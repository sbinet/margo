package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("please connect to localhost:7777")
	http.HandleFunc("/", rootHandle)
	log.Fatal(http.ListenAndServe(":7777", nil))
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	clients++
	time.Sleep(2 * time.Second)
	fmt.Fprintf(w, "time:  %v\n", time.Now().UTC())
	fmt.Fprintf(w, "conns: %v\n", clients)
}

var clients = 0
