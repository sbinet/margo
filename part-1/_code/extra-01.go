package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fname := os.Args[1]
	f, err := os.Open(fname) // HLxxx
	if err != nil {
		log.Fatalf("could not open %q: %v\n", fname, err)
	}
	defer f.Close() // HLxxx

	words := 0
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := strings.Split(scan.Text(), " ")
		for _, word := range line {
			if word != "" {
				words++
			}
		}
	}

	err = scan.Err() // HLxxx
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error: %v\n", err)
	}
	fmt.Printf("%d\n", words)
}
