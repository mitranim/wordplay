package main

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

	writeFile(SRC_FILE, entries.Bytes())
}
