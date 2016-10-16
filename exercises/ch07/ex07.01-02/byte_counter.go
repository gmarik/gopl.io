package main

import (
	"bufio"
	"io"
)

// Exercise 7.1: Using the ideas from ByteCounter, implement counters for
// words and for lines. You will find bufio.ScanWords useful.
//
// Exercise 7.2: Write a function CountingWriter with the signature below
// that, given an io.Writer, returns a new Writer that wraps the original,
// and a pointer to an int64 variable that at any moment contains the number
// of bytes written to the new Writer.
//
//			func CountingWriter(w io.Writer) (io.Writer, *int64)
//

// ByteCounter acts as a io.Writer in order to count bytes written
// It's pretty much useless since Write returns number of bytes written
type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	var off int
	for {
		n, _, err := bufio.ScanWords(p[off:], true)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		off += n
		*c += WordCounter(1) // convert int to
	}
	return len(p), nil
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	var off int
	for {
		n, _, err := bufio.ScanLines(p[off:], true)
		if err != nil {
			return 0, err
		}
		off += n
		if n == 0 {
			break
		}
		*c += LineCounter(1) // convert int to
	}
	return len(p), nil
}

type CountWriter struct {
	w     io.Writer
	count int64
}

func (cw *CountWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	if err != nil {
		return n, err
	}

	cw.count += int64(n)

	return n, nil
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	cw := &CountWriter{w: w}
	return cw, &cw.count
}
