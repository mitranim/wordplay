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

		entry.ReplaceAuthor(
			[]string{`LandRaider`, `VengefulAncient`, `vengefulancient`},
			`LR`,
		)

		entry.ReplaceAuthor(
			[]string{`LeoJo`, `LeoJo231094`, `El Jay`, `.el.jay.`},
			`LJ`,
		)

		entry.ReplaceAuthor([]string{`Yury`}, `Y`)

		entry.ReplaceAuthor([]string{``, `Mitranim`, `mitranim`}, `M`)
	}

	dupes := entries.Dupes()
	if len(dupes) > 0 {
		fmt.Println(`duplicates:`)
		grepr.Println(dupes)
		os.Exit(1)
	}

	WriteFile(path, entries.String())
}
