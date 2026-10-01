package svg

import (
	"strings"
	"testing"
)

// What a drawing could not do is said out loud, and what it says has to be
// readable by whoever is going to go and fix the drawing.

func TestWarningSaysWhichElementAndWhat(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10"><image href="a.png"/><path d="M0 0 X1 1"/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ws := img.Warnings()
	if len(ws) < 2 {
		t.Fatalf("got %d warnings, want one for the image and one for the path: %s", len(ws), ws)
	}
	for _, w := range ws {
		if w.Element == "" {
			t.Errorf("a warning came with no element: %+v", w)
		}
		if w.What == "" {
			t.Errorf("a warning came with nothing said about it: %+v", w)
		}
		if !strings.Contains(w.String(), w.Element) {
			t.Errorf("%q does not name the element it came from", w.String())
		}
	}
	// A drawing that came out as it should says so rather than nothing at all.
	img, err = Parse(`<svg viewBox="0 0 10 10"><rect width="10" height="10"/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a drawing with nothing wrong said %q", got)
	}
}

func TestErrorSaysWhyADrawingCouldNotBeRead(t *testing.T) {
	// A drawing that could not be read at all is an error rather than a
	// warning, and what it says names the package so a stack of them can be read
	// from the bottom up.
	for _, src := range []string{``, `<rect/>`, `<svg><g></svg>`} {
		_, err := Parse(src)
		if err == nil {
			t.Errorf("parse %q: it read a drawing out of it", src)
			continue
		}
		if !strings.HasPrefix(err.Error(), "antui/svg: ") {
			t.Errorf("parse %q: the error is %q, want it to say where it came from", src, err)
		}
	}
	// The root of a drawing is the one tag that must be there, and a file whose
	// root is something else is not a drawing.
	_, err := Parse(`<rect width="1" height="1"/>`)
	var typed *Error
	if !asError(err, &typed) {
		t.Fatalf("the error is %v, want one of this package's own", err)
	}
	if !strings.Contains(typed.What, "<svg>") {
		t.Errorf("the error says %q, want it to say which tag was expected", typed.What)
	}
}

// asError is whether err is of this package's error type, so that a caller can
// tell a drawing it could not read from any other failure.
func asError[T error](err error, target *T) bool {
	e, ok := err.(T)
	if ok {
		*target = e
	}
	return ok
}
