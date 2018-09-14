// margo-hello is a simple Hello World command to test everything is correctly setup
package main

import (
	"fmt"
	"os"
)

func main() {
	who := "MarGo"
	if len(os.Args) > 1 {
		who = os.Args[1]
	}
	fmt.Printf("Hello %s!\n", who)
}
