package main

import (
	"strings"
	"testing"
	"time"
)

func Test_charset(t *testing.T) {
	testCharset(t, new(charset), "")
	testCharset(t, charsetSpace, " \t\v")
	testCharset(t, charsetNewline, "\r\n")
	testCharset(t, charsetWhitespace, " \t\v\r\n")
	testCharset(t, charsetDelim, " \t\v\r\n#()[];,")
}

func testCharset(t *testing.T, set *charset, chars string) {
	for i := 0; i <= 256; i++ {
		if strings.ContainsRune(chars, rune(i)) {
			continue
		}

		if set.has(i) {
			t.Fatalf("charset shouldn't contain %#0.2x", i)
		}
	}

	for _, char := range chars {
		if !set.hasRune(char) {
			t.Fatalf("charset should contain %#0.2x", char)
		}
	}
}

func Test_findIsoTime(t *testing.T) {
	const src = `a206523 (back up code (work in progress), 2020-10-26T11:51:32+03:00)`

	exp := time.Date(2020, 10, 26, 11, 51, 32, 0, tryTimeLoc(time.LoadLocation(`Europe/Moscow`)))
	res := findIsoTime(src)

	if !exp.Equal(res) {
		t.Fatalf("expected to parse %q into %q, got %q", src, exp, res)
	}
}
