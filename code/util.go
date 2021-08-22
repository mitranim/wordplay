package main

import (
	"context"
	"fmt"
	h "net/http"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	x "github.com/mitranim/gax"
	e "github.com/pkg/errors"
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

func init() {
	spew.Config.Indent = "  "
	spew.Config.ContinueOnMethod = true
}

func bytesToStringAlloc(bytes []byte) string   { return string(bytes) }
func stringToBytesAlloc(input string) []byte   { return []byte(input) }
func bytesToMutableString(input []byte) string { return *(*string)(unsafe.Pointer(&input)) }

type charset = boolCharset

// Fixed size because it's simpler and we only need ASCII support.
type boolCharset [128]bool

func (self *boolCharset) has(val int) bool      { return val < len(self) && self[val] }
func (self *boolCharset) hasByte(val byte) bool { return self.has(int(val)) }
func (self *boolCharset) hasRune(val rune) bool { return self.has(int(val)) }

func (self *boolCharset) str(str string) *boolCharset {
	for _, char := range str {
		self[char] = true
	}
	return self
}

// ASCII only. Substitute for missing `uint128`. Methods use pointer receiver
// because the value is larger than a machine word, and copying has a
// measurable cost in hotpaths.
type bitCharset [128 / 8]byte

func (self *bitCharset) hasByte(val byte) bool { return self.has(int(val)) }
func (self *bitCharset) hasRune(val rune) bool { return self.has(int(val)) }

func (self *bitCharset) has(val int) bool {
	return isAscii(val) && (self[byteInd(val)]&(byteBit(val))) != 0
}

func (self *bitCharset) add(val int) {
	reqAscii(val)
	self[byteInd(val)] |= byteBit(val)
}

func (self *bitCharset) del(val int) {
	reqAscii(val)
	self[byteInd(val)] ^= byteBit(val)
}

func (self *bitCharset) str(str string) *bitCharset {
	for _, char := range str {
		self.add(int(char))
	}
	return self
}

func byteInd(val int) int  { return val >> 3 }
func byteBit(val int) byte { return byte(1 << (val & 7)) }
func isAscii(val int) bool { return val >= 0 && val < 128 }

func reqAscii(val int) {
	if !isAscii(val) {
		panic(e.Errorf(`%q is not ASCII`, val))
	}
}

var (
	charsetSpace      = new(charset).str(" \t\v")
	charsetNewline    = new(charset).str("\r\n")
	charsetWhitespace = new(charset).str(" \t\v\r\n")
	charsetDelim      = new(charset).str(" \t\v\r\n#()[];,")
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

// Significantly faster than using `strings.HasPrefix` and/or
// `utf8.DecodeRuneInString`.
func headChar(str string) (char rune, size int) {
	if len(str) >= 2 && str[0] == '\r' && str[1] == '\n' {
		return '\n', 2
	}
	if len(str) >= 1 && (str[0] == '\r' || str[0] == '\n') {
		return '\n', 1
	}

	for i, val := range str {
		if i == 0 {
			char = val
			size = len(str)
		} else {
			size = i
			break
		}
	}

	return
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
