package svg

import "bytes"

// sniffSVG says whether the bytes are a drawing rather than a photograph: the
// first tag they open with — past the byte-order mark, an XML declaration, a
// comment or a doctype in front of it — is an `<svg>`. The mark may be the one
// UTF-8 writes or the one a UTF-16 file opens with, which says the drawing is
// written two bytes to a character and has to be read that way before its tags
// mean anything. Only the front is looked at, the same as every other format
// this reads by what it begins with, so what a photograph carries in its
// metadata after the picture itself is never mistaken for the picture.
func sniffSVG(data []byte) bool {
	if len(data) >= 2 && (data[0] == 0xFF && data[1] == 0xFE || data[0] == 0xFE && data[1] == 0xFF) {
		data = []byte(utf16ToUTF8(string(data)))
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	for {
		switch {
		case bytes.HasPrefix(data, []byte("<!--")):
			data = past(data, "-->")
		case bytes.HasPrefix(data, []byte("<?")):
			data = past(data, ">")
		case bytes.HasPrefix(data, []byte("<!")):
			data = pastDeclaration(data)
		case bytes.HasPrefix(data, []byte("<svg")):
			rest := data[len("<svg"):]
			return len(rest) == 0 || rest[0] == '>' || rest[0] == '/' || isSpace(rest[0])
		default:
			return false
		}
	}
}

// pastDeclaration is everything after a doctype, which holds its whole
// internal subset between brackets when it has one and only ends at the
// bracket's own close.
func pastDeclaration(data []byte) []byte {
	lt := bytes.Index(data, []byte(">"))
	lb := bytes.Index(data, []byte("["))
	if lb >= 0 && (lt < 0 || lb < lt) {
		return past(data, "]>")
	}
	return past(data, ">")
}

// past is everything after the first end marker, or nothing when the text
// never says it — a declaration left open is not a drawing this can read.
func past(data []byte, end string) []byte {
	i := bytes.Index(data, []byte(end))
	if i < 0 {
		return nil
	}
	return data[i+len(end):]
}
