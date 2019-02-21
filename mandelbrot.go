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
//Takes:
//	arg: A set of complex coordinates
//Returns:
//	mag: The magnitude of that point
func magnitude(arg complex128) (mag float64) {
	mag = math.Sqrt(real(arg)*real(arg) + imag(arg)*imag(arg))
	return
}

//Since the ultimate purpose is to compare the magnitude to 2, what if we just compared the square of the magnitude to 4?
//Saves a square root call, which should provide some benefit over a million points or so.
//Takes:
//	arg: A set of complex coordinates
//Returns:
//	mag: The square of the magnitude of that point
func magnitudeSquared(arg complex128) (mag float64) {
	mag = real(arg)*real(arg) + imag(arg)*imag(arg)
	return
}

//Core Mandelbrot calculation: x2 = x1^2 + arg
//Default seed should be zero, setting seed equal to arg is another valid approach
//Takes:
//	arg: the complex coordinates of the point,
//	seed: the seed for the calculation which should generally be zero. Setting the see equal to the arg is another valid approach
//	maxIterations: the maximum number of iterations. If we exceed this without diverging we assume convergence. Raising this value increases the resolution along the edges of the set
//Returns:
//	converges: whether or not the point converges
//	iterations:	If the point diverges, it returns the number of iterations required to establish divergence
//				If the point converges, it sets the iterations to -1 to distinguish it from points that were never calculated which have iterations of 0
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
//After a point is calculated, deposit it in the calculated channel
//When finished it announces on the finished channel
//Takes:
//	to_calculate: a channel from which it consumes uncalculated data points
//	calculated: a channel on which to send calculated data points
//	finished: a channel on which to announce completion
//	max_iterations: the maximum number of iterations to attempt before concluding convergence
//Returns:
//	Nothing
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
//Takes:
//	calculated: a channel from which it consumes calculated data points
//	finished: a gif to which it writes the calculated data points
//	completed: a channel on which to announce completion
//Returns:
//	Nothing
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

//Find all the points in every frame of our gif
//Takes:
//	to_be_calculated: a map from coordinates to data points
//	starting_coordinate: the center coordinate of the region to examine
//	number_frames: the number of frames for which points need to be found
//	zoom_factor: the zoom factor between frames
//	frame_dimension: the length of one side of the frame, which is assumed to be square
func callingAllPoints(to_be_calculated map[complex128]data_point, starting_coordinate complex128, number_frames float64, zoom_factor float64, frame_dimension float64) {
	//Split the complex coordinate a+bi into its real and imaginary coefficients
	a := real(starting_coordinate)
	b := imag(starting_coordinate)
	//For each frame in the gif
	for i := float64(1); i <= number_frames; i++ {
		//figure out how far apart our points are
		x_offset := float64(zoom_factor * math.Pow(.9, i-1))
		//x_offset and y_offset should always be the same, but we're leaving both in just in case
		//Note that if they ever differ, the python script will need to be revised
		y_offset := x_offset
		//For every x value in the frame...
		x := a - frame_dimension*x_offset/2
		for foo := float64(0); foo < frame_dimension; foo++ {
			//find all the y-values in the frame
			y := b - frame_dimension*y_offset/2
			for bar := float64(0); bar < frame_dimension; bar++ {
				//for every x,y point add this zoom level to the list of zoom levels that use it
				data := to_be_calculated[complex(float64(x), float64(y))]
				data.zoom_levels = append(data.zoom_levels, int(i))
				data.coordinate = complex(float64(x), float64(y))
				//add the point to the map
				to_be_calculated[complex(float64(x), float64(y))] = data
				y += y_offset
			}
			x += x_offset
		}
	}
}

