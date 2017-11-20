package main

import (
	"flag"
	"fmt"
	"sync"
)

var (
	sum = 0
	mu  sync.RWMutex
)

func main() {
	nn := flag.Int("xy", 2, "xy")
	flag.Parse()
	nx := *nn
	ny := *nn
	fmt.Printf("n=%d\n", compute(nx, ny))
}

func compute(nx, ny int) int {
	var pt Point
	var wg sync.WaitGroup
	wg.Add(2)

	r := Point{pt.x + 1, pt.y}
	d := Point{pt.x, pt.y + 1}
	go inspect(r, &wg, nx, ny)
	go inspect(d, &wg, nx, ny)

	wg.Wait()
	return sum
}

type Point struct {
	x, y int
}

func inspect(pt Point, wg *sync.WaitGroup, nx, ny int) {
	defer wg.Done()
	if pt.x < nx {
		r := Point{pt.x + 1, pt.y}
		wg.Add(1)
		go inspect(r, wg, nx, ny)
	}
	if pt.y < ny {
		d := Point{pt.x, pt.y + 1}
		wg.Add(1)
		go inspect(d, wg, nx, ny)
	}

	if pt.y == ny && pt.x == nx {
		mu.Lock()
		sum++
		mu.Unlock()
	}
}
