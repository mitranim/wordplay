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

	dupes := entries.Dupes()

	if len(dupes) > 0 {
		fmt.Println(`duplicates:`)
		grepr.Println(dupes)
		os.Exit(1)
	}

	entries.ReplaceAuthors(map[string]string{
		`LandRaider`:      `LR`,
		`VengefulAncient`: `LR`,
		`vengefulancient`: `LR`,
		`LeoJo`:           `LJ`,
		`LeoJo231094`:     `LJ`,
		`El Jay`:          `LJ`,
		`.el.jay.`:        `LJ`,
		`Yury`:            `Y`,
		`Mitranim`:        `M`,
		`mitranim`:        `M`,
		``:                `M`,
	})

	WriteFile(path, entries.String())
}
