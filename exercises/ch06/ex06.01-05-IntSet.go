package main

import (
	"bytes"
	"fmt"
)

// TODO: Problem: Testing is somewhat difficult and fragile:
// even though test cases aren't pointers underlying slices of IntSet are shared
// which affects/conflicts with test cases
// Solution: always copy the set being modified, but this requires Copy be implemented in advance
// Question: are there language patterns testing idioms that would provide non-shared test cases or it's something developer has to take care
//
// TODO: split test cases per operation and target edge cases
//

// Exercise 6.1: Implement these additional methods
// 		func (*IntSet) Len() int      // return the number of elements
// 		func (*IntSet) Remove(x int)  // remove x from the set
// 		func (*IntSet) Clear()        // remove all elements from the set
// 		func (*IntSet) Copy() *IntSet // return a copy of the set
//
// Exercise 6.2: Define a variadic (*IntSet).AddAll(...int) method that
// allows a list of values to be added, such as s.AddAll(1, 2, 3).
//
// Exercise 6.3: (*IntSet).UnionWith computes the union of two sets using |,
// the word-parallel bitwise OR operator. Implement methods for
// IntersectWith, DifferenceWith, and SymmetricDifference for the
// corresponding set operations. (The symmetric difference of two sets
// contains the elements present in one set or the other but not both.)
//
// Exercise 6.4: Add a method Elems that returns a slice containing the
// elements of the set, suitable for iterating over with a range loop.”

// Exercise 6.5: The type of each word used by IntSet is uint64, but 64-bit
// arithmetic may be inefficient on a 32-bit platform. Modify the program to
// use the uint type, which is the most efficient unsigned integer type for
// the platform. Instead of dividing by 64, define a constant holding the
// effective size of uint in bits, 32 or 64. You can use the perhaps
// too-clever expression 32 << (^uint(0) >> 63) for this purpose.

// Truth tables
//
// -----------  union | inters | diff | symdiff aka (xor)
// | p  |  q  | p | q | p & q  | p&^q | (p|q)&^(p&q) aka (p|q)&(^p|^q)
// |----|-----|-------|--------|------|---------
// | 0  |  0  |   0   |   0    |  0   | 0 = 0 & 1
// | 0  |  1  |   1   |   0    |  0   | 1 = 1 & 1
// | 1  |  0  |   1   |   0    |  1   | 1 = 1 & 1
// | 1  |  1  |   1   |   1    |  0   | 0 = 1 & 0
//
// https://en.wikipedia.org/wiki/Set_(mathematics)#Basic_operations
//
// Truth tables: https://en.wikipedia.org/wiki/Truth_table
// diff: aka material non-implication https://en.wikipedia.org/wiki/Material_nonimplication
// symdiff: aka exclusive or https://en.wikipedia.org/wiki/Exclusive_or

// hack to get the size of the uint(which is platform dependant and behind the scenes is either uint32 or uint64)
// 1. default to 32 bit platform
// 2. but if ^(uint(0) >> 63 gives 1
// 3. then it's 64 bit platform so double 32 ( 32 << 1)
const wordSize = 32 << (^uint(0) >> 63)

type IntSet struct {
	words []uint
}

func NewIntSet(xs ...int) *IntSet {
	var x IntSet
	x.AddAll(xs...)

	return &x
}

func (s *IntSet) wordBit(x int) (int, uint) {
	return x / wordSize, uint(x % wordSize)
}

func (s *IntSet) Has(x int) bool {
	word, bit := s.wordBit(x)
	return word < len(s.words) && (s.words[word]&(1<<bit)) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := s.wordBit(x)

	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}

	s.words[word] |= (1 << bit)
}

func (s *IntSet) AddAll(xs ...int) {
	for _, x := range xs {
		s.Add(x)
	}
}
func (s *IntSet) Len() int {
	var sum int

	for _, w := range s.words {
		for i := 0; i < wordSize; i += 1 {
			if (w & (1 << uint(i))) != 0 {
				sum += 1
			}
		}
	}

	return sum
}

func (s *IntSet) Remove(x int) bool {
	word, bit := s.wordBit(x)
	if word >= len(s.words) {
		return false
	}
	if s.words[word]&(1<<bit) == 0 {
		return false
	}

	s.words[word] &= ^(1 << bit)

	return true
}

func (s *IntSet) Copy() *IntSet {
	var x IntSet

	if s.words != nil {
		x.words = make([]uint, len(s.words))
		copy(x.words, s.words)
	}

	return &x
}

