package main

// #include "my-cgo-lib.h"
import "C"

func main() {
	C.Lib()
}
