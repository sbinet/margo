# Introduction to REST API with Go

You know how to create web servers in `Go`.

In this hands-on session, we'll see how one can create a simple web server
in `Go` exposing a REST API to a book database.

Let's say we have a nifty database of books we own at home.
For simplicity's sake, let's say it is a simple:

```go
var books = []Book{
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

type Book struct {
	Title   string
	Authors []string
	Pages   int
}
```

We don't have many books, but we do have good taste.

Ok, what should our (JSON) REST API web server do?
It should probably expose these endpoints:

- `"/"` displays a welcome message and the number of books in the library
- `"/books"` displays a listing of all the books in the library
- `"/book/<index>"` displays a particular book
- `"/create"` creates a new book entry in the library

Let's write this basic structure.

```go
var mu sync.RWMutex // to protect access to the books database

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
```

Run this server in a terminal and then, in another terminal:

```sh
$> curl localhost:7777
<h1>Welcome to the Library</h1>
We have 4 books.

$> curl localhost:7777/books
<h1>Book index</h1>
<ol>
	<li>The Go Programming Language</li>
	<li>Advanced Programming in the UNIX Environment</li>
	<li>The Practice of Programming</li>
	<li>The C Programming Language</li>
</ol>

$> curl localhost:7777/book/
Enter a book index. (no index given)

$> curl localhost:7777/book/0
main.Book{Title:"The Go Programming Language", Authors:[]string{"Alan A. A. Donovan", "Brian W. Kernighan"}, Pages:380}

$> curl localhost:7777/book/4
incorrect index 4

$> curl localhost:7777/book/foo
invalid index 'foo'

```

Ok, the basic functionalities are here, with a basic model of `Book`.

## Sending JSON

Now, instead of printing an enumerated list of books, let's have the `"/books"`
endpoint print back the whole list as JSON.

It's actually rather easy:

```go
func booksHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ") // for pretty printing
	err := enc.Encode(books)
	if err != nil {
		log.Printf("error encoding JSON: %v\n", err)
	}
}
```

Now, requesting `"/books"` gives back:

```sh
$> curl localhost:7777/books
[
 {
  "Title": "The Go Programming Language",
  "Authors": [
   "Alan A. A. Donovan",
   "Brian W. Kernighan"
  ],
  "Pages": 380
 },
 {
  "Title": "Advanced Programming in the UNIX Environment",
  "Authors": [
   "W. Richard Stevens",
   "Stephen A. Rago"
  ],
  "Pages": 1024
 },
 {
  "Title": "The Practice of Programming",
  "Authors": [
   "Brian W. Kernighan",
   "Rob Pike"
  ],
  "Pages": 267
 },
 {
  "Title": "The C Programming Language",
  "Authors": [
   "Brian W. Kernighan",
   "Dennis Ritchie"
  ],
  "Pages": 274
 }
]
```

That's nice, but we can do better: this doesn't look very JSON idiomatic.
Usually, JSON keys are not uppercased.
This is easily fixed by using `struct-tags`:

```go
type Book struct {
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
	Pages   int      `json:"pages"`
}
```

Now, requesting `"/books"` prints back:

```sh
$> curl localhost:7777/books
[
 {
  "title": "The Go Programming Language",
  "authors": [
   "Alan A. A. Donovan",
   "Brian W. Kernighan"
  ],
  "pages": 380
 },
 {
  "title": "Advanced Programming in the UNIX Environment",
  "authors": [
   "W. Richard Stevens",
   "Stephen A. Rago"
  ],
  "pages": 1024
 },
 {
  "title": "The Practice of Programming",
  "authors": [
   "Brian W. Kernighan",
   "Rob Pike"
  ],
  "pages": 267
 },
 {
  "title": "The C Programming Language",
  "authors": [
   "Brian W. Kernighan",
   "Dennis Ritchie"
  ],
  "pages": 274
 }
]
```

One more thing.
We should declare the correct `"Content-Type"` for our response's header,
so clients can expect JSON:

```go
func booksHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// ...
}
```

Let's also modify `"/book/<index>"` to print out JSON:

```go
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
```

## Create a new book entry

Let's add the new and final `"/create"` endpoint to add new book entries:

```go
func main() {
	fmt.Printf("please connect to localhost:7777\n")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/book/", bookHandler)
	http.HandleFunc("/create", createHandler)
	log.Fatal(http.ListenAndServe(":7777", nil))
}
```

and the `createHandler`:

```go
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
```

Run the server and, in another terminal:

```sh
$> curl -d '{"title":"Effective Go", "pages":50, "authors": ["Go core team"]}' localhost:7777/create

$> curl localhost:7777/books
[
 {
  "title": "The Go Programming Language",
  "authors": [
   "Alan A. A. Donovan",
   "Brian W. Kernighan"
  ],
  "pages": 380
 },
 {
  "title": "Advanced Programming in the UNIX Environment",
  "authors": [
   "W. Richard Stevens",
   "Stephen A. Rago"
  ],
  "pages": 1024
 },
 {
  "title": "The Practice of Programming",
  "authors": [
   "Brian W. Kernighan",
   "Rob Pike"
  ],
  "pages": 267
 },
 {
  "title": "The C Programming Language",
  "authors": [
   "Brian W. Kernighan",
   "Dennis Ritchie"
  ],
  "pages": 274
 },
 {
  "title": "Effective Go",
  "authors": [
   "Go core team"
  ],
  "pages": 50
 }
]

$> curl localhost:7777
<h1>Welcome to the Library</h1>
We have 5 books.
```

And voila, a very simple RESTful API web server.
