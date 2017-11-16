package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

// START OMIT
func main() {
	const origin = "http://localhost/"
	ws, err := websocket.Dial("ws://localhost:12345/echo", "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte("hello there\n")); err != nil {
		log.Fatal(err)
	}

	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %q\n", msg[:n])
}

// END OMIT
