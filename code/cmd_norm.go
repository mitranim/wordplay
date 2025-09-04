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
	entries.Norm()

	reds := gg.Filter(entries, Entry.HasRedundantAuthor)
	if len(reds) > 0 {
		fmt.Println(`redundant authors in entries:`)
		fmt.Println(reds.String())
		os.Exit(1)
	}

	dupes := entries.Dupes()
	if len(dupes) > 0 {
		fmt.Println(`duplicates:`)
		grepr.Println(dupes)
		os.Exit(1)
	}

	entries.ReplaceAuthors(map[string]string{
		`LandRaider`:      `VA`,
		`LR`:              `VA`,
		`VengefulAncient`: `VA`,
		`vengefulancient`: `VA`,
		`LeoJo`:           `LJ`,
		`LeoJo231094`:     `LJ`,
		`El Jay`:          `LJ`,
		`.el.jay.`:        `LJ`,
		`Yury`:            `Y`,
		`Kayez`:           `K`,
		`Kaeyz`:           `K`,
		`kayez`:           `K`,
		`kaeyz`:           `K`,
		`_tigy_`:          `T`,
		`Pablo`:           `P`,
		`Mitranim`:        `M`,
		`mitranim`:        `M`,
		``:                `M`,
	})

	WriteFile(path, entries.String())
}
