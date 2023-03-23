package main

import (
	"fmt"
	"os"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

func CmdNorm() {
	NormFile(`../readme.md`)
	NormFile(`../readme_ru.md`)
}

func NormFile(path string) {
	entries := ParseEntries(gg.ReadFile[string](path))

	for ind := range entries {
		entry := &entries[ind]

		if entry.Author == `LandRaider` || entry.Author == `VengefulAncient` {
			entry.Author = `LR`
		}

		if entry.Author == `LeoJo` || entry.Author == `LeoJo231094` {
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
