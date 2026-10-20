package css

import "testing"

// loadSheet parses a sheet and loads it into a table the way SetStyle would,
// so a test can apply to it without a file on disk.
func loadSheet(t *testing.T, sheet string) *CSSClasses {
	t.Helper()
	c := NewTable()
	c.sheet = parseOK(t, sheet)
	return c
}

// TestHydraRepaintsTheLoadedSheet checks the skin reaches the widgets: the
// accent and the near-black glass land on the roles the sheet styles, and the
// file's own layout is left alone.
func TestHydraRepaintsTheLoadedSheet(t *testing.T) {
	c := loadSheet(t, `button { padding: 9px 20px; font-size: 14px; }`)

	if err := c.Hydra(); err != nil {
		t.Fatalf("Hydra: %v", err)
	}
	vp := Viewport{Width: 520, Height: 620}

	btn := c.GetStyle(StyleQuery{Tag: RoleButton, State: StateNone, Viewport: vp})
	if got := btn.Background; got.R() != 0x57 || got.G() != 0xC1 || got.B() != 0xFF {
		t.Errorf("button background = %v, want the accent", got)
	}
	if in := c.GetStyle(StyleQuery{Tag: RoleInput, State: StateNone, Viewport: vp}); in.Background.R() != 0x02 {
		t.Errorf("input background = %v, want the field glass", in.Background)
	}
	if !btn.Has("padding") || !btn.Has("font-size") {
		t.Error("the sheet's padding and font-size were lost")
	}
}

// TestHydraAppliesToEveryTable checks that skinning one table does not stop
// the next one.
func TestHydraAppliesToEveryTable(t *testing.T) {
	loadSheet(t, `button { color: #111111; }`).Hydra()
	other := loadSheet(t, `button { color: #111111; }`)

	if err := other.Hydra(); err != nil {
		t.Fatalf("second table Hydra: %v", err)
	}
	st := other.GetStyle(StyleQuery{Tag: RoleButton, State: StateNone, Viewport: Viewport{Width: 360, Height: 360}})
	if got := st.Background.R(); got != 0x57 {
		t.Errorf("second table's button background = %#02x, want the accent", got)
	}
}
