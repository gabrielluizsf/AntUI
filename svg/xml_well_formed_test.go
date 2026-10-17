package svg

import (
	"strings"
	"testing"
)

// A drawing that is not well-formed XML is refused rather than read into a tree
// that differs from the one declared, and a drawing that is well-formed — even
// in the corners other readers get wrong — is read. These check both halves:
// the spellings XML allows and the ones it does not.

func TestADrawingOpensPastItsByteOrderMark(t *testing.T) {
	e, err := parseDocument("\ufeff<rect/>")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if e.Name != "rect" {
		t.Errorf("name = %q, want rect", e.Name)
	}
}

func TestAClosingTagMayHoldSpaceBeforeTheBracket(t *testing.T) {
	for _, src := range []string{`<g></g>`, `<g></g >`, `<g></g  >`} {
		if _, err := parseDocument(src); err != nil {
			t.Errorf("parse %s: %v", src, err)
		}
	}
}

func TestADrawingEndsAtItsRootElement(t *testing.T) {
	for _, bad := range []string{
		`<svg></svg><svg></svg>`,
		`<svg/>trailing`,
		`<svg></svg></g>`,
		"<svg/>\ntext",
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where the drawing carries on past its root", bad)
		}
	}
}

func TestADrawingMayHoldACommentOrInstructionAfterItsRoot(t *testing.T) {
	for _, good := range []string{
		`<svg></svg><!-- done -->`,
		`<svg></svg><?rendered by a tool?>`,
		"  <svg></svg>  \n\t",
		`<svg></svg><!-- and -->`,
	} {
		if _, err := parseDocument(good); err != nil {
			t.Errorf("parse %s: %v", good, err)
		}
	}
}

func TestAnAttributeMayNotBeWrittenTwiceOnOneTag(t *testing.T) {
	// A name repeated exactly is an error; a name in another case is another
	// name, as XML reads one, and says nothing about the first.
	for _, bad := range []string{
		`<rect width="1" width="2"/>`,
		`<rect hidden="1" hidden="2"/>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where the attribute is repeated", bad)
		}
	}
	if _, err := parseDocument(`<rect x="1" X="2"/>`); err != nil {
		t.Errorf("parse: %v, where the two names are different names", err)
	}
}

func TestAnAttributeNameIsReadAsTheFileSpellsIt(t *testing.T) {
	// `viewBox` and `viewbox` are two different names, and only the one XML
	// spells means the viewBox: a file that writes the other has said nothing
	// about it, the same as it says nothing to a browser.
	e, err := parseDocument(`<svg viewbox="0 0 1 1" width="4"/>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if e.attr("viewBox") != "" {
		t.Errorf("a lowercase viewbox was read as the viewBox: %q", e.attr("viewBox"))
	}
	if got := e.attr("viewbox"); got != "0 0 1 1" {
		t.Errorf("the name as written is kept: %q", got)
	}
}

func TestAnAttributeValueMayNotHoldALessThan(t *testing.T) {
	for _, bad := range []string{
		`<rect x="a<b"/>`,
		`<rect x='a<b'/>`,
		`<rect x=a<b/>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where the value holds a `<`", bad)
		}
	}
}

func TestATagOrAttributeNameIsOneXMLAllows(t *testing.T) {
	for _, bad := range []string{
		`<1rect/>`,
		`<re"ct/>`,
		`<g!/>`,
		`<rect 1x="2"/>`,
		`<rect x="1" st.roke<=2/>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where the name is not one XML allows", bad)
		}
	}
}

func TestNamesXMLAllowsAreRead(t *testing.T) {
	for _, good := range []string{
		`<rect x="1"/>`,
		`<linearGradient id="g1"/>`,
		`<rect xlink:href="#a" xmlns:xlink="http://www.w3.org/1999/xlink"/>`,
		`<text xml:space="preserve"> </text>`,
		`<a_b.c-d:e xmlns:a_b.c-d="http://www.w3.org/2000/svg"/>`,
	} {
		if _, err := parseDocument(good); err != nil {
			t.Errorf("parse %s: %v", good, err)
		}
	}
}

// The name a tag is closed by has to be its own, and the message says both.
func TestAMismatchedCloseSaysWhatClosedIt(t *testing.T) {
	_, err := parseDocument(`<g></rect>`)
	if err == nil {
		t.Fatal("parse: no error")
	}
	if !strings.Contains(err.Error(), "closed by </rect>") {
		t.Errorf("error = %q, want it to name </rect>", err)
	}
}

func TestACharacterXMLKeepsOutOfADocumentIsAnError(t *testing.T) {
	// XML allows tab, line break and carriage return and nothing else below a
	// space, and every byte has to be part of the UTF-8 the document is
	// written in. A control character or a stray byte is an error wherever it
	// stands — in the text, in a value, in a comment or in a CDATA section —
	// because a reader that passed one on would paint a document the file does
	// not hold.
	for _, bad := range []string{
		"<svg>\x01</svg>",
		"<rect x=\"\x0b\"/>",
		"<!--\x0c--><svg/>",
		"<svg><![CDATA[\x00]]></svg>",
		"<svg>\xff</svg>",
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %q: no error, where a character XML keeps out stands", bad)
		}
	}
	if _, err := parseDocument("<svg>\t\n\r</svg>"); err != nil {
		t.Errorf("parse: %v, where the characters are the three XML allows", err)
	}
}
