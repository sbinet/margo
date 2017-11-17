// +build OMIT

package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func init() {
	// hack to make this program runnable on the playground.
	go func() {
		time.Sleep(10 * time.Second)
		log.Fatalf("time's up")
	}()
}

func main() {
	boring("boring!") // HL
}

// START OMIT
func boring(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

// STOP OMIT
