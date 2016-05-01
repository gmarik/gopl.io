package main

import (
	"bytes"
	"testing"
)

func TestOutline(t *testing.T) {
	html := `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>The Go Programming Language</title>

<link type="text/css" rel="stylesheet" href="/lib/godoc/style.css">

<link rel="search" type="application/opensearchdescription+xml" title="godoc" href="/opensearch.xml" />

<link rel="stylesheet" href="/lib/godoc/jquery.treeview.css">
`

	// set the output
	buf := &bytes.Buffer{}
	out = buf
	outline(bytes.NewReader([]byte(html)))

	got := buf.String()
	exp := `  <html []>
    <head []>
      <meta [http-equiv content]/>
      <meta [name content]/>
      <meta [name content]/>
      <title []>
      </title>
      <link [type rel href]/>
      <link [rel type title href]/>
      <link [rel href]/>
    </head>
    <body []/>
  </html>
`
	if got != exp {
		// t.Errorf("\nExp: %v\nGot: %v", []byte(exp), []byte(got))
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}

func TestStop(t *testing.T) {
	html := `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>The Go Programming Language</title>

<link id="stop" type="text/css" rel="stylesheet" href="/lib/godoc/style.css">

<link rel="search" type="application/opensearchdescription+xml" title="godoc" href="/opensearch.xml" />

<link rel="stylesheet" href="/lib/godoc/jquery.treeview.css">
`

	// set the output
	buf := &bytes.Buffer{}
	out = buf
	stopAt = "stop"
	outline(bytes.NewReader([]byte(html)))

	got := buf.String()
	exp := `  <html []>
    <head []>
      <meta [http-equiv content]/>
      <meta [name content]/>
      <meta [name content]/>
      <title []>
      </title>
`
	if got != exp {
		// t.Errorf("\nExp: %v\nGot: %v", []byte(exp), []byte(got))
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}
