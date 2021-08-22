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
		_ = entries.Bytes()
	}
}
