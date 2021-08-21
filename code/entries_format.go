package main

// Assumes entries have been normalized. TODO enforce normalization: sort by
// author first.
func FormatEntries(entries []Entry) (buf []byte) {
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

func FormatEntriesNew(entries []Entry) (buf []byte) {
	for _, entry := range entries {
		buf = entry.Append(buf)
	}
	return buf
}

func isBlank(str string) bool {
	for _, char := range str {
		switch char {
		case ' ', '\t', '\v', '\r', '\n':
		default:
			return false
		}
	}
	return true
}
