package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	h "net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
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
	SRC_FILE          = `../readme.md`
	SHORT_SNIPPET_LEN = 64
)

func init() {
	spew.Config.Indent = "  "
	// spew.Config.ContinueOnMethod = true
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

func makeCmd(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

/*
Runs a command for side effects, connecting its stdout and stderr to the parent
process.
*/
func runCmd(command string, args ...string) {
	try.To(makeCmd(command, args...).Run())
}

/*
Runs a command and returns its stdout. Stderr is connected to the parent
process.
*/
func runCmdOut(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	return bytesToMutableString(bytes.TrimSpace(try.ByteSlice(cmd.Output())))
}

var reNewline = regexp.MustCompile(`(?:\r\n|\r|\n)`)

// Seems missing from the standard library.
func splitLines(str string) []string {
	return reNewline.Split(str, -1)
}

// WTF I shouldn't have to write this.
func timeMin(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	if b.Before(a) {
		return b
	}
	return a
}

func tryTimeLoc(val *time.Location, err error) *time.Location {
	try.To(err)
	return val
}

func tryTime(val time.Time, err error) time.Time {
	try.To(err)
	return val
}

var reIsoTime = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:Z|\+\d+:\d+|\+\d+|\+\w+)`)

func findIsoTime(str string) time.Time {
	defer try.Detailf(`failed to parse time from %q`, str)
	return tryTime(time.Parse(time.RFC3339, reIsoTime.FindString(str)))
}

// Why do I have to write this?
func writeFile(path string, val []byte) {
	try.To(os.MkdirAll(filepath.Dir(path), os.ModePerm))
	try.To(os.WriteFile(path, val, os.ModePerm))
}

// Why do I have to write this?
func writeFileStr(path string, val string) {
	writeFile(path, stringToBytesAlloc(val))
}
