package main

import (
	"reflect"
	"testing"
)

func TestIntSet_Basic(t *testing.T) {
	cases := []struct {
		desc    string
		set     IntSet
		adds    []int
		hasnots []int
		len     int
		string  string
	}{
		{
			desc:    "empty",
			set:     IntSet{},
			adds:    nil,
			hasnots: []int{2},
			len:     0,
			string:  "{}",
		},
		{
			desc:    "single element, 0",
			set:     *NewIntSet(0),
			adds:    []int{0},
			hasnots: []int{2},
			len:     1,
			string:  "{0}",
		},
		{
			desc:    "single element",
			set:     *NewIntSet(1),
			adds:    []int{1},
			hasnots: []int{2},
			len:     1,
			string:  "{1}",
		},
		{
			desc:    "few elements",
			set:     *NewIntSet(1, 2, 124),
			adds:    []int{1, 2, 124},
			hasnots: []int{0, 5, 7, 8, 500},
			len:     3,
			string:  "{1 2 124}",
		},
	}

	for _, c := range cases {
		t.Run(c.desc+"/Add", func(t *testing.T) {
			var got IntSet
			for _, v := range c.adds {
				got.Add(v)
			}
			if !reflect.DeepEqual(c.set, got) {
				t.Errorf("\nExp: %v\nGot: %v", c.set, got)
			}
		})

		t.Run(c.desc+"/AddAll", func(t *testing.T) {
			var got IntSet
			got.AddAll(c.adds...)
			if !reflect.DeepEqual(c.set, got) {
				t.Errorf("\nExp: %v\nGot: %v", c.set, got)
			}
		})

		t.Run(c.desc+"/Len", func(t *testing.T) {
			var (
				got = c.set.Len()
				set = c.len
			)

			if set != got {
				t.Errorf("\nExp: %v\nGot: %v", set, got)
			}
		})

		t.Run(c.desc+"/Copy", func(t *testing.T) {
			var (
				got = c.set.Copy()
				set = &c.set
			)

			if !reflect.DeepEqual(set, got) {
				t.Errorf("\nExp: %v\nGot: %v", set, got)
			}
		})

		t.Run(c.desc+"/Has", func(t *testing.T) {
			for _, h := range c.adds {
				if false == c.set.Has(h) {
					t.Errorf("Exp to have %v", h)
				}
			}
		})

		t.Run(c.desc+"/Has/not", func(t *testing.T) {
			for _, h := range c.hasnots {
				if true == c.set.Has(h) {
					t.Errorf("Exp to not have %v", h)
				}
			}
		})

		t.Run(c.desc+"/Remove", func(t *testing.T) {
			if c.set.Len() != c.len {
				t.Fatalf("invalid state: %v", c.set)
			}

			set := c.set.Copy()
			for _, h := range c.adds {
				if false == set.Remove(h) {
					t.Errorf("Exp to remove %v", h)
				}
			}
		})

		t.Run(c.desc+"/Remove/not", func(t *testing.T) {
			for _, h := range c.hasnots {
				if true == c.set.Remove(h) {
					t.Errorf("Exp to not remove %v", h)
				}
			}
		})

		t.Run(c.desc+"/Clear", func(t *testing.T) {
			// pre-clear validation
			set := c.set.Copy()
			if c.set.Len() != c.len {
				t.Fatalf("invalid state: %v", c.set)
			}

			set.Clear()

			if set.Len() != 0 {
				t.Errorf("Exp to be emtpy, but is %v", set)
			}
		})

		t.Run(c.desc+"/Elems", func(t *testing.T) {
			var (
				got = c.set.Elems()
				exp = c.adds
			)

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %#v\nGot: %#v", exp, got)
			}
		})

		t.Run(c.desc+"/String", func(t *testing.T) {
			var (
				got = c.set.String()
				exp = c.string
			)

			if exp != got {
				t.Errorf("\nExp: %v\nGot: %v", exp, got)
			}
		})
	}
}
func TestIntSet_Ops(t *testing.T) {
	tcases := []struct {
		desc string

		A, B IntSet

		union, intersect IntSet
		differ, symmdiff IntSet
	}{
		{
			desc:      "A,B =0",
			A:         IntSet{},
			B:         IntSet{},
			union:     IntSet{},
			intersect: IntSet{},
			differ:    IntSet{},
			symmdiff:  IntSet{},
		},

		{
			desc:      "A>0, B=0",
			A:         *NewIntSet(1, 2, 3, 124),
			B:         *NewIntSet(),
			union:     *NewIntSet(1, 2, 3, 124),
			intersect: IntSet{},
			differ:    IntSet{},
			symmdiff:  IntSet{},
		},

		{
			desc:      "A=0,B>0",
			A:         *NewIntSet(),
			B:         *NewIntSet(1, 2, 3, 124),
			union:     *NewIntSet(1, 2, 3, 124),
			intersect: IntSet{},
			differ:    IntSet{},
			symmdiff:  IntSet{},
		},
		{
			desc:      "A=1, B>1",
			A:         *NewIntSet(5),
			B:         *NewIntSet(1, 2, 3, 124),
			union:     *NewIntSet(1, 2, 3, 5, 124),
			intersect: IntSet{},
			differ:    *NewIntSet(5),
			symmdiff:  *NewIntSet(1, 2, 3, 5, 124),
		},

		{
			desc:      "A+,B+, A > B",
			A:         *NewIntSet(1, 2, 3, 124),
			B:         *NewIntSet(5, 12),
			union:     *NewIntSet(1, 2, 3, 5, 12, 124),
			intersect: *NewIntSet(),
			differ:    *NewIntSet(1, 2, 3, 124),
			symmdiff:  *NewIntSet(1, 2, 3, 5, 12, 124),
		},

		{
			desc:      "A = B",
			A:         *NewIntSet(1, 2, 3, 124),
			B:         *NewIntSet(1, 2, 3, 124),
			union:     *NewIntSet(1, 2, 3, 124),
			intersect: *NewIntSet(1, 2, 3, 124),
			differ:    IntSet{},
			symmdiff:  IntSet{},
		},

		{
			desc:      "A > B, with intersection",
			A:         *NewIntSet(1, 2, 3, 124),
			B:         *NewIntSet(2, 13, 4),
			union:     *NewIntSet(1, 2, 3, 4, 13, 124),
			intersect: *NewIntSet(2),
			differ:    *NewIntSet(1, 3, 124),
			symmdiff:  *NewIntSet(1, 3, 4, 13, 124),
		},
	}

	for _, tc := range tcases {
		t.Run("UnionWith/"+tc.desc, func(t *testing.T) {
			var (
				exp = &tc.union
				got = tc.A.Copy()
			)

			got.UnionWith(&tc.B)

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %v\nGot: %v", exp, got)
			}
		})

		t.Run("IntersectWith/"+tc.desc, func(t *testing.T) {
			var (
				exp = &tc.intersect
				got = tc.A.Copy()
			)

			got.IntersectWith(&tc.B)

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %#v\nGot: %#v", exp, got)
			}
		})

		t.Run("DifferenceWith/"+tc.desc, func(t *testing.T) {
			var (
				exp = &tc.differ
				got = tc.A.Copy()
			)

			got.DifferenceWith(&tc.B)

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %#v\nGot: %#v", exp, got)
			}
		})

		t.Run("SymmDiffWith/"+tc.desc, func(t *testing.T) {
			var (
				exp = &tc.symmdiff
				got = tc.A.Copy()
			)

			got.SymmDiffWith(&tc.B)

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %#v\nGot: %#v", exp, got)
			}
		})
	}
}
