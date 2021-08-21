package main

type Entries []Entry

func (self Entries) Append(val Entry) Entries {
	return append(self, val)
}

func (self Entries) AppendMany(vals ...Entry) Entries {
	return append(self, vals...)
}

type Entry struct {
	Author   string
	Phrase   string
	Meanings []string
	Tags     []string

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
		buf = appendJoined(buf, "; ", self.Meanings)
		buf = append(buf, ")"...)
	}
	return buf
}

func (self Entry) AppendTags(buf []byte) []byte {
	if self.HasTags() {
		buf = append(buf, "["...)
		buf = appendJoined(buf, "; ", self.Tags)
		buf = append(buf, "]"...)
	}
	return buf
}

func (self Entry) AppendAuthor(buf []byte) []byte {
	if len(self.Author) > 0 {
		buf = append(buf, " © "...)
		buf = append(buf, self.Author...)
	}
	return buf
}

func (self Entry) AppendAuthorOld(buf []byte) []byte {
	if len(self.Author) > 0 {
		buf = append(buf, "# "...)
		buf = append(buf, self.Author...)
		buf = appendNewlines(buf)
	}
	return buf
}

func (self Entry) GetAuthor() string { return self.Author }
func (self Entry) HasMeanings() bool { return len(self.Meanings) > 0 }
func (self Entry) HasTags() bool     { return len(self.Tags) > 0 }

func (self *Entry) appendMeaning(val string) {
	self.Meanings = append(self.Meanings, val)
}

func (self *Entry) appendTag(val string) {
	self.Tags = append(self.Tags, val)
}

// A simplistic "ordered map" for lists of entries.
type EntryMap struct {
	Keys []string
	Map  map[string]Entries
}

func (self EntryMap) Ungroup() Entries {
	total := 0
	for _, list := range self.Map {
		total += len(list)
	}

	out := make(Entries, 0, total)
	for _, key := range self.Keys {
		out = out.AppendMany(self.Map[key]...)
	}
	return out
}

func GroupEntries(entries Entries, fun func(Entry) string) EntryMap {
	grouped := EntryMap{Map: map[string]Entries{}}
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

		grouped.Map[key] = grouped.Map[key].Append(entry)
	}
	return grouped
}

func GroupEntriesByAuthor(entries Entries) EntryMap {
	return GroupEntries(entries, Entry.GetAuthor)
}

// Unused.
func IntersperseEntriesByAuthor(entries Entries) Entries {
	grouped := GroupEntriesByAuthor(entries)
	ungrouped := make(Entries, 0, len(entries))
	_ = grouped

	// for len(ungrouped) < len(entries) {
	// 	for _, key := range grouped.Keys {

	// 	}
	// }

	return ungrouped
}
