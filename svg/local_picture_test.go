package svg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFetchPictureRefusesADevice keeps a path from naming the memory of the
// machine: /dev/zero read whole never ends, and a pipe never answers.
func TestFetchPictureRefusesADevice(t *testing.T) {
	if _, err := os.Stat("/dev/zero"); err != nil {
		t.Skip("there is no /dev/zero here")
	}
	for _, target := range []string{"/dev/zero", "file:///dev/zero"} {
		if _, err := fetchPicture(target); err == nil {
			t.Errorf("%s was read, want it refused as no plain file", target)
		}
	}
}

// TestFetchPictureRefusesAFileTooBig holds a file on disk to the size a
// server's answer is held to.
func TestFetchPictureRefusesAFileTooBig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "big.png")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(maxPictureBytes + 1); err != nil {
		t.Fatal(err)
	}
	f.Close()
	if _, err := fetchPicture(path); err == nil || !strings.Contains(err.Error(), "bigger than") {
		t.Fatalf("a file past the limit was read with %v, want it refused", err)
	}
}

// TestNetworkDrawingMayNotReadThisMachine is the escape of CVE-2023-38633: a
// drawing that came from a server names a file on this machine by `file:`,
// itself or through a drawing written in place inside it, which carries the
// server in its chain. Neither is read. (A bare path in such a drawing
// resolves against the server, the way a browser reads it, so it never
// reaches the disk to begin with.)
func TestNetworkDrawingMayNotReadThisMachine(t *testing.T) {
	dir := t.TempDir()
	secret := filepath.Join(dir, "secret.png")
	if err := os.WriteFile(secret, pngClaiming(1, 1), 0o600); err != nil {
		t.Fatal(err)
	}
	href := "file://" + filepath.ToSlash(secret)
	src := `<svg xmlns="http://www.w3.org/2000/svg" width="4" height="4"><image href="` + href + `" width="4" height="4"/></svg>`
	for _, from := range []struct {
		path  string
		chain []string
	}{
		{"https://example.com/a.svg", nil},
		{"data:image/svg+xml,inner", []string{"https://example.com/a.svg"}},
	} {
		img, err := parseText(src, from.path, from.chain)
		if err != nil {
			t.Fatal(err)
		}
		if pic := img.Root.find("image"); pic != nil && pic.Pic != nil {
			t.Errorf("a drawing read through %v read %s", append(from.chain, from.path), href)
		}
		said := false
		for _, w := range img.Warnings() {
			said = said || strings.Contains(w.What, "may not read")
		}
		if !said {
			t.Errorf("reading %s through %v said %v, want it refused", href, append(from.chain, from.path), img.Warnings())
		}
	}
	if got := resolveHref("https://example.com/a.svg", secret); !strings.HasPrefix(got, "https://example.com/") {
		t.Errorf("a bare path in a drawing from a server resolved to %q, want it on the server", got)
	}
}

// TestDiskDrawingStillReadsBesideIt keeps the rule from costing what it
// protects: a drawing read off the disk still reads the picture beside it.
func TestDiskDrawingStillReadsBesideIt(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "beside.png"), pngClaiming(1, 1), 0o600); err != nil {
		t.Fatal(err)
	}
	src := `<svg xmlns="http://www.w3.org/2000/svg" width="4" height="4"><image href="beside.png" width="4" height="4"/></svg>`
	img, err := parseText(src, filepath.Join(dir, "a.svg"), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range img.Warnings() {
		if strings.Contains(w.What, "may not read") {
			t.Fatalf("a drawing on disk was refused the picture beside it: %v", w)
		}
	}
}
