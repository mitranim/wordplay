package main

import (
	"fmt"
	"os"

	"github.com/mitranim/repr"
)

func cmdNorm() {
	normFile(`../readme.md`)
	normFile(`../readme_ru.md`)
}

func normFile(path string) {
	entries := ParseEntries(readFileString(path))

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

	writeFile(path, entries.Bytes())
}
