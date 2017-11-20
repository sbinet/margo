package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
)

var (
	nx  = 3
	ny  = 3
	exp = 0
)

var (
	sum = 0
	mu  sync.RWMutex
)

func main() {
	nnx := flag.Int("xy", nx, "xy")
	flag.Parse()
	nx = *nnx
	ny = *nnx
	exp = (nx + 1) * (ny + 1) * 2
	log.Printf("goroutines: %v\n", exp)
	n := start()
	fmt.Printf("n=%d\n", n)
}

type Point struct {
	x, y int
}

func start() int {
	var pt Point
	var wg sync.WaitGroup
	wg.Add(2)

	right := Point{pt.x + 1, pt.y}
	left := Point{pt.x, pt.y + 1}
	go inspect(right, &wg)
	go inspect(left, &wg)

	wg.Wait()
	return sum
}

func inspect(pt Point, wg *sync.WaitGroup) {
	defer wg.Done()
	if pt.x < nx {
		right := Point{pt.x + 1, pt.y}
		wg.Add(1)
		go inspect(right, wg)
	}
	if pt.y < ny {
		left := Point{pt.x, pt.y + 1}
		wg.Add(1)
		go inspect(left, wg)
	}

	if pt.y == ny && pt.x == nx {
		mu.Lock()
		sum++
		mu.Unlock()
	}
}
