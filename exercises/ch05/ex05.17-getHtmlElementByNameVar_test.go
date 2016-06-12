package main

// Exercise 5.17: Write a variadic function ElementsByTagName that, given an HTML node tree and zero or more names, returns all the elements that match one of those names

import (
	"bytes"
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func ElementsByTagName(n *html.Node, names ...string) []*html.Node {
	var nodes []*html.Node

	var visit func(*html.Node)

	//depthFirst
	visit = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if n.Type == html.ElementNode && oneOf(c.Data, names...) {
				nodes = append(nodes, c)
			}
			visit(c)
		}
	}

	visit(n)

	return nodes
}

func oneOf(name string, names ...string) bool {
	for _, n := range names {
		if name == n {
			return true
		}
	}
	return false
}

func TestElementsByTagName(t *testing.T) {

	doc, err := html.Parse(bytes.NewReader([]byte(htmlSnip)))
	if err != nil {
		t.Fatal(err)
	}

	var got []string

	for _, node := range ElementsByTagName(doc, "head", "a") {
		got = append(got, node.Data)
	}

	exp := []string{"head", "a", "a"}

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}

var htmlSnip = `<html>
<head>
</head>
<body>
	<img />
	<a />
	<a />
	<br />
</body>
</html>
`
