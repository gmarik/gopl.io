package main

// Exercise 5.16: Write a variadic version of strings.Join.

import (
	"reflect"
	"testing"
)

func join(args ...string) string {
	s := ""
	for i := 0; i < len(args); i += 1 {
		s += args[i]
		if i+1 < len(args) {
			s += ","
		}
	}

	return s
}

func TestJoin(t *testing.T) {

	cases := []struct {
		in  []string
		exp string
	}{
		{
			[]string{},
			"",
		},
		{
			[]string{"a"},
			"a",
		},
		{
			[]string{"a", "b"},
			"a,b",
		},
		{
			[]string{"a", "b", "c"},
			"a,b,c",
		},
	}

	for _, c := range cases {
		exp := c.exp
		got := join(c.in...)

		if !reflect.DeepEqual(got, exp) {
			t.Errorf("\nExp: %v\nGot: %v", exp, got)
		}
	}
}
