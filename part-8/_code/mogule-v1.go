package main

import (
	"log"

	"github.com/sbinet-vgo/api"
)

func main() {
	var g api.Greeter
	log.Printf("%s", g.Greet("Bob"))
}
