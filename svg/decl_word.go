package svg

// declWord answers where a word of the declaration stands in the order the
// grammar gives them — a version, then an encoding, then a standalone marker —
// or -1 for a word the declaration may not hold at all. The words are read as
// XML spells them: `VERSION` is not the word `version`.
func declWord(name string) int {
	switch name {
	case "version":
		return 0
	case "encoding":
		return 1
	case "standalone":
		return 2
	}
	return -1
}
