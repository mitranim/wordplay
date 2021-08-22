package main

import (
	"io"
	"strings"

	"github.com/mitranim/try"
	e "github.com/pkg/errors"
)

const (
	SHORT_SNIPPET_LEN = 64
)

func ParseEntries(src string) Entries {
	defer try.Detail(`failed to parse string into entries`)

	parser := MakeParser(src)
	parser.Parse()
	return parser.Entries
}

type Parser struct {
	Source  string
	Entries Entries
	cursor  int
	entry   Entry
}

func MakeParser(content string) Parser {
	return Parser{Source: content}
}

func (self *Parser) Parse() {
	for self.more() {
		if !self.scanned((*Parser).any) {
			self.fail(e.New(`unrecognized content`))
		}
	}
}

func (self *Parser) any() {
	switch {
	case self.scanned((*Parser).whitespace):
	case self.scanned((*Parser).heading):
	case self.scanned((*Parser).entryQuoted):
	case self.scanned((*Parser).entryUnquoted):
	}
}

func (self *Parser) whitespace() { self.bytesWith(charsetWhitespace) }
func (self *Parser) space()      { self.bytesWith(charsetSpace) }

func (self *Parser) heading() {
	if !self.scannedByte('#') {
		return
	}

	if !self.scanned((*Parser).space) {
		self.fail(e.New(`malformed header: expected '#' followed by space and author name`))
	}

	start := self.cursor
	self.nonNewline()

	author := strings.TrimSpace(self.from(start))
	if len(author) == 0 {
		self.fail(e.New(`malformed header: expected '#' followed by space and author name`))
	}

	self.entry.Author = author
	self.delimWhitespace()
}

func (self *Parser) entryQuoted() {
	if !self.scannedByte('"') {
		return
	}

	start := self.cursor

loop:
	for {
		char, size := self.headChar()
		self.reqMore(char, size, '"')

		switch char {
		case '"':
			phrase := strings.TrimSpace(self.from(start))
			if len(phrase) == 0 {
				self.fail(e.New(`quoted phrase is empty`))
			}

			self.entry.Phrase = phrase
			self.mov(size)
			break loop

		default:
			self.mov(size)
		}
	}

	self.entryRest()
}

func (self *Parser) nonDelim() { self.charsWithout(charsetDelim) }

func (self *Parser) entryUnquoted() {
	start := self.cursor

	if !self.scanned((*Parser).nonDelim) {
		return
	}

	self.entry.Phrase = self.from(start)
	self.entryRest()
}

func (self *Parser) entryRest() {
	self.entryMeanings()
	self.entryTags()
	self.entryFlush()
	self.delimWhitespace()
}

func (self *Parser) entryMeanings() {
	self.space()
	if !self.scannedByte('(') {
		return
	}

	depth := 1
	start := self.cursor

loop:
	for {
		char, size := self.headChar()
		self.reqMore(char, size, ')')

		switch char {
		case '(':
			depth++
			self.mov(size)

		case ')':
			depth--

			if depth == 0 {
				self.appendMeaning(self.from(start))
				self.mov(size)
				break loop
			}

			// Should be impossible.
			if depth < 0 {
				self.fail(e.New(`mismatched closing ")"`))
			}

			self.mov(size)

		case ';':
			self.appendMeaning(self.from(start))
			self.mov(size)
			start = self.cursor

		default:
			self.mov(size)
		}
	}
}

func (self *Parser) appendMeaning(val string) {
	val = strings.TrimSpace(val)
	if len(val) == 0 {
		self.fail(e.New(`unexpected empty meaning`))
	}
	self.entry.appendMeaning(val)
}

func (self *Parser) entryTags() {
	self.space()
	if !self.scannedByte('[') {
		return
	}

	start := self.cursor

loop:
	for {
		char, size := self.headChar()
		self.reqMore(char, size, ']')

		switch char {
		case '[':
			self.fail(e.New(`unexpected nested "["`))

		case ';':
			self.appendTag(self.from(start))
			self.mov(size)
			start = self.cursor

		case ']':
			self.appendTag(self.from(start))
			self.mov(size)
			break loop

		default:
			self.mov(size)
		}
	}
}

func (self *Parser) appendTag(val string) {
	val = strings.TrimSpace(val)
	if len(val) == 0 {
		self.fail(e.New(`unexpected empty tag`))
	}
	self.entry.appendTag(val)
}

func (self *Parser) entryFlush() {
	self.Entries = append(self.Entries, self.entry)
	self.entry = Entry{Author: self.entry.Author}
}

func (self *Parser) delimWhitespace() {
	self.space()

	if !self.more() {
		return
	}
	if !self.scannedNewline() {
		self.fail(e.New(`expected at least two newlines or EOF`))
	}

	if !self.more() {
		return
	}
	if !self.scannedNewline() {
		self.fail(e.New(`expected at least two newlines or EOF`))
	}

	self.whitespace()
}

func (self *Parser) newline() {
	char, size := self.headChar()
	if char == '\n' {
		self.mov(size)
	}
}

func (self *Parser) nonNewline() { self.charsWithout(charsetNewline) }

func (self *Parser) more() bool { return self.cursor < len(self.Source) }

func (self *Parser) rest() string {
	if self.more() {
		return self.Source[self.cursor:]
	}
	return ""
}

func (self *Parser) from(start int) string {
	if start < 0 {
		start = 0
	}
	if start < self.cursor {
		return self.Source[start:self.cursor]
	}
	return ""
}

func (self *Parser) headChar() (rune, int) {
	return headChar(self.rest())
}

func (self *Parser) mov(size int) { self.cursor += size }

func (self *Parser) end() { self.cursor = len(self.Source) }

func (self *Parser) scanned(fun func(*Parser)) bool {
	start := self.cursor
	fun(self)
	return self.cursor > start
}

func (self *Parser) scannedNewline() bool {
	return self.scanned((*Parser).newline)
}

func (self *Parser) scannedByte(char byte) bool {
	if self.more() && self.Source[self.cursor] == char {
		self.cursor++
		return true
	}
	return false
}

func (self *Parser) bytesWith(set charset) {
	for self.more() && set.hasByte(self.Source[self.cursor]) {
		self.cursor++
	}
}

func (self *Parser) charsWith(set charset) {
	for i, char := range self.rest() {
		if !set.hasRune(char) {
			self.cursor += i
			return
		}
	}
	self.end()
}

func (self *Parser) charsWithout(set charset) {
	for i, char := range self.rest() {
		if set.hasRune(char) {
			self.cursor += i
			return
		}
	}
	self.end()
}

func (self *Parser) reqMore(char rune, size int, delim rune) {
	if size == 0 {
		self.failEof(delim)
	}

	// Manually inline `self.failNewline` to avoid weird perf regression (WTF).
	if charsetNewline.hasRune(char) {
		self.fail(e.Errorf(`expected closing %q, found newline`, delim))
	}
}

func (self *Parser) failNewline(char rune, delim rune) {
	if charsetNewline.hasRune(char) {
		self.fail(e.Errorf(`expected closing %q, found newline`, delim))
	}
}

func (self *Parser) failEof(delim rune) {
	self.fail(e.Wrapf(io.EOF, `expected closing %q, found EOF`, delim))
}

func (self *Parser) fail(err error) {
	panic(ParseErr{
		Source:  self.Source,
		Cursor:  self.cursor,
		Snippet: snippet(self.rest(), SHORT_SNIPPET_LEN),
		Cause:   err,
	})
}
