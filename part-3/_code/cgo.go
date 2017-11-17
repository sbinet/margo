package main

// #cgo LDFLAGS: -lm
// #include <math.h>
import "C"

func main() {
	println("C.sqrt(2) = ", C.sqrt(2))
}
