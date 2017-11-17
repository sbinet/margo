package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	v := 1.0
	for i := 0; i < 10; i++ {
		v -= (v*v - x) / (2 * v)
	}
	return v
}

func main() {
	fmt.Println(Sqrt(2))
}
