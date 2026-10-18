package svg

// isStrongLTR says whether a rune runs left to right on its own.
func isStrongLTR(r rune) bool {
	switch {
	case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z':
		return true
	case r >= 0x00C0 && r <= 0x02AF: // Latin-1, Latin Extended, IPA
		return true
	case r >= 0x0370 && r <= 0x058F: // Greek, Cyrillic, Armenian
		return true
	case r >= 0x0900 && r <= 0x097F: // Devanagari
		return true
	case r >= 0x0E00 && r <= 0x0E7F: // Thai
		return true
	case r >= 0x3040 && r <= 0x30FF: // Hiragana, Katakana
		return true
	case r >= 0x4E00 && r <= 0x9FFF: // CJK unified ideographs
		return true
	case r >= 0xAC00 && r <= 0xD7AF: // Hangul syllables
		return true
	}
	return false
}

// isStrongRTL says whether a rune runs right to left on its own.
func isStrongRTL(r rune) bool {
	switch {
	case r >= 0x0590 && r <= 0x08FF: // Hebrew, Arabic, Syriac, Thaana, NKo
		return true
	case r >= 0xFB1D && r <= 0xFDFF: // Hebrew and Arabic presentation forms
		return true
	case r >= 0xFE70 && r <= 0xFEFF: // Arabic presentation forms-B
		return true
	}
	return false
}
