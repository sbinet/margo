package main

import (
        "fmt"
        "log"
        "net/http"
)

func main() {
        fmt.Println("please connect to localhost:7777/hello")
        http.HandleFunc("/hello", HelloServer)
        log.Fatal(http.ListenAndServe(":7777", nil))
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
        log.Println(req.URL)
        fmt.Fprintf(w, "Hello, world!\nURL = %s\n", req.URL)
}
