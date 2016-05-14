package main

import "testing"

// Exercise 5.9: Write a function expand(s string, f func(string) string) string that replaces each substring “$foo” within s by the text returned by f("foo").
//

// see os.Expand for full implementation
func expand(s string, mapping func(string) string) string {
	buf := make([]byte, 0, len(s))

	i := 0
	for j := 0; j < len(s); j += 1 {
		if s[j] == '$' && j+1 < len(s) {
			buf = append(buf, s[i:j]...)
			name, width := getName(s[j+1:])
			buf = append(buf, mapping(name)...)
			j += width
			i = j + 1
		}
	}

	return string(buf)
}

func getName(s string) (_ string, width int) {
	for i := 0; i < len(s); i += 1 {
		if 'a' <= s[i] && s[i] <= 'z' {
			width += 1
			continue
		}
		break
	}
	return s[:width], width
}

func TestExpand(t *testing.T) {

	m := func(s string) string {
		return map[string]string{
			"hello": "Hello",
			"wo":    "world",
		}[s]
	}

	exp := "Hello world"
	got := expand("$hello $wo", m)

	if got != exp {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}

}

func Test_getName(t *testing.T) {

	exp := "name"
	got, _ := getName("name asdf")
	if got != exp {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}

}
