package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/mitranim/try"
	e "github.com/pkg/errors"
)

const TIMES_FILE = `../fixtures/entry_times.json`

func init() { commands[`times`] = cmdTimes }

func cmdTimes() {
	times := TimeMap{}

	for _, val := range readEntries(SRC_FILE) {
		times.AddMin(val.Phrase, entryTimeMin(val))
	}

	writeFile(TIMES_FILE, try.ByteSlice(json.Marshal(times)))
}

type TimeMap map[string]time.Time

func (self TimeMap) AddMin(key string, inst time.Time) {
	prev, ok := self[key]

	if ok {
		self[key] = timeMin(prev, inst)
	} else {
		self[key] = inst
	}
}

func readEntries(path string) Entries {
	return ParseEntries(bytesToMutableString(try.ByteSlice(os.ReadFile(path))))
}

func entryTimeMin(entry Entry) (inst time.Time) {
	defer func() {
		if inst.IsZero() {
			panic(e.Errorf(`missing time for entry %#v`, entry))
		}
	}()

	for _, line := range splitLines(gitLogSearch(entry.Phrase)) {
		if inst.IsZero() {
			inst = findIsoTime(line)
		} else {
			inst = timeMin(inst, findIsoTime(line))
		}
	}

	return
}

func gitLogSearch(search string) string {
	return runCmdOut("git", "log", "--format=reference", "--date=iso-strict", "-S", search)
}
