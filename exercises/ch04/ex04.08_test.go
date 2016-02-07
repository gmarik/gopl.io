package main

// Exercise 4.8: Modify charcount to count letters, digits, and so on in their Unicode categories, using functions like unicode.IsLetter.

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"testing"

	"strings"
	"unicode"
	"unicode/utf8"
)

type stats struct {
	counts  map[rune]int
	cat     map[string]int
	utflen  []int
	invalid int
}

var (
	input = `hello 123 ©Ω√ world`

	exp_stats = stats{
		map[rune]int{
			100:   1,
			101:   1,
			104:   1,
			108:   3,
			111:   2,
			114:   1,
			119:   1,
			169:   1,
			32:    3,
			49:    1,
			50:    1,
			51:    1,
			63743: 1,
			8730:  1,
			937:   1,
		},
		map[string]int{
			"lower":   10,
			"digit":   3,
			"number":  3,
			"graphic": 19,
			"letter":  11,
		},
		[]int{1: 16, 2: 2, 3: 2, 4: 0},
		0,
	}
)

func charcount(rd io.Reader) (stats, error) {
	counts := make(map[rune]int)    // counts of Unicode characters
	cat := make(map[string]int)     // counts categories
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters

	in := bufio.NewReader(rd)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			return stats{}, err
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++

		tests := []struct {
			name string
			f    func(rune) bool
		}{
			{"control", unicode.IsControl},
			{"digit", unicode.IsDigit},
			{"graphic", unicode.IsGraphic},
			{"letter", unicode.IsLetter},
			{"lower", unicode.IsLower},
			{"mark", unicode.IsMark},
			{"number", unicode.IsNumber},
		}

		for _, s := range tests {
			if s.f(r) {
				cat[s.name]++
			}
		}

	}
	return stats{counts, cat, utflen[:], invalid}, nil
}

type sortable struct {
	m []rune
}

func (s *sortable) Len() int {
	return len(s.m)
}

func (s *sortable) Swap(i, j int) {
	s.m[i], s.m[j] = s.m[j], s.m[i]
}

func (s *sortable) Less(i, j int) bool {
	return string(s.m[i]) < string(s.m[j])
}

func printCounts(s stats) string {
	counts, utflen, invalid := s.counts, s.utflen, s.invalid

	ordered := make([]rune, 0, len(counts))
	for r, _ := range counts {
		ordered = append(ordered, r)
	}

	sort.Sort(&sortable{ordered})
	w := &bytes.Buffer{}

	fmt.Fprintf(w, "rune\tcount\n")
	for _, c := range ordered {
		n := counts[c]
		fmt.Fprintf(w, "%q\t%d\n", c, n)
	}
	fmt.Fprint(w, "\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Fprintf(w, "%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Fprintf(w, "\n%d invalid UTF-8 characters\n", invalid)
	}

	return w.String()
}

func TestCounts(t *testing.T) {
	exp := exp_stats
	got, err := charcount(strings.NewReader(input))
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}

func TestPrintCounts(t *testing.T) {

	stats, err := charcount(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	got := printCounts(stats)
	exp := `rune	count
' '	3
'1'	1
'2'	1
'3'	1
'd'	1
'e'	1
'h'	1
'l'	3
'o'	2
'r'	1
'w'	1
'©'	1
'Ω'	1
'√'	1
'\uf8ff'	1

len	count
1	16
2	2
3	2
4	0
`
	if got != exp {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}
