package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// fontRule writes a stylesheet declaring the fixture Roboto family beside the
// rule that asks for it, and hands back the path a table reads it from. When
// faces is false only the regular file is declared, so a weight or a slant it
// does not carry has to be faked on top of it.
func fontRule(t *testing.T, faces bool, rules string) string {
	t.Helper()
	dir, err := filepath.Abs(filepath.Join("..", "testdata", "fonts"))
	if err != nil {
		t.Fatalf("Abs: %v", err)
	}
	slash := filepath.ToSlash(dir)
	body := `@font-face { font-family: "Roboto"; src: url("` + slash + `/Roboto-Regular.ttf"); }` + "\n"
	if faces {
		body += `@font-face { font-family: "Roboto"; src: url("` + slash + `/Roboto-Bold.ttf"); font-weight: 700; }` + "\n"
		body += `@font-face { font-family: "Roboto"; src: url("` + slash + `/Roboto-Italic.ttf"); font-style: italic; }` + "\n"
	}
	body += rules
	path := filepath.Join(t.TempDir(), "app.css")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write sheet: %v", err)
	}
	return path
}

// TestCSSFontFamilyDrawsInTheFileItNames checks what a font-family reaches:
// the @font-face file the family declares is the face the run draws in, a
// line is as tall as that face, and only the weight or slant the file does not
// carry is faked on top of it.
func TestCSSFontFamilyDrawsInTheFileItNames(t *testing.T) {
	pin := canvas.DefaultFace()
	canvas.SetDefaultFace(nil) // the built-in, so the file's own height is the one measured
	defer canvas.SetDefaultFace(pin)

	for _, c := range []struct {
		name   string
		faces  bool
		rules  string
		file   string
		bold   bool
		italic bool
	}{
		{"a family holding a face for every weight and slant", true,
			`button { font-family: "Roboto"; }`, "Roboto-Regular.ttf", false, false},
		{"a weight the file already carries is not faked", true,
			`button { font-family: Roboto; font-weight: 700; }`, "Roboto-Bold.ttf", false, false},
		{"a slant the file already leans is not faked", true,
			`button { font-family: Roboto; font-style: italic; }`, "Roboto-Italic.ttf", false, false},
		{"a weight no file of that size has is faked", false,
			`button { font-family: Roboto; font-weight: 700; }`, "Roboto-Regular.ttf", true, false},
		{"a slant no file leans is faked", false,
			`button { font-family: Roboto; font-style: italic; }`, "Roboto-Regular.ttf", false, true},
		{"no family named keeps the face the canvas had", true,
			`button { font-weight: 700; }`, "", true, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			win := blankWin(t, 240, 60)
			cs := newCSSStyle(win)
			if err := cs.classes.SetStyle(fontRule(t, c.faces, c.rules)); err != nil {
				t.Fatalf("SetStyle: %v", err)
			}
			// The face the match should have taken, for saying the run
			// draws in that file rather than in another.
			var matched *canvas.Face
			if c.file != "" {
				f, err := canvas.LoadFace(
					filepath.Join("..", "testdata", "fonts", c.file),
					float64(css.DefaultFontSize))
				if err != nil {
					t.Fatalf("LoadFace: %v", err)
				}
				matched = f
			}
			st := cs.baseStyle(css.RoleButton, State{})
			o := cs.textOpts(st)
			if got, want := o.Face != nil, matched != nil; got != want {
				t.Fatalf("a face was drawn with: %v, want one: %v", got, want)
			}
			if o.Bold != c.bold || o.Italic != c.italic {
				t.Errorf("bold = %v, italic = %v; want %v, %v", o.Bold, o.Italic, c.bold, c.italic)
			}
			if matched == nil {
				return
			}
			if got, want := o.Face.Width("A"), matched.Width("A"); got != want {
				t.Errorf("width of A = %d, want the matched file's %d", got, want)
			}
			if got, line := cs.textHeight(st), o.Face.Height(); got != line {
				t.Errorf("line height = %d, want the face's %d", got, line)
			}
			if got := cs.textHeight(st); got == canvas.TextHeight() {
				t.Errorf("line height = the built-in font's %d, want the face's own", got)
			}
		})
	}
}
