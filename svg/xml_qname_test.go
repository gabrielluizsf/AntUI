package svg

import "testing"

func TestANameMayBeWrittenWithOnePrefix(t *testing.T) {
	// A name parts in two at most, a prefix in front of a local part, whether
	// the name is a tag's or the element binds the prefix it wears.
	for _, src := range []string{
		`<a_b.c-d:e/>`,
		`<circle xmlns:s="http://www.w3.org/2000/svg" s:fill="red"/>`,
	} {
		if _, err := parseDocument(src); err != nil {
			t.Errorf("parse %s: %v", src, err)
		}
	}
}

func TestANameXMLDoesNotAllowIsAnError(t *testing.T) {
	// Two colons, a colon with nothing on one side of it, and a name opening
	// with a digit are not names, whatever a drawing means by them.
	for _, src := range []string{
		`<a:b:c/>`,
		`<:a/>`,
		`<a:/>`,
		`<1a/>`,
		`<rect a:b:c="1"/>`,
		`<rect :a="1"/>`,
		`<rect a:="1"/>`,
		`<rect 1a="1"/>`,
	} {
		if _, err := parseDocument(src); err == nil {
			t.Errorf("%s was read, want an error", src)
		}
	}
}

func TestATargetXMLBelongsToTheDeclaration(t *testing.T) {
	// XML reserves the target `xml` for the declaration that opens a document,
	// whatever case a file spells it in; an instruction with another target is
	// just an instruction and is read past.
	if _, err := parseDocument(`<?XML version="1.0"?><svg/>`); err == nil {
		t.Error("an instruction wearing the target XML was read, want an error")
	}
	if _, err := parseDocument(`<svg><?xMl?></svg>`); err == nil {
		t.Error("an instruction wearing the target xMl was read, want an error")
	}
	for _, src := range []string{
		`<?xml version="1.0"?><svg/>`,
		`<?xml-stylesheet type="text/xsl"?><svg/>`,
		`<svg><?php echo 1; ?></svg>`,
	} {
		if _, err := parseDocument(src); err != nil {
			t.Errorf("parse %s: %v", src, err)
		}
	}
}
