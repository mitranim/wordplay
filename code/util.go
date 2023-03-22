package main

import (
	"fmt"
	"path/filepath"

	"github.com/mitranim/gg"
)

const (
	SRC_FILE          = `../readme.md`
	SHORT_SNIPPET_LEN = 64
)

// Fixed size because it's simpler and we only need ASCII support.
// Used by pointer because large size = slow copying.
// Simpler and faster than bitset.
type charset [128]bool

func (self *charset) has(val int) bool      { return val < len(self) && self[val] }
func (self *charset) hasByte(val byte) bool { return self.has(int(val)) }
func (self *charset) hasRune(val rune) bool { return self.has(int(val)) }

func (self *charset) hasRunes(src string) bool {
	for _, char := range src {
		if self.hasRune(char) {
			return true
		}
	}
	return false
}

func (self *charset) str(str string) *charset {
	for _, char := range str {
		self[char] = true
	}
	return self
}

func (self *charset) union(set *charset) *charset {
	for ind, ok := range set {
		if ok {
			self[ind] = true
		}
	}
	return self
}

var (
	charsetSpace      = new(charset).str(" \t\v")
	charsetNewline    = new(charset).str("\r\n")
	charsetPunct      = new(charset).str(`#()[];,`)
	charsetWhitespace = new(charset).union(charsetSpace).union(charsetNewline)
	charsetDelim      = new(charset).union(charsetWhitespace).union(charsetPunct)
)

func counter(n int) []struct{} { return make([]struct{}, n) }

func sep(ptr *string, sep string) {
	if len(*ptr) > 0 {
		*ptr += sep
	}
}

func spf(ptr *string, pattern string, args ...interface{}) {
	*ptr += fmt.Sprintf(pattern, args...)
}

func snippet(input string, limit int) string {
	for ind, char := range input {
		switch char {
		case '\n', '\r':
			return input[:ind]
		}

		if ind > limit {
			return input[:ind] + "…"
		}
	}
	return input
}

func leadingNewlineSize(str string) int {
	if len(str) >= 2 && str[0] == '\r' && str[1] == '\n' {
		return 2
	}
	if len(str) >= 1 && (str[0] == '\r' || str[0] == '\n') {
		return 1
	}
	return 0
}

func appendNewlineIfNeeded(buf []byte) []byte {
	if len(buf) > 0 {
		return appendNewline(buf)
	}
	return buf
}

func appendNewline(buf []byte) []byte {
	return append(buf, '\n')
}

func appendJoined(buf []byte, sep string, vals []string) []byte {
	for ind, val := range vals {
		if ind > 0 {
			buf = append(buf, sep...)
		}
		buf = append(buf, val...)
	}
	return buf
}

func writeFile[A gg.Text](path string, val A) {
	gg.MkdirAll(filepath.Dir(path))
	gg.WriteFile(path, val)
}

// Like `utf8.DecodeRuneInString`, but much faster in Go < 1.17, and without
// `utf8.RuneError`.
func headChar(str string) (char rune, size int) {
	for ind, val := range str {
		if ind == 0 {
			char = val
			size = len(str)
		} else {
			size = ind
			break
		}
	}
	return
}
