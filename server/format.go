package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

const (
	SHORT_SNIPPET_LEN = 64
	LONG_SNIPPET_LEN  = 1024
)

type FormatError struct {
	Line    int
	Column  int
	Content []byte
	Snippet []byte
	Reason  error
}

func (self FormatError) Error() string {
	return self.formatError(false)
}

func (self FormatError) formatError(expand bool) string {
	format := ``
	args := []interface{}{}

	if self.Line > 0 || self.Column > 0 {
		format += `%v:%v`
		args = append(args, self.Line, self.Column)
	}

	if len(self.Snippet) > 0 {
		if format != "" {
			format += `: `
		}
		format += `%q`
		args = append(args, self.Snippet)
	}

	if self.Reason != nil {
		if format != "" {
			format += `: `
		}
		if expand {
			format += `%+v`
		} else {
			format += `%v`
		}
		args = append(args, self.Reason)
	}

	return fmt.Sprintf(format, args...)
}

func (self FormatError) Format(fms fmt.State, verb rune) {
	switch verb {
	case 'v':
		if fms.Flag('#') {
			spew.Fdump(fms, self)
			return
		}
		if fms.Flag('+') {
			fms.Write(stringToBytesAlloc(self.formatError(true)))
			return
		}
		fms.Write(stringToBytesAlloc(self.formatError(true)))
	default:
		fms.Write(stringToBytesAlloc(self.formatError(true)))
	}
}

type Entry struct {
	Author   string
	Phrase   string
	Meanings []string
	Tags     []string
	Line     int
	Column   int
}

func (entry Entry) FormatAppendPhrase(buf []byte) []byte {
	if regWhitespaceAnywhere.MatchString(entry.Phrase) {
		buf = append(buf, `"`...)
		buf = append(buf, entry.Phrase...)
		buf = append(buf, `"`...)
	} else {
		buf = append(buf, entry.Phrase...)
	}
	return buf
}

func (entry Entry) FormatAppendMeanings(buf []byte) []byte {
	if len(entry.Meanings) > 0 {
		buf = append(buf, "("...)
		for i, val := range entry.Meanings {
			if i > 0 {
				buf = append(buf, "; "...)
			}
			buf = append(buf, val...)
		}
		buf = append(buf, ")"...)
	}
	return buf
}

func (entry Entry) FormatAppendTags(buf []byte) []byte {
	if len(entry.Tags) > 0 {
		buf = append(buf, "["...)
		for i, val := range entry.Tags {
			if i > 0 {
				buf = append(buf, "; "...)
			}
			buf = append(buf, val...)
		}
		buf = append(buf, "]"...)
	}
	return buf
}

func (entry *Entry) Clean() {
	entry.Author = strings.TrimSpace(entry.Author)
	entry.Phrase = strings.TrimSpace(entry.Phrase)
	stringSliceTrimSpace(entry.Meanings)
	stringSliceTrimSpace(entry.Tags)
}

func ValidateEntryAddition(entries []Entry, newEntry Entry) error {
	if isBlank(newEntry.Phrase) {
		return errors.New(`phrase must be non-empty`)
	}

	if isBlank(newEntry.Author) {
		return errors.New(`author must be non-empty`)
	}

	for _, entry := range entries {
		if entry.Phrase == newEntry.Phrase {
			return errors.Errorf(`redundant entry; phrase %q is already present`, entry.Phrase)
		}
	}

	return nil
}

// A simplistic "ordered map" for lists of entries.
type EntryMap struct {
	Keys []string
	Map  map[string][]Entry
}

func (self EntryMap) Ungroup() []Entry {
	total := 0
	for _, list := range self.Map {
		total += len(list)
	}

	out := make([]Entry, 0, total)
	for _, key := range self.Keys {
		out = append(out, self.Map[key]...)
	}
	return out
}

func GroupEntries(entries []Entry, fun func(Entry) string) EntryMap {
	grouped := EntryMap{Map: map[string][]Entry{}}
	for _, entry := range entries {
		key := fun(entry)

		keyKnown := false
		for _, ordKey := range grouped.Keys {
			if ordKey == key {
				keyKnown = true
				break
			}
		}
		if !keyKnown {
			grouped.Keys = append(grouped.Keys, key)
		}

		grouped.Map[key] = append(grouped.Map[key], entry)
	}
	return grouped
}

func GroupEntriesByAuthor(entries []Entry) EntryMap {
	return GroupEntries(entries, func(entry Entry) string { return entry.Author })
}

