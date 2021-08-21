package main

import (
	"context"
	"testing"
)

func TestReadBackingFile(t *testing.T) {
	t.Skip()
	if testing.Short() {
		t.Skip()
	}

	content, version := ReadBackingFile(context.Background())
	_ = ParseEntries(content)

	if version == "" {
		t.Fatal("expected a non-empty commit hash")
	}
}
