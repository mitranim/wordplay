package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitranim/repr"
	"github.com/mitranim/try"
)

var tSrc = bytesToMutableString(bytes.TrimSpace(try.ByteSlice(os.ReadFile(SRC_FILE))))

func TestParseAndFormat(t *testing.T) {
	source := tSrc
	entries := ParseEntries(source)
	output := strings.TrimSpace(entries.String())

	if testing.Verbose() {
		fmt.Printf("source:  %v\n", writeTempFile(t, "source", stringToBytesAlloc(source)))
		fmt.Printf("output:  %v\n", writeTempFile(t, "output", stringToBytesAlloc(output)))
		fmt.Printf("entries: %v\n", writeTempFile(t, "entries", repr.Bytes(entries)))
	}

	if source != output {
		if testing.Verbose() {
			t.Fatalf("mismatch of source and formatted output")
		} else {
			t.Fatal("mismatch of source and formatted output; run test in verbose mode for details")
		}
	}
}

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

func writeTempFile(t *testing.T, subpath string, content []byte) string {
	tempDir := os.TempDir()
	if tempDir == "" {
		t.Fatal("failed to create temporary directory: got empty path")
	}

	path := filepath.Join(tempDir, `wordplay_testing`, subpath)
	writeFile(path, content)
	return path
}