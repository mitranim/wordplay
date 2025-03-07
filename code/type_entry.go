package main

import (
	"strings"

	"github.com/mitranim/gg"
)

type Entries []Entry

// Implement `fmt.Stringer`.
func (self Entries) String() string { return gg.ToString(self.AppendTo(nil)) }

// Implement `gg.AppenderTo`.
func (self Entries) AppendTo(buf []byte) []byte {
	for _, val := range self {
		buf = val.AppendTo(buf)
	}
	return buf
}

func (self Entries) Dupes() (out []string) {
	counts := make(map[string]int)

	// TODO consider case-insensitive matching.
	return gg.MapCompact(self, func(val Entry) (_ string) {
		key := val.Pk()
		count := counts[key]
		counts[key]++
		if count == 1 {
			return key
		}
		return
	})
}

func (self Entries) ReplaceAuthors(src map[string]string) {
	for ind := range self {
		self[ind].ReplaceAuthor(src)
	}
}

func (self Entries) Norm() { gg.EachPtr(self, (*Entry).Norm) }

type Entry struct {
	Author   string
	Phrase   string
	Meanings []string
	Tags     []string // Preferably ISO 639-1 or ISO 639-2 codes.
}

// Implement `gg.Pker`.
func (self Entry) Pk() string { return self.Phrase }

// Implement `fmt.Stringer`.
func (self Entry) String() string { return gg.ToString(self.AppendTo(nil)) }

// Implement `gg.AppenderTo`.
func (self Entry) AppendTo(src []byte) []byte {
	buf := gg.Buf(src)
	if buf.Len() > 0 {
		buf.AppendNewline()
	}
	buf = self.AppendPhrase(buf)
	buf = self.AppendMeanings(buf)
	buf = self.AppendTags(buf)
	buf = self.AppendAuthor(buf)
	buf.AppendNewline()
	return buf
}

func (self Entry) AppendPhrase(src []byte) []byte {
	return append(src, self.Phrase...)
}

func (self Entry) AppendMeanings(src []byte) []byte {
	buf := gg.Buf(src)
	if self.HasMeanings() {
		buf.AppendString(` (`)
		buf = AppendJoined(buf, `; `, self.Meanings)
		buf.AppendString(`)`)
	}
	return buf
}

func (self Entry) AppendTags(src []byte) []byte {
	buf := gg.Buf(src)
	if self.HasTags() {
		buf.AppendString(` [`)
		buf = AppendJoined(buf, `; `, self.Tags)
		buf.AppendString(`]`)
	}
	return buf
}

func (self Entry) AppendAuthor(src []byte) []byte {
	buf := gg.Buf(src)
	if len(self.Author) > 0 {
		buf.AppendString(` © `)
		buf.AppendString(self.Author)
	}
	return buf
}

func (self Entry) HasMeanings() bool { return len(self.Meanings) > 0 }
func (self Entry) HasTags() bool     { return len(self.Tags) > 0 }

func (self Entry) HasRedundantAuthor() bool {
	return hasAuthorSign(self.Author) ||
		hasAuthorSign(self.Phrase) ||
		gg.Some(self.Meanings, hasAuthorSign) ||
		gg.Some(self.Tags, hasAuthorSign)
}

func (self *Entry) AddMeaning(val string) {
	val = strings.TrimSpace(val)
	if len(val) == 0 {
		panic(gg.Errv(`unexpected empty meaning`))
	}
	self.Meanings = append(self.Meanings, val)
}

func (self *Entry) AddTag(val string) {
	val = strings.TrimSpace(val)
	if len(val) == 0 {
		panic(gg.Errv(`unexpected empty tag`))
	}
	self.Tags = append(self.Tags, val)
}

func (self *Entry) ReplaceAuthor(src map[string]string) {
	val, ok := src[self.Author]
	if ok {
		self.Author = val
	}
}

func (self *Entry) Norm() {
	self.Author = StrNorm(self.Author)
	self.Phrase = StrNorm(self.Phrase)
	gg.MapMut(self.Meanings, StrNorm)
	gg.MapMut(self.Tags, StrNorm)
}

func hasAuthorSign(src string) bool { return strings.Contains(src, `©`) }
