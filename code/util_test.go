package main

import (
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func Test_Charset(t *testing.T) {
	defer gtest.Catch(t)

	testCharset(Charset(nil), ``)
	testCharset(CharsetSpace, " \t\v")
	testCharset(CharsetNewline, "\r\n")
	testCharset(CharsetWhitespace, " \t\v\r\n")
	testCharset(CharsetDelim, "()[]©\r\n")
}

func testCharset(set Charset, chars string) {
	for ind := range gg.Iter(256) {
		if strings.ContainsRune(chars, rune(ind)) {
			continue
		}

		if set.HasInt(ind) {
			panic(gg.Errf(`charset must not contain int %#0.2x (%v); charset: %v`, ind, ind, set))
		}
	}

	for _, char := range chars {
		if !set.HasRune(char) {
			panic(gg.Errf(`charset must contain rune %#0.2x (%q); charset: %v`, char, char, set))
		}
	}
}
