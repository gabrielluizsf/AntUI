package svg

import (
	"testing"
	"unicode/utf16"
)

// A drawing that declares itself is held to the declaration: the version has
// to be one this reads and the encoding one this can turn into characters, and
// the declaration may open the document and nowhere else. What it says about
// how it is written is also done rather than believed — a drawing in UTF-16
// arrives as bytes and is read as the text it stands for.

func TestADrawingHoldsToTheVersionItDeclares(t *testing.T) {
	for _, good := range []string{
		`<?xml version="1.0"?><svg/>`,
		`<?xml version='1.1' ?><svg/>`,
		`<?xml version="1.0" encoding="UTF-8"?><svg/>`,
		`<?xml version="1.0" encoding="utf-8" standalone="yes"?><svg/>`,
	} {
		if _, err := parseDocument(good); err != nil {
			t.Errorf("parse %s: %v", good, err)
		}
	}
	for _, bad := range []string{
		`<?xml version="2.0"?><svg/>`,
		`<?xml encoding="UTF-8"?><svg/>`,
		`<?xml version="1.0" encoding="ISO-8859-1"?><svg/>`,
		`<?xml version="1.0" encoding="Windows-1252"?><svg/>`,
		`<?xml version="1.0"`,
		`<?xml VERSION="1.0"?><svg/>`,
		`<?xml version="1.0" foo="bar"?><svg/>`,
		`<?xml encoding="UTF-8" version="1.0"?><svg/>`,
		`<?xml version="1.0" version="1.0"?><svg/>`,
		`<?xml version="1.0" standalone="maybe"?><svg/>`,
		`<?xml version=1.0?><svg/>`,
		`<?xml version="1.0" junk?><svg/>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where the declaration is not one this reads", bad)
		}
	}
}

func TestTheDeclarationOnlyOpensTheDocument(t *testing.T) {
	for _, bad := range []string{
		`<svg><?xml version="1.0"?></svg>`,
		`<!-- note --><?xml version="1.0"?><svg/>`,
		`<svg></svg><?xml version="1.0"?>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where the declaration does not open the document", bad)
		}
	}
	// An instruction whose name merely starts with xml is a plain one and
	// stands where any instruction may.
	if _, err := parseDocument(`<?xml-stylesheet href="a.css"?><svg/>`); err != nil {
		t.Errorf("parse: %v", err)
	}
}

func TestACommentOrInstructionThatWasNeverClosedStopsTheReading(t *testing.T) {
	for _, bad := range []string{
		`<svg><!-- never closed<rect/></svg>`,
		`<svg><? never closed<rect/></svg>`,
		`<svg><style><![CDATA[ never closed</style>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where it was never closed", bad)
		}
	}
}

func TestACommentMayNotHoldARunOfDashes(t *testing.T) {
	for _, bad := range []string{
		`<svg><!-- a -- b --><rect/></svg>`,
		`<svg><!-- a---><rect/></svg>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where the comment holds a run of dashes", bad)
		}
	}
	if _, err := parseDocument(`<svg><!----><!-- fine --></svg>`); err != nil {
		t.Errorf("parse: %v", err)
	}
}

func TestADeclarationInsideTheDrawingIsAnError(t *testing.T) {
	if _, err := parseDocument(`<svg><!DOCTYPE svg><rect/></svg>`); err == nil {
		t.Error("parse: no error, where a doctype stands inside the drawing")
	}
}

func TestTheWritingMayNotHoldACDataClose(t *testing.T) {
	if _, err := parseDocument(`<svg>a]]>b</svg>`); err == nil {
		t.Error("parse: no error, where the writing holds a `]]>`")
	}
	if _, err := parseDocument(`<svg><style><![CDATA[a]b]]></style></svg>`); err != nil {
		t.Errorf("parse: %v", err)
	}
}

// utf16Of turns text into a UTF-16 file: its byte-order mark, then one unit
// per two bytes, big endian or little.
func utf16Of(s string, le bool) string {
	units := utf16.Encode([]rune(s))
	out := []byte{0xFE, 0xFF}
	if le {
		out = []byte{0xFF, 0xFE}
	}
	for _, u := range units {
		if le {
			out = append(out, byte(u), byte(u>>8))
		} else {
			out = append(out, byte(u>>8), byte(u))
		}
	}
	return string(out)
}

func TestADrawingInUTF16IsReadAsWhatItStandsFor(t *testing.T) {
	const src = `<svg width="2" height="2"><rect width="2" height="2" fill="#00ff00"/></svg>`
	for _, file := range []string{utf16Of(src, false), utf16Of(src, true)} {
		_, cv := painted(t, file, 2, 2)
		if got := cv.At(1, 1); got.G() != 0xff {
			t.Errorf("pixel = %v, want the green the drawing declares", got)
		}
	}
}
