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
	for _, bad := range []string{
		`<rect width="1" width="2"/>`,
		`<rect hidden hidden/>`,
		`<rect x="1" X="2"/>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where the attribute is repeated", bad)
		}
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
		`<rect xlink:href="#a"/>`,
		`<text xml:space="preserve"> </text>`,
		`<a_b.c-d:e/>`,
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
