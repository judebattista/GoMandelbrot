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
	If the magnitude of a given iteration for point c is greater than 2, then the sequence tend to infinity
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
	iterations = 0
	converges = true
	return
}

func testCheckConvergence() {
	for foo := 0.0; foo < 2; foo += 0.2 {
		for bar := 0.0; bar < 2; bar += 0.2 {
			var coords = complex(foo, bar)
			converges, iterations := checkConvergence(coords, 0, 10)
			if converges {
				fmt.Printf("%v + %vi converges: %v after %v iterations.\n", real(coords), imag(coords), converges, iterations)
			}
		}
	}
}

func calculator(to_calculate chan data_point, calculated chan data_point, finished chan bool) {
	for current_point := range to_calculate {
		current_point.converges, current_point.iterations = checkConvergence(current_point.coordinate, 0+0i, 2000)
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
	//testCheckConvergence()

	number_frames := float64(1)

	starting_coordinate := -.170337 + -1.06506i
	a := real(starting_coordinate)
	b := imag(starting_coordinate)
	zoom_factor := 0.001
	frame_dimension := float64(4096)

	gif := gif{int(number_frames), make([]frame, int(number_frames))}

	to_be_calculated := make(map[complex128]data_point)

	//Can we parallelize this?
	for i := float64(1); i <= number_frames; i++ {
		radius := float64((frame_dimension * i * zoom_factor) / 2)
		x_offset := float64(i * zoom_factor)
		y_offset := float64(i * zoom_factor)
		for x := a - (radius); x < a+(radius-1); x += x_offset {
			for y := b - (radius); y < b+(radius-1); y += y_offset {
				data := to_be_calculated[complex(float64(x), float64(y))]
				data.zoom_levels = append(data.zoom_levels, int(i))
				data.coordinate = complex(float64(x), float64(y))
				to_be_calculated[complex(float64(x), float64(y))] = data
			}
		}
	}

	num_threads := 8
	to_calculate := make(chan data_point)
	calculated := make(chan data_point)
	finished := make(chan bool)

	//Spin up go routines first
	//If a channel does not have listeners, we get a deadlock
	for i := 0; i < num_threads; i++ {
		go calculator(to_calculate, calculated, finished)
	}

	go collector(calculated, gif, finished)

	//Dump everything from to_be_calculated into the to_calculate channel
	for _, v := range to_be_calculated {
		to_calculate <- v
	}

	//Close the to_calculate channel in order to permit the range syntax to function
	close(to_calculate)

	//Collect one finish from each of the calculator threads
	for i := 0; i < num_threads; i++ {
		<-finished
	}

	close(calculated)

	//Collect a final finish from the collector thread
	<-finished

	//Write the data sets to files
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
