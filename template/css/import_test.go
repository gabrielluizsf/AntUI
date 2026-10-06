package css

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFiles lays a stylesheet tree out in one directory and returns it.
func writeFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("MkdirAll(%s): %v", path, err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatalf("WriteFile(%s): %v", path, err)
		}
	}
	return dir
}

// parseMain reads main.css out of a directory of files.
func parseMain(t *testing.T, files map[string]string) *Sheet {
	t.Helper()
	sh, err := ParseFile(filepath.Join(writeFiles(t, files), "main.css"))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	return sh
}

// hasWarn reports whether the sheet holds a warning naming what.
func hasWarn(sh *Sheet, what string) bool {
	for _, w := range sh.Warn {
		if strings.Contains(w, what) {
			return true
		}
	}
	return false
}

// TestImportReadsTheFileBesideIt checks that an @import brings the other
// sheet's rules in where the statement stands: before the rules that follow
// it, so the sheet's own declarations still win, and with nothing left
// unexplained in the warnings.
func TestImportReadsTheFileBesideIt(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css":  "@import \"theme.css\";\n.t { width: 10px; }",
		"theme.css": ".t { width: 20px; }\n.u { width: 30px; }",
	})
	if len(sh.Warn) != 0 {
		t.Fatalf("warnings = %v, want none", sh.Warn)
	}
	square := Viewport{Width: 800, Height: 800}
	if v, ok := mediaWidth(sh, square, "u"); !ok || v != 30 {
		t.Errorf(".u from the imported file = %d, %v; want 30, true", v, ok)
	}
	if v, ok := mediaWidth(sh, square, "t"); !ok || v != 10 {
		t.Errorf(".t = %d, %v; want 10, true — the import stands first, so "+
			"the rule after it still wins", v, ok)
	}
}

// TestImportResolvesBesideTheFileThatImportedIt checks that each file reads
// its own imports from its own directory, not from the sheet's root.
func TestImportResolvesBesideTheFileThatImportedIt(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css":  `@import "sub/a.css";`,
		"sub/a.css": `@import "../b.css";`,
		"b.css":     `.n { width: 40px; }`,
	})
	if len(sh.Warn) != 0 {
		t.Fatalf("warnings = %v, want none", sh.Warn)
	}
	if v, ok := mediaWidth(sh, Viewport{Width: 800, Height: 800}, "n"); !ok || v != 40 {
		t.Errorf(".n = %d, %v; want 40, true — the second import resolves "+
			"beside the file that wrote it", v, ok)
	}
}

// TestImportWithoutAFileToResolveToIsReported reads an @import in text with
// no file behind it: it says so and leaves the rest of the sheet standing.
func TestImportWithoutAFileToResolveToIsReported(t *testing.T) {
	sh, err := Parse(`@import "a.css"; button { color: red; }`)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if !hasWarn(sh, "no file to resolve it against") {
		t.Errorf("warnings = %v, want one about the URL having nowhere to go", sh.Warn)
	}
	if st := sh.Style("button", nil, StateNone, 800); !st.Has("color") {
		t.Error("the rule after the @import should still be read")
	}
}

// TestImportOfFileThatIsNotThereIsReported checks a local file that does not
// open: a warning, and a sheet that keeps its own rules.
func TestImportOfFileThatIsNotThereIsReported(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css": `@import "missing.css"; .m { width: 5px; }`,
	})
	if !hasWarn(sh, `ignoring @import of "missing.css"`) {
		t.Errorf("warnings = %v, want one naming the file", sh.Warn)
	}
	if v, ok := mediaWidth(sh, Viewport{Width: 800, Height: 800}, "m"); !ok || v != 5 {
		t.Errorf(".m = %d, %v; want 5, true", v, ok)
	}
}

