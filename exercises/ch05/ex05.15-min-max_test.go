package main

import (
	"fmt"
	"testing"
)

func Min(min int, vals ...int) int {
	switch l := len(vals); l {
	case 0:
		return min
	default:
		for _, v := range vals {
			if v < min {
				min = v
			}
		}

		return min
	}
}

func Min0(args ...int) int {

	switch l := len(args); l {
	case 0:
		panic(fmt.Errorf("Arguments expected"))
	case 1:
		return args[0]
	default:
		min := args[0]
		for _, v := range args[1:] {
			if v < min {
				min = v
			}
		}

		return min
	}
}

func TestMin(t *testing.T) {
	{
		exp := -3
		got := Min(1, 2, -3, 3, 4)
		if got != exp {
			t.Errorf("\nExp: %v\nGot: %v", exp, got)
		}
	}
	{
		exp := 1
		got := Min(1)
		if got != exp {
			t.Errorf("\nExp: %v\nGot: %v", exp, got)
		}
	}
}

func TestMin0(t *testing.T) {
	exp := -3
	got := Min0(1, 2, -3, 3, 4)
	if got != exp {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}

func TestMinNoArgs(t *testing.T) {
	var recovered error
	func() {
		defer func() {
			r := recover()
			var ok bool
			if recovered, ok = r.(error); !ok {
				recovered = fmt.Errorf("Unexpected error, %#v", r)
			}
		}()

		Min0()
	}()

	if recovered == nil || recovered.Error() != "Arguments expected" {
		t.Fatalf("Error expected, got %#v", recovered)
	}
}
