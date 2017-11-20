package main

import "C"

//export Add
func Add(i, j int) int {
	return i + j
}

func main() {}
