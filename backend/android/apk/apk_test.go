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

func TestValidPackage(t *testing.T) {
	for _, good := range []string{"com.example.app", "org.antui.hello", "a.b", "com.x2.y_3"} {
		if err := validPackage(good); err != nil {
			t.Errorf("validPackage(%q): %v", good, err)
		}
	}
	for _, bad := range []struct{ in, why string }{
		{"hello", "one part"},
		{"com..app", "empty part"},
		{"com.2fast", "starts with a digit"},
		{"com.my-app", "a hyphen"},
		// The platform refuses these, not us: the id becomes a package name
		// in generated Java, and "class" is not a package name.
		{"com.example.class", "a Java keyword"},
		{"com.new.app", "a Java keyword"},
	} {
		if err := validPackage(bad.in); err == nil {
			t.Errorf("validPackage(%q) was accepted, and should fail on %s", bad.in, bad.why)
		}
	}
}

func TestConfigDefaults(t *testing.T) {
	dir := t.TempDir()
	c := Config{Package: "com.example.app", Dir: dir}
	if err := c.check(); err != nil {
		t.Fatal(err)
	}
	if c.Label != filepath.Base(dir) {
		t.Errorf("Label = %q, want the directory's name", c.Label)
	}
	if c.Lib != defaultLib {
		t.Errorf("Lib = %q, want %q", c.Lib, defaultLib)
	}
	if c.MinSDK != sdk.MinSDK || c.TargetSDK != sdk.PlayTarget {
		t.Errorf("SDK levels = %d..%d, want %d..%d",
			c.MinSDK, c.TargetSDK, sdk.MinSDK, sdk.PlayTarget)
	}
	if len(c.ABIs) != len(sdk.Ship) {
		t.Errorf("ABIs = %v, want %v", c.ABIs, sdk.Ship)
	}
	if c.Strip == nil || !*c.Strip {
		t.Error("Strip should default on")
	}
	if !c.Keystore.Debug() {
		t.Error("the default keystore should be the debug one")
	}

	// A minimum above the target is a build that cannot mean anything.
	bad := Config{Package: "com.example.app", Dir: dir, MinSDK: 40, TargetSDK: 30}
	if err := bad.check(); err == nil {
		t.Error("minSdk above targetSdk was accepted")
	}
}

// The manifest has three attributes that are the difference between an app
// that starts and one that dies at launch with a message about something
// else. They are checked by name because that is what they are worth.
func TestManifestLoadBearingAttributes(t *testing.T) {
	c := Config{Package: "com.example.app", Dir: t.TempDir(), Lib: "antui"}
	if err := c.check(); err != nil {
		t.Fatal(err)
	}
	m := c.Manifest()
	// hasCode and the activity's name depend on whether the shim is in, and
	// are checked by the two tests below. These four do not.
	for _, want := range []string{
		`android:extractNativeLibs="false"`,
		`android:exported="true"`,
		`android:name="android.app.lib_name" android:value="antui"`,
		`android.intent.category.LAUNCHER`,
	} {
		if !strings.Contains(m, want) {
			t.Errorf("the manifest is missing %s:\n%s", want, m)
		}
	}
	// Debuggable is off unless asked for: the store refuses an upload with it.
	if strings.Contains(m, "debuggable") {
		t.Error("debuggable appeared in a manifest that did not ask for it")
	}
}

func TestManifestEscapes(t *testing.T) {
	c := Config{
		Package:     "com.example.app",
		Dir:         t.TempDir(),
		Label:       `Bob's "Big" <Adventure> & Co`,
		Permissions: []string{"android.permission.INTERNET"},
	}
	if err := c.check(); err != nil {
		t.Fatal(err)
	}
	m := c.Manifest()
	if strings.Contains(m, `Bob's "Big"`) {
		t.Errorf("the label went in unescaped, which makes the manifest not XML:\n%s", m)
	}
	if !strings.Contains(m, "&amp;") || !strings.Contains(m, "&lt;Adventure&gt;") {
		t.Errorf("the label is not escaped as XML:\n%s", m)
	}
	if !strings.Contains(m, `<uses-permission android:name="android.permission.INTERNET" />`) {
		t.Errorf("the permission is missing:\n%s", m)
	}
}

