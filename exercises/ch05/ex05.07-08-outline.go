package main

// Exercise 5.7: Develop startElement and endElement into a general HTML pretty-printer. Print comment nodes, text nodes, and the attributes of each element (<a href='...'>). Use short forms like <img/> instead of <img></img> when an element has no children. Write a test to ensure that the output can be parsed successfully. (See Chapter 11.)
// Exercise 5.8: Modify forEachNode so that the pre and post functions return a boolean result indicating whether to continue the traversal. Use it to write a function ElementByID with the following signature that finds the first HTML element with the specified id attribute. The function should stop the traversal as soon as a match is found.

// Outline prints the outline of an HTML document tree.

import (
	"fmt"
	"io"
	"os"

	"strings"

	"golang.org/x/net/html"
)

var (
	out    io.Writer = os.Stdout
	stopAt string
)

func main() {
	outline(os.Stdin)
}

func outline(r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	stopWhenFound := func(n *html.Node, depth int, hasChildren bool) bool {
		for _, an := range n.Attr {
			if an.Key == "id" && an.Val == stopAt && stopAt != "" {
				return true
			}
		}

		return startElement(n, depth, hasChildren)
	}

	forEachNode(doc, 0, stopWhenFound, endElement)

	return nil
}

func forEachNode(n *html.Node, depth int, pre, post func(*html.Node, int, bool) bool) bool {
	if pre != nil {
		stop := pre(n, depth, n.FirstChild != nil)
		if stop {
			return stop
		}

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		stop := forEachNode(c, depth+1, pre, post)
		if stop {
			return stop
		}
	}

	if post != nil {
		stop := post(n, depth, n.FirstChild != nil)
		if stop {
			return stop
		}
	}
	return false
}

func startElement(n *html.Node, depth int, hasChildren bool) bool {
	if n.Type == html.ElementNode {
		var attrNames []string
		for _, an := range n.Attr {
			attrNames = append(attrNames, an.Key)
		}
		var terminator string
		if !hasChildren {
			terminator = "/"
		}
		fmt.Fprintf(out, "%*s<%s [%s]%s>\n", depth*2, "", n.Data, strings.Join(attrNames, " "), terminator)
	}

	return false
}

func endElement(n *html.Node, depth int, hasChildren bool) bool {
	if hasChildren && n.Type == html.ElementNode {
		fmt.Fprintf(out, "%*s</%s>\n", depth*2, "", n.Data)
	}
	return false
}
