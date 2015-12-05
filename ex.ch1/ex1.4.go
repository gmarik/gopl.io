// Exercise 1.4: Modify dup2 to print the names of all files in which each
// duplicated line occurs.
//

// Running: go run ex.ch1/ex1.4.go -- ex.ch1/*.go|sort -n

// See page 10.
//!+

// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
)

type stats struct {
	count int
	files map[string]bool
}

func main() {
	counts := make(map[string]*stats)

	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, stats := range counts {
		if stats.count > 0 {
			fmt.Printf("%d\t%s\t%v\n", stats.count, line, stats.files)
		}
	}
}

func countLines(f *os.File, counts map[string]*stats) {
	input := bufio.NewScanner(f)

	for input.Scan() {
		t := input.Text()

		st := counts[t]

		if nil == st {
			// TODO: no zero value for maps? no zero value for reference types?
			st = &stats{files: map[string]bool{}}
			counts[t] = st
		}

		st.count += 1
		st.files[f.Name()] = true
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!-
