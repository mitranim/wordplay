package main

import (
	"fmt"
	"io"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

type ParseErr struct {
	Source  string
	Cursor  int
	Snippet string
	Cause   error
}

func (self ParseErr) Error() string {
	return self.fmt(false)
}

func (self ParseErr) Unwrap() error { return self.Cause }

func (self ParseErr) Format(fms fmt.State, verb rune) {
	switch verb {
	case 'v':
		if fms.Flag('#') {
			gg.Write(fms, grepr.String(self))
			return
		}
		if fms.Flag('+') {
			gg.Write(fms, self.fmt(true))
			return
		}
		gg.Write(fms, self.fmt(true))
	default:
		gg.Write(fms, self.fmt(true))
	}
}

func (self ParseErr) fmt(expand bool) (out string) {
	row, col := rowCol(self.Source, self.Cursor)
	spf(&out, `<row:col> %v:%v`, row, col)

	if self.Cause != nil {
		sep(&out, `: `)
		if expand {
			spf(&out, `%+v`, self.Cause)
		} else {
			spf(&out, `%v`, self.Cause)
		}
	}

	if len(self.Snippet) > 0 {
		sep(&out, `; following text: `)
		spf(&out, `%q`, self.Snippet)
	}
	return
}

func rowCol(src string, pos int) (row int, col int) {
	var chars int
	for ind, char := range src {
		chars++
		if chars == pos {
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

func errNewline(delim rune) error {
	return gg.Errf(`expected closing %q, found newline`, delim)
}

func errEof(delim rune) error {
	return gg.Wrapf(io.EOF, `expected closing %q, found EOF`, delim)
}
