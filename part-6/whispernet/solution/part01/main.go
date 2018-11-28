// This program reads from standard input and writes JSON-encoded messages to
// standard output. For example, this input line:
//	Hello!
// Produces this output:
//	{"Body":"Hello!"}
//
package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Message struct {
	Body string
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for sc.Scan() {
		msg := Message{Body: sc.Text()}
		err := enc.Encode(msg)
		if err != nil {
			log.Fatalf("could not encode message %#v: %v", msg, err)
		}
	}

	if err := sc.Err(); err != nil {
		log.Fatalf("error while scanning: %v", err)
	}
}
