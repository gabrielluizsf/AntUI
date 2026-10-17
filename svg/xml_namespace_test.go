package svg

import "testing"

// An element is drawn through the namespace it stands in, not through the
// letters of its local name: a prefix or a default bound anywhere but the SVG
// namespace is a drawing within the drawing, and what is written in it is not
// this package's to paint. A prefix nothing declares is read as the SVG one,
// because that is how the drawings that arrive here spell one.

const greenRect = `<rect width="4" height="4" fill="#00ff00"/>`

func TestAnElementOutsideTheSVGNamespaceIsNotPainted(t *testing.T) {
	cases := []string{
		`<svg width="4" height="4"><g xmlns="http://example.com"><rect width="4" height="4" fill="#ff0000"/></g>` + greenRect + `</svg>`,
		`<svg width="4" height="4" xmlns:x="http://example.com"><x:g><rect width="4" height="4" fill="#ff0000"/></x:g>` + greenRect + `</svg>`,
		`<svg width="4" height="4"><x:circle xmlns:x="http://example.com" r="9"/><x:rect width="4" height="4" fill="#ff0000"/>` + greenRect + `</svg>`,
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
		`<svg width="4" height="4"><q:rect width="4" height="4" fill="#00ff00"/></svg>`,
		`<svg width="4" height="4"><svg:rect width="4" height="4" fill="#00ff00"/></svg>`,
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
