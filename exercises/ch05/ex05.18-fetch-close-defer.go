package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	for _, v := range os.Args[1:] {
		local, _, err := fetch(v)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Downloaded %q to %q\n", v, local)
	}
}

func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	n, err = io.Copy(f, resp.Body)
	// Close file, but prefer error from Copy, if any.
	//
	// if closeErr := f.Close(); err == nil {
	// 	err = closeErr
	// }
	return local, n, err
}
