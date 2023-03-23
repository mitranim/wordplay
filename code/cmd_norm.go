package main

import (
	"fmt"
	"os"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

func cmdNorm() {
	normFile(`../readme.md`)
	normFile(`../readme_ru.md`)
}

func normFile(path string) {
	entries := ParseEntries(gg.ReadFile[string](path))

	for ind := range entries {
		entry := &entries[ind]

		if entry.Author == `LandRaider` {
			entry.Author = `LR`
		}

		if entry.Author == `LeoJo` {
			entry.Author = `LJ`
		}

		if entry.Author == `` || entry.Author == `Mitranim` {
			entry.Author = `M`
		}
	}

	dupes := entries.Dupes()
	if len(dupes) > 0 {
		fmt.Println(`duplicates:`)
		grepr.Println(dupes)
		os.Exit(1)
	}

	WriteFile(path, entries.String())
}
