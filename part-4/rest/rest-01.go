package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var (
	mu    sync.RWMutex
	books = []Book{
		{
			Title:   "The Go Programming Language",
			Authors: []string{"Alan A. A. Donovan", "Brian W. Kernighan"},
			Pages:   380,
		},
		{
			Title:   "Advanced Programming in the UNIX Environment",
			Authors: []string{"W. Richard Stevens", "Stephen A. Rago"},
			Pages:   1024,
		},
		{
			Title:   "The Practice of Programming",
			Authors: []string{"Brian W. Kernighan", "Rob Pike"},
			Pages:   267,
		},
		{
			Title:   "The C Programming Language",
			Authors: []string{"Brian W. Kernighan", "Dennis Ritchie"},
			Pages:   274,
		},
	}
)

type Book struct {
	Title   string
	Authors []string
	Pages   int
}

func main() {
	fmt.Printf("please connect to localhost:7777\n")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/book/", bookHandler)
	log.Fatal(http.ListenAndServe(":7777", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	fmt.Fprintf(w, "<h1>Welcome to the Library</h1>\nWe have %d books.\n", len(books))
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	fmt.Fprintf(w, "<h1>Book index</h1>\n")

	fmt.Fprintf(w, "<ol>\n")
	for _, book := range books {
		fmt.Fprintf(w, "\t<li>%s</li>\n", html.EscapeString(book.Title))
	}
	fmt.Fprintf(w, "</ol>\n")
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	path := html.EscapeString(r.URL.Path[len("/book/"):])
	switch path {
	case "":
		fmt.Fprintf(w, "Enter a book index. (no index given)\n")
	default:
		idx, err := strconv.Atoi(path)
		if err != nil {
			log.Printf("incorrect index %q: %v\n", path, err)
			http.Error(w, "invalid index '"+path+"'", http.StatusNotFound)
			return
		}
		if idx >= len(books) {
			log.Printf("incorrect index %d (>=%d)\n", idx, len(books))
			http.Error(w, "incorrect index "+path, http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%#v\n", books[idx])
	}
}
