package svg

import "testing"

func TestAnEntityTheDrawingDefinesIsReadIntoItsText(t *testing.T) {
	// A drawing may declare its own names in the doctype's internal subset,
	// and every use of one stands for the text it was given — in an attribute
	// value and in the writing of an element alike.
	const decl = `<!DOCTYPE svg [<!ENTITY color "red"><!ENTITY word "a &amp; b">]>`
	e, err := parseDocument(decl + `<rect fill="&color;"/>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := e.attr("fill"); got != "red" {
		t.Errorf("fill = %q, want red", got)
	}
	g, err := parseDocument(decl + `<g>&word;</g>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got, want := g.Text, "a & b"; got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

func TestAnEntityHoldsOtherEntities(t *testing.T) {
	// What an entity stands for may name another in turn, as long as the names
	// do not come back around to themselves.
	e, err := parseDocument(`<!DOCTYPE svg [<!ENTITY a "1"><!ENTITY b "&a;2"><!ENTITY c "&b;3">]><rect x="&c;"/>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := e.attr("x"); got != "123" {
		t.Errorf("x = %q, want 123", got)
	}
}

func TestAnEntityThatComesBackAroundIsAnError(t *testing.T) {
	// An entity whose text names itself, directly or by way of another, would
	// never stop being read.
	for _, src := range []string{
		`<!DOCTYPE svg [<!ENTITY a "&a;">]><rect x="&a;"/>`,
		`<!DOCTYPE svg [<!ENTITY a "&b;"><!ENTITY b "&a;">]><rect x="&a;"/>`,
	} {
		if _, err := parseDocument(src); err == nil {
			t.Errorf("%s was read, want an error", src)
		}
	}
}

func TestAnEntityAnotherDocumentHoldsIsAnError(t *testing.T) {
	// The text of an entity written in another document is not here to read,
	// and nothing in the package fetches it.
	src := `<!DOCTYPE svg [<!ENTITY e SYSTEM "http://example.com/e.xml">]><rect x="&e;"/>`
	if _, err := parseDocument(src); err == nil {
		t.Error("an entity of another document was read, want an error")
	}
}

func TestAnEntityXMLDefinesIsNotAFileToDefine(t *testing.T) {
	// The five entities XML defines for every document, and a name a file
	// gives twice, are not a file's to say.
	for _, src := range []string{
		`<!DOCTYPE svg [<!ENTITY lt "x">]><rect x="1"/>`,
		`<!DOCTYPE svg [<!ENTITY amp "x">]><rect x="1"/>`,
		`<!DOCTYPE svg [<!ENTITY a "1"><!ENTITY a "2">]><rect x="&a;"/>`,
	} {
		if _, err := parseDocument(src); err == nil {
			t.Errorf("%s was read, want an error", src)
		}
	}
}

func TestAnEntityTheDrawingNeverClosedIsAnError(t *testing.T) {
	for _, src := range []string{
		`<rect x="&amp"/>`,
		`<rect x="&no-semicolon"/>`,
		`<g>&amp</g>`,
	} {
		if _, err := parseDocument(src); err == nil {
			t.Errorf("%s was read, want an error", src)
		}
	}
}

func TestACharacterReferenceStandsForItsCharacter(t *testing.T) {
	e, err := parseDocument(`<g>&#65;&#x42;&#x1F600;</g>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got, want := e.Text, "AB\U0001F600"; got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

func TestAReferenceToACharacterXMLKeepsOutIsAnError(t *testing.T) {
	// XML keeps the control characters and the unpaired surrogates out of a
	// document entirely, and a reference to one of them names no character a
	// reader may be handed.
	for _, src := range []string{
		`<rect x="&#0;"/>`,
		`<rect x="&#xD800;"/>`,
		`<rect x="&#x110000;"/>`,
		`<rect x="&#;"/>`,
		`<rect x="&#X41;"/>`,
		`<g>&#x0;</g>`,
	} {
		if _, err := parseDocument(src); err == nil {
			t.Errorf("%s was read, want an error", src)
		}
	}
}

func TestCDATAKeepsItsEntitiesAsWritten(t *testing.T) {
	// What a CDATA section holds is the writing itself, where an ampersand is
	// a character and starts nothing.
	e, err := parseDocument(`<style><![CDATA[ .a { content: "&amp;" } ]]></style>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got, want := e.Text, `.a { content: "&amp;" }`; got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}