// Assemble has one rule that is not a preference: a native library has to go
// in uncompressed, because the manifest says the loader maps it out of the
// APK rather than unpacking it, and a deflated entry cannot be mapped.
func TestAssembleStoresLibraries(t *testing.T) {
	dir := t.TempDir()

	base := filepath.Join(dir, "base.apk")
	f, err := os.Create(base)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	// A deflated manifest, as aapt2 writes it, and a deflated resources.arsc,
	// which is the case Assemble has to rewrite.
	for _, name := range []string{"AndroidManifest.xml", "resources.arsc"} {
		w, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate})
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(strings.Repeat("x", 1000)))
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	f.Close()

	lib := filepath.Join(dir, "libantui.so")
	if err := os.WriteFile(lib, []byte("not really an elf"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(dir, "out.apk")
	cfg := &Config{Package: "com.example.app", Dir: dir}
	if err := cfg.check(); err != nil {
		t.Fatal(err)
	}
	if err := Assemble(base, map[sdk.ABI]string{sdk.Arm64: lib}, "antui", out, cfg); err != nil {
		t.Fatal(err)
	}

	zr, err := zip.OpenReader(out)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	methods := map[string]uint16{}
	for _, e := range zr.File {
		methods[e.Name] = e.Method
	}
	if m, ok := methods["lib/arm64-v8a/libantui.so"]; !ok {
		t.Fatalf("the library is not in the APK: %v", methods)
	} else if m != zip.Store {
		t.Error("the library is compressed, and could not be mapped out of the APK")
	}
	if methods["resources.arsc"] != zip.Store {
		t.Error("resources.arsc is compressed, which API 30 and up refuse")
	}
	// The manifest keeps whatever aapt2 chose; it is not ours to change.
	if methods["AndroidManifest.xml"] != zip.Deflate {
		t.Error("the manifest was recompressed rather than copied across")
	}
	// The shim is in by default, and a dex is read rather than mapped, so it
	// may be compressed.
	if m, ok := methods["classes.dex"]; !ok {
		t.Errorf("the Java shim is not in the APK: %v", methods)
	} else if m != zip.Deflate {
		t.Error("classes.dex is stored, which only wastes space")
	}
}

// An app that needs nothing from Java gets no Java: no dex in the package,
// and the manifest names the platform's own activity.
func TestNoShimLeavesNoJava(t *testing.T) {
	off := false
	c := Config{Package: "com.example.app", Dir: t.TempDir(), Shim: &off}
	if err := c.check(); err != nil {
		t.Fatal(err)
	}
	m := c.Manifest()
	if !strings.Contains(m, `android:hasCode="false"`) {
		t.Errorf("hasCode is not false without a dex:\n%s", m)
	}
	if !strings.Contains(m, `android:name="android.app.NativeActivity"`) {
		t.Errorf("the manifest does not name the platform's activity:\n%s", m)
	}
}

// And an app that does gets both, which have to agree: a manifest that says
// there is no code with a dex beside it loads nothing at all.
func TestShimIsNamedAndDeclared(t *testing.T) {
	c := Config{Package: "com.example.app", Dir: t.TempDir()}
	if err := c.check(); err != nil {
		t.Fatal(err)
	}
	m := c.Manifest()
	if !strings.Contains(m, `android:hasCode="true"`) {
		t.Errorf("hasCode is not true with a dex:\n%s", m)
	}
	if !strings.Contains(m, `android:name="dev.antui.AntuiActivity"`) {
		t.Errorf("the manifest does not name the shim:\n%s", m)
	}
}
