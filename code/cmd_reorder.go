package main

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mitranim/try"
)

func init() { commands[`reorder`] = cmdReorder }

func cmdReorder() {
	var times TimeMap
	json.Unmarshal(try.ByteSlice(os.ReadFile(TIMES_FILE)), &times)

	var timed TimedEntries
	for _, val := range readEntries(SRC_FILE) {
		timed = append(timed, TimedEntry{val, times[val.Phrase]})
	}
	writeFileStr(`../fixtures/entries_timed_unsorted`, spew.Sdump(timed))

	sort.Stable(sort.Reverse(timed))
	writeFileStr(`../fixtures/entries_timed_sorted`, spew.Sdump(timed))

	entries := untimedEntries(timed)
	writeFile(`../fixtures/readme.md`, bytes.TrimSpace(FormatEntries(entries)))
}

type TimedEntries []TimedEntry

func (self TimedEntries) Len() int { return len(self) }

func (self TimedEntries) Swap(a, b int) { self[a], self[b] = self[b], self[a] }

func (self TimedEntries) Less(a, b int) bool {
	if self[a].Time.Before(self[b].Time) {
		return true
	}
	if self[a].Time.After(self[b].Time) {
		return false
	}
	if self[a].Author < self[b].Author {
		return true
	}
	if self[a].Author > self[b].Author {
		return false
	}
	return false
}

type TimedEntry struct {
	Entry
	Time time.Time
}

func untimedEntries(vals TimedEntries) Entries {
	out := make(Entries, 0, len(vals))
	for _, val := range vals {
		out = append(out, val.Entry)
	}
	return out
}
