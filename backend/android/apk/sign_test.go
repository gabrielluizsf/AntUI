package apk

import (
	"archive/zip"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// toolchain is the SDK, or a skip. A machine with no Android SDK is a real
// machine — most CI runners are one until they are told otherwise — and the
// right answer there is that these do not run, not that they fail.
func toolchain(t *testing.T) *sdk.Toolchain {
	t.Helper()
	tc, err := sdk.Find(sdk.Require{})
	if err != nil {
		t.Skipf("no Android SDK here: %v", err)
	}
	if tc.BuildTools.Dir == "" {
		t.Skip("no build-tools installed")
	}
	return tc
}

// A package this library signs has to be one Android will install, and the
// only thing that can say so is apksigner.
//
// It is worth checking here rather than only on a device because the failure
// is so far from its cause: an APK with a v1 signature and no v2 installs on
// Android 6 and is refused from Android 7 onwards, with an error about the
// package being corrupt.
func TestSignedPackageVerifies(t *testing.T) {
	tc := toolchain(t)
	if tc.JDK.Dir == "" {
		t.Skip("no JDK, and apksigner is a Java program")
	}
	dir := t.TempDir()

	// The smallest thing that is a zip. apksigner does not care what is in
	// it; what is being checked is the signing, not the contents.
	unsigned := filepath.Join(dir, "unsigned.apk")
	writeTestZip(t, unsigned, map[string]string{
		"AndroidManifest.xml": "not really a manifest",
		"classes.dex":         "not really a dex",
	})

	out := filepath.Join(dir, "signed.apk")
	k := Keystore{Path: filepath.Join(dir, "test.jks"), Alias: "test",
		StorePass: "testtest", KeyPass: "testtest"}
	if err := Create(tc, k, "CN=Test", 30); err != nil {
		t.Fatalf("could not make a key: %v", err)
	}
	if err := SignAPK(tc, k, 21, unsigned, out); err != nil {
		t.Fatalf("signing failed: %v", err)
	}

	// The minimum is given rather than read: apksigner works it out from the
	// manifest, and the manifest in this zip is a placeholder. What is being
	// checked is the signature and not the package.
	body, err := exec.Command(tc.BuildTools.Apksigner(), "verify", "-v",
		"--min-sdk-version", "21", out).CombinedOutput()
	if err != nil {
		t.Fatalf("apksigner refused it: %v\n%s", err, body)
	}
	report := string(body)
	// v1 alone has not been enough since Android 7, and v2 is what makes a
	// package installable on everything current.
	for _, want := range []string{
		"Verified using v1 scheme (JAR signing): true",
		"Verified using v2 scheme (APK Signature Scheme v2): true",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("apksigner did not say %q:\n%s", want, report)
		}
	}
	// And v4 is off on purpose: it writes a second file beside the package,
	// and leaving a stray .idsig next to whatever the caller asked for is a
	// worse surprise than a slower install.
	if _, err := os.Stat(out + ".idsig"); err == nil {
		t.Error("signing left an .idsig beside the package")
	}
}

// Resources are linked by aapt2 into a table, and the table is what the
// manifest's references point at. A package whose label is a reference to an
// entry that is not there installs and shows a blank name.
func TestLinkedResourcesHoldTheIcon(t *testing.T) {
	tc := toolchain(t)
	if tc.Platform.Dir == "" {
		t.Skip("no platform installed, and aapt2 link needs android.jar")
	}
	dir := t.TempDir()

	cfg := &Config{
		Package: "com.example.linked",
		Label:   "Linked",
		Dir:     dir,
		Icon:    source(t, dir, 256),
	}
	out := filepath.Join(dir, "base.apk")
	if err := Link(tc, cfg, out); err != nil {
		t.Fatalf("link failed: %v", err)
	}

	// aapt2 reads its own output, which is the point: this is checking that
	// what was produced is what the platform's own tool understands.
	body, err := exec.Command(tc.BuildTools.Aapt2(), "dump", "resources", out).CombinedOutput()
	if err != nil {
		t.Fatalf("aapt2 could not read what it wrote: %v\n%s", err, body)
	}
	table := string(body)
	for _, want := range []string{"mipmap/ic_launcher", "mipmap/ic_launcher_round"} {
		if !strings.Contains(table, want) {
			t.Errorf("the resource table has no %s:\n%s", want, first(table, 40))
		}
	}

	// And the manifest that came out points at them.
	badging, err := exec.Command(tc.BuildTools.Aapt2(), "dump", "badging", out).CombinedOutput()
	if err != nil {
		t.Fatalf("aapt2 dump badging: %v\n%s", err, badging)
	}
	if !strings.Contains(string(badging), "application-label:'Linked'") {
		t.Errorf("the label did not survive linking:\n%s", first(string(badging), 20))
	}
	if !strings.Contains(string(badging), "application-icon") {
		t.Errorf("the icon did not survive linking:\n%s", first(string(badging), 20))
	}
}

// writeTestZip makes a zip with the named entries in it.
func writeTestZip(t *testing.T, path string, entries map[string]string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	z := zip.NewWriter(f)
	for name, body := range entries {
		w, err := z.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
}

// first is the beginning of a long output, for an error message that should
// not be a screenful.
func first(text string, lines int) string {
	parts := strings.SplitN(text, "\n", lines+1)
	if len(parts) > lines {
		parts = parts[:lines]
	}
	return strings.Join(parts, "\n")
}
