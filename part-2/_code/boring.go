// +build OMIT

package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	boring("boring!")
}

// START OMIT
func boring(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Second)
	}
}

// STOP OMIT

func init() {
	// hack to make this program runnable on the playground.
	go func() {
		time.Sleep(10 * time.Second)
		log.Fatalf("time's up")
	}()
}
