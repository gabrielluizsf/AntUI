package sdk

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseVersion(t *testing.T) {
	for _, c := range []struct {
		in    string
		parts []int
		pre   string
	}{
		{"37.0.0", []int{37, 0, 0}, ""},
		{"30.0.16138531", []int{30, 0, 16138531}, ""},
		{"30.0.16138531-beta3", []int{30, 0, 16138531}, "beta3"},
		{"36.1.0-rc1", []int{36, 1, 0}, "rc1"},
		{"21", []int{21}, ""},
		{"", nil, ""},
		{"VanillaIceCream", nil, ""},
	} {
		v := ParseVersion(c.in)
		if v.Pre != c.pre || len(v.Parts) != len(c.parts) {
			t.Fatalf("ParseVersion(%q) = %v %q, want %v %q", c.in, v.Parts, v.Pre, c.parts, c.pre)
		}
		for i := range c.parts {
			if v.Parts[i] != c.parts[i] {
				t.Fatalf("ParseVersion(%q) = %v, want %v", c.in, v.Parts, c.parts)
			}
		}
	}
}

func TestVersionCompare(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want int
	}{
		{"37.0.0", "36.1.0", 1},
		{"36.1.0", "36.0.0", 1},
		{"37.0.0", "37.0.0", 0},
		// A release beats the preview of the same number, which is the whole
		// reason Compare looks at Pre at all.
		{"37.0.0", "37.0.0-rc2", 1},
		{"37.0.0-rc1", "37.0.0-rc2", -1},
		// Fewer parts is smaller when the shared ones are equal.
		{"37.0", "37.0.1", -1},
		{"37.0", "37.0.0", 0},
	} {
		if got := ParseVersion(c.a).Compare(ParseVersion(c.b)); got != c.want {
			t.Errorf("%s vs %s = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

// A release is chosen over a newer preview, so that installing a beta by
// accident does not silently become what everything is built with.
func TestPickPrefersRelease(t *testing.T) {
	got, ok := pick([]string{"37.0.0-rc1", "36.1.0", "37.0.0-rc2"}, ParseVersion)
	if !ok || got != "36.1.0" {
		t.Fatalf("pick = %q %v, want 36.1.0", got, ok)
	}
	got, ok = pick([]string{"37.0.0-rc1", "36.1.0", "37.0.0"}, ParseVersion)
	if !ok || got != "37.0.0" {
		t.Fatalf("pick = %q %v, want 37.0.0", got, ok)
	}
	// With nothing but previews, the newest preview is better than nothing.
	got, ok = pick([]string{"37.0.0-rc1", "37.0.0-rc2"}, ParseVersion)
	if !ok || got != "37.0.0-rc2" {
		t.Fatalf("pick = %q %v, want 37.0.0-rc2", got, ok)
	}
	if _, ok := pick(nil, ParseVersion); ok {
		t.Error("pick of nothing reported success")
	}
}

func TestParsePlatform(t *testing.T) {
	for _, c := range []struct {
		name            string
		api, minor, ext int
		ok              bool
	}{
		{"android-30", 30, 0, 0, true},
		{"android-36.1", 36, 1, 0, true},
		{"android-37.2", 37, 2, 0, true},
		{"android-34-ext11", 34, 0, 11, true},
		{"android-VanillaIceCream", 0, 0, 0, false},
		{"sources", 0, 0, 0, false},
	} {
		p, ok := parsePlatform(c.name)
		if ok != c.ok {
			t.Fatalf("parsePlatform(%q) ok = %v, want %v", c.name, ok, c.ok)
		}
		if ok && (p.API != c.api || p.Minor != c.minor || p.Ext != c.ext) {
			t.Errorf("parsePlatform(%q) = %d.%d ext%d, want %d.%d ext%d",
				c.name, p.API, p.Minor, p.Ext, c.api, c.minor, c.ext)
		}
	}
}

func TestABI(t *testing.T) {
	for _, c := range []struct {
		in     string
		abi    ABI
		goarch string
		bits   int
	}{
		{"arm64-v8a", Arm64, "arm64", 64},
		{"arm64", Arm64, "arm64", 64},
		{"armeabi-v7a", Arm, "arm", 32},
		{"arm", Arm, "arm", 32},
		{"x86_64", X64, "amd64", 64},
		{"amd64", X64, "amd64", 64},
		{"386", X86, "386", 32},
	} {
		abi, err := ParseABI(c.in)
		if err != nil {
			t.Fatalf("ParseABI(%q): %v", c.in, err)
		}
		if abi != c.abi || abi.GOARCH() != c.goarch || abi.Bits() != c.bits {
			t.Errorf("ParseABI(%q) = %s/%s/%d, want %s/%s/%d",
				c.in, abi, abi.GOARCH(), abi.Bits(), c.abi, c.goarch, c.bits)
		}
	}
	if _, err := ParseABI("mips"); err == nil {
		t.Error("ParseABI(mips) should fail: the ABI has been gone for years")
	}
	// The 32-bit ARM triple the compiler uses is not the one binutils uses,
	// and getting it wrong finds no compiler at all.
	if got := Arm.Triple(); got != "armv7a-linux-androideabi" {
		t.Errorf("Arm.Triple() = %q", got)
	}
	if Arm.GOARM() != "7" || Arm64.GOARM() != "" {
		t.Error("GOARM is only set for 32-bit ARM")
	}
	// A bundle without a 64-bit ABI is refused by the store.
	has64 := false
	for _, a := range Ship {
		if a.Bits() == 64 {
			has64 = true
		}
	}
	if !has64 {
		t.Error("Ship carries no 64-bit ABI")
	}
}

func TestFindRootIgnoresAnEmptyDirectory(t *testing.T) {
	// Pointing ANDROID_HOME at a directory that exists but holds no SDK is a
	// common broken setup, and taking it would fail much later.
	dir := t.TempDir()
	t.Setenv("ANDROID_HOME", dir)
	t.Setenv("ANDROID_SDK_ROOT", "")
	if isSDK(dir) {
		t.Fatal("an empty directory was taken for an SDK")
	}
	if err := os.Mkdir(filepath.Join(dir, "platform-tools"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !isSDK(dir) {
		t.Fatal("a directory with platform-tools was not taken for an SDK")
	}
}

func TestMissingSaysHowToFixIt(t *testing.T) {
	m := &Missing{Root: "/sdk"}
	m.add("build-tools", "build-tools;37.0.0")
	m.add("ndk", "ndk;30.0.16138531")
	m.add("a JDK", "")
	want := `sdkmanager "build-tools;37.0.0" "ndk;30.0.16138531"`
	if got := m.Fix(); got != want {
		t.Errorf("Fix() = %q, want %q", got, want)
	}
	if m.Error() == "" {
		t.Error("Missing has no message")
	}
}

// The rest needs a real SDK, so it runs on a machine that has one and is
// skipped everywhere else.
func TestFindOnThisMachine(t *testing.T) {
	t.Setenv("ANDROID_HOME", os.Getenv("ANDROID_HOME"))
	tc, err := Find(Require{NDK: true, MinPlatform: MinSDK})
	if err != nil {
		t.Skipf("no usable SDK here: %v", err)
	}
	if !exists(tc.BuildTools.Aapt2()) {
		t.Errorf("aapt2 is not at %s", tc.BuildTools.Aapt2())
	}
	if !exists(tc.Platform.Jar()) {
		t.Errorf("android.jar is not at %s", tc.Platform.Jar())
	}
	for _, abi := range All {
		cc, err := tc.NDK.Clang(abi, MinSDK)
		if err != nil {
			t.Errorf("%s: %v", abi, err)
			continue
		}
		if !exists(cc) {
			t.Errorf("%s: no compiler at %s", abi, cc)
		}
	}
}

// TestDoctor prints the report for this machine. It asserts almost nothing —
// what it is for is `go test -v`, which is the quickest way to see what an
// install actually looks like from here.
func TestDoctor(t *testing.T) {
	d := Check(Require{NDK: true})
	t.Log("\n" + d.String())
	if len(d.Findings) == 0 {
		t.Error("the doctor found nothing at all to say")
	}
}
