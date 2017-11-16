package main

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

// START OMIT

func main() {
	r, err := gzip.NewReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		log.Fatal(err)
	}
}

// END OMIT