// On i7-8750H 2.2—3.9 GHz, this takes 2.4 milliseconds for 30 kilobytes of
// content. Could be better, but unlikely to be a bottleneck.
func ParseEntries(content []byte) ([]Entry, error) {
	state := MakeParseState(content)
	err := state.Parse()
	return state.Entries, err
}

// Assumes entries have been normalized. TODO enforce normalization: sort by
// author first.
func FormatEntries(entries []Entry) []byte {
	var buf []byte
	var author string

	for _, entry := range entries {
		if entry.Author != author {
			buf = append(buf, "#"...)
			if entry.Author != "" {
				buf = append(buf, " "...)
				buf = append(buf, entry.Author...)
			}
			buf = append(buf, "\n\n"...)
			author = entry.Author
		}

		buf = entry.FormatAppendPhrase(buf)

		if len(entry.Meanings) > 0 {
			buf = append(buf, " "...)
			buf = entry.FormatAppendMeanings(buf)
		}

		buf = append(buf, "\n\n"...)
	}

	return buf
}

func FormatEntriesNew(entries []Entry) []byte {
	var buf []byte

	for _, entry := range entries {
		buf = entry.FormatAppendPhrase(buf)

		if len(entry.Meanings) > 0 {
			buf = append(buf, " "...)
			buf = entry.FormatAppendMeanings(buf)
		}

		if len(entry.Tags) > 0 {
			buf = append(buf, " "...)
			buf = entry.FormatAppendTags(buf)
		}

		if entry.Author != "" {
			buf = append(buf, " © "...)
			buf = append(buf, entry.Author...)
		}

		buf = append(buf, "\n\n"...)
	}

	return buf
}

type ParseState struct {
	Entries     []Entry
	Content     []byte
	Line        int
	Column      int
	EntryLine   int
	EntryColumn int
	Author      string
	Phrase      string
	Meanings    []string
}

func MakeParseState(content []byte) ParseState {
	return ParseState{
		Content: content,
		Line:    1,
		Column:  0,
	}
}

var regWhitespaceAnywhere = regexp.MustCompile(`\s`)
var regWhitespace = regexp.MustCompile(`^\s*`)
var regSameLineWhitespace = regexp.MustCompile(`^[^\S\r\n]*`)
var regHead = regexp.MustCompile(`^# ([^\r\n]*)`)
var regEntryQuoted = regexp.MustCompile(`^"([^"\r\n]*)"`)
var regEntryUnquoted = regexp.MustCompile(`^(\S+)`)
var regEntryMeaning = regexp.MustCompile(`^[^\S\r\n]*\(([^()\r\n]*)\)`)
var regNonEmptyLine = regexp.MustCompile(`^[^\S\r\n]*\S`)
var semicolonBytes = []byte{';'}

func (state *ParseState) Parse() error {
	for len(state.Content) > 0 {
		cursor := regWhitespace.FindIndex(state.Content)[1]

		if cursor > 0 {
			state.AdvanceCursor(cursor)
			continue
		}

		if state.Content[0] == '#' {
			err := state.ParseHead()
			if err != nil {
				return err
			}
			continue
		}

		if state.Content[0] == '"' {
			err := state.ParseEntryQuoted()
			if err != nil {
				return err
			}
			continue
		}

		err := state.ParseEntryUnquoted()
		if err != nil {
			return err
		}
	}

	return nil
}

func (state *ParseState) ParseDelim() error {
	cursor := regWhitespace.FindIndex(state.Content)[1]

	if len(state.Content) > cursor {
		lines, _ := lineColumnDelta(state.Content[:cursor])
		if lines < 2 {
			return state.Error(errors.New(`expected at least one empty line`))
		}
	}

	state.AdvanceCursor(cursor)
	return nil
}

func (state *ParseState) ParseHead() error {
	indexes := regHead.FindSubmatchIndex(state.Content)
	if len(indexes) == 0 {
		return state.Error(errors.New(`malformed header: expected '#' followed by space and author name`))
	}

	state.Author = bytesToStringAlloc(bytes.TrimSpace(state.Content[indexes[2]:indexes[3]]))
	state.AdvanceCursor(indexes[1])
	state.ParseDelim()
	return nil
}

func (state *ParseState) ParseEntryQuoted() error {
	indexes := regEntryQuoted.FindSubmatchIndex(state.Content)
	if len(indexes) == 0 {
		return state.Error(errors.New(`malformed entry: expected a phrase in quotes`))
	}

	phrase := bytesToStringAlloc(bytes.TrimSpace(state.Content[indexes[2]:indexes[3]]))
	if len(phrase) == 0 {
		// TODO determine correct column.
		return state.Error(errors.New(`phrase appears empty`))
	}

	state.EntryLine = state.Line
	state.EntryColumn = state.Column
	state.Phrase = phrase
	state.AdvanceCursor(indexes[1])

	state.ParseEntryMeanings()
	return nil
}

