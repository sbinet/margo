// This program listens on the host and port specified by the -listen flag.
// For each incoming connection, it launches a goroutine that reads and decodes
// JSON-encoded messages from the connection and prints them to standard
// output.
//
// You can test this program by running it in one terminal:
// 	$ part3 -listen=localhost:8000
// And running part2 in another terminal:
// 	$ part2 -dial=localhost:8000
// Lines typed in the second terminal should appear as JSON objects in the
// first terminal.
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var listenAddr = flag.String("listen", "localhost:8000", "host:port to listen on")

type Message struct {
	Body string
}

func main() {
	flag.Parse()

	srv, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("could not listen on %q: %v", *listenAddr, err)
	}

	for {
		c, err := srv.Accept()
		if err != nil {
			log.Fatalf("could not accept connection: %v", err)
		}
		go serve(c)
	}

}

func serve(c net.Conn) {
	defer c.Close()

	dec := json.NewDecoder(c)

	for {
		var msg Message
		err := dec.Decode(&msg)
		if err != nil {
			if err == io.EOF {
				log.Printf("client disconnected %v", c.RemoteAddr())
				return
			}
			log.Fatalf("could not decode message: %v", err)
		}
		fmt.Printf("message> %#v\n", msg)
	}
}
