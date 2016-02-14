package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// downloading( for some reason boolean flag has to be the last one)
// go run exercises/ch04/ex04.12_xkcd.go  --dl-from 1642 -index-path="$PWD/out" -dl
//
// search mode otherwise:
// go run exercises/ch04/ex04.12_xkcd.go Gravitational Waves

var (
	isDownload = flag.Bool("dl", false, "Download mode")
	dlFrom     = flag.Int("dl-from", 0, "Download from id")
	dlTo       = flag.Int("dl-to", 1643, "Download to id")
	indexPath  = flag.String("index-path", "./xkcd", "path do dir with saved json files")
)

func main() {
	flag.Parse()

	if *isDownload {
		if err := os.MkdirAll(*indexPath, os.ModeDir|os.ModePerm); err != nil {
			log.Fatal(err)
		}

		log.Printf("Downloading: %d-%d to %q\n", *dlFrom, *dlTo, *indexPath)
		download(*dlFrom, *dlTo, *indexPath)
		return
	}

	query := strings.Join(flag.Args(), " ")

	log.Println("Searching for: ", query)
	if err := search(query, *indexPath); err != nil {
		log.Fatal(err)
	}
}

type xkcd struct {
	Day   string `json:"day"`
	Month string `json:"month"`
	Year  string `json:"year"`

	Num int `json:"num"`

	Link       string `json:"link"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
}

func match(query string, info *xkcd) bool {
	for _, s := range []*string{&info.Transcript, &info.Alt, &info.Title} {
		if strings.Index(*s, query) >= 0 {
			return true
		}
	}

	return false
}

func search(query, indexPath string) error {
	matches, err := filepath.Glob(path.Join(indexPath, "*.json"))
	if err != nil {
		return err
	}

	var info xkcd

	for _, fn := range matches {
		data, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Println("Error reading:", fn)
			continue
		}

		if err := json.Unmarshal(data, &info); err != nil {
			log.Println("Error parsing:", fn)
			continue
		}

		if match(query, &info) {
			log.Printf("[%d] %q (%s) %s, %s\n\n", info.Num, info.Title, info.Img, info.Transcript, info.Alt)
		}
	}
	return nil
}

type result struct {
	id  int
	err error

	url  string
	body io.ReadCloser
}

func (m *result) Filename() string {
	return fmt.Sprintf("xkcd-%04d.json", m.id)
}

func downloader(sinkCh, downloadCh, persistCh chan result) {
	for r := range downloadCh {
		resp, err := http.Get(r.url)
		if err != nil {
			r.err = fmt.Errorf("Error downloading %q: %q\n", r.url, err)
			sinkCh <- r
			continue
		}
		if resp.StatusCode != http.StatusOK {
			r.err = fmt.Errorf("Error downloading %q: %d response", r.url, resp.StatusCode)
			sinkCh <- r
			continue
		}

		r.body = resp.Body
		persistCh <- r
	}
}

func writer(sinkCh, persistCh chan result, indexPath string) {
	for r := range persistCh {
		f, err := os.Create(path.Join(indexPath, r.Filename()))
		if err != nil {
			r.err = fmt.Errorf("Error creating %q: %q", r.url, err)
			sinkCh <- r
			continue
		}

		_, err = io.Copy(f, r.body)
		if err != nil {
			log.Printf("Error writing %q: %q", r.Filename(), err)
			sinkCh <- r
			continue
		}
		r.body.Close()
		f.Close()

		sinkCh <- r
	}
}

func download(from, to int, indexPath string) error {
	var (
		downloadCh = make(chan result)
		persistCh  = make(chan result)
		sinkCh     = make(chan result)
	)

	for i := 0; i < 20; i += 1 {
		go downloader(sinkCh, downloadCh, persistCh)
		go writer(sinkCh, persistCh, indexPath)
	}

	go func() {
		for i := from; i <= to; i += 1 {
			// for i := 1; i < 2; i += 1 {
			downloadCh <- result{id: i, url: fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)}
		}
	}()

	var count int

	for r := range sinkCh {
		if r.err != nil {
			log.Println(r.err)
		} else {
			log.Printf("Saved %s\n", r.url)
		}

		if count == (to - from) {
			break
		}
		count += 1
	}
	log.Printf("Done")

	return nil
}