func (state *ParseState) ParseEntryUnquoted() error {
	indexes := regEntryUnquoted.FindSubmatchIndex(state.Content)
	if len(indexes) == 0 {
		return state.Error(errors.New(`malformed entry: expected a word`))
	}

	state.EntryLine = state.Line
	state.EntryColumn = state.Column
	state.Phrase = bytesToStringAlloc(state.Content[indexes[2]:indexes[3]])
	state.AdvanceCursor(indexes[1])

	state.ParseEntryMeanings()
	return nil
}

// FIXME parse tags.
func (state *ParseState) ParseEntryMeanings() error {
	if regNonEmptyLine.Match(state.Content) {
		indexes := regEntryMeaning.FindSubmatchIndex(state.Content)
		if len(indexes) == 0 {
			return state.Error(errors.New(`expected meanings in parens (nested parens are not allowed)`))
		}

		inner := bytes.TrimSpace(state.Content[indexes[2]:indexes[3]])
		if len(inner) == 0 {
			// TODO determine correct column.
			return state.Error(errors.New(`expected non-empty meanings between parens`))
		}

		var meanings []string
		for _, chunk := range bytes.Split(inner, semicolonBytes) {
			meaning := bytesToStringAlloc(bytes.TrimSpace(chunk))
			if len(meaning) == 0 {
				// TODO determine correct column.
				return state.Error(errors.New(`one of the meanings is empty`))
			}
			meanings = append(meanings, meaning)
		}

		state.Meanings = meanings
		state.AdvanceCursor(indexes[1])
	}

	state.Entries = append(state.Entries, Entry{
		Author:   state.Author,
		Line:     state.EntryLine,
		Column:   state.EntryColumn,
		Phrase:   state.Phrase,
		Meanings: state.Meanings,
	})

	state.EntryLine = 0
	state.EntryColumn = 0
	state.Phrase = ""
	state.Meanings = nil

	state.ParseDelim()
	return nil
}

func (state *ParseState) AdvanceCursor(end int) {
	head := state.Content[:end]
	tail := state.Content[end:]
	lines, columns := lineColumnDelta(head)

	state.Content = tail

	if lines == 0 {
		state.Column += columns
	} else {
		state.Line += lines
		state.Column = columns + 1
	}
}

func (state *ParseState) consumeSameLineWhitespace() {
	cursor := regSameLineWhitespace.FindIndex(state.Content)[1]
	if cursor > 0 {
		state.AdvanceCursor(cursor)
	}
}

func (state ParseState) Error(reason error) FormatError {
	return FormatError{
		Line:    state.Line,
		Column:  state.Column,
		Content: state.Content,
		Snippet: snippet(state.Content, SHORT_SNIPPET_LEN),
		Reason:  reason,
	}
}

func IntersperseEntriesByAuthor(entries []Entry) []Entry {
	grouped := GroupEntriesByAuthor(entries)
	ungrouped := make([]Entry, 0, len(entries))
	_ = grouped

	// for len(ungrouped) < len(entries) {
	// 	for _, key := range grouped.Keys {

	// 	}
	// }

	return ungrouped
}

func snippet(input []byte, limit int) []byte {
	for i, char := range input {
		switch char {
		case '\n', '\r':
			return input[:i]
		}
		if i > limit {
			out := make([]byte, 0, i+3)
			out = append(out, input[:i]...)
			out = append(out, "…"...)
			return out
		}
	}
	return input
}

func lineColumnDelta(input []byte) (int, int) {
	var lines int
	var columns int

	for i, char := range input {
		if char == '\r' && len(input) > i+1 && input[i+1] == '\n' {
			continue
		}

		if char == '\r' || char == '\n' {
			lines++
			columns = 0
			continue
		}

		columns++
	}

	return lines, columns
}

func isBlank(str string) bool {
	for _, char := range str {
		switch char {
		case ' ', '\t', '\v', '\r', '\n':
		default:
			return false
		}
	}
	return true
}

func bytesToStringAlloc(bytes []byte) string   { return string(bytes) }
func stringToBytesAlloc(input string) []byte   { return []byte(input) }
func bytesToMutableString(input []byte) string { return *(*string)(unsafe.Pointer(&input)) }

func stringSliceTrimSpace(slice []string) {
	for i := range slice {
		slice[i] = strings.TrimSpace(slice[i])
	}
}
