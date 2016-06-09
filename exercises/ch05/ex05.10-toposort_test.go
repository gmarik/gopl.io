// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 136.

//TODO: Exercise 5.10: Rewrite topoSort to use maps instead of slices and eliminate the initial sort. Verify that the results, though nondeterministic, are valid topological orderings.
//Exercise 5.11: The instructor of the linear algebra course decides that calculus is now a prerequisite. Extend the topoSort function to report cycles

package main

import (
	"bytes"
	"fmt"
	"sort"
	"testing"
)

// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

// The toposort program prints the nodes of a DAG in topological order.
//
func topoSortSlice(m map[string][]string) []string {
	var (
		order    []string
		visitAll func(items []string)

		seen = make(map[string]bool)
	)

	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}

func topoSort(m map[string][]string) []string {
	var (
		order    []string
		visitAll func([]string)

		seen = make(map[string]bool)
	)

	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	visitAll(keys)
	return order
}

func Test_TopoSortSlice(t *testing.T) {
	var (
		w = &bytes.Buffer{}
	)

	for i, course := range topoSortSlice(prereqs) {
		fmt.Fprintf(w, "%d:\t%s\n", i+1, course)
	}

	exp := `1:	intro to programming
2:	discrete math
3:	data structures
4:	algorithms
5:	linear algebra
6:	calculus
7:	formal languages
8:	computer organization
9:	compilers
10:	databases
11:	operating systems
12:	networks
13:	programming languages
`
	got := w.String()

	if got != exp {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}

}

func Test_TopoSort(t *testing.T) {
	t.Skip("TODO: clarify the requirements and implement")
	return

	var (
		w = &bytes.Buffer{}
	)

	for i, course := range topoSort(prereqs) {
		fmt.Fprintf(w, "%d:\t%s\n", i+1, course)
	}

	exp := `1:	intro to programming
2:	discrete math
3:	data structures
4:	algorithms
5:	linear algebra
6:	calculus
7:	formal languages
8:	computer organization
9:	compilers
10:	databases
11:	operating systems
12:	networks
13:	programming languages
`
	got := w.String()

	if got != exp {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}
