// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Findlinks3 crawls the web, starting with the URLs on the command line.
//
// Exercise 5.13: Modify crawl to make local copies of the pages it finds, creating directories as necessary. Don’t make copies of pages that come from a different domain. For example, if the original page comes from golang.org, save all files from there, but exclude ones from vimeo.com

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"golang.org/x/net/html"
)

//!+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

//!-breadthFirst

//!+crawl
func crawl(url string) []string {
	fmt.Println(url)
	_url, rc, err := Fetch(url)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer rc.Close()

	buf := &bytes.Buffer{}
	teeReader := io.TeeReader(rc, buf)

	list, err := Extract(_url, teeReader)
	if err != nil {
		log.Println(err)
	}

	// if 0 < strings.Index(_url.String(), "golang.org") {
	if err := Save(_url, buf); err != nil {
		log.Println(err)
	}
	// }

	return list
}

// stores content of the url to disk
// TODO: detect if index.html has to be created
// TODO: only save content from certain host
// TODO: make it more modular
func Save(u *url.URL, buf io.Reader) error {
	p := path.Join("/tmp/save/", u.Path)
	if err := os.MkdirAll(p, os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path.Join(p, "index.html"))
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, buf); err != nil {
		return err
	}
	return f.Close()
}

//!-crawl

//!+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}

func Fetch(url string) (*url.URL, io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	return resp.Request.URL, resp.Body, nil
}

// Extract parses the response as HTML,
// and returns the links in the HTML document.
func Extract(u *url.URL, r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing as HTML: %v", err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := u.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

//!-Extract

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
