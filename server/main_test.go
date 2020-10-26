package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/mitranim/repr"
)

const SRC = "../readme.md"
const TEMPDIR_NAME = "wordplay_testing"
const FILE_READ_WRITE_MODE = 0600
const DIR_READ_WRITE_MODE = 0700

func TestParseAndFormat(t *testing.T) {
	content, err := ioutil.ReadFile(SRC)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	entries, err := ParseEntries(content)
	if err != nil {
		t.Fatalf("%v:%v", SRC, err)
	}

	printed := FormatEntries(entries)
	source := bytes.TrimSpace(content)
	output := bytes.TrimSpace(printed)

	if testing.Verbose() {
		fmt.Printf("source: %v\n", writeTempFile(t, "source", source))
		fmt.Printf("output: %v\n", writeTempFile(t, "output", output))
		fmt.Printf("output new: %v\n", writeTempFile(t, "output_new", bytes.TrimSpace(FormatEntriesNew(entries))))
		fmt.Printf("entries: %v\n", writeTempFile(t, "entries", repr.Bytes(entries)))
	}

	if bytesToStringAlloc(source) != bytesToStringAlloc(output) {
		if testing.Verbose() {
			t.Fatalf("mismatch of source and formatted output")
		} else {
			t.Fatal("mismatch of source and formatted output; run test in verbose mode for details")
		}
	}
}

func TestGroup(t *testing.T) {
	entries := tbEntries(t)
	grouped := GroupEntriesByAuthor(entries)

	ordKeys := make([]string, len(grouped.Keys))
	copy(ordKeys, grouped.Keys)
	sort.Strings(ordKeys)

	mapKeys := make([]string, 0, len(grouped.Map))
	for key := range grouped.Map {
		mapKeys = append(mapKeys, key)
	}
	sort.Strings(mapKeys)

	if !reflect.DeepEqual(ordKeys, mapKeys) {
		t.Fatalf("mismatch between ordered keys and actual map keys: %#v vs %#v",
			ordKeys, mapKeys)
	}

	for key, group := range grouped.Map {
		for _, entry := range group {
			if entry.Author != key {
				t.Fatalf(`author mismatch in grouped entry: author = %v, entry.Author = %v`,
					key, entry.Author)
			}
		}
	}

	count := 0
	for _, group := range grouped.Map {
		count += len(group)
	}
	if count != len(entries) {
		t.Fatalf(`total amount of grouped entries doesn't match amount of source entries; expected %v, counted %v`,
			len(entries), count)
	}
}

func TestUngroup(t *testing.T) {
	entries := tbEntries(t)
	grouped := GroupEntriesByAuthor(entries)
	ungrouped := grouped.Ungroup()

	if !reflect.DeepEqual(entries, ungrouped) {
		if testing.Verbose() && false {
			t.Fatalf("mismatch between parsed and grouped-ungrouped entries\nparsed:\n%v\nungrouped:\n%v",
				repr.String(entries), repr.String(ungrouped))
		} else {
			t.Log("grouped.Keys:", repr.String(grouped.Keys))
			t.Fatal("mismatch between parsed and grouped-ungrouped entries")
		}
	}
}

func TestAppendEntry(t *testing.T) {
	// Sorting would change the order.
	const firstAuthor = "3 1623f82d4c287874bec41230f8c3e6838d4d"
	const secondAuthor = "2 9ceea2c8039da1aa931395534064059a2e77"
	const lastAuthor = "1 b76ef9342a9431f3cbd85a775e8b630bdb6d"

	unsorted := append([]Entry{
		Entry{
			Author:   firstAuthor,
			Phrase:   "5ba2388d355805ac6cc2c37edb90aef56da0",
			Meanings: []string{"8ca04ae07a87fb5206a6c774cb1a857df307"},
		},
		Entry{
			Author:   secondAuthor,
			Phrase:   "bf6d6e13ffe18617e904bac818b598a91bd6",
			Meanings: []string{"f7b2897b7b101f9edb094828ed588da36fe6"},
		},
	}, append(tbEntries(t), Entry{
		Author:   lastAuthor,
		Phrase:   "5f478fa6662b8a7d4be27f10bc5c4e2ea92a",
		Meanings: []string{"979b8a0b47fc310e8d1c62895d9600dffe91"},
	})...)

	sorted := GroupEntriesByAuthor(unsorted).Ungroup()

	if !reflect.DeepEqual(unsorted, sorted) {
		t.Fatal(`expected entries with unique authors to be unaffected by author grouping`)
	}

	inserted := Entry{
		Author:   secondAuthor,
		Phrase:   "d1389aa76cbdbc33dee9235a8abdd9ae8d5a",
		Meanings: []string{"4640b6b52a6da70f70085366f727b20610b4"},
	}
	unsorted = append(unsorted, inserted)

	sorted = GroupEntriesByAuthor(unsorted).Ungroup()

	insertionIndex := -1
	for i, entry := range sorted {
		if reflect.DeepEqual(entry, inserted) {
			insertionIndex = i
			break
		}
	}

	const expectedIndex = 2
	if insertionIndex != expectedIndex {
		t.Fatalf(`expected new entry with known author to be colocated with another; expected index: %v; found index: %v`,
			expectedIndex, insertionIndex)
	}
}

func TestDeduplication(t *testing.T) {

}

func TestReadBackingFile(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	content, version, err := ReadBackingFile(context.Background())
	if err != nil {
		t.Fatalf("%+v", err)
	}

	_, err = ParseEntries(content)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if version == "" {
		t.Fatal("expected a non-empty commit hash")
	}
}

func tbEntries(tb testing.TB) []Entry {
	content, err := ioutil.ReadFile(SRC)
	if err != nil {
		tb.Fatalf("%+v", err)
	}

	entries, err := ParseEntries(content)
	if err != nil {
		tb.Fatalf("%+v", err)
	}

	return entries
}

func bn(b *testing.B) []struct{} { return make([]struct{}, b.N) }

func BenchmarkParse(b *testing.B) {
	content, err := ioutil.ReadFile(SRC)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for range bn(b) {
		_, err := ParseEntries(content)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFormat(b *testing.B) {
	entries := tbEntries(b)

	b.ResetTimer()

	for range bn(b) {
		_ = FormatEntries(entries)
	}
}

func BenchmarkGroup(b *testing.B) {
	entries := tbEntries(b)

	b.ResetTimer()

	for range bn(b) {
		_ = GroupEntriesByAuthor(entries)
	}
}

func BenchmarkUngroup(b *testing.B) {
	entries := tbEntries(b)
	grouped := GroupEntriesByAuthor(entries)

	b.ResetTimer()

	for range bn(b) {
		_ = grouped.Ungroup()
	}
}

func writeTempFile(t *testing.T, subpath string, content []byte) string {
	tempDir := os.TempDir()
	if tempDir == "" {
		t.Fatal("failed to create temporary directory: OS API returned empty path")
	}

	dir := filepath.Join(tempDir, TEMPDIR_NAME)
	err := os.MkdirAll(dir, DIR_READ_WRITE_MODE)
	if err != nil {
		t.Fatalf("failed to create temporary directory: %+v", err)
	}

	filePath := filepath.Join(dir, subpath)
	err = ioutil.WriteFile(filePath, content, FILE_READ_WRITE_MODE)
	if err != nil {
		t.Fatalf("failed to write temporary file: %+v", err)
	}

	return filePath
}
