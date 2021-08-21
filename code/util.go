package main

import (
	"context"
	"fmt"
	h "net/http"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	x "github.com/mitranim/gax"
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

func isSpace(char rune) bool {
	switch char {
	case ' ', '\t', '\v':
		return true
	default:
		return false
	}
}

func isNewline(char rune) bool {
	switch char {
	case '\r', '\n':
		return true
	default:
		return false
	}
}

func isWhitespace(char rune) bool {
	return isSpace(char) || isNewline(char)
}

func isNonNewline(char rune) bool { return !isNewline(char) }

func isDelimPunct(char rune) bool {
	switch char {
	case '#', '(', ')', '[', ']', ';':
		return true
	default:
		return false
	}
}

func isDelim(char rune) bool {
	return isWhitespace(char) || isDelimPunct(char)
}

func isNonDelim(char rune) bool { return !isDelim(char) }

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

func strHas(str string, fun func(rune) bool) bool {
	for _, char := range str {
		if fun(char) {
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
