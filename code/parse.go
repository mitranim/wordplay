package main

import (
	"strings"
	"unicode/utf8"

	"github.com/mitranim/try"
	e "github.com/pkg/errors"
)

func ParseEntries(src string) Entries {
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
	defer try.Detail(`failed to parse entries`)
	defer try.Trans(self.err)

	for self.more() {
		if !self.scanned((*Parser).any) {
			panic(e.New(`unrecognized content`))
		}
	}
}

func (self *Parser) any() {
	switch {
	case self.scanned((*Parser).whitespace):
	case self.scanned((*Parser).entryQuoted):
	case self.scanned((*Parser).entryUnquoted):
	}
}

func (self *Parser) entryQuoted() {
	if !self.scannedByte('"') {
		return
	}

	start := self.cursor
	self.singleLineUntil('"')

	phrase := strings.TrimSpace(self.from(start))
	self.cursor += len(`"`)

	if len(phrase) == 0 {
		panic(e.New(`quoted phrase is empty`))
	}

	self.entry.Phrase = phrase
	self.entryRest()
}

func (self *Parser) nonDelim() { self.charsWithout(charsetDelim) }

func (self *Parser) entryUnquoted() {
	start := self.cursor
	if !self.scanned((*Parser).nonDelim) {
		return
	}

	phrase := strings.TrimSpace(self.from(start))
	if len(phrase) == 0 {
		panic(e.New(`phrase is empty`))
	}

	self.entry.Phrase = phrase
	self.entryRest()
}

func (self *Parser) entryRest() {
	self.entryMeanings()
	self.entryTags()
	self.entryAuthor()
	self.entryFlush()
	self.delimWhitespace()
}

func (self *Parser) entryMeanings() {
	self.space()
	if !self.scannedByte('(') {
		return
	}

	start := self.cursor
	cursor := self.cursor

	for i, char := range self.rest() {
		if charsetNewline.hasRune(char) {
			self.cursor = start + i + utf8.RuneLen(char)
			panic(errNewline(')'))
		}

		if char == ';' {
			self.cursor = start + i
			self.appendMeaning(self.from(cursor))
			self.cursor += len(`;`)
			cursor = self.cursor
			continue
		}

		if char == ')' {
			self.cursor = start + i
			self.appendMeaning(self.from(cursor))
			self.cursor += len(`)`)
			return
		}

		if char == '(' {
			self.cursor = start + i + len(`(`)
			panic(e.New(`unexpected nested "("`))
		}
	}

	self.end()
	panic(errEof(')'))
}

func (self *Parser) appendMeaning(val string) {
	// defer self.detail()
	self.entry.appendMeaning(val)
}

func (self *Parser) entryTags() {
	self.space()
	if !self.scannedByte('[') {
		return
	}

	start := self.cursor
	cursor := self.cursor

	for i, char := range self.rest() {
		if charsetNewline.hasRune(char) {
			self.cursor = start + i + utf8.RuneLen(char)
			panic(errNewline(']'))
		}

		if char == ';' {
			self.cursor = start + i
			self.appendTag(self.from(cursor))
			self.cursor += len(`;`)
			cursor = self.cursor
			continue
		}

		if char == ']' {
			self.cursor = start + i
			self.appendTag(self.from(cursor))
			self.cursor += len(`]`)
			return
		}

		if char == '[' {
			self.cursor = start + i + len(`[`)
			panic(e.New(`unexpected nested "["`))
		}
	}

	self.end()
	panic(errEof(']'))
}

func (self *Parser) appendTag(val string) {
	// defer self.detail()
	self.entry.appendTag(val)
}

func (self *Parser) entryAuthor() {
	self.space()
	if !self.scannedChar('©') {
		return
	}

	start := self.cursor
	self.nonNewline()

	author := strings.TrimSpace(self.from(start))
	if len(author) == 0 {
		panic(e.New(`expected "©" to be followed by author name`))
	}

	self.entry.Author = author
}

func (self *Parser) entryFlush() {
	self.Entries = append(self.Entries, self.entry)
	self.entry = Entry{}
}

func (self *Parser) delimWhitespace() {
	self.space()

	// nolint:staticcheck
	if self.more() && !self.scannedNewline() || self.more() && !self.scannedNewline() {
		panic(e.New(`expected at least two newlines or EOF`))
	}

	self.whitespace()
}

func (self *Parser) whitespace() { self.bytesWith(charsetWhitespace) }
func (self *Parser) space()      { self.bytesWith(charsetSpace) }
func (self *Parser) newline()    { self.cursor += leadingNewlineSize(self.rest()) }
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

func (self *Parser) end() { self.cursor = len(self.Source) }

func (self *Parser) headByte() byte {
	return self.Source[self.cursor]
}

func (self *Parser) scanned(fun func(*Parser)) bool {
	start := self.cursor
	fun(self)
	return self.cursor > start
}

func (self *Parser) scannedNewline() bool {
	return self.scanned((*Parser).newline)
}

func (self *Parser) scannedByte(char byte) bool {
	if self.more() && self.headByte() == char {
		self.cursor++
		return true
	}
	return false
}

func (self *Parser) scannedChar(val rune) bool {
	char, size := headChar(self.rest())
	if size > 0 && val == char {
		self.cursor += size
		return true
	}
	return false
}

func (self *Parser) bytesWith(set *charset) {
	for self.more() && set.hasByte(self.headByte()) {
		self.cursor++
	}
}

func (self *Parser) charsWithout(set *charset) {
	for i, char := range self.rest() {
		if set.hasRune(char) {
			self.cursor += i
			return
		}
	}
	self.end()
}

func (self *Parser) singleLineUntil(delim rune) {
	for i, char := range self.rest() {
		if charsetNewline.hasRune(char) {
			self.cursor += i
			panic(errNewline(delim))
		}

		if char == delim {
			self.cursor += i
			return
		}
	}

	self.end()
	panic(errEof(delim))
}

func (self *Parser) err(err error) error {
	return ParseErr{
		Source:  self.Source,
		Cursor:  self.cursor,
		Snippet: snippet(self.rest(), SHORT_SNIPPET_LEN),
		Cause:   err,
	}
}
