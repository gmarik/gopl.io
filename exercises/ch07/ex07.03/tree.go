package main

import (
	"bytes"
	"fmt"
)

// Exercise 7.3: Write a String method for the *tree type in
// gopl.io/ch4/treesort (§4.4) that reveals the sequence of values in the
// tree.
//
type tree struct {
	value       int
	left, right *tree
}

func (t *tree) String() string {
	if t == nil {
		return "<nil>"
	}

	var buf bytes.Buffer

	var printer = func(depth int, kind int, value int) {
		padding := bytes.Repeat([]byte(" "), depth*2)
		fmt.Fprintf(&buf, "%s%d\n", padding, value)
	}

	t.visit(0, 0, printer)

	return buf.String()
}

func (t *tree) visit(depth int, kind int, visitor func(depth int, kind int, value int)) {
	if t == nil {
		return
	}

	visitor(depth, kind, t.value)

	t.left.visit(depth+1, 1, visitor)
	t.right.visit(depth+1, 2, visitor)
}
