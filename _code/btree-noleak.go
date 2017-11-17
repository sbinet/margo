package main

import (
	"fmt"
	"log"

	"golang.org/x/tour/tree"
)

// STARTWALK OMIT

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	quit := make(chan int)
	defer close(quit)
	walk(t, ch, quit)
	close(ch)
}

func walk(t *tree.Tree, ch, quit chan int) {
	if t == nil {
		return
	}
	walk(t.Left, ch, quit)
	select {
	case ch <- t.Value:
	case <-quit:
		return
	}
	walk(t.Right, ch, quit)
}

// ENDWALK OMIT

// STARTSAME OMIT

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	quit := make(chan int)
	ch1 := make(chan int)
	ch2 := make(chan int)
	defer close(quit)

	go func() { walk(t1, ch1, quit); close(ch1) }()
	go func() { walk(t2, ch2, quit); close(ch2) }()

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2
		if !ok1 || !ok2 {
			return ok1 == ok2
		}
		if v1 != v2 {
			return false
		}
	}
}

// ENDSAME OMIT

func main() {
	ch := make(chan int)
	go Walk(tree.New(1), ch)
	for i := range ch {
		fmt.Printf("%d\n", i)
	}

	ok1 := Same(tree.New(1), tree.New(1))
	if !ok1 {
		log.Fatalf("FAILED ok1=%v", ok1)
	}

	ok2 := Same(tree.New(1), tree.New(2))
	if ok2 {
		log.Fatalf("FAILED: ok2=%v", ok2)
	}
	fmt.Printf("OK\n")
}
