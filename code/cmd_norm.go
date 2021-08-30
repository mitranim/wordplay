package main

import (
	"fmt"
	"os"

	"github.com/mitranim/repr"
)

func cmdNorm() {
	entries := ParseEntries(readFileString(SRC_FILE))

	for i := range entries {
		entry := &entries[i]

		if entry.Author == "LandRaider" {
			entry.Author = "LR"
		}

		if entry.Author == "LeoJo" {
			entry.Author = "LJ"
		}

		if entry.Author == "" || entry.Author == "Mitranim" {
			entry.Author = "M"
		}
	}

	dupes := entries.Dupes()
	if len(dupes) > 0 {
		fmt.Println(`duplicates:`)
		repr.Println(dupes)
		os.Exit(1)
	}

	writeFile(SRC_FILE, entries.Bytes())
}
