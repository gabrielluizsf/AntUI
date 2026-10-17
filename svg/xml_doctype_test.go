package svg

import "testing"

func TestADoctypeDeclaresWhatTheDrawingMayLeanon(t *testing.T) {
	// A doctype may name its root, point at an external DTD, and hold an
	// internal subset. Of all that, what a drawing leans on is its own names:
	// every other declaration is read past, since which elements a file holds
	// is a matter of validation this renderer does not keep to.
	for _, src := range []string{
		`<!DOCTYPE svg><svg/>`,
		`<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "svg11.dtd"><svg/>`,
		`<!DOCTYPE svg SYSTEM "svg.dtd"><svg/>`,
		`<!DOCTYPE svg [<!ELEMENT svg ANY>]><svg/>`,
		`<!DOCTYPE svg [<!ATTLIST svg x CDATA "a>b">]><svg/>`,
		`<!DOCTYPE svg [<!-- a note -->]><svg/>`,
		`<!DOCTYPE svg PUBLIC "a[b" "c>d"><svg/>`,
	} {
		if _, err := parseDocument(src); err != nil {
			t.Errorf("parse %s: %v", src, err)
		}
	}
}

func TestADoctypeThatIsNeverClosedIsAnError(t *testing.T) {
	for _, src := range []string{
		`<!DOCTYPE svg [<!ELEMENT svg ANY><svg/>`,
		`<!DOCTYPE svg`,
		`<!DOCTYPE svg "unclosed><svg/>`,
		`<!DOCTYPE svg [<!ELEMENT svg "a><svg/>`,
	} {
		if _, err := parseDocument(src); err == nil {
			t.Errorf("%s was read, want an error", src)
		}
	}
}

func TestOnlyOneDoctypeOpensADocument(t *testing.T) {
	if _, err := parseDocument(`<!DOCTYPE svg><!DOCTYPE svg><svg/>`); err == nil {
		t.Error("two doctypes were read, want an error")
	}
}

func TestADeclarationThatIsNotADoctypeDoesNotOpenADocument(t *testing.T) {
	// The prolog holds the doctype and nothing else; an element declaration or
	// a stray word has no doctype to belong to.
	for _, src := range []string{
		`<!ELEMENT svg ANY><svg/>`,
		`<!ENTITY a "1"><svg/>`,
		`<!DOCTYPE svg [junk]><svg/>`,
	} {
		if _, err := parseDocument(src); err == nil {
			t.Errorf("%s was read, want an error", src)
		}
	}
}

func TestANameTheDoctypeLeansOnIsReadInItsPlace(t *testing.T) {
	// A parameter entity — %name; — holds pieces of the doctype itself, and
	// what it holds may be an entity declaration of the drawing's or the text
	// of one.
	for _, tc := range []struct {
		src, fill string
	}{
		{`<!DOCTYPE svg [<!ENTITY % p "red"><!ENTITY c "%p;">]><rect fill="&c;"/>`, "red"},
		{`<!DOCTYPE svg [<!ENTITY % e "<!ENTITY c 'red'>">%e;]><rect fill="&c;"/>`, "red"},
	} {
		e, err := parseDocument(tc.src)
		if err != nil {
			t.Fatalf("parse %s: %v", tc.src, err)
		}
		if got := e.attr("fill"); got != tc.fill {
			t.Errorf("fill = %q, want %q", got, tc.fill)
		}
	}
}

func TestANameTheDoctypeDoesNotDefineIsAnError(t *testing.T) {
	// A parameter entity that no declaration holds, one whose text another
	// document holds, and one whose text comes back around to itself are each
	// an error: the doctype is read once, and nothing fetches the rest.
	for _, src := range []string{
		`<!DOCTYPE svg [%nope;]><rect x="1"/>`,
		`<!DOCTYPE svg [<!ENTITY % e SYSTEM "http://example.com/e.dtd">%e;]><rect x="1"/>`,
		`<!DOCTYPE svg [<!ENTITY % a "%b;"><!ENTITY % b "%a;">%a;]><rect x="1"/>`,
	} {
		if _, err := parseDocument(src); err == nil {
			t.Errorf("%s was read, want an error", src)
		}
	}
}
