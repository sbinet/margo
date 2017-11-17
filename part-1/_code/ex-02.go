package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	sum := 0
	for _, arg := range os.Args[1:] {
		v, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error** %v\n", err)
			os.Exit(1)
		}
		sum += v
	}
	fmt.Printf("sum= %v\n", sum)
}
