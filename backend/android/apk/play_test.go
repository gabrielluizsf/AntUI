// The builder is host-side tooling. Its tests are platform-independent on
// Linux and macOS, but on Windows a git checkout turns the golden files
// under testdata into CRLF, so byte comparisons against them fail. These
// are Android tests; let the Android workflow run them.
//
//go:build !windows

package apk

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// clean is a config with nothing about it the store would object to. Each
// test below spoils exactly one thing, so that what it is testing is the
// difference and not the whole config.
func clean(t *testing.T) Config {
	t.Helper()
	dir := t.TempDir()
	icon := filepath.Join(dir, "icon.png")
	if err := os.WriteFile(icon, []byte("not really a png"), 0o644); err != nil {
		t.Fatal(err)
	}
	return Config{
		Package:   "org.example.app",
		Label:     "App",
		Dir:       dir,
		TargetSDK: sdk.PlayTarget,
		ABIs:      sdk.Ship,
		Icon:      icon,
		Keystore: Keystore{
			Path: filepath.Join(dir, "release.jks"), Alias: "release",
			StorePass: "s3cret", KeyPass: "s3cret",
		},
	}
}

func TestPlayCleanConfig(t *testing.T) {
	if ps := clean(t).PlayProblems(); len(ps) != 0 {
		t.Fatalf("a clean config has problems: %v", ps)
	}
}

func TestPlayProblems(t *testing.T) {
	for _, tc := range []struct {
		name  string
		spoil func(*Config)
		want  string
		fatal bool
	}{
		{"32-bit only", func(c *Config) { c.ABIs = []sdk.ABI{sdk.Arm} },
			"no 64-bit ABI", true},
		{"no arm64", func(c *Config) { c.ABIs = []sdk.ABI{sdk.X64} },
			"no arm64-v8a", false},
		{"old target", func(c *Config) { c.TargetSDK = sdk.PlayTarget - 1 },
			"the store requires", true},
		{"debuggable", func(c *Config) { c.Debuggable = true },
			"debuggable", true},
		{"debug key", func(c *Config) { c.Keystore = Keystore{} },
			"shared debug key", true},
		{"example id", func(c *Config) { c.Package = "com.example.app" },
			"com.example", true},
		{"no icon", func(c *Config) { c.Icon = "" },
			"no icon", false},
		{"unstripped", func(c *Config) { off := false; c.Strip = &off },
			"symbols in it", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := clean(t)
			tc.spoil(&cfg)
			ps := cfg.PlayProblems()
			found := false
			for _, p := range ps {
				if strings.Contains(p.Text, tc.want) {
					found = true
					if p.Fatal != tc.fatal {
						t.Errorf("%q: fatal = %v, want %v", p.Text, p.Fatal, tc.fatal)
					}
					if p.Fix == "" {
						t.Errorf("%q says what is wrong but not what to do", p.Text)
					}
				}
			}
			if !found {
				t.Fatalf("no problem mentioning %q; got %v", tc.want, ps)
			}
			if Fatal(ps) != tc.fatal {
				t.Errorf("Fatal = %v, want %v", Fatal(ps), tc.fatal)
			}
		})
	}
}

// A config is what was meant; the file is what gets uploaded. These are the
// things only the file can answer.
func TestPlayProblemsInPackage(t *testing.T) {
	for _, tc := range []struct {
		name    string
		ext     string
		entries map[string]uint16 // path → compression method
		want    []string
		fatal   bool
	}{
		{"nothing to run", ".apk", map[string]uint16{
			"classes.dex": zip.Deflate,
		}, []string{"no native library"}, true},
		{"32-bit only", ".apk", map[string]uint16{
			"lib/armeabi-v7a/libantui.so": zip.Store,
		}, []string{"no 64-bit library"}, true},
		{"compressed library", ".apk", map[string]uint16{
			"lib/arm64-v8a/libantui.so": zip.Deflate,
		}, []string{"stored compressed"}, true},
		{"a good package", ".apk", map[string]uint16{
			"lib/arm64-v8a/libantui.so": zip.Store,
		}, nil, false},
		// A bundle is not a package: its libraries are deflated on purpose,
		// because bundletool decides how they are stored when it splits it.
		{"bundle without symbols", ".aab", map[string]uint16{
			"base/lib/arm64-v8a/libantui.so": zip.Deflate,
		}, []string{"carries no symbols"}, false},
		{"bundle with symbols", ".aab", map[string]uint16{
			"base/lib/arm64-v8a/libantui.so":             zip.Deflate,
			bundleSymbols + "/arm64-v8a/libantui.so.sym": zip.Deflate,
		}, nil, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "app"+tc.ext)
			writeZip(t, path, tc.entries)

			ps, err := PlayProblemsIn(path)
			if err != nil {
				t.Fatal(err)
			}
			if len(ps) != len(tc.want) {
				t.Fatalf("got %v, wanted %d problem(s)", ps, len(tc.want))
			}
			for i, want := range tc.want {
				if !strings.Contains(ps[i].Text, want) {
					t.Errorf("problem %d is %q, wanted one mentioning %q",
						i, ps[i].Text, want)
				}
			}
			if Fatal(ps) != tc.fatal {
				t.Errorf("Fatal = %v, want %v", Fatal(ps), tc.fatal)
			}
		})
	}
}

func writeZip(t *testing.T, path string, entries map[string]uint16) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	z := zip.NewWriter(f)
	for name, method := range entries {
		w, err := z.CreateHeader(&zip.FileHeader{Name: name, Method: method})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte("elf, more or less")); err != nil {
			t.Fatal(err)
		}
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
}