// TestImportLeavesTheNetworkAlone checks that a URL this canvas does not go
// and get — another host, a data: document, a protocol-relative one — is
// reported and skipped rather than read.
func TestImportLeavesTheNetworkAlone(t *testing.T) {
	for _, ref := range []string{
		`url("https://example.com/a.css")`,
		`"data:text/css,.a{width:1px}"`,
		`url(//example.com/a.css)`,
	} {
		t.Run(ref, func(t *testing.T) {
			sh, err := Parse(`@import ` + ref + `; button { color: red; }`)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if !hasWarn(sh, "only local files are read") {
				t.Errorf("warnings = %v, want one about the URL being remote", sh.Warn)
			}
			if len(sh.Rules()) != 1 {
				t.Errorf("rules = %d, want only the one in this sheet", len(sh.Rules()))
			}
		})
	}
}

// TestImportReadsEachFileOnce checks the second import of a file already
// read: dropped with a warning, so a sheet cannot hold it twice.
func TestImportReadsEachFileOnce(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css": `@import "a.css"; @import "a.css";`,
		"a.css":    `.a { width: 20px; }`,
	})
	if !hasWarn(sh, "already read") {
		t.Errorf("warnings = %v, want one about the file being read already", sh.Warn)
	}
	if len(sh.Rules()) != 1 {
		t.Errorf("rules = %d, want 1 — the second import of the same file is dropped",
			len(sh.Rules()))
	}
}

// TestImportLoopEnds checks two files importing each other: the chain stops
// where it meets a file it has already read, and every file in it is styled.
func TestImportLoopEnds(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css": `@import "a.css"; .m { width: 10px; }`,
		"a.css":    `@import "b.css"; .a { width: 20px; }`,
		"b.css":    `@import "a.css"; .b { width: 30px; }`,
	})
	if !hasWarn(sh, "already read") {
		t.Errorf("warnings = %v, want one about the file being read already", sh.Warn)
	}
	if len(sh.Rules()) != 3 {
		t.Errorf("rules = %d, want 3 — main.css, a.css and b.css, each once",
			len(sh.Rules()))
	}
}

// TestImportOfTheSheetItselfIsStopped checks that a stylesheet importing its
// own name does not read itself: the file is on the list before it starts.
func TestImportOfTheSheetItselfIsStopped(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css": `@import "./main.css"; .m { width: 10px; }`,
	})
	if !hasWarn(sh, "already read") {
		t.Errorf("warnings = %v, want one about the file being read already", sh.Warn)
	}
	if len(sh.Rules()) != 1 {
		t.Errorf("rules = %d, want 1", len(sh.Rules()))
	}
}

// TestImportChainIsBounded checks a long chain of files, each a different
// one: the seen list cannot stop it, so the depth limit does.
func TestImportChainIsBounded(t *testing.T) {
	files := map[string]string{"main.css": `@import "f00.css";`}
	for i := 0; i < 40; i++ {
		files[fmt.Sprintf("f%02d.css", i)] = fmt.Sprintf(`@import "f%02d.css"; .f { width: 1px; }`, i+1)
	}
	sh := parseMain(t, files)
	if !hasWarn(sh, `"f16.css": more than 16 files deep`) {
		t.Errorf("warnings = %v, want one about the chain's depth", sh.Warn)
	}
	if len(sh.Rules()) != 16 {
		t.Errorf("rules = %d, want 16 — the first sixteen files read, and no more",
			len(sh.Rules()))
	}
}

// TestImportIsGatedByItsQuery checks the media query an @import names: the
// rules of the file it reads only reach the sheet where the query holds.
func TestImportIsGatedByItsQuery(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css": `@import url("a.css") screen and (min-width: 600px);`,
		"a.css":    `.g { width: 10px; }`,
	})
	if len(sh.Warn) != 0 {
		t.Fatalf("warnings = %v, want none", sh.Warn)
	}
	wide := Viewport{Width: 800, Height: 800}
	narrow := Viewport{Width: 400, Height: 400}
	if v, ok := mediaWidth(sh, wide, "g"); !ok || v != 10 {
		t.Errorf(".g at 800 = %d, %v; want 10, true", v, ok)
	}
	if _, ok := mediaWidth(sh, narrow, "g"); ok {
		t.Error(".g at 400 should not apply — the query the import named does not hold")
	}
}

