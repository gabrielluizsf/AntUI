package svg

import (
	"fmt"
	"strings"
	"testing"
)

// laughs is a doctype whose names each hold ten of the one before, levels
// deep: half a kilobyte that reads into ten to the levels of text.
func laughs(levels int, mark string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<?xml version="1.0"?><!DOCTYPE svg [<!ENTITY %s0 "ha">`, mark)
	for i := 1; i <= levels; i++ {
		fmt.Fprintf(&b, `<!ENTITY %s%d "`, mark, i)
		for range 10 {
			fmt.Fprintf(&b, "&%s%d;", mark, i-1)
		}
		b.WriteString(`">`)
	}
	return b.String()
}

// TestParseRefusesABillionLaughs is CVE-2003-1564 played against this reader:
// a doctype that only refuses names that come back to themselves reads nine
// levels of ten into a gigabyte, and the drawing has to be refused well before.
func TestParseRefusesABillionLaughs(t *testing.T) {
	src := laughs(9, "l") + `]><svg xmlns="http://www.w3.org/2000/svg"><text>&l9;</text></svg>`
	if _, err := Parse(src); err == nil || !strings.Contains(err.Error(), "more than") {
		t.Fatalf("a billion laughs parsed with %v, want it refused for its size", err)
	}
}

// TestParseRefusesLaughsInAnAttribute is the same bomb set off where an
// attribute's value is read rather than a text's.
func TestParseRefusesLaughsInAnAttribute(t *testing.T) {
	src := laughs(9, "l") + `]><svg xmlns="http://www.w3.org/2000/svg"><rect id="&l9;"/></svg>`
	if _, err := Parse(src); err == nil {
		t.Fatal("a billion laughs in an attribute parsed, want it refused for its size")
	}
}

// TestParseRefusesParameterLaughs is the bomb built out of the doctype's own
// %names;, which are read into the doctype rather than into the drawing.
func TestParseRefusesParameterLaughs(t *testing.T) {
	var b strings.Builder
	b.WriteString(`<!DOCTYPE svg [<!ENTITY % p0 "ha">`)
	for i := 1; i <= 9; i++ {
		fmt.Fprintf(&b, `<!ENTITY %% p%d "`, i)
		for range 10 {
			fmt.Fprintf(&b, "%%p%d;", i-1)
		}
		b.WriteString(`">`)
	}
	b.WriteString(`<!ENTITY t "%p9;">]><svg xmlns="http://www.w3.org/2000/svg"/>`)
	if _, err := Parse(b.String()); err == nil {
		t.Fatal("a billion laughs of parameter names parsed, want it refused for its size")
	}
}

// TestParseStillReadsTheEntitiesADrawingUses keeps the limit from costing what
// it protects: a doctype naming a colour and a namespace is read as before.
func TestParseStillReadsTheEntitiesADrawingUses(t *testing.T) {
	src := `<!DOCTYPE svg [<!ENTITY red "#ff0000"><!ENTITY ns "http://www.w3.org/2000/svg">]>
		<svg xmlns="&ns;" viewBox="0 0 4 4"><rect width="4" height="4" fill="&red;"/></svg>`
	img, err := Parse(src)
	if err != nil {
		t.Fatalf("a drawing leaning on its doctype was refused: %v", err)
	}
	if rect := img.Root.find("rect"); rect == nil {
		t.Fatal("the rectangle is not in the drawing")
	}
}