func (s *IntSet) Clear() {
	s.words = s.words[:0]
}

func (s *IntSet) Elems() []int {
	var elems []int

	for i, w := range s.words {
		var bit uint = 1
		for j := 0; j < wordSize; j += 1 {
			if (w & bit) != 0 {
				elems = append(elems, i*wordSize+j)
			}
			bit = bit << 1
		}
	}

	return elems
}

func (s *IntSet) String() string {
	var buf bytes.Buffer

	buf.WriteString("{")

	for j, w := range s.words {
		if w == 0 {
			continue
		}

		for i := 0; i < wordSize; i += 1 {
			if w&(1<<uint(i)) != 0 {
				if buf.Len() != len("{") {
					buf.WriteString(" ")
				}
				fmt.Fprintf(&buf, "%d", j*wordSize+i)
			}
		}
	}

	buf.WriteString("}")

	return buf.String()
}

//	 -----------------------------------------
//   | len(a) | len(b) |  a unionwith b
//   |--------|--------|----------------------
// 1.|   *    |   0    |	a
// 2.|   0    |   *    |	b
// 3.|   <b   |   *    |  a |= b; + b
// 4.|   =b   |   *    |  a |= b
// 5.|   >b   |   *    |  a[:len(b)] |= b

func (a *IntSet) UnionWith(b *IntSet) {
	if len(b.words) == 0 {
		return
	}

	for i := range b.words {
		if i >= len(a.words) {
			break
		}
		a.words[i] |= b.words[i]
	}

	if len(a.words) < len(b.words) {
		a.words = append(a.words, b.words[len(a.words):]...)
	}
}

//	 -----------------------------------------
//   | len(a) | len(b) |  a intersect b
//   |--------|--------|----------------------
// 1.|   *    |   0    |	0
// 2.|   0    |   *    |	0
// 3.|   <b   |   *    |  a &= b
// 4.|   =b   |   *    |  a &= b
// 5.|   >b   |   *    |  a[:len(b)] &= b, cut the reset

func (a *IntSet) IntersectWith(b *IntSet) {
	if len(a.words) == 0 || len(b.words) == 0 {
		a.words = nil
		return
	}

	var n = len(b.words)

	if len(a.words) < n {
		n = len(a.words)
	} else {
		a.words = a.words[:n]
	}

	for i := 0; i < n; i++ {
		var ridx = n - i - 1
		a.words[ridx] &= b.words[ridx]

		if a.words[ridx] == 0 {
			a.words = a.words[:ridx]
		}
	}

	if len(a.words) == 0 {
		a.words = nil
	}
}

//	 -----------------------------------------
//   | len(a) | len(b) |  a diff b
//   |--------|--------|----------------------
// 1.|   *    |   0    |	0
// 2.|   0    |   *    |	0
// 3.|   <=b  |   *    |  a &= ^b
// 5.|   >b   |   *    |  a[:len(b)] &= ^b, keep the rest

func (a *IntSet) DifferenceWith(b *IntSet) {
	if len(a.words) == 0 || len(b.words) == 0 {
		a.words = nil
		return
	}

	var n = len(b.words)
	if len(a.words) < len(b.words) {
		n = len(a.words)
	}

	for i := 0; i < n; i++ {
		var ridx = n - i - 1
		a.words[ridx] &= ^b.words[ridx]

		if a.words[ridx] == 0 {
			a.words = a.words[:ridx]
		}
	}

	if len(a.words) == 0 {
		a.words = nil
	}
}

//	 -----------------------------------------
//   | len(a) | len(b) |  a symmdiff b
//   |--------|--------|----------------------
// 1.|   *    |   0    | 0
// 2.|   0    |   *    | 0
// 3.|   <b   |   *    | a  op b
// 4.|   =b   |   *    | a  op b
// 5.|   >b   |   *    | a[:len(b)] op b

func (a *IntSet) SymmDiffWith(b *IntSet) {
	if len(a.words) == 0 || len(b.words) == 0 {
		a.words = nil
		return
	}

	var n = len(b.words)
	if len(a.words) < len(b.words) {
		n = len(a.words)
		a.words = append(a.words, b.words[n:]...)
	}

	for i := 0; i < n; i++ {
		var ridx = n - i - 1

		q, p := a.words[ridx], b.words[ridx]

		a.words[ridx] = (p | q) & ^(p & q)

		if a.words[ridx] == 0 {
			a.words = a.words[:ridx]
		}
	}

	if len(a.words) == 0 {
		a.words = nil
	}
}
