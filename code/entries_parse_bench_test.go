package main

import (
	"testing"
)

func BenchmarkParse(b *testing.B) {
	for range counter(b.N) {
		_ = ParseEntries(tSrc)
	}
}

func BenchmarkFormat(b *testing.B) {
	entries := ParseEntries(tSrc)

	b.ResetTimer()

	for range counter(b.N) {
		_ = FormatEntries(entries)
	}
}

func BenchmarkGroup(b *testing.B) {
	entries := ParseEntries(tSrc)

	b.ResetTimer()

	for range counter(b.N) {
		_ = GroupEntriesByAuthor(entries)
	}
}

func BenchmarkUngroup(b *testing.B) {
	entries := ParseEntries(tSrc)
	grouped := GroupEntriesByAuthor(entries)

	b.ResetTimer()

	for range counter(b.N) {
		_ = grouped.Ungroup()
	}
}
