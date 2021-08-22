package main

import (
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
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
			spew.Fdump(fms, self)
			return
		}
		if fms.Flag('+') {
			io.WriteString(fms, self.fmt(true))
			return
		}
		io.WriteString(fms, self.fmt(true))
	default:
		io.WriteString(fms, self.fmt(true))
	}
}

func (self ParseErr) fmt(expand bool) (out string) {
	row, col := rowCol(self.Source, self.Cursor)
	spf(&out, `%v:%v`, row, col)

	if self.Cause != nil {
		sep(&out, `: `)
		if expand {
			spf(&out, `%+v`, self.Cause)
		} else {
			spf(&out, `%v`, self.Cause)
		}
	}

	if len(self.Snippet) > 0 {
		sep(&out, `; found: `)
		spf(&out, `%q`, self.Snippet)
	}
	return
}

func rowCol(str string, cursor int) (row int, col int) {
	for i, char := range str {
		if char == '\r' && i < len(str)-2 && str[i+1] == '\n' {
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
