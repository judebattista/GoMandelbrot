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

//find out how far away from the origin our complex coordinate is
func magnitude(arg complex128) (mag float64) {
	mag = math.Sqrt(real(arg)*real(arg) + imag(arg)*imag(arg))
	return
}

func magnitudeSquared(arg complex128) (mag float64) {
	mag = real(arg)*real(arg) + imag(arg)*imag(arg)
	return
}

type data_point struct {
	coordinate  complex128
	converges   bool
	iterations  int
	zoom_levels []int
}

type frame struct {
	m  map[complex128]data_point
	id int
}

type gif struct {
	num_frames int
	frames     []frame
}

//x2 = x1^2 + arg
//Default seed to zero, setting seed equal to arg is another valid approach
func checkConvergence(arg complex128, seed complex128, maxIterations int) (converges bool, iterations int) {
	currentTerm := seed
	var lastTerm complex128
	for ndx := 0; ndx < maxIterations; ndx++ {
		lastTerm = currentTerm
		currentTerm = (lastTerm * lastTerm) + arg
		//If we ever find ourselves more than 2 from the origin, we diverge
		//if magnitude(currentTerm) > 2 {
		//But... we can save a step if we square both sides!
		if magnitudeSquared(currentTerm) > 4 {
			iterations = ndx + 1
			converges = false
			return
		}
	}
	iterations = -1
	converges = true
	return
}

func calculator(to_calculate chan data_point, calculated chan data_point, finished chan bool, max_iterations int) {
	for current_point := range to_calculate {
		current_point.converges, current_point.iterations = checkConvergence(current_point.coordinate, 0+0i, max_iterations)
		calculated <- current_point
	}
	finished <- true
}

func collector(calculated chan data_point, finished gif, completed chan bool) {
	for current_point := range calculated {
		for _, value := range current_point.zoom_levels {
			//Insert comment here
			//targetMap := finished.frames[value-1].m
			if finished.frames[value-1].m == nil {
				finished.frames[value-1].m = make(map[complex128]data_point)
			}
			finished.frames[value-1].m[current_point.coordinate] = current_point
		}
	}
	completed <- true
}

func main() {
	//determine zoom
	//determine set of points
	//find any previously calculated points
	//check each remaining point for convergence
	//write the zoom level to an image

	starting_coordinate := -0.7463 + 0.1102i
	a := real(starting_coordinate)
	b := imag(starting_coordinate)
	//Everything interesting happens between -2 and 2 on both axes
	//If the starting coordinate is not 0+0i, the offset needs to be changed to include that window
	biggest_coord_offset := float64(.01)

	//IMPORTANT: If these values change, they must also be changed in the python script
	//TODO: Put common values in a config file and read it into both Go and Python scripts
	frame_dimension := float64(1024)
	number_frames := float64(10)
	max_iterations := 100

	zoom_factor := (2 * biggest_coord_offset) / frame_dimension

	gif := gif{int(number_frames), make([]frame, int(number_frames))}

	to_be_calculated := make(map[complex128]data_point)

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
