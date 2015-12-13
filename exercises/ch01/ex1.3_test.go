// Exercise 1.3: Experiment to measure the difference in running time between our
// potentially inefficient versions and the one that uses strings.Join.
// (Section 1.6 illustrates part of the time package, and Section 11.4 shows how
// to write benchmark tests for systematic performance evaluation.)

// $ go test ex1.3_test.go -bench .
// testing: warning: no tests to run
// PASS
// BenchmarkConcat1_10-2  	 1000000	      1484 ns/op
// BenchmarkConcat2_10-2  	 1000000	      1485 ns/op
// BenchmarkConcat3_10-2  	 3000000	       449 ns/op
// BenchmarkConcat1_100-2 	  100000	     22604 ns/op
// BenchmarkConcat2_100-2 	  100000	     22186 ns/op
// BenchmarkConcat3_100-2 	  500000	      2959 ns/op
// BenchmarkConcat1_1000-2	    2000	    691194 ns/op
// BenchmarkConcat2_1000-2	    2000	    707082 ns/op
// BenchmarkConcat3_1000-2	   50000	     28483 ns/op
// ok  	command-line-arguments	16.894s

// See page 8.

// Echo3 prints its command-line arguments.
package main

import (
	"strings"
	"testing"
)

//!+

func concat1(args []string) string {
	var s, sep string
	for i := 1; i < len(args); i++ {
		s += sep + args[i]
		sep = " "
	}
	return s
}

func concat2(args []string) string {
	s, sep := "", ""
	for _, arg := range args[1:] {
		s += sep + arg
		sep = " "
	}
	return s
}

func concat3(args []string) string {
	return strings.Join(args[1:], " ")
}

var (
	short = [10]string{}
	mid   = [100]string{}
	big   = [1000]string{}
)

func bench(concat func([]string) string, n int, slice []string) {
	for i := 0; i < n; i += 1 {
		concat(slice)
	}
}

func BenchmarkConcat1_10(b *testing.B) {
	bench(concat1, b.N, short[:])
}
func BenchmarkConcat2_10(b *testing.B) {
	bench(concat2, b.N, short[:])
}
func BenchmarkConcat3_10(b *testing.B) {
	bench(concat3, b.N, short[:])
}

func BenchmarkConcat1_100(b *testing.B) {
	bench(concat1, b.N, mid[:])
}
func BenchmarkConcat2_100(b *testing.B) {
	bench(concat2, b.N, mid[:])
}
func BenchmarkConcat3_100(b *testing.B) {
	bench(concat3, b.N, mid[:])
}

func BenchmarkConcat1_1000(b *testing.B) {
	bench(concat1, b.N, big[:])
}
func BenchmarkConcat2_1000(b *testing.B) {
	bench(concat2, b.N, big[:])
}
func BenchmarkConcat3_1000(b *testing.B) {
	bench(concat3, b.N, big[:])
}

//!-
