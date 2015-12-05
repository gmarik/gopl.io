// Exercise 1.2: Modify the echo program to print the index and value of each
// of its arguments, one per line.

// See page 6.
//!+

// Echo2 prints its command-line arguments.
package main

import (
	"fmt"
	"os"
)

func main() {
	for i, arg := range os.Args[1:] {
		fmt.Printf("%d %s\n", i, arg)
	}
}

//!-
