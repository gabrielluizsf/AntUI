package css

import (
	"os"
	"path/filepath"
	"testing"
)

// load writes css to a temp file and loads it into the table.
func load(t *testing.T, css string) *CSSClasses {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "app.css")
	if err := os.WriteFile(path, []byte(css), 0o644); err != nil {
		t.Fatal(err)
	}
	c := NewTable()
	if err := c.SetStyle(path); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSetStyleLoadsFile(t *testing.T) {
	c := load(t, `
		body { background-color: #000000; }
		button.primary { background-color: #3E63DD; color: #FFFFFF; }
		@media (min-width: 720px) {
			button { width: 90%; }
		}
	`)
	c.Button = "primary"

	st := c.GetStyle(RoleButton, nil, StateNone, 360)
	if !st.Has("background-color") {
		t.Error("expected background-color to apply from .primary")
	}
	if st.Background.R() != 0x3E || st.Background.G() != 0x63 || st.Background.B() != 0xDD {
		t.Errorf("unexpected background %v", st.Background)
	}
	if !st.Has("color") || st.Color.R() != 0xFF || st.Color.G() != 0xFF || st.Color.B() != 0xFF {
		t.Errorf("unexpected color %v", st.Color)
	}

	// Media query is off below 720 wide.
	if st.Has("width") {
		t.Errorf("width should not apply at 360 wide, got %v", st.Width)
	}

	stWide := c.GetStyle(RoleButton, nil, StateNone, 800)
	if !stWide.Has("width") {
		t.Error("expected width to apply inside the media query")
	}
	if got := stWide.Width.Px(800); got != 720 {
		t.Errorf("expected 90%% of 800 = 720, got %d", got)
	}
}

func TestSetStyleBadFile(t *testing.T) {
	c := NewTable()
	if err := c.SetStyle(filepath.Join(t.TempDir(), "missing.css")); err == nil {
		t.Fatal("expected an error loading a missing file")
	}
}

func TestGetStyleCaches(t *testing.T) {
	c := load(t, "button { background-color: #123456; }")
	a := c.GetStyle(RoleButton, nil, StateNone, 360)
	b := c.GetStyle(RoleButton, nil, StateNone, 360)
	if !sameStyle(a, b) {
		t.Error("two reads should return equal styles")
	}
}

// sameStyle relaxes the comparison to the fields this test cares about.
func sameStyle(a, b Style) bool { return a.Background == b.Background }

func TestSetPropertyWins(t *testing.T) {
	c := load(t, "button { background-color: #111111; }")
	c.SetProperty(RoleButton, nil, "background-color", "#222222")

	st := c.GetStyle(RoleButton, nil, StateNone, 360)
	if st.Background.R() != 0x22 {
		t.Errorf("SetProperty should win over the file, got %v", st.Background)
	}

	val, ok := c.Property(RoleButton, nil, "background-color", 360)
	if !ok || val != "#222222" {
		t.Errorf("Property should answer the override, got %q ok=%v", val, ok)
	}
}

func TestDisplayNone(t *testing.T) {
	c := load(t, "button { display: none; }")
	st := c.GetStyle(RoleButton, nil, StateNone, 360)
	if st.Display != 1 {
		t.Error("expected display:none to set Display=1")
	}
}

func TestWarnings(t *testing.T) {
	c := load(t, "button { background-color: #333333; }\n@media (orientation: portrait) { button { color: red; } }")
	if len(c.Warnings()) == 0 {
		t.Error("expected a warning for the unsupported media condition")
	}
}
