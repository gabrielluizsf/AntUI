package svg

import "testing"

func TestSniffSVGKnowsADrawingFromAPhotograph(t *testing.T) {
	// Only the front matters, and what comes before the tag — a byte-order
	// mark, an xml declaration, a comment, a doctype — is read past to find it.
	drawings := []string{
		`<svg width="1" height="1"/>`,
		"\ufeff<svg/>",
		`<?xml version="1.0"?><svg/>`,
		"<!-- a note --><svg>",
		`<!DOCTYPE svg [<!ENTITY a "1">]><svg/>`,
	}
	for _, s := range drawings {
		if !sniffSVG([]byte(s)) {
			t.Errorf("sniff %q: said a photograph", s)
		}
	}
	photographs := []string{
		"",
		"not a drawing",
		"<?xml version=\"1.0\"?><html/>",
		"<svgfoo/>",
		"<!DOCTYPE svg [<!ENTITY a \"1\">", // a doctype that was never closed
	}
	for _, s := range photographs {
		if sniffSVG([]byte(s)) {
			t.Errorf("sniff %q: said a drawing", s)
		}
	}
}

func TestSniffSVGReadsADrawingWrittenTwoBytesToACharacter(t *testing.T) {
	// A UTF-16 file opens with its own byte-order mark, and its tags only mean
	// anything once the bytes behind the mark are read as the characters they
	// stand for.
	const src = `<svg width="1" height="1"/>`
	for _, le := range []bool{false, true} {
		if !sniffSVG([]byte(utf16Of(src, le))) {
			t.Errorf("sniff a UTF-16 drawing (le=%v): said a photograph", le)
		}
	}
	if sniffSVG([]byte(utf16Of("<html/>", false))) {
		t.Error("sniff a UTF-16 document that is no drawing: said a drawing")
	}
}
