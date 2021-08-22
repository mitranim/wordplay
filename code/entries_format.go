package main

// Assumes entries have been normalized. TODO enforce normalization: sort by
// author first.
func FormatEntriesOld(entries Entries) (buf []byte) {
	var author string

	for _, entry := range entries {
		if entry.Author != author {
			author = entry.Author
			buf = entry.AppendAuthorOld(buf)
		}

		buf = entry.AppendOld(buf)
	}

	return buf
}

func FormatEntries(entries Entries) (buf []byte) {
	for _, entry := range entries {
		buf = entry.Append(buf)
	}
	return buf
}
