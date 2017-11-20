package main

import (
	"flag"
	"fmt"
)

func main() {
	nn := flag.Int("xy", 2, "NxN size of the grid")
	flag.Parse()

	nx := *nn
	ny := *nn

	fmt.Printf("n=%d\n", compute(nx, ny))
}

func compute(nx, ny int) int {
	// ...
	return 0
}
