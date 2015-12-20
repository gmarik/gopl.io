package main

// Exercise 2.3: Rewrite PopCount to use a loop instead of a single
// expression.  Compare the performance of the two versions. (Section 11.4
// shows how to compare the performance of different implementations
// systematically.)
//
// Exercise 2.4: Write a version of PopCount that counts bits by shifting its
// argument through 64 bit positions, testing the rightmost bit each time.
// Compare its performance to the table-lookup version.
//
// Exercise 2.5: The expression x&(x-1) clears the rightmost non-zero bit of x.
// Write a version of PopCount that counts bits by using this fact, and assess
// its performance.

//
// Running: $ go test -bench=. exercises/ch02/ex02.03-popcount_test.go

import "testing"

// pc[i] is the population count of i.
var pc [256]byte

func init() {

	// Precompute bits table
	// n |bit rep| bit set
	// -------------------
	// 0 |    0000| 0
	// 1 |    0001| 1
	// 2 |    0010| 1
	// 3 |    0011| 2
	// 4 |    0100| 1
	// 5 |    0101| 2
	// 6 |    0110| 2
	// 7 |    0111| 3
	// 8 |    1000| 1
	// 9 |    1001| 2
	//10 |    1010| 2
	//11 |    1011| 3
	//12 |    1100| 2
	//
	// Example:
	// i = 11 = 1011b
	// i/2 = 11/2 = 11 >> 1 = 1011 >> 1 = 0101b = 5
	// i & 1 = 11 & 1 = 1011b & 0001b = 1
	// pc[11] = pc[5] + byte(i & 1) = 2 + 1 = 3
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

//
// PopCount* Explanation:
//
// uint64 - 64 bit unsigned integer
// 64 = 8 * 8 bits (8 bytes)
// XXXXXXXX XXXXXXXX XXXXXXXX ... XXXXXXXX XXXXXXXX
// TotalBitsSet = each_byte(x).map(&bits_set).inject(0, &:+)

// PopCount returns the population count (number of set bits) of x.
func PopCountSingleExpr(x uint64) int {
	return int(
		pc[byte(x>>(0*8))] +
			pc[byte(x>>(1*8))] +
			pc[byte(x>>(2*8))] +
			pc[byte(x>>(3*8))] +
			pc[byte(x>>(4*8))] +
			pc[byte(x>>(5*8))] +
			pc[byte(x>>(6*8))] +
			pc[byte(x>>(7*8))])
}

// PopCount returns the population count (number of set bits) of x.
func PopCountLoop(x uint64) int {
	var sum byte
	for i := 0; i < 64; i += 8 {
		sum += pc[byte(x>>uint(i))]
	}
	return int(sum)

}

func PopCount64Shifts(x uint64) int {
	var sum byte
	for i := 0; i < 64; i += 1 {
		sum, x = sum+byte(x&1), x>>1
	}
	return int(sum)
}

// Explanation:
// clears righmost non-zero bit:
// 0101
// 0100
// 0000
func PopCountClearLastBit(x uint64) int {
	var sum byte
	for ; x > 0; x = x & (x - 1) {
		sum += 1
	}
	return int(sum)
}

func TestPopCounts(t *testing.T) {
	const (
		BitsCount = 32
		V         = 0x0123456789ABCDEF
	)

	fs := []func(uint64) int{PopCountClearLastBit, PopCount64Shifts, PopCountLoop, PopCountSingleExpr}

	for i, f := range fs {
		if BitsCount != f(V) {
			t.Errorf("Invalid count value %v expected %v [tcase: %d]", f(V), BitsCount, i)
		}
	}
}

func BenchmarkClearLastBit(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		PopCountClearLastBit(0x0123456789ABCDEF)
	}
}

func Benchmark64Shifts(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		PopCount64Shifts(0x0123456789ABCDEF)
	}
}

func BenchmarkPopCountLoop(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		PopCountLoop(0x0123456789ABCDEF)
	}
}

func BenchmarkPopCountSingleExpr(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		PopCountSingleExpr(0x0123456789ABCDEF)
	}
}
