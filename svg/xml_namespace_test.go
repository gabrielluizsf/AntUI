package svg

import "testing"

// An element is drawn through the namespace it stands in, not through the
// letters of its local name: a prefix or a default bound anywhere but the SVG
// namespace is a drawing within the drawing, and what is written in it is not
// this package's to paint. A prefix nothing declares is an error, as it is to
// an XML Namespaces reader: a drawing wearing one is not the drawing it
// declares, and guessing at what it stood for would paint what the file did
// not say.

const greenRect = `<rect width="4" height="4" fill="#00ff00"/>`

func TestAnElementOutsideTheSVGNamespaceIsNotPainted(t *testing.T) {
	cases := []string{
		`<svg width="4" height="4"><g xmlns="http://example.com"><rect width="4" height="4" fill="#ff0000"/></g>` + greenRect + `</svg>`,
		`<svg width="4" height="4" xmlns:x="http://example.com"><x:g><rect width="4" height="4" fill="#ff0000"/></x:g>` + greenRect + `</svg>`,
		`<svg width="4" height="4" xmlns:x="http://example.com"><x:circle r="9"/><x:rect width="4" height="4" fill="#ff0000"/>` + greenRect + `</svg>`,
	}
	for _, src := range cases {
		_, cv := painted(t, src, 4, 4)
		if got := cv.At(2, 2); got.G() != 0xff || got.R() != 0 {
			t.Errorf("pixel = %v, want only the green rect of the two to paint", got)
		}
	}
}

func TestARootOutsideTheSVGNamespaceIsAnError(t *testing.T) {
	if _, err := Parse(`<svg xmlns="http://example.com"/>`); err == nil {
		t.Error("parse: no error, where the root stands in a foreign namespace")
	}
}

func TestAPrefixOrDefaultInTheSVGNamespaceDraws(t *testing.T) {
	cases := []string{
		`<svg width="4" height="4" xmlns:q="http://www.w3.org/2000/svg"><q:rect width="4" height="4" fill="#00ff00"/></svg>`,
		`<svg width="4" height="4" xmlns="http://www.w3.org/2000/svg">` + greenRect + `</svg>`,
		`<svg width="4" height="4"><q:rect xmlns:q="http://www.w3.org/2000/svg" width="4" height="4" fill="#00ff00"/></svg>`,
		`<svg width="4" height="4"><svg:rect xmlns:svg="http://www.w3.org/2000/svg" width="4" height="4" fill="#00ff00"/></svg>`,
	}
	for _, src := range cases {
		_, cv := painted(t, src, 4, 4)
		if got := cv.At(2, 2); got.G() != 0xff {
			t.Errorf("pixel = %v, want the green the drawing declares", got)
		}
	}
}

func TestANameWrittenInTheSpellingOfAnApiIsRefused(t *testing.T) {
	// The `{uri}name` spelling carries names between APIs, not between the
	// tags of a file: `{` is not a name XML allows, and a file writing one is
	// not the drawing it declares.
	if _, err := parseDocument(`<svg><{http://www.w3.org/2000/svg}rect/></svg>`); err == nil {
		t.Error("parse: no error, where the tag is written in the spelling of an API")
	}
}

func TestAnElementInAForeignNamespaceStillClosesRight(t *testing.T) {
	for _, bad := range []string{
		`<svg><x:g xmlns:x="http://example.com"></g></svg>`,
		`<svg><x:g xmlns:x="http://example.com"></x:g></svg2>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where the close is not the open", bad)
		}
	}
	if _, err := parseDocument(`<svg><x:g xmlns:x="http://example.com"></x:g></svg>`); err != nil {
		t.Errorf("parse: %v", err)
	}
}

func TestAPrefixNothingDeclaresIsAnError(t *testing.T) {
	for _, bad := range []string{
		`<svg><x:g/></svg>`,
		`<svg xmlns:x="http://example.com"><y:g/></svg>`,
		`<svg><rect x:href="#a"/></svg>`,
		`<svg><x:g></x:g2></svg>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where a prefix is worn that nothing declares", bad)
		}
	}
}

func TestAnAttributeIsHeldToTheNamespaceItsPrefixStandsFor(t *testing.T) {
	// An attribute written against XLink or XML is kept under the name the
	// drawing reads it by, whatever prefix a file spelled it with; an
	// attribute of any other namespace is another document's and is let go.
	e, err := parseDocument(`<svg xmlns:xl="http://www.w3.org/1999/xlink" xmlns:o="urn:other">` +
		`<image xl:href="#a" o:tweak="1" xml:space="preserve" x="1"/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	img := e.Kids[0]
	if got := img.attr("xlink:href"); got != "#a" {
		t.Errorf("xlink:href = %q, want it kept under the name the drawing reads", got)
	}
	if got := img.attr("xml:space"); got != "preserve" {
		t.Errorf("xml:space = %q, want it kept", got)
	}
	if img.hasAttr("o:tweak") {
		t.Error("an attribute of another namespace was kept")
	}
}

func TestThePrefixesXMLKeepsForItselfAreNotAFileToBind(t *testing.T) {
	for _, bad := range []string{
		`<svg xmlns:xml="urn:other"/>`,
		`<svg xmlns:xmlns="urn:other"/>`,
	} {
		if _, err := parseDocument(bad); err == nil {
			t.Errorf("parse %s: no error, where a prefix XML keeps for itself is bound elsewhere", bad)
		}
	}
	if _, err := parseDocument(`<svg xmlns:xml="http://www.w3.org/XML/1998/namespace"/>`); err != nil {
		t.Errorf("parse: %v, where the xml prefix is bound to the XML namespace", err)
	}
}

func TestAPrefixDeclaredToNoNamespaceIsAnError(t *testing.T) {
	// A prefix has to stand for something: `xmlns=""` empties the default and
	// is how a file says an element is in no namespace, but a prefix declared
	// to nothing names nothing and is not a binding a reader can take.
	if _, err := parseDocument(`<svg xmlns:x=""><x:g/></svg>`); err == nil {
		t.Error("parse: no error, where a prefix is declared to no namespace")
	}
	if _, err := parseDocument(`<svg xmlns=""/>`); err != nil {
		t.Errorf("parse: %v, where only the default is emptied", err)
	}
}

func TestAnEmptiedDefaultIsNotTheSVGNamespace(t *testing.T) {
	// `xmlns=""` undeclares the default, so the element is in no namespace at
	// all — not the SVG one — and what is written in it is not this package's
	// to paint, the same as a browser draws none of it.
	_, cv := painted(t, `<svg width="4" height="4"><g xmlns="">`+greenRect+`</g></svg>`, 4, 4)
	if got := cv.At(2, 2); got.G() != 0 {
		t.Errorf("pixel = %v, want an emptied default namespace painted nowhere", got)
	}
	if _, err := Parse(`<svg xmlns=""/>`); err == nil {
		t.Error("parse: no error, where the root empties the default namespace")
	}
}
