package main

/*
	real values: -2.0 to 1.0
	imag values: -1.5 to 1.5

*/
import (
	"math"
	"math/cmplx"
)

//find out how far away from the origin our complex coordinate is
func magnitude(arg complex128) (mag float64) {
	mag = math.Sqrt(real(arg)*real(arg) + imag(arg)*imag(arg))
	return
}

//x2 = x1^2 + arg
//Default seed to zero, setting seed equal to arg is another valid approach
func checkConvergence(arg complex128, seed complex128, maxIterations int, epsilon float64) (converges bool, iterations int) {
	var delta float64
	currentTerm := seed
	var lastTerm complex128
	for ndx := 0; ndx < maxIterations; ndx++ {
		lastTerm = currentTerm
		currentTerm = (lastTerm * lastTerm) + arg
		//This is probably not the right approach
		delta = cmplx.Abs(currentTerm - lastTerm)
		if delta < epsilon {
			iterations = ndx + 1
			converges = true
			return
		}
	}
	iterations = maxIterations + 1
	converges = false
	return
}

func main() {

}