// TestImportedQueryFoldsWithTheFile'sOwn checks that a query on the import
// and a @media inside the file it reads are both in force at once.
func TestImportedQueryFoldsWithTheFilesOwn(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css": `@import "a.css" (min-width: 600px);`,
		"a.css":    `@media (max-width: 600px) { .f { width: 10px; } }`,
	})
	if len(sh.Warn) != 0 {
		t.Fatalf("warnings = %v, want none", sh.Warn)
	}
	for _, c := range []struct {
		vp   Viewport
		want bool
	}{
		{Viewport{Width: 600, Height: 600}, true},
		{Viewport{Width: 800, Height: 800}, false},
		{Viewport{Width: 400, Height: 400}, false},
	} {
		_, ok := mediaWidth(sh, c.vp, "f")
		if ok != c.want {
			t.Errorf(".f at %dx%d applies = %v, want %v", c.vp.Width, c.vp.Height, ok, c.want)
		}
	}
}

// TestImportInsideMediaTakesTheOuterQuery checks an @import inside an
// @media: the file it reads lands inside that media too.
func TestImportInsideMediaTakesTheOuterQuery(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css": `@media (min-width: 600px) { @import "a.css"; }`,
		"a.css":    `.o { width: 10px; }`,
	})
	if len(sh.Warn) != 0 {
		t.Fatalf("warnings = %v, want none", sh.Warn)
	}
	if _, ok := mediaWidth(sh, Viewport{Width: 800, Height: 800}, "o"); !ok {
		t.Error(".o at 800 should apply — the @media around the import holds")
	}
	if _, ok := mediaWidth(sh, Viewport{Width: 400, Height: 400}, "o"); ok {
		t.Error(".o at 400 should not apply — the @media around the import does not")
	}
}

// TestImportClausesAreReported checks the layer() and supports() clauses:
// neither is folded into the cascade, both are said so, and the file the
// statement names is still read.
func TestImportClausesAreReported(t *testing.T) {
	sh := parseMain(t, map[string]string{
		"main.css": `@import "a.css" layer(base);` + "\n" +
			`@import "a.css" supports(display: grid);`,
		"a.css": `.c { width: 10px; }`,
	})
	if !hasWarn(sh, "layer()") {
		t.Errorf("warnings = %v, want one about layer()", sh.Warn)
	}
	if !hasWarn(sh, "supports()") {
		t.Errorf("warnings = %v, want one about supports()", sh.Warn)
	}
	if !hasWarn(sh, "already read") {
		t.Errorf("warnings = %v, want one for the second import of the same file", sh.Warn)
	}
	if v, ok := mediaWidth(sh, Viewport{Width: 800, Height: 800}, "c"); !ok || v != 10 {
		t.Errorf(".c = %d, %v; want 10, true — the clause is dropped, the file is not",
			v, ok)
	}
}

// TestImportStatementsWithNothingToReadAreTolerated checks the @import forms
// that name no file or never end: each is a warning, and none of them is a
// compile error.
func TestImportStatementsWithNothingToReadAreTolerated(t *testing.T) {
	for _, sheet := range []string{
		`@import;`,
		`@import "";`,
		`@import "a.css`,
		`@import url(a.css`,
		`@import screen;`,
	} {
		t.Run(sheet, func(t *testing.T) {
			sh, err := Parse(sheet + "\nbutton { color: red; }")
			if err != nil {
				t.Fatalf("Parse(%q): %v", sheet, err)
			}
			if len(sh.Warn) == 0 {
				t.Error("warnings = none, want one about the @import")
			}
		})
	}
}
