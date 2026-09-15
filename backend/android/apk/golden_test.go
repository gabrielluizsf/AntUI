// The builder is host-side tooling. Its tests are platform-independent on
// Linux and macOS, but on Windows a git checkout turns the golden files
// under testdata into CRLF, so byte comparisons against them fail. These
// are Android tests; let the Android workflow run them.
//
//go:build !windows

package apk

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/canvas"
)

// -update rewrites the golden files instead of comparing against them.
var update = flag.Bool("update", false, "rewrite the golden files")

// The manifest, kept whole, so that a change to it is something a person
// reads rather than something that happens.
//
// **The golden is the manifest this library writes, and not what aapt2 makes
// of it.** A golden for aapt2's binary XML, or for apksigner's signature
// block, would be a golden for somebody else's program: it breaks on every
// build-tools upgrade, and it says nothing at all about whether the code
// here is right. What is worth pinning is the part that is ours — every
// attribute, in order, with the values a package's behaviour depends on.
//
// TestManifestLoadBearingAttributes checks the handful that matter one at a
// time and says why each one matters, which is the better test of the two.
// This one catches what that cannot: an attribute quietly added, dropped or
// reordered by a change somewhere else.
func TestManifestGolden(t *testing.T) {
	dir := t.TempDir()
	yes, no := true, false
	cfg := Config{
		Package:     "com.example.golden",
		Label:       "Golden",
		Dir:         dir,
		VersionCode: 7,
		VersionName: "2.1",
		MinSDK:      21,
		TargetSDK:   sdk.PlayTarget,
		ABIs:        sdk.Ship,
		Permissions: []string{
			"android.permission.INTERNET",
			"android.permission.CAMERA",
		},
		Features: []Feature{{Name: "android.hardware.camera", Required: false}},
		Queries:  []string{"android.intent.action.TTS_SERVICE"},
		Links: []DeepLink{
			{Scheme: "golden"},
			{Scheme: "https", Host: "example.com", Path: "/open", Verify: true},
		},
		Orientation:    "portrait",
		Icon:           source(t, dir, 64),
		IconBackground: canvas.RGB(0x3E, 0x63, 0xDD),
		Shim:           &yes,
		Backup:         &no,
	}
	if err := cfg.check(); err != nil {
		t.Fatal(err)
	}
	got := cfg.Manifest()

	golden := filepath.Join("testdata", "manifest.xml")
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("wrote %s", golden)
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("%v\n(run: go test ./antui/backend/android/apk -run Golden -update)", err)
	}
	if got != string(want) {
		t.Errorf("the manifest changed.\n\n--- what is in %s ---\n%s\n"+
			"--- what came out ---\n%s\n"+
			"If the change is meant, run:\n"+
			"    go test ./antui/backend/android/apk -run Golden -update\n"+
			"and read the diff before committing it.", golden, want, got)
	}
}

// And the same for the one file this library encodes by hand rather than
// asking a tool for.
//
// BundleConfig.pb is nineteen bytes of protocol buffer written here, field
// by field, and a wrong field number produces a perfectly well-formed
// message that bundletool reads as empty. That happened, and it cost an
// afternoon: the version went in field 1 and it is field 2.
//
// The bytes are checked against a golden rather than against bundletool
// because bundletool is not on every machine — but the golden itself was
// made by comparing with what bundletool writes, which is the only way to
// know it is right. See bundle.go.
func TestBundleConfigGolden(t *testing.T) {
	got := bundleConfig()
	golden := filepath.Join("testdata", "bundleconfig.pb")
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, got, 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("wrote %s (%d bytes)", golden, len(got))
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("%v\n(run: go test ./antui/backend/android/apk -run Golden -update)", err)
	}
	if string(got) != string(want) {
		t.Errorf("BundleConfig.pb changed\n  was % x\n  now % x", want, got)
	}
}
