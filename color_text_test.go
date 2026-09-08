package antui

import (
	"encoding/json"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

func TestColorText(t *testing.T) {
	for _, tc := range []struct {
		text  string
		want  canvas.Color
		wrote string // what it comes back out as, when different
	}{
		{"#3E63DD", canvas.RGB(0x3E, 0x63, 0xDD), ""},
		{"3E63DD", canvas.RGB(0x3E, 0x63, 0xDD), "#3E63DD"},
		{"#803E63DD", canvas.RGBA(0x3E, 0x63, 0xDD, 0x80), ""},
		{"#FF3E63DD", canvas.RGB(0x3E, 0x63, 0xDD), "#3E63DD"},
		// Six digits mean opaque, not transparent.
		{"#000000", canvas.Black, ""},
		{"#FFFFFF", canvas.White, ""},
	} {
		got, err := canvas.ParseColor(tc.text)
		if err != nil {
			t.Errorf("ParseColor(%q): %v", tc.text, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseColor(%q) = %#08x, want %#08x", tc.text, uint32(got), uint32(tc.want))
		}
		want := tc.wrote
		if want == "" {
			want = tc.text
		}
		if got.String() != want {
			t.Errorf("%q came back as %q, want %q", tc.text, got.String(), want)
		}
	}

	for _, bad := range []string{"", "#", "red", "#FFF", "#GGGGGG", "#1234567"} {
		if c, err := canvas.ParseColor(bad); err == nil {
			t.Errorf("ParseColor(%q) = %v, wanted an error", bad, c)
		}
	}
}

func TestColorJSON(t *testing.T) {
	type thing struct {
		Background canvas.Color `json:"background"`
	}
	body, err := json.Marshal(thing{canvas.RGB(0x3E, 0x63, 0xDD)})
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `{"background":"#3E63DD"}` {
		t.Fatalf("marshalled to %s", body)
	}
	var back thing
	if err := json.Unmarshal(body, &back); err != nil {
		t.Fatal(err)
	}
	if back.Background != canvas.RGB(0x3E, 0x63, 0xDD) {
		t.Fatalf("came back as %#08x", uint32(back.Background))
	}
}