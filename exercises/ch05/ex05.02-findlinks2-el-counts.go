package main

// Exercise 5.2: Write a function to populate a mapping from element names—p, div, span, and so on—to the number of elements with that name in an HTML document tree.
import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks: %v\n", err)
		os.Exit(1)
	}

	m := make(map[string]uint)
	visit(m, doc)
	for name, count := range m {
		fmt.Printf("%q: %d\n", name, count)
	}
}

func visit(m map[string]uint, n *html.Node) {
	if n.Type == html.ElementNode {
		m[n.Data] += 1
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visit(m, c)
	}
}
