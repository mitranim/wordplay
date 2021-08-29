package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/mitranim/try"
)

const (
	SRC_FILE          = `../readme.md`
	SHORT_SNIPPET_LEN = 64
)

// nolint:deadcode,unused
func bytesToStringAlloc(bytes []byte) string { return string(bytes) }

func stringToBytesAlloc(input string) []byte { return []byte(input) }

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

// Why do I have to write this?
func readFile(path string) []byte {
	return try.ByteSlice(os.ReadFile(path))
}

// Why do I have to write this?
func readFileString(path string) string {
	return bytesToMutableString(readFile(path))
}

// Why do I have to write this?
func writeFile(path string, val []byte) {
	try.To(os.MkdirAll(filepath.Dir(path), os.ModePerm))
	try.To(os.WriteFile(path, val, os.ModePerm))
}

// Why do I have to write this?
// nolint:deadcode,unused
func writeFileStr(path string, val string) {
	writeFile(path, stringToBytesAlloc(val))
}

// Like `utf8.DecodeRuneInString` but much faster and without `utf8.RuneError`.
func headChar(str string) (char rune, size int) {
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
