package main

import "io"

// Exercise 7.4: The strings.NewReader function returns a value that
// satisfies the io.Reader interface (and others) by reading from its
// argument, a string. Implement a simple version of NewReader yourself, and
// use it to make the HTML parser (§5.2) take input from a string.

// Exercise 7.5: The LimitReader function in the io package accepts an
// io.Reader r and a number of bytes n, and returns another Reader that reads
// from r but reports an end-of-file condition after n bytes. Implement it.
//
// 		func LimitReader(r io.Reader, n int64) io.Reader”
//

// Reader implements Reader interface on a string
type Reader struct {
	i int // the index reader is at
	s string
}

// NewReader constructs Reader from a string
func NewReader(s string) *Reader {
	return &Reader{s: s}
}

func (r *Reader) Read(buf []byte) (int, error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}

	n := copy(buf, r.s[r.i:])
	r.i += n

	return n, nil
}

// LimitReader implements both Reader and limiting logic
type LimitReader struct {
	n int64 // max number bytes to read
	r io.Reader
}

// NewLimitReader constructs Reader from a string limited to n bytes
func NewLimitReader(r io.Reader, n int64) *LimitReader {
	return &LimitReader{r: r, n: n}
}

func (r *LimitReader) Read(buf []byte) (int, error) {
	if r.n <= 0 {
		return 0, io.EOF
	}

	if int64(len(buf)) > r.n {
		buf = buf[:r.n]
	}

	n, err := r.r.Read(buf)
	r.n -= int64(n)

	return n, err
}
