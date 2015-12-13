// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Alexa top https://gist.github.com/jesstess/1392819

// Try fetchall with longer argument lists, such as samples from the top million
// web sites available at alexa.com. How does the program behave if a web site
// just doesn’t respond? (Section 8.9 describes mechanisms for coping in such
// cases.)

// Running:
// cat ex.ch1/alexa-top-100.txt |xargs go run ex.ch1/ex01.11-fetchall.go
// Get http://bp.blogspot.com: dial tcp: lookup bp.blogspot.com: no such host
// Get http://1e100.net: dial tcp: lookup 1e100.net: no such host
// 0.67s    11884  http://conduit.com
// 0.73s     3256  http://pornhub.com
// 0.81s     3256  http://xhamster.com
// 0.84s     3256  http://redtube.com
// 0.86s     3256  http://livejasmin.com
// 0.87s    35332  http://fc2.com
// 0.87s     3278  http://thepiratebay.org
// 0.90s     3256  http://tube8.com
// 0.98s     3256  http://xvideos.com
// 1.03s     9811  http://imageshack.us
// 1.09s   213313  http://digg.com
// 1.40s    63789  http://photobucket.com
// 1.40s   106716  http://go.com
// ...
// 7.25s   429621  http://sohu.com
// 8.26s        0  http://youku.com
// Get https://www.taobao.com/: net/http: TLS handshake timeout
// 17.90s   613957  http://qq.com
// Get http://rapidshare.com: dial tcp 10.11.12.13:80: i/o timeout
// Get http://yieldmanager.com: dial tcp 208.67.66.24:80: i/o timeout
// Get http://cnzz.com: dial tcp 42.156.162.55:80: i/o timeout
// 30.02s elapsed

// See page 17.
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch("http://"+url, ch) // start a goroutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

//!-
