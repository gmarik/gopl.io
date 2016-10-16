package main

import (
	"io"
	"reflect"
	"testing"
)

var (
	_    io.Writer = new(ByteCounter)
	w, _           = CountingWriter(nil)
	_    io.Writer = w
)

func Test_Counters(t *testing.T) {

	tcases := []struct {
		data  []byte
		bytec ByteCounter
		wordc WordCounter
		linec LineCounter
	}{
		{
			data:  []byte(""),
			bytec: ByteCounter(0),
			wordc: WordCounter(0),
			linec: LineCounter(0),
		},
		{
			data:  []byte("hello world"),
			bytec: ByteCounter(11),
			wordc: WordCounter(2),
			linec: LineCounter(1),
		},
		{
			data:  []byte("h e l l o\nworld"),
			bytec: ByteCounter(15),
			wordc: WordCounter(6),
			linec: LineCounter(2),
		},
	}

	for _, tc := range tcases {
		t.Run("ByteCounter", func(t *testing.T) {
			var (
				got ByteCounter
				exp = tc.bytec
			)

			n, err := got.Write(tc.data)
			if err != nil {
				t.Fatal("Unexpected error:", err)
			}

			if n != len(tc.data) {
				t.Errorf("\nExp: %v\nGot: %v", len(tc.data), n)
			}

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %v\nGot: %v", exp, got)
			}
		})

		t.Run("WordCounter", func(t *testing.T) {
			var (
				got WordCounter
				exp = tc.wordc
			)

			n, err := got.Write(tc.data)
			if err != nil {
				t.Fatal("Unexpected error:", err)
			}

			if n != len(tc.data) {
				t.Errorf("\nExp: %v\nGot: %v", len(tc.data), n)
			}

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %v\nGot: %v", exp, got)
			}
		})

		t.Run("LineCounter", func(t *testing.T) {
			var (
				got LineCounter
				exp = tc.linec
			)

			n, err := got.Write(tc.data)
			if err != nil {
				t.Fatal("Unexpected error:", err)
			}

			if n != len(tc.data) {
				t.Errorf("\nExp: %v\nGot: %v", len(tc.data), n)
			}

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %v\nGot: %v", exp, got)
			}
		})
	}

	t.Run("CountingWriter", func(t *testing.T) {
		var (
			cw, count = CountingWriter(new(WordCounter))
			exp       int64
		)

		for _, tc := range tcases {
			_, err := cw.Write(tc.data)
			if err != nil {
				t.Fatal("Unexpected error:", err)
			}

			var got = *count
			exp += int64(tc.bytec)

			if exp != got {
				t.Errorf("\nExp: %v\nGot: %v", exp, got)
			}
		}
	})
}
