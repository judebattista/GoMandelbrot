package main

import (
	"fmt"
	"math"
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

func main() {
	//determine zoom
	//determine set of points
	//find any previously calculated points
	//check each remaining point for convergence
	//write the zoom level to an image

	number_frames := float64(5)

	starting_coordinate := 0 + 0i
	a := real(starting_coordinate)
	b := imag(starting_coordinate)
	zoom_factor := 0.5
	frame_dimension := float64(256)

	//gif := gif{number_frames, []frame{}}

	to_be_calculated := make(map[complex128]data_point)

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

	for k, v := range to_be_calculated {
		if (k == 0+0i) || (k == 0+.5i) || (k == 0+1i) || (k == 0+1.5i) || (k == 0+2i) || (k == 0+2.5i) {
			fmt.Println(k, v)
		}
	}

	testCheckConvergence()
}
