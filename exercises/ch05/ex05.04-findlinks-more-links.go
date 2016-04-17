package main

// Exercise 5.4: Extend the visit function so that it extracts other kinds of links from the document, such as images, scripts, and style sheets.
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

	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

var linkMap = map[string]string{
	"a":    "href",
	"link": "href",

	"img":    "src",
	"script": "src",
	"iframe": "src",
}

func visit(links []string, n *html.Node) []string {
	for k, v := range linkMap {
		if n.Data == k && n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == v {
					links = append(links, a.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}
