package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	verbose := flag.Bool("v", false, "enable/disable verbose mode") // HLflag
	flag.Parse()                                                    // HLflag

	sum := 0
	if *verbose { // HLflag
		fmt.Printf("%v\n", sum)
	}
	for _, arg := range flag.Args() { // HLflag
		v, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error** %v\n", err)
			os.Exit(1)
		}
		sum += v

		if *verbose { // HLflag
			fmt.Printf("+ %v -> %v\n", v, sum)
		}
	}
	if *verbose { // HLflag
		fmt.Printf("===============\n")
	}
	fmt.Printf("sum= %v\n", sum)
}
