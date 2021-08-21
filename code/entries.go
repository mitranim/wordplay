package main

import "unsafe"

type Entry struct {
	Author string
	Phrase string

	Meaning     string
	Meanings    []string
	MeaningsArr [2]string

	Tag     string
	Tags    []string
	TagsArr [2]string

	// Unused.
	Row int
	Col int
}

func (self Entry) Append(buf []byte) []byte {
	buf = self.AppendContent(buf)
	buf = self.AppendAuthor(buf)
	buf = appendNewlines(buf)
	return buf
}

func (self Entry) AppendOld(buf []byte) []byte {
	buf = self.AppendContent(buf)
	buf = appendNewlines(buf)
	return buf
}

func (self Entry) AppendContent(buf []byte) []byte {
	if len(self.Phrase) > 0 {
		buf = self.AppendPhrase(buf)
	}

	if self.HasMeanings() {
		buf = append(buf, " "...)
		buf = self.AppendMeanings(buf)
	}

	if self.HasTags() {
		buf = append(buf, " "...)
		buf = self.AppendTags(buf)
	}

	return buf
}

func (self Entry) AppendPhrase(buf []byte) []byte {
	if strHas(self.Phrase, isWhitespace) {
		buf = append(buf, `"`...)
		buf = append(buf, self.Phrase...)
		buf = append(buf, `"`...)
	} else {
		buf = append(buf, self.Phrase...)
	}
	return buf
}

func (self Entry) AppendMeanings(buf []byte) []byte {
	if self.HasMeanings() {
		buf = append(buf, "("...)
		buf = appendJoinedWith(buf, "; ", self.Meaning, self.Meanings)
		buf = append(buf, ")"...)
	}
	return buf
}

func (self Entry) AppendTags(buf []byte) []byte {
	if self.HasTags() {
		buf = append(buf, "["...)
		buf = appendJoinedWith(buf, "; ", self.Tag, self.Tags)
		buf = append(buf, "]"...)
	}
	return buf
}

func (self Entry) AppendAuthor(buf []byte) []byte {
	if len(self.Author) == 0 {
		return buf
	}

	buf = append(buf, " © "...)
	buf = append(buf, self.Author...)
	return buf
}

func (self Entry) AppendAuthorOld(buf []byte) []byte {
	if len(self.Author) == 0 {
		return buf
	}

	buf = append(buf, "# "...)
	buf = append(buf, self.Author...)
	buf = appendNewlines(buf)
	return buf
}

func (self Entry) GetAuthor() string { return self.Author }

func (self Entry) HasMeanings() bool {
	return len(self.Meaning) > 0 || len(self.Meanings) > 0
}

func (self Entry) HasTags() bool {
	return len(self.Tag) > 0 || len(self.Tags) > 0
}

func (self *Entry) Clone() Entry {
	out := *self

	if unsafe.Pointer(&self.MeaningsArr) == *(*unsafe.Pointer)(unsafe.Pointer(&self.Meanings)) {
		out.Meanings = out.MeaningsArr[:]
	}

	if unsafe.Pointer(&self.TagsArr) == *(*unsafe.Pointer)(unsafe.Pointer(&self.Tags)) {
		out.Tags = out.TagsArr[:]
	}

	return out
}

func (self *Entry) Norm() {
	if self.MeaningsArr != [2]string{} {
		self.Meanings = self.MeaningsArr[:]
	}
	if self.TagsArr != [2]string{} {
		self.Tags = self.TagsArr[:]
	}
}

func (self *Entry) appendMeaning(val string) {
	strAppendSwap(&self.Meaning, &self.Meanings, &self.MeaningsArr, val)
}

func (self *Entry) appendTag(val string) {
	strAppendSwap(&self.Tag, &self.Tags, &self.TagsArr, val)
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
	return GroupEntries(entries, Entry.GetAuthor)
}

// Unused.
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
