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
	iterations = maxIterations + 1
	converges = true
	return
}

func main() {
	//determine zoom
	//determine set of points
	//find any previously calculated points
	//check each remaining point for convergence
	//write the zoom level to an image
}
