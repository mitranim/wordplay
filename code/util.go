package main

import (
	"context"
	"fmt"
	h "net/http"
	"strings"
	"unicode/utf8"
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

func stringSliceMut(vals []string, fun func(string) string) {
	for i := range vals {
		vals[i] = fun(vals[i])
	}
}

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

func isNonWhitespace(char rune) bool { return !isWhitespace(char) }

func isMeaningsDelim(char rune) bool {
	return isNewline(char) || char == ')'
}

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

func headChar(str string) (rune, int) {
	if strings.HasPrefix(str, "\r\n") {
		return '\n', len("\r\n")
	}
	if len(str) > 0 && (str[0] == '\r' || str[0] == '\n') {
		return '\n', 1
	}
	return utf8.DecodeRuneInString(str)
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

func assert(ok bool) {
	if !ok {
		panic("internal error: failed a condition that should never be failed, see the stacktrace")
	}
}

func strAppendSwap(one *string, many *[]string, arr *[2]string, val string) {
	if len(*one) == 0 {
		if len(*many) == 0 {
			*one = val
		} else {
			*many = append(*many, val)
			*arr = [2]string{}
		}
	} else {
		if cap(*many) == 0 {
			*many = arr[:]
			(*many)[0] = *one
			(*many)[1] = val
		} else {
			*many = append(*many, *one, val)
			*arr = [2]string{}
		}
		*one = ""
	}
}

func appendJoinedWith(buf []byte, sep string, val string, vals []string) []byte {
	if len(val) > 0 {
		buf = append(buf, val...)
		if len(vals) > 0 {
			buf = append(buf, sep...)
		}
	}

	buf = appendJoined(buf, sep, vals)
	return buf
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
