package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// Echo the data received on the WebSocket.
func echoServer(ws *websocket.Conn) { io.Copy(ws, ws) }

func main() {
	http.Handle("/echo", websocket.Handler(echoServer))
	log.Fatal(http.ListenAndServe(":12345", nil))
}
