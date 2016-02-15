package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Exercise 4.13: The JSON-based web service of the Open Movie Database lets you search https://omdbapi.com/ for a movie by name and download its poster image. Write a tool poster that downloads the poster image for the movie named on the command line.

func main() {
	c := &http.Client{Timeout: 2 * time.Second}

	req, err := http.NewRequest("GET", "http://www.omdbapi.com/?plot=short&r=json", nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Set("t", strings.Join(os.Args[1:], " "))
	req.URL.RawQuery = q.Encode()

	log.Printf("Querying: %q", req.URL.String())

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Response status: ", resp.StatusCode)
	}

	var search SearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&search); err != nil {
		log.Fatal(err)
	}

	if search.Error != "" {
		log.Printf("%#v", search)
		log.Fatal(search.Error)
	}

	log.Println("Poster URL:", search.Poster)
}

type SearchResponse struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	Rated      string `json:"Rated"`
	Released   string `json:"Released"`
	Runtime    string `json:"Runtime"`
	Genre      string `json:"Genre"`
	Director   string `json:"Director"`
	Writer     string `json:"Writer"`
	Actors     string `json:"Actors"`
	Plot       string `json:"Plot"`
	Language   string `json:"Language"`
	Country    string `json:"Country"`
	Awards     string `json:"Awards"`
	Poster     string `json:"Poster"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	Response   string `json:"Response"`
	Error      string `json:"Error"`
}
