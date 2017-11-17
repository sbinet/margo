package main

import (
	"fmt"
	"math"
)

// STARTABSER OMIT

type Abser interface {
	Abs() float64
}

func main() {
	var a Abser
	f := MyFloat(-math.Sqrt2)
	v := Vertex{3, 4}

	a = f                // a MyFloat implements Abser
	fmt.Println(a.Abs()) // =1.4142135623730951

	a = &v               // a *Vertex implements Abser
	fmt.Println(a.Abs()) // =5
}

type MyFloat float64

type Vertex struct {
	X, Y float64
}

// ENDABSER OMIT

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
