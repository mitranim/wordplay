package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/mitranim/gg/gtest"
)

var tSrc = strings.TrimSpace(gg.ReadFile[string](SRC_FILE))

func Test_parse_and_format(t *testing.T) {
	defer gtest.Catch(t)

	source := tSrc
	entries := ParseEntries(source)
	output := strings.TrimSpace(entries.String())

	if testing.Verbose() {
		fmt.Printf("source:  %v\n", writeTempFile(`source`, source))
		fmt.Printf("output:  %v\n", writeTempFile(`output`, output))
		fmt.Printf("entries: %v\n", writeTempFile(`entries`, grepr.String(entries)))
	}

	if source != output {
		if testing.Verbose() {
			t.Fatal(`mismatch of source and formatted output`)
		} else {
			t.Fatal(`mismatch of source and formatted output; run test in verbose mode for details`)
		}
	}
}

func Benchmark_parse(b *testing.B) {
	defer gtest.Catch(b)

	for range counter(b.N) {
		_ = ParseEntries(tSrc)
	}
}

func Benchmark_format(b *testing.B) {
	defer gtest.Catch(b)

	entries := ParseEntries(tSrc)
	b.ResetTimer()

	for range counter(b.N) {
		_ = entries.Bytes()
	}
}

func writeTempFile(subpath, content string) string {
	tempDir := os.TempDir()
	gtest.NotZero(tempDir, `need non-empty path for temporary directory`)

	path := filepath.Join(tempDir, `wordplay_testing`, subpath)
	writeFile(path, content)
	return path
}
