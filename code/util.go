package main

import (
	"fmt"
	"io"
	"path/filepath"
	r "reflect"
	"strings"

	"github.com/mitranim/gg"
	"github.com/mitranim/rf"
)

const (
	SRC_FILE          = `../readme.md`
	SHORT_SNIPPET_LEN = 64
)

// Performs much better than equivalent map-based set.
type Charset []bool

func (self Charset) HasInt(val int) bool   { return val < len(self) && self[val] }
func (self Charset) HasByte(val byte) bool { return self.HasInt(int(val)) }
func (self Charset) HasRune(val rune) bool { return self.HasInt(int(val)) }

func (self Charset) HasRunes(src string) bool {
	for _, char := range src {
		if self.HasRune(char) {
			return true
		}
	}
	return false
}

func (self Charset) Add(val int) Charset {
	diff := val - len(self)
	if diff >= 0 {
		self = gg.GrowLen(self, diff+1)
	}
	self[val] = true
	return self
}

func (self Charset) AddStr(src string) Charset {
	for _, char := range src {
		self = self.Add(int(char))
	}
	return self
}

func (self Charset) AddFrom(src Charset) Charset {
	for ind, ok := range src {
		if ok {
			self = self.Add(ind)
		}
	}
	return self
}

var (
	CharsetSpace      = Charset(nil).AddStr(" \t\v")
	CharsetNewline    = Charset(nil).AddStr("\r\n")
	CharsetWhitespace = Charset(nil).AddFrom(CharsetSpace).AddFrom(CharsetNewline)
	CharsetDelim      = Charset(nil).AddStr(`()[]©`).AddFrom(CharsetNewline)
)

func Snippet(src string, limit uint) string {
	return gg.Ellipsis(UntilNewline(src), limit)
}

func UntilNewline[A gg.Text](src A) A {
	ind := strings.IndexAny(gg.ToString(src), "\r\n")
	if ind >= 0 {
		return src[:ind]
	}
	return src
}

func LeadingNewlineSize(str string) int {
	if len(str) >= 2 && str[0] == '\r' && str[1] == '\n' {
		return 2
	}
	if len(str) >= 1 && (str[0] == '\r' || str[0] == '\n') {
		return 1
	}
	return 0
}

func AppendJoined(buf []byte, sep string, vals []string) []byte {
	for ind, val := range vals {
		if ind > 0 {
			buf = append(buf, sep...)
		}
		buf = append(buf, val...)
	}
	return buf
}

func WriteFile[A gg.Text](path string, val A) {
	gg.MkdirAll(filepath.Dir(path))
	gg.WriteFile(path, val)
}

// Like `utf8.DecodeRuneInString`, but much faster in Go < 1.17, and without
// `utf8.RuneError`.
func HeadChar(str string) (char rune, size int) {
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

func ReqField[Out, Src any](src Src, off uintptr) Out {
	typ := gg.Type[Src]()
	field := rf.TypeOffsetFields(typ)[off][0]
	val := r.ValueOf(&src).Elem().FieldByIndex(field.Index).Interface().(Out)

	if gg.IsZero(val) {
		panic(gg.Errf(`unexpected zero value of field %q of %v`, field.Name, typ))
	}
	return val
}

func SortReverse[A gg.Lesser[A]](val []A) {
	gg.Sort(val)
	gg.Reverse(val)
}

func JoinLinesSparse[A gg.Text](src ...A) A {
	return gg.ToText[A](gg.JoinOpt(src, gg.Newline+gg.Newline))
}

// TODO move to `gg`. Needs tests. TODO another version that takes char index.
func RowCol(src string, byteInd int) (row int, col int) {
	for ind, char := range src {
		if ind >= byteInd {
			break
		}

		if char == '\r' && ind < len(src)-2 && src[ind+1] == '\n' {
			continue
		}

		if char == '\r' || char == '\n' {
			row++
			col = 0
			continue
		}

		col++
	}

	row++
	col++
	return
}

func Unquote(src string) string {
	const quoteAscii = '"'
	const quoteLeft = `“`
	const quoteRight = `”`

	size := len(src)

	if size > 1 && src[0] == quoteAscii && src[size-1] == quoteAscii {
		inner := src[1 : size-1]

		// This avoids stripping opening and closing quotes from a string that
		// actually contains multiple pairs of quotes, breaking two pairs.
		if !strings.ContainsRune(inner, quoteAscii) {
			return inner
		}

		return src
	}

	if strings.HasPrefix(src, quoteLeft) && strings.HasSuffix(src, quoteRight) {
		inner := src[len(quoteLeft) : size-len(quoteRight)]

		// This avoids stripping opening and closing quotes from a string that
		// actually contains multiple pairs of quotes, breaking two pairs.
		if !strings.Contains(inner, quoteLeft) && !strings.Contains(inner, quoteRight) {
			return inner
		}

		return src
	}

	return src
}

// TODO move to `gg`.
func AppendNewlineOpt[A ~string](val A) A {
	if len(val) > 0 && !gg.HasNewlineSuffix(val) {
		return val + A(gg.Newline)
	}
	return val
}

// Permissive version of `fmt.Fprintln`: does nothing if output is nil.
// Also doesn't automatically space-out adjacent strings.
// TODO move to `gg`.
func Fprintln(out io.Writer, msg ...any) {
	if out != nil {
		gg.Write(out, AppendNewlineOpt(gg.Str(msg...)))
	}
}

// Permissive version of `fmt.Fprintf`: does nothing if output is nil.
// TODO move to `gg`.
func Fprintf(out io.Writer, pat string, arg ...any) {
	if out != nil {
		fmt.Fprintf(out, pat, gg.NoEscUnsafe(arg)...)
	}
}

var StrNorm = strings.NewReplacer(
	"\x00", ``,
	"\u0000", ``,
	"\u00a0", ` `,
	`’`, `'`,
	`“`, `"`,
	`”`, `"`,
).Replace
