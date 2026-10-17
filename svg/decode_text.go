package svg

import "strings"

// decodeText takes a document as it sits in a file and hands back the text the
// parser reads: UTF-8, with the byte-order mark Windows editors leave at the
// front taken off, or UTF-16 — which arrives behind a mark of its own and is
// turned into UTF-8 here rather than read as characters it never wrote.
func decodeText(src string) string {
	switch {
	case strings.HasPrefix(src, "\ufeff"):
		return strings.TrimPrefix(src, "\ufeff")
	case strings.HasPrefix(src, "\xFF\xFE"), strings.HasPrefix(src, "\xFE\xFF"):
		return utf16ToUTF8(src)
	}
	return src
}
