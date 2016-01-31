package main

import "testing"
import "reflect"
import "unicode"
import "unicode/utf8"

// Exercise 4.3: Rewrite reverse to use an array pointer instead of a slice.
// Exercise 4.4: Write a version of rotate that operates in a single pass.
// Exercise 4.5: Write an in-place function to eliminate adjacent duplicates in a []string slice.
// Exercise 4.6: Write an in-place function that squashes each run of adjacent Unicode spaces (see unicode.IsSpace) in a UTF-8-encoded []byte slice into a single ASCII space.
// Exercise 4.7: Modify reverse to reverse the characters of a []byte slice that represents a UTF-8-encoded string, in place. Can you do it without allocating new memory?

// TODO: how to avoid duplication?
func reverse5(s *[5]int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func reverseByte(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// TODO: not sure if that's what authors intended with "single pass"
func rotate(s []int, r int) {
	if r == 0 {
		return
	}

	var (
		l = len(s)
		n = r % l
	)

	if n == 0 {
		return
	}

	if r < 0 {
		n = l + n
	}

	buf := make([]int, n)

	copy(buf, s[l-n:])
	copy(s[n:], s[0:l-n])
	copy(s[0:n], buf)
}

func dedup(s []string) []string {
	i := 0
	for {
		if i+1 >= len(s) {
			break
		}
		if s[i] == s[i+1] {
			if i+1 < len(s) {
				copy(s[i+1:], s[i+2:])
			}
			s = s[:len(s)-1]
		} else {
			i += 1
		}
	}

	return s
}

func squash(str string, p func(r rune) bool) string {
	s := []byte(str)

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRune(s[i:])

		if unicode.IsSpace(r) {
			spaceSize := size
			for {
				r, size := utf8.DecodeRune(s[i+spaceSize:])
				if !unicode.IsSpace(r) {
					break
				}
				spaceSize += size
			}

			s[i] = byte(' ')

			copy(s[i+1:], s[i+spaceSize:])
			s = s[:len(s)+1-spaceSize]
			i += 1
		} else {
			i += size
		}
	}

	return string(s)
}

func reverseUnicode(in string) string {
	var (
		s = []byte(in)
		i = 0
		l = len(s)
	)
	for {
		if i >= l {
			break
		}
		_, size := utf8.DecodeRune(s[i:])
		reverseByte(s[i : i+size])
		i += size
	}
	reverseByte(s)
	return string(s)
}

func TestReverseArrayPtr(t *testing.T) {

	var a = [5]int{5, 4, 3, 2, 1}
	reverse5(&a)

	exp := [5]int{1, 2, 3, 4, 5}
	got := a

	if got != exp {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}

func TestRotateSinglePass(t *testing.T) {

	cases := []struct {
		n        int
		got, exp []int
	}{
		{
			0,
			[]int{0, 1, 2, 3, 4, 5},
			[]int{0, 1, 2, 3, 4, 5},
		},
		{
			6,
			[]int{0, 1, 2, 3, 4, 5},
			[]int{0, 1, 2, 3, 4, 5},
		},
		{
			2,
			[]int{0, 1, 2, 3, 4, 5},
			[]int{4, 5, 0, 1, 2, 3},
		},
		{
			11,
			[]int{0, 1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5, 0},
		},
		{
			0,
			[]int{0, 1, 2, 3, 4, 5},
			[]int{0, 1, 2, 3, 4, 5},
		},
		{
			-2,
			[]int{0, 1, 2, 3, 4, 5},
			[]int{2, 3, 4, 5, 0, 1},
		},
		{
			-11,
			[]int{0, 1, 2, 3, 4, 5},
			[]int{5, 0, 1, 2, 3, 4},
		},
	}

	for _, c := range cases {
		got, exp := c.got, c.exp

		rotate(got, c.n)

		if !reflect.DeepEqual(got, exp) {
			t.Errorf("\nExp: %v\nGot: %v\nn: %v", exp, got, c.n)
		}
	}

}

func TestDedup(t *testing.T) {
	cases := []struct {
		exp, got []string
	}{
		{
			[]string{"a", "b"},
			[]string{"a", "a", "b", "b", "b"},
		},
		{
			[]string{"a", "b"},
			[]string{"a", "b", "b", "b"},
		},
		{
			[]string{"a"},
			[]string{"a", "a"},
		},
		{
			[]string{"a"},
			[]string{"a"},
		},
	}

	for _, v := range cases {
		exp := v.exp
		got := dedup(v.got)

		if !reflect.DeepEqual(got, exp) {
			t.Errorf("\nExp: %v\nGot: %v", exp, got)
		}
	}
}

func TestSquash(t *testing.T) {
	exp := "hello world !"
	got := squash("hello  world \u0085 \v \n !", unicode.IsSpace)

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}

func TestReverseUnicode(t *testing.T) {

	exp := "Ω∆\u4e16 ,olleH"
	got := reverseUnicode("Hello, \u4e16∆Ω")

	if exp != got {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}
