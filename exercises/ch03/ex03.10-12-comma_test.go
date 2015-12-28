package main

import "bytes"
import "fmt"
import "testing"
import "strings"

// Exercise 3.10: Write a non-recursive version of comma, using bytes.Buffer
// instead of string concatenation.
//

// Exercise 3.11: Enhance comma so that it deals correctly with floating-point
// numbers and an optional sign.
//

// Exercise 3.12: Write a function that reports whether two strings are anagrams
// of each other, that is, they contain the same letters in a different order.

func Anagrams(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i += 1 {

		if a[i] != b[len(b)-i-1] {
			return false
		}
	}

	return true
}

func commaFloat(s string) (string, error) {
	if len(s) < 2 {
		return s, nil
	}

	var sign = ""
	if s[0] == '+' || s[0] == '-' {
		sign, s = string(s[0]), s[1:]
	}

	var decimal = ""

	if i := strings.Index(s, "."); i > -1 {
		s, decimal = s[0:i], s[i:]
	}

	s, err := commaNonRecursive(s)
	if err != nil {
		return s, err
	}

	return sign + s + decimal, nil
}

// Reverse returns its argument string reversed rune-wise left to right.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func commaNonRecursive(s string) (string, error) {
	s = Reverse(s)
	var b bytes.Buffer

	for i, v := range s {
		// i	     |  0 1 2 3 4 5 6 7 8 1
		// i%3     |  0 1 2 0 1 2 0 1 2 0
		// (i+1)%3 |  1 2 0 1 2 0 0 1 2 0
		if i%3 == 0 && i > 0 && i < len(s) {
			b.WriteRune(',')
		}

		if _, err := b.WriteRune(v); err != nil {
			return "", err
		}
	}

	return Reverse(b.String()), nil
}

func ExampleTestNonRecursive() {
	var s string

	s, _ = commaNonRecursive("12345678")
	fmt.Println(s)

}

func TestCommaNonRecursive(t *testing.T) {

	cases := []struct {
		in  string
		out string
	}{
		{"123", "123"},
		{"1234", "1,234"},
		{"12345", "12,345"},
		{"123456", "123,456"},
		{"1234567", "1,234,567"},
		{"-123", "-123"},
		{"-1234", "-1,234"},
		{"+12345", "+12,345"},
		{"+123456", "+123,456"},
		{"+1234567", "+1,234,567"},
		{"-123.1", "-123.1"},
		{"-1234.1234", "-1,234.1234"},
		{"+12345.1234", "+12,345.1234"},
		{"+123456.1234", "+123,456.1234"},
		{"+1234567.1234", "+1,234,567.1234"},
	}

	for _, c := range cases {
		got, err := commaFloat(c.in)
		if err != nil {
			t.Error(err)
		}

		exp := c.out

		if got != exp {
			t.Errorf("\nExp: %v\nGot: %v", exp, got)
		}
	}
}

func TestAnagrams(t *testing.T) {
	cases := []struct {
		a, b       string
		isAanagram bool
	}{
		{"", "", true},
		{"a", "a", true},
		{"ab", "ab", false},
		{"ab", "ab", false},
		{"aba", "aba", true},
		{"abac", "abac", false},
		{"ababa", "ababa", true},
		{"ababa", "cbaba", false},
		{"a", "ok", false},
	}

	for _, v := range cases {
		exp := v.isAanagram
		got := Anagrams(v.a, v.b)

		if got != exp {
			t.Errorf("\nExp: %v\nGot: %v\n for %s %s", exp, got, v.a, v.b)
		}
	}
}
