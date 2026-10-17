package svg

import "strings"

// decodeText takes a document as it sits in a file and hands back the text the
// parser reads: the line breaks a text file may hold in any of three ways read
// as one, then UTF-8, with the byte-order mark Windows editors leave at the
// front taken off, or UTF-16 — which arrives behind a mark of its own and is
// turned into UTF-8 here rather than read as characters it never wrote. A
// drawing that says `x="1<CR><LF>2"` and one that says `x="1<LF>2"` are the
// same drawing to XML, and they are read as the same drawing here.
func decodeText(src string) string {
	switch {
	case strings.HasPrefix(src, "\ufeff"):
		src = strings.TrimPrefix(src, "\ufeff")
	case strings.HasPrefix(src, "\xFF\xFE"), strings.HasPrefix(src, "\xFE\xFF"):
		src = utf16ToUTF8(src)
	}
	return normalizeBreaks(src)
}
