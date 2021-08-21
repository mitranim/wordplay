package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitranim/repr"
	"github.com/mitranim/try"
	e "github.com/pkg/errors"
)

const SRC = "../readme.md"
const TEMPDIR_NAME = "wordplay_testing"

var tSrc = bytesToMutableString(bytes.TrimSpace(try.ByteSlice(os.ReadFile(SRC))))

func TestParseAndFormat(t *testing.T) {
	source := tSrc
	entries := ParseEntries(source)
	output := bytesToMutableString(bytes.TrimSpace(FormatEntries(entries)))

	if testing.Verbose() {
		fmt.Printf("source: %v\n", writeTempFile(t, "source", stringToBytesAlloc(source)))
		fmt.Printf("output: %v\n", writeTempFile(t, "output", stringToBytesAlloc(output)))
		fmt.Printf("output new: %v\n", writeTempFile(t, "output_new", FormatEntriesNew(entries)))
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

func writeTempFile(t *testing.T, subpath string, content []byte) string {
	tempDir := os.TempDir()
	if tempDir == "" {
		t.Fatal("failed to create temporary directory: OS API returned empty path")
	}

	dir := filepath.Join(tempDir, TEMPDIR_NAME)
	err := os.MkdirAll(dir, os.ModePerm)
	try.To(e.Wrap(err, `failed to create temporary directory`))

	filePath := filepath.Join(dir, subpath)
	err = os.WriteFile(filePath, content, os.ModePerm)
	try.To(e.Wrap(err, `failed to write temporary file`))

	return filePath
}
