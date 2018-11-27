package main

import (
	"encoding/json"
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
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
	Pages   int      `json:"pages"`
}

func main() {
	fmt.Printf("please connect to localhost:7777\n")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/book/", bookHandler)
	http.HandleFunc("/create", createHandler)
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	err := enc.Encode(books)
	if err != nil {
		log.Printf("error encoding JSON: %v\n", err)
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	path := html.EscapeString(r.URL.Path[len("/book/"):])
	if path == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "[]\n")
		return
	}

	idx, err := strconv.Atoi(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "[]\n")
		return
	}

	if idx >= len(books) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "[]\n")
		return
	}

	err = json.NewEncoder(w).Encode(books[idx])
	if err != nil {
		log.Printf("error encoding book[%d]: %v\n", idx, err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid http request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Printf("error decoding book: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	mu.Lock()
	books = append(books, book)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}
