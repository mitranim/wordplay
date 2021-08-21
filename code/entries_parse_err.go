package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

type ParseErr struct {
	Row     int
	Col     int
	Source  string
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
			fms.Write(stringToBytesAlloc(self.fmt(true)))
			return
		}
		fms.Write(stringToBytesAlloc(self.fmt(true)))
	default:
		fms.Write(stringToBytesAlloc(self.fmt(true)))
	}
}

func (self ParseErr) fmt(expand bool) (out string) {
	spf(&out, `%v:%v`, self.Row+1, self.Col+1)

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
