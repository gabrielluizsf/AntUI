package apk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/canvas"
)

func TestConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ConfigName)
	want := Config{
		Package:        "com.example.hello",
		Label:          "Hello",
		Dir:            dir,
		VersionCode:    3,
		VersionName:    "1.2",
		MinSDK:         21,
		TargetSDK:      sdk.PlayTarget,
		ABIs:           []sdk.ABI{sdk.Arm64, sdk.Arm},
		Permissions:    []string{"android.permission.INTERNET"},
		Queries:        []string{"android.intent.action.TTS_SERVICE"},
		Icon:           filepath.Join(dir, "icon.png"),
		IconBackground: canvas.RGB(0x3E, 0x63, 0xDD),
		Links:          []DeepLink{{Scheme: "https", Host: "example.com", Verify: true}},
		Keystore:       Keystore{Path: filepath.Join(dir, "release.jks"), Alias: "release"},
	}
	if err := want.Save(path); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	// The paths went out relative and came back absolute, which is the point
	// of them: the file moves and still means the same thing.
	for _, f := range []struct{ name, got, want string }{
		{"Dir", got.Dir, want.Dir},
		{"Icon", got.Icon, want.Icon},
		{"Keystore.Path", got.Keystore.Path, want.Keystore.Path},
	} {
		if f.got != f.want {
			t.Errorf("%s = %q, want %q", f.name, f.got, f.want)
		}
	}
	if got.IconBackground != want.IconBackground {
		t.Errorf("IconBackground = %v, want %v", got.IconBackground, want.IconBackground)
	}
	if len(got.ABIs) != 2 || got.ABIs[0] != sdk.Arm64 {
		t.Errorf("ABIs = %v", got.ABIs)
	}
	if len(got.Links) != 1 || !got.Links[0].Verify {
		t.Errorf("Links = %v", got.Links)
	}

	// The paths in the file are relative, or it is not a file that travels.
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), dir) {
		t.Errorf("the file holds an absolute path:\n%s", body)
	}
	// And a colour is written as a colour.
	if !strings.Contains(string(body), `"#3E63DD"`) {
		t.Errorf("the background is not written as a colour:\n%s", body)
	}
}

// A password in a file beside the source is a password in the repository.
func TestConfigNeverWritesPasswords(t *testing.T) {
	path := filepath.Join(t.TempDir(), ConfigName)
	cfg := Config{
		Package: "com.example.hello",
		Keystore: Keystore{
			Path: "release.jks", Alias: "release",
			StorePass: "hunter2", KeyPass: "hunter2",
		},
	}
	if err := cfg.Save(path); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "hunter2") {
		t.Fatalf("the config file holds a password:\n%s", body)
	}

	t.Setenv("ANTUIAPK_STOREPASS", "from the environment")
	t.Setenv("ANTUIAPK_KEYPASS", "from the environment")
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Keystore.StorePass != "from the environment" {
		t.Errorf("StorePass = %q", got.Keystore.StorePass)
	}
	if got.Keystore.KeyPass != "from the environment" {
		t.Errorf("KeyPass = %q", got.Keystore.KeyPass)
	}
}

// A misspelt key is the quietest way for a config file to be wrong: the app
// simply comes out without the permission, and nothing anywhere says so.
func TestConfigRejectsUnknownKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), ConfigName)
	if err := os.WriteFile(path, []byte(`{"permission": ["a"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("a misspelt key was accepted")
	}
	if !strings.Contains(err.Error(), "permission") {
		t.Errorf("the error does not name the key: %v", err)
	}
}
