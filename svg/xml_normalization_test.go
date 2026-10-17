package svg

import "testing"

func TestALineBreakIsReadAsOneWhateverWroteIt(t *testing.T) {
	// The same drawing saved by three editors holds three ways of ending a
	// line, and XML reads them all as one line feed before anything else.
	for _, src := range []string{
		"<g xml:space=\"preserve\">a\nb</g>",
		"<g xml:space=\"preserve\">a\r\nb</g>",
		"<g xml:space=\"preserve\">a\rb</g>",
	} {
		e, err := parseDocument(src)
		if err != nil {
			t.Fatalf("parse %q: %v", src, err)
		}
		if got, want := e.Text, "a\nb"; got != want {
			t.Errorf("%q gave text %q, want %q", src, got, want)
		}
	}
}

func TestWhitespaceInsideAnAttributeValueIsOneSpace(t *testing.T) {
	// A tab or a line break written literally in a value is normalized to a
	// space; a character reference holds the character itself and is not.
	for _, tc := range []struct{ src, want string }{
		{"<rect x=\"1\t2\"/>", "1 2"},
		{"<rect x=\"1\n2\"/>", "1 2"},
		{"<rect x=\"1\r\n2\"/>", "1 2"},
	} {
		e, err := parseDocument(tc.src)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.src, err)
		}
		if got := e.attr("x"); got != tc.want {
			t.Errorf("%q gave x = %q, want %q", tc.src, got, tc.want)
		}
	}
	e, err := parseDocument(`<rect x="1&#10;2&#9;3"/>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got, want := e.attr("x"), "1\n2\t3"; got != want {
		t.Errorf("a character reference gave x = %q, want %q", got, want)
	}
}
