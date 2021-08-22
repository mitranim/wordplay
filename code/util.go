package main

import (
	"context"
	"fmt"
	"io"
	h "net/http"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	x "github.com/mitranim/gax"
	"github.com/mitranim/try"
)

type (
	Rew = h.ResponseWriter
	Req = h.Request
	Ctx = context.Context
)

var (
	E  = x.E
	AP = x.AP
)

const (
	SHORT_SNIPPET_LEN = 64
)

func init() {
	spew.Config.Indent = "  "
	spew.Config.ContinueOnMethod = true
}

func bytesToStringAlloc(bytes []byte) string   { return string(bytes) }
func stringToBytesAlloc(input string) []byte   { return []byte(input) }
func bytesToMutableString(input []byte) string { return *(*string)(unsafe.Pointer(&input)) }

// Fixed size because it's simpler and we only need ASCII support.
// Used by pointer because large size = slow copying.
// Simpler and faster than bitset.
type charset [128]bool

func (self *charset) has(val int) bool      { return val < len(self) && self[val] }
func (self *charset) hasByte(val byte) bool { return self.has(int(val)) }
func (self *charset) hasRune(val rune) bool { return self.has(int(val)) }

func (self *charset) str(str string) *charset {
	for _, char := range str {
		self[char] = true
	}
	return self
}

func (self *charset) union(set *charset) *charset {
	for i, ok := range set {
		if ok {
			self[i] = true
		}
	}
	return self
}

var (
	charsetSpace      = new(charset).str(" \t\v")
	charsetNewline    = new(charset).str("\r\n")
	charsetPunct      = new(charset).str("#()[];,")
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
	for i, char := range input {
		switch char {
		case '\n', '\r':
			return input[:i]
		}

		if i > limit {
			return input[:i] + "…"
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

func strHas(str string, set *charset) bool {
	for _, char := range str {
		if set.hasRune(char) {
			return true
		}
	}
	return false
}

func appendNewlines(buf []byte) []byte {
	return append(buf, "\n\n"...)
}

func appendJoined(buf []byte, sep string, vals []string) []byte {
	for i, val := range vals {
		if i > 0 {
			buf = append(buf, sep...)
		}
		buf = append(buf, val...)
	}
	return buf
}

func writeString(out io.Writer, val string) {
	try.Int(io.WriteString(out, val))
}
