package main

import (
	"strings"

	"github.com/mitranim/try"
	e "github.com/pkg/errors"
)

const (
	SHORT_SNIPPET_LEN = 64
	LONG_SNIPPET_LEN  = 1024
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

	cursor int // must be == (row + 1) * (col + 1)
	row    int
	col    int
	entry  Entry
}

func MakeParser(content string) Parser {
	// return Parser{Source: content, Entries: make(Entries, 0, 2048)}
	return Parser{Source: content}
}

func (self *Parser) Parse() {
	for self.more() {
		if !self.scanned((*Parser).any) {
			panic(self.err(e.New(`unrecognized content`)))
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

func (self *Parser) whitespace() { self.chars(isWhitespace) }
func (self *Parser) space()      { self.chars(isSpace) }

func (self *Parser) heading() {
	if !self.scannedChar('#') {
		return
	}

	if !self.scannedBy(isSpace) {
		panic(self.err(e.New(`malformed header: expected '#' followed by space and author name`)))
	}

	start := self.cursor
	self.nonNewline()

	author := strings.TrimSpace(self.from(start))
	if len(author) == 0 {
		panic(self.err(e.New(`malformed header: expected '#' followed by space and author name`)))
	}

	self.entry.Author = author
	self.delimWhitespace()
}

func (self *Parser) entryQuoted() {
	row, col := self.row, self.col
	if !self.scannedChar('"') {
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
				panic(self.err(e.New(`quoted phrase is empty`)))
			}

			self.entry.Phrase = phrase
			self.entry.Row = row
			self.entry.Col = col
			self.scan(char, size)
			break loop

		default:
			self.scan(char, size)
		}
	}

	self.entryRest()
}

func (self *Parser) nonDelim() { self.chars(isNonDelim) }

func (self *Parser) entryUnquoted() {
	start, row, col := self.cursor, self.row, self.col

	if !self.scanned((*Parser).nonDelim) {
		return
	}

	self.entry.Phrase = self.from(start)
	self.entry.Row = row
	self.entry.Col = col

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
	if !self.scannedChar('(') {
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
			self.scan(char, size)

		case ')':
			depth--

			if depth == 0 {
				self.appendMeaning(self.from(start))
				self.scan(char, size)
				break loop
			}

			if depth < 0 {
				panic(self.err(e.New(`mismatched closing ")"`)))
			}

			self.scan(char, size)

		case ';':
			self.appendMeaning(self.from(start))
			self.scan(char, size)
			start = self.cursor

		default:
			self.scan(char, size)
		}
	}
}

func (self *Parser) appendMeaning(val string) {
	val = strings.TrimSpace(val)
	if len(val) == 0 {
		panic(self.err(e.New(`unexpected empty meaning`)))
	}
	self.entry.appendMeaning(val)
}

// Untested!
func (self *Parser) entryTags() {
	self.space()
	if !self.scannedChar('[') {
		return
	}

	start := self.cursor

loop:
	for {
		char, size := self.headChar()
		self.reqMore(char, size, ']')

		switch char {
		case '[':
			panic(self.err(e.New(`unexpected nested "["`)))

		case ';':
			self.appendTag(self.from(start))
			self.scan(char, size)
			start = self.cursor

		case ']':
			self.appendTag(self.from(start))
			self.scan(char, size)
			break loop

		default:
			self.scan(char, size)
		}
	}
}

func (self *Parser) appendTag(val string) {
	val = strings.TrimSpace(val)
	if len(val) == 0 {
		panic(self.err(e.New(`unexpected empty tag`)))
	}
	self.entry.appendTag(val)
}

func (self *Parser) entryFlush() {
	self.Entries = self.Entries.Append(self.entry)
	self.entry = Entry{Author: self.entry.Author}
}

func (self *Parser) delimWhitespace() {
	self.space()

	if !self.more() {
		return
	}
	if !self.scannedNewline() {
		panic(self.err(e.New(`expected at least two newlines or EOF`)))
	}

	if !self.more() {
		return
	}
	if !self.scannedNewline() {
		panic(self.err(e.New(`expected at least two newlines or EOF`)))
	}

	self.whitespace()
}

func (self *Parser) more() bool { return self.left() > 0 }

func (self *Parser) rest() string {
	if self.more() {
		return self.Source[self.cursor:]
	}
	return ""
}

func (self *Parser) left() int {
	return len(self.Source) - self.cursor
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

func (self *Parser) scanned(fun func(*Parser)) bool {
	start := self.cursor
	fun(self)
	return self.cursor > start
}

func (self *Parser) scannedNewline() bool {
	return self.scanned((*Parser).newline)
}

func (self *Parser) scannedChar(char rune) bool {
	head, size := self.headChar()
	if size > 0 && head == char {
		self.scan(head, size)
		return true
	}
	return false
}

func (self *Parser) scannedBy(fun func(rune) bool) bool {
	start := self.cursor
	self.chars(fun)
	return self.cursor > start
}

func (self *Parser) newline() {
	char, size := self.headChar()
	if char == '\n' {
		self.scan(char, size)
	}
}

func (self *Parser) nonNewline() { self.chars(isNonNewline) }

func (self *Parser) scan(char rune, size int) {
	if size > 0 {
		self.cursor += size

		// Newlines must be normalized via `headChar`.
		if char == '\n' {
			self.row++
			self.col = 0
		} else {
			self.row = 0
			self.col++
		}
	}
}

func (self *Parser) scanChar() {
	self.scan(self.headChar())
}

func (self *Parser) scanChars(count int) {
	for range counter(count) {
		self.scanChar()
	}
}

func (self *Parser) chars(fun func(rune) bool) {
	for {
		char, size := self.headChar()
		if size > 0 && fun(char) {
			self.scan(char, size)
		} else {
			return
		}
	}
}

func (self *Parser) isNextChar(char rune) bool {
	head, size := self.headChar()
	return size > 0 && head == char
}

func (self *Parser) isNextPrefix(prefix string) bool {
	return strings.HasPrefix(self.rest(), prefix)
}

func (self *Parser) reqMore(char rune, size int, delim rune) {
	if size == 0 {
		panic(self.err(e.Errorf(`expected closing %q, found EOF`, delim)))
	}

	if isNewline(char) {
		panic(self.err(e.Errorf(`expected closing %q, found newline`, delim)))
	}
}

func (self Parser) err(err error) ParseErr {
	return ParseErr{
		Row:     self.row,
		Col:     self.col,
		Source:  self.Source,
		Snippet: snippet(self.rest(), SHORT_SNIPPET_LEN),
		Cause:   err,
	}
}
