package main

import (
	"log"

	"github.com/sbinet-vgo/api/v2"
)

func main() {
	log.Printf("%s", api.Welcome("Bob!"))
}
