package main

import (
	"reflect"
	"testing"
)

func TestTree_String(t *testing.T) {
	tcases := []struct {
		string string
		tree   *tree
	}{
		{
			tree:   nil,
			string: "<nil>",
		},

		{
			tree: &tree{
				value: 1,
				left: &tree{value: 2,
					right: &tree{value: 22},
				},
				right: &tree{value: 3,
					left: &tree{value: 31},
				},
			},
			string: "1\n  2\n    22\n  3\n    31\n",
		},
	}

	for _, tc := range tcases {
		t.Run("", func(t *testing.T) {
			var (
				exp = tc.string
				got = tc.tree.String()
			)

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %q\nGot: %q", exp, got)
			}
		})
	}
}
