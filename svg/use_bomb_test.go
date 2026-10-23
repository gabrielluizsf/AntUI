package svg

import (
	"fmt"
	"strings"
	"testing"
)

// TestParseStopsAUseBomb is librsvg's CVE-2019-20446 played through `<use>`:
// each group is ten uses of the one before, six deep, a million shapes out of
// a few hundred bytes. The drawing is read, short of the shapes past the
// limit, and says why they are missing.
func TestParseStopsAUseBomb(t *testing.T) {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10"><defs><rect id="u0" width="1" height="1"/>`)
	for i := 1; i <= 6; i++ {
		fmt.Fprintf(&b, `<g id="u%d">`, i)
		for range 10 {
			fmt.Fprintf(&b, `<use href="#u%d"/>`, i-1)
		}
		b.WriteString(`</g>`)
	}
	b.WriteString(`</defs><use href="#u6"/></svg>`)
	img, err := Parse(b.String())
	if err != nil {
		t.Fatalf("the use bomb was refused outright: %v", err)
	}
	if img.built > maxBuiltNodes {
		t.Errorf("the drawing built %d nodes, want no more than %d", img.built, maxBuiltNodes)
	}
	said := false
	for _, w := range img.Warnings() {
		said = said || strings.Contains(w.What, "left out")
	}
	if !said {
		t.Errorf("the drawing said %v, want it to say the rest was left out", img.Warnings())
	}
}
