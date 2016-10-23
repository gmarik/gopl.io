package main

import (
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

var (
	_ io.Reader = (*Reader)(nil)
)

func TestReader(t *testing.T) {
	tcases := []struct {
		string string
		reader io.Reader
	}{
		{
			string: "hello world",
			reader: NewReader("hello world"),
		},
		{
			string: "hell",
			reader: NewLimitReader(NewReader("hello world"), 4),
		},
		{
			string: "hello",
			reader: NewLimitReader(NewReader("hello"), 500),
		},
	}

	for _, tc := range tcases {
		t.Run("", func(t *testing.T) {
			var (
				data, err = ioutil.ReadAll(tc.reader)
				exp       = tc.string
				got       = string(data)
			)

			if err != nil {
				t.Fatal("Unexpected error:", err)
			}

			if !reflect.DeepEqual(exp, got) {
				t.Errorf("\nExp: %v\nGot: %v", exp, got)
			}
		})

		t.Run("EOF", func(t *testing.T) {
			var (
				_, err = tc.reader.Read(nil)
			)

			if err != io.EOF {
				t.Fatalf("EOF expected, got: %v", err)
			}
		})
	}

}
