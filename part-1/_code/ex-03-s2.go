package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	v := 1.0
	old := 1.0
	delta := 1e10
	for delta > 1e-12 {
		v, old = v-(v*v-x)/(2*v), v
		delta = math.Abs(old - v)
	}
	return v
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(math.Sqrt(2))
}