//Take all the points in the gif object and write them to a series of csv files, each representing one frame of the gif
//These files will be written to a gif by our test.py script
//Takes:
//	gif: A gif object containing all the data to be written to files
//Returns:
//	Nothing
func writeGifFiles(gif gif) {
	for i, v := range gif.frames {
		//Try and create each file
		file_name := fmt.Sprintf("frame%02d.txt", i)
		file, err := os.Create(file_name)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		//Dump all the relevant points into the file
		for _, point := range v.m {
			fmt.Fprintf(file, "%v, %v, %v\n", real(point.coordinate), imag(point.coordinate), point.iterations)
		}
		file.Close()
	}
}

//Create all the go routines responsible for calculating individual points
//Takes:
//	num_threads: The number of go routines to create
//	to_calculate: The channel the go routines will pull uncalculated points from
// 	calculated: The channel the go routines will send calculated points too
//	finished: The channel on which each go routine will announce completion
//	max_iterations: The maximum number of iterations over a point before we assume convergence
func spinUpCalculators(num_threads int, to_calculate chan data_point, calculated chan data_point, finished chan bool, max_iterations int) {
	for i := 0; i < num_threads; i++ {
		go calculator(to_calculate, calculated, finished, max_iterations)
	}
}

//Create our collector go routine responsible for taking all the calculated data points and writing them to a map based on their coordinates
//Note that since collector writes to a map, there should only be a single instance
//Takes:
//	calculated: The channel the go routine will pull calculated points from
//	gif: The gif which holds the map of points to which the go routine will write the calculated points
//	finished: The channel on which the go routine will announce completion
func spinUpCollector(calculated chan data_point, gif gif, finished chan bool) {
	go collector(calculated, gif, finished)
}

//Dump all the points from to_be_calculated into to_calculate
//Takes:
//	to_be_calculated: The map from coordinates to data points that holds all the newly initialized data points
//	to_calculate: The channel into which the data points will be dumped
func feedCalculators(to_be_calculated map[complex128]data_point, to_calculate chan data_point) {
	for _, v := range to_be_calculated {
		to_calculate <- v
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

	//IMPORTANT: If the following three values change, they must also be changed in the python script
	//	frame_dimension, number_frames, max_iterations
	//TODO: Put common values in a config file and read it into both Go and Python scripts
	//How many points are in each frame?
	frame_dimension := float64(1024)
	//How many frames are in the gif?
	number_frames := float64(30)
	//How many times will we iterate over a point before deciding that it diverges.
	//The higher the value, the more detail is revealed along the fractal edges.
	max_iterations := 1000

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
	//Channel to supply the points to the calculators
	to_calculate := make(chan data_point)
	//Channel to receive the points after their convergence has been calculated
	calculated := make(chan data_point)
	//Channel to announce completion
	finished := make(chan bool)

	//IMPORTANT: Spin up the listeners BEFORE you start using the channels
	//This was the source of our perplexing deadlock
	//Spin up all of our calculators
	spinUpCalculators(num_threads, to_calculate, calculated, finished, max_iterations)

	//As the calculators dump points into calculated, collector will pull them out and write them to gif
	spinUpCollector(calculated, gif, finished)

	//Start feeding the calculators by dumping the data points from the map to the channel
	//Each point dumped into the to_calculate channel will be pulled by one of the calculators and evaluated for divergence.
	feedCalculators(to_be_calculated, to_calculate)

	//We're fresh outta input, so close the channel
	//Once it's empty the calculators will all terminate automatically thanks to the range syntax
	close(to_calculate)

	//Wait for every calculator to send it's finished message
	for i := 0; i < num_threads; i++ {
		<-finished
	}

	//Normally you do not want to close a channel outside of that channel's sender
	//In this case, since we have multiple senders, none of them can authoritatively close it
	//So we have to wait for all of them the announce they are finished and then close it in main
	//close the calculated channel since no more inputs are coming.
	close(calculated)

	//Wait for one final finished message, presumably from collector
	//This could be risky since we are not guaranteed that all the calculators will finish and announce before collector
	//In practice, since all the calculators feed the channel that collector depends on, it is a reasonable risk.
	<-finished

	//Output the csv files that will be turned into a beeeyootiful gif. Or garbage. Pick your parameters wisely!
	writeGifFiles(gif)
}
