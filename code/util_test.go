package main

import (
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func Test_charset(t *testing.T) {
	defer gtest.Catch(t)

	testCharset(new(charset), ``)
	testCharset(charsetSpace, " \t\v")
	testCharset(charsetNewline, "\r\n")
	testCharset(charsetWhitespace, " \t\v\r\n")
	testCharset(charsetDelim, " \t\v\r\n#()[];,")
}

func testCharset(set *charset, chars string) {
	for ind := range gg.Iter(256) {
		if strings.ContainsRune(chars, rune(ind)) {
			continue
		}

		if set.has(ind) {
			panic(gg.Errf(`charset must not contain %#0.2x`, ind))
		}
	}

	for _, char := range chars {
		if !set.hasRune(char) {
			panic(gg.Errf(`charset must contain %#0.2x`, char))
		}
	}
}
