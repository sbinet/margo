package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// STARTITEM OMIT
type Word struct {
	word  string
	count int
}

type Words []Word

func (p Words) Len() int { return len(p) }

func (p Words) Less(i, j int) bool {
	ii := p[i]
	jj := p[j]
	switch {
	case ii.count == jj.count:
		return ii.word < jj.word
	}
	return ii.count < jj.count
}

func (p Words) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// ENDITEM OMIT

func main() {
	fname := os.Args[1]
	f, err := os.Open(fname) // HLxxx
	if err != nil {
		log.Fatalf("could not open %q: %v\n", fname, err)
	}
	defer f.Close() // HLxxx

	// START OMIT
	nwords := 0
	nlines := 0
	stats := make(map[string]int)

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := strings.Split(scan.Text(), " ")
		nlines++ // HLxxx
		for _, word := range line {
			if word != "" {
				nwords++
				stats[word]++ // HLxxx
			}
		}
	}
	// END OMIT

	err = scan.Err() // HLxxx
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error: %v\n", err)
	}
	fmt.Printf("#lines: %d\n", nlines)
	fmt.Printf("#words: %d\n", nwords)

	// STARTSORT OMIT
	words := make(Words, 0, len(stats))
	for w, n := range stats {
		words = append(words, Word{word: w, count: n}) // HLxxx
	}
	sort.Sort(sort.Reverse(words)) // HLxxx

	fmt.Printf("\npopcon:\n")
	for i, v := range words[:5] {
		fmt.Printf("#%d: %q (%d)\n", i+1, v.word, v.count)
	}
	// ENDSORT OMIT
}
