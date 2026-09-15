package shim

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every class the shim has is reached only from native code, so R8 has to be
// told to keep the package whole. If a class ever moves out of it, the rules
// stop covering it and nothing else would notice.
func TestProGuardCoversEveryClass(t *testing.T) {
	rules := string(ProGuard)
	if !strings.Contains(rules, "-keep class "+Package+".** { *; }") {
		t.Fatalf("the rules do not keep %s whole", Package)
	}
	for _, name := range javaFiles(t) {
		src := read(t, name)
		if want := "package " + Package + ";"; !strings.Contains(src, want) {
			t.Errorf("%s is not in %s, so the rules do not cover it", name, Package)
		}
	}
}

// A class name in a FindClass call is a string, and a wrong one is not a
// compile error: it is a NoClassDefFoundError the first time that code path
// runs, which may be on a device, in a feature nobody tried. This checks
// every such name against the Java that is actually here.
func TestEveryClassLookedUpExists(t *testing.T) {
	have := map[string]bool{}
	for _, name := range javaFiles(t) {
		have[strings.TrimSuffix(name, ".java")] = true
	}

	// dev/antui/X, as JNI spells a class name.
	ref := regexp.MustCompile(`"` + strings.ReplaceAll(Package, ".", "/") + `/([A-Za-z0-9_$]+)"`)
	found := 0
	root := ".." // antui/backend/android, which is everything that calls into the shim
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		switch filepath.Ext(path) {
		case ".go", ".c", ".h":
		default:
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, m := range ref.FindAllStringSubmatch(string(body), -1) {
			found++
			if !have[m[1]] {
				t.Errorf("%s looks up %s/%s, and there is no such class",
					path, Package, m[1])
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if found == 0 {
		t.Fatal("no class lookups found at all; the search is broken, not the code")
	}
	t.Logf("%d lookups across %d classes", found, len(have))
}

func javaFiles(t *testing.T) []string {
	t.Helper()
	names, err := fs.Glob(Sources, "*.java")
	if err != nil || len(names) == 0 {
		t.Fatalf("no Java in the shim: %v", err)
	}
	return names
}

func read(t *testing.T, name string) string {
	t.Helper()
	body, err := fs.ReadFile(Sources, name)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}
