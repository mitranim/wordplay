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

func (self ParseErr) Error() string { return self.fmt(false) }

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

func (self ParseErr) fmt(expand bool) string {
	row, col := RowCol(self.Source, self.Cursor)

	var buf gg.Buf
	buf.AppendString(`<row:col> `)
	buf.AppendInt(row)
	buf.AppendString(`:`)
	buf.AppendInt(col)

	if len(self.Snippet) > 0 {
		buf.AppendString(`; followed by: `)
		buf.AppendString(gg.Quote(self.Snippet))
	}

	if self.Cause != nil {
		buf.AppendString(`; `)
		if expand {
			buf.Fprintf(`%+v`, self.Cause)
		} else {
			buf.AppendAny(self.Cause)
		}
	}

	return buf.String()
}

func errNewline(delim rune) error {
	return gg.Errf(`expected closing %q, found newline`, delim)
}

func errEof(delim rune) error {
	return gg.Wrapf(io.EOF, `expected closing %q, found EOF`, delim)
}
