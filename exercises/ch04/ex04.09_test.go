package main

import (
	"bufio"
	"io"
	"reflect"
	// "sort"
	"testing"

	"strings"
)

var text = `
func wordfreq(rd io.Reader) (map[string]int, error) {
	freq := make(map[string]int)
	input := bufio.NewScanner(rd)

	for input.Split(bufio.ScanWords) {
		if err := input.Err(); err != nil {
			return nil, err
		}

		freq[input.Text()] += 1
	}
	return freq, nil
}
`

// Exercise 4.9: Write a program wordfreq to report the frequency of each word in an input text file. Call input.Split(bufio.ScanWords) before the first call to Scan to break the input into words instead of lines.”

func wordfreq(rd io.Reader) (map[string]int, error) {
	freq := make(map[string]int)
	input := bufio.NewScanner(rd)
	input.Split(bufio.ScanWords)

	for input.Scan() {
		if err := input.Err(); err != nil {
			return nil, err
		}

		freq[input.Text()] += 1
	}
	return freq, nil
}

func TestWordfreq(t *testing.T) {

	exp := map[string]int{
		"input.Split(bufio.ScanWords)": 1,
		"(map[string]int,":             1,
		"bufio.NewScanner(rd)":         1,
		"freq[input.Text()]":           1,
		"make(map[string]int)":         1,
		"input.Err();":                 1,
		"io.Reader)":                   1,
		"wordfreq(rd":                  1,

		"for":    1,
		"err":    3,
		"!=":     1,
		"freq":   1,
		":=":     3,
		"1":      1,
		"nil,":   1,
		"freq,":  1,
		"input":  1,
		"nil":    2,
		"error)": 1,
		"if":     1,
		"return": 2,
		"+=":     1,
		"func":   1,
		"}":      3,
		"{":      3,
	}

	got, err := wordfreq(strings.NewReader(text))

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}

}
