package main

/*
	real values: -2.0 to 1.0
	imag values: -1.5 to 1.5
	If the magnitude of a given iteration for point c is greater than 2, then the sequence tend to infinity
*/
import (
	"math"
)

//find out how far away from the origin our complex coordinate is
func magnitude(arg complex128) (mag float64) {
	mag = math.Sqrt(real(arg)*real(arg) + imag(arg)*imag(arg))
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
		//This is probably not the right approach
		if magnitude(currentTerm) > 2 {
			iterations = ndx + 1
			converges = false
			return
		}
	}
	converges = true
	return
}

func main() {
	//determine zoom
	//determine set of points
	//find any previously calculated points
	//check each remaining point for convergence
	//write the zoom level to an image

	number_frames := 30

	starting_coordinate := 0 + 0i
	zoom_factor := 2
	fram_dimension := 256

	gif := gif{number_frames, []frame{}}

	to_be_calculated := make(map[complex128]data_point)

	for i := 0; i < number_frames; i++ {

	}
}
