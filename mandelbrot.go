package main

import (
	"fmt"
	"log"
	"math"
	"os"
)

/*
	real values: -2.0 to 1.0
	imag values: -1.5 to 1.5
	If the magnitude of a given iteration for point c is greater than 2, then the sequence tends to infinity
*/

//Represents a complex coordinate along with the mandelbrot-relevant info
//The zoom_levels slice will contain a list of all the zoom levels that include the point
type data_point struct {
	coordinate  complex128
	converges   bool
	iterations  int
	zoom_levels []int
}

//Represents one frame of a gif
//Maps a set of coordinates to the corresponding data_point and tracks its order in the gif
type frame struct {
	m  map[complex128]data_point
	id int
}

//Represents an animated gif. Has the data for each of its frames and the total number of frames
type gif struct {
	num_frames int
	frames     []frame
}

//find out how far away from the origin our complex coordinate is
func magnitude(arg complex128) (mag float64) {
	mag = math.Sqrt(real(arg)*real(arg) + imag(arg)*imag(arg))
	return
}

//Since the ultimate purpose is to compare the magnitude to 2, what if we just compared the square of the magnitude to 4?
//Saves a square root call, which should provide some benefit over a million points or so.
func magnitudeSquared(arg complex128) (mag float64) {
	mag = real(arg)*real(arg) + imag(arg)*imag(arg)
	return
}

//Core Mandelbrot calculation: x2 = x1^2 + arg
//Default seed should be zero, setting seed equal to arg is another valid approach
func checkConvergence(arg complex128, seed complex128, maxIterations int) (converges bool, iterations int) {
	//Start at the seed value
	currentTerm := seed
	var lastTerm complex128
	for ndx := 0; ndx < maxIterations; ndx++ {
		//Since the sequence is defined recursively we need to store the current term before we calculate the next one
		lastTerm = currentTerm
		//Calculate the next term
		currentTerm = (lastTerm * lastTerm) + arg
		//If we ever find ourselves more than 2 from the origin, we diverge
		//if magnitude(currentTerm) > 2 {
		//But... we can save a step if we square both sides!
		//If it diverges we want to know how long it took to diverge. This represents how close it is to being in the set
		if magnitudeSquared(currentTerm) > 4 {
			iterations = ndx + 1
			converges = false
			return
		}
	}
	//If the point converges, then we set the iterations to -1
	//This is to differentiate them from uncalculated points which have 0 iterations
	iterations = -1
	converges = true
	return
}

//Go routine to calculate the convergence of a point.
//It pulls from the to_calculate channel until said channel is empty and closed
func calculator(to_calculate chan data_point, calculated chan data_point, finished chan bool, max_iterations int) {
	//Get points from the channel until it is empty and closed
	for current_point := range to_calculate {
		//Set the mandelbrot info for the point
		current_point.converges, current_point.iterations = checkConvergence(current_point.coordinate, 0+0i, max_iterations)
		//Send the point to the calculated channel for collection
		calculated <- current_point
	}
	//Once the to_calculate channel runs try, let finished know that this go routine is done.
	finished <- true
}

//Function to collect the calculated points from the channel and store them in a map
//Since it writes to a map, this should not be parallelized
func collector(calculated chan data_point, finished gif, completed chan bool) {
	for current_point := range calculated {
		for _, value := range current_point.zoom_levels {
			//We could not find an easy way to automatically initialize the map on creation,
			//so we need a nil check before we access it
			if finished.frames[value-1].m == nil {
				finished.frames[value-1].m = make(map[complex128]data_point)
			}
			//store the calculated point
			finished.frames[value-1].m[current_point.coordinate] = current_point
		}
	}
	completed <- true
}

func callingAllPoints(to_be_calculated map[complex128]data_point, starting_coordinate complex128, number_frames float64, zoom_factor float64, frame_dimension float64) {
	a := real(starting_coordinate)
	b := imag(starting_coordinate)
	for i := float64(1); i <= number_frames; i++ {
		x_offset := float64(zoom_factor * math.Pow(.9, i-1))
		//x_offset and y_offset should always be the same, but we're leaving both in just in case
		//Note that if they ever differ, the python script will need to be revised
		y_offset := x_offset
		x := a - frame_dimension*x_offset/2
		for foo := float64(0); foo < frame_dimension; foo++ {
			y := b - frame_dimension*y_offset/2
			for bar := float64(0); bar < frame_dimension; bar++ {
				data := to_be_calculated[complex(float64(x), float64(y))]
				data.zoom_levels = append(data.zoom_levels, int(i))
				data.coordinate = complex(float64(x), float64(y))
				to_be_calculated[complex(float64(x), float64(y))] = data
				y += y_offset
			}
			x += x_offset
		}
	}
}

func main() {
	//determine zoom
	//determine set of points
	//find any previously calculated points
	//check each remaining point for convergence
	//write the zoom level to an image

	starting_coordinate := -0.7463 + 0.1102i
	//Everything interesting happens between -2 and 2 on both axes
	//If the starting coordinate is not 0+0i, the offset needs to be changed to include that window
	biggest_coord_offset := float64(.01)

	//IMPORTANT: If these values change, they must also be changed in the python script
	//TODO: Put common values in a config file and read it into both Go and Python scripts
	//How many points are in each frame?
	frame_dimension := float64(1024)
	//How many frames are in the gif?
	number_frames := float64(10)
	//How many times will we iterate over a point before deciding that it diverges.
	//The higher the value, the more detail is revealed along the fractal edges.
	max_iterations := 100

	//How much will we change our resolution per frame?
	zoom_factor := (2 * biggest_coord_offset) / frame_dimension

	//Create the object that represents our gif
	gif := gif{int(number_frames), make([]frame, int(number_frames))}

	//map to hold all the points prior to calculation
	to_be_calculated := make(map[complex128]data_point)

	//
	callingAllPoints(to_be_calculated, starting_coordinate, number_frames, zoom_factor, frame_dimension)

	//how many threads does your processor support? Probably roughly twice the number of cores
	num_threads := 8
	to_calculate := make(chan data_point)
	calculated := make(chan data_point)
	finished := make(chan bool)

	for i := 0; i < num_threads; i++ {
		go calculator(to_calculate, calculated, finished, max_iterations)
	}

	go collector(calculated, gif, finished)

	for _, v := range to_be_calculated {
		to_calculate <- v
	}

	close(to_calculate)

	for i := 0; i < num_threads; i++ {
		<-finished
	}

	close(calculated)

	<-finished

	for i, v := range gif.frames {
		file_name := fmt.Sprintf("frame%02d.txt", i)
		file, err := os.Create(file_name)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		for _, point := range v.m {
			fmt.Fprintf(file, "%v, %v, %v\n", real(point.coordinate), imag(point.coordinate), point.iterations)
		}
		file.Close()
	}
}
