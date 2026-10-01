package svg

import (
	"strings"
	"testing"
)

// A drawing is read from text and painted as it was written. These check the
// reading: that a tag is found, that what it was given is kept, and that text
// written between tags comes back the way it went in.

func TestParseDocumentReadsATag(t *testing.T) {
	e, err := parseDocument(`<rect x="1" y="2" width="3" height="4"/>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if e.Name != "rect" {
		t.Errorf("name = %q, want rect", e.Name)
	}
	for name, want := range map[string]string{"x": "1", "y": "2", "width": "3", "height": "4"} {
		if got := e.attr(name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
}

func TestParseDocumentReadsChildren(t *testing.T) {
	e, err := parseDocument(`<g><rect/><circle/></g>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(e.Kids) != 2 {
		t.Fatalf("got %d children, want 2", len(e.Kids))
	}
	if e.Kids[0].Name != "rect" || e.Kids[1].Name != "circle" {
		t.Errorf("children are %q and %q, want rect and circle", e.Kids[0].Name, e.Kids[1].Name)
	}
}

func TestParseDocumentNested(t *testing.T) {
	e, err := parseDocument(`<g><g><rect/></g></g>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := e.find("rect"); got == nil {
		t.Fatal("the rect is not to be found")
	}
}

func TestParseDocumentAttributeStyles(t *testing.T) {
	// An attribute may be in either sort of quote, and one written on its own
	// has its own name for a value.
	for _, src := range []string{
		`<rect fill='red'/>`,
		`<rect fill="red"/>`,
		`<rect fill=red/>`,
	} {
		e, err := parseDocument(src)
		if err != nil {
			t.Fatalf("parse %s: %v", src, err)
		}
		if got := e.attr("fill"); got != "red" {
			t.Errorf("%s gave fill = %q, want red", src, got)
		}
	}
	e, err := parseDocument(`<rect hidden/>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !e.hasAttr("hidden") {
		t.Error("an attribute written on its own was not kept")
	}
}

func TestParseDocumentTextAndWhitespace(t *testing.T) {
	// Whitespace around a child element is there to make the file readable and
	// says nothing about the drawing, so it is not kept as text.
	e, err := parseDocument("<text>\n  hello\n  <tspan/>\n</text>")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !strings.HasPrefix(e.Text, "hello") {
		t.Errorf("text = %q, want it to start at hello", e.Text)
	}
	if strings.Contains(e.Text, "\n") {
		t.Errorf("text = %q, want the line breaks taken out", e.Text)
	}
}

func TestParseDocumentCommentsAndDeclarations(t *testing.T) {
	// The things a file says about itself rather than drawing — a comment, a
	// doctype, an xml declaration — are read past, and a comment may sit
	// anywhere at all.
	for _, src := range []string{
		`<?xml version="1.0"?><svg><!-- a note --><rect/></svg>`,
		`<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "svg11.dtd"><svg><rect/></svg>`,
		`<svg><rect/><!-- between --><circle/></svg>`,
		`<svg><g><!-- inside --></g></svg>`,
	} {
		if _, err := parseDocument(src); err != nil {
			t.Errorf("parse %s: %v", src, err)
		}
	}
}

func TestParseDocumentCDATA(t *testing.T) {
	e, err := parseDocument(`<style><![CDATA[ .a { fill: red } ]]></style>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := e.Text; got != ".a { fill: red }" {
		t.Errorf("text = %q, want the CDATA as it was written", got)
	}
}

func TestParseDocumentNamespacedTags(t *testing.T) {
	// A document may spell its tags with a prefix or in a namespace written out
	// in full, and both read the same as the bare name.
	for _, src := range []string{
		`<circle/>`,
		`<svg:circle/>`,
		`<ns:circle/>`,
		`<circle xmlns="http://www.w3.org/2000/svg"/>`,
	} {
		e, err := parseDocument(src)
		if err != nil {
			t.Fatalf("parse %s: %v", src, err)
		}
		if e.Name != "circle" {
			t.Errorf("%s gave name %q, want circle", src, e.Name)
		}
	}
}

func TestParseDocumentUnescapesEntities(t *testing.T) {
	e, err := parseDocument(`<text label="a &amp; b &lt; c &#65; &#x42;"/>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got, want := e.attr("label"), "a & b < c A B"; got != want {
		t.Errorf("label = %q, want %q", got, want)
	}
}

func TestParseDocumentLeavesAnUnknownEntityAlone(t *testing.T) {
	// An entity this package does not know is left as it was written, because
	// guessing at what it meant would draw something other than what the file
	// says.
	e, err := parseDocument(`<text label="a &unknown; b"/>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got, want := e.attr("label"), "a &unknown; b"; got != want {
		t.Errorf("label = %q, want %q", got, want)
	}
}

func TestParseDocumentRejectsWhatCannotBeRead(t *testing.T) {
	// A tag that is never closed has no shape to paint, which is an error
	// rather than a warning.
	for _, src := range []string{
		``,
		`   `,
		`not a drawing`,
		`<g><rect/>`,
		`<rect`,
		`<rect x="unclosed/>`,
	} {
		if _, err := parseDocument(src); err == nil {
			t.Errorf("parse %q: it read a drawing that is not one", src)
		}
	}
}
