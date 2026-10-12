package template

// wrapWord picks where a word that overflows the line ends it. It returns the
// index to resume at and true when the line breaks there; -1 and false keep
// the word on the current line.
func wrapWord(text string, wordStart, wordEnd, wordW, brk, w, maxRef int, breakWord bool) (int, bool) {
	if brk != -1 {
		return brk, true
	}
	if breakWord {
		return cutWord(text, wordStart, wordEnd, maxRef-w), true
	}
	if w == 0 {
		// No prior content and no wrap point: the overlong word fills this
		// line and the scan keeps going.
		return -1, false
	}
	// Prior content but no fitting break: the word spills onto its own
	// overlong line, as CSS overflow-wrap normal allows.
	return wordEnd, true
}
