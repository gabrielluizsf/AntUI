package svg

import "strings"

// normalizeBreaks reads the ways a file may end a line — a lone return, the
// pair a Windows editor writes, or the plain line feed every other one does —
// as the one line feed XML reads them all as, before anything else in the text
// is looked at. Without it the same drawing would hold different characters
// depending on which editor saved it.
func normalizeBreaks(s string) string {
	if !strings.ContainsRune(s, '\r') {
		return s
	}
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
}
