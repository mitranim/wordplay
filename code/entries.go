package main

type Entries []Entry

func (self Entries) Bytes() (buf []byte) {
	for _, val := range self {
		buf = val.Append(buf)
	}
	return
}

func (self Entries) String() string {
	return bytesToMutableString(self.Bytes())
}

type Entry struct {
	Author   string
	Phrase   string
	Meanings []string
	Tags     []string
}

func (self *Entry) Append(buf []byte) []byte {
	buf = self.AppendPhrase(buf)
	buf = self.AppendMeanings(buf)
	buf = self.AppendTags(buf)
	buf = self.AppendAuthor(buf)
	buf = appendNewlines(buf)
	return buf
}

func (self *Entry) AppendPhrase(buf []byte) []byte {
	if strHas(self.Phrase, charsetWhitespace) {
		buf = append(buf, `"`...)
		buf = append(buf, self.Phrase...)
		buf = append(buf, `"`...)
	} else {
		buf = append(buf, self.Phrase...)
	}
	return buf
}

func (self *Entry) AppendMeanings(buf []byte) []byte {
	if self.HasMeanings() {
		buf = append(buf, " ("...)
		buf = appendJoined(buf, "; ", self.Meanings)
		buf = append(buf, ")"...)
	}
	return buf
}

func (self *Entry) AppendTags(buf []byte) []byte {
	if self.HasTags() {
		buf = append(buf, " ["...)
		buf = appendJoined(buf, "; ", self.Tags)
		buf = append(buf, "]"...)
	}
	return buf
}

func (self *Entry) AppendAuthor(buf []byte) []byte {
	if len(self.Author) > 0 {
		buf = append(buf, " © "...)
		buf = append(buf, self.Author...)
	}
	return buf
}

func (self *Entry) HasMeanings() bool { return len(self.Meanings) > 0 }
func (self *Entry) HasTags() bool     { return len(self.Tags) > 0 }

func (self *Entry) appendMeaning(val string) {
	self.Meanings = append(self.Meanings, val)
}

func (self *Entry) appendTag(val string) {
	self.Tags = append(self.Tags, val)
}
