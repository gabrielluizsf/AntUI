//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// session is a desktop's environment for a test: nothing is in it unless
// the test put it there.
func session(vars map[string]string) func(string) string {
	return func(key string) string { return vars[key] }
}

// freshTheme puts the package's reading back to a desktop that has never
// been asked, so a test sees what it set up rather than what the machine
// running it said.
func freshTheme(t *testing.T) {
	t.Helper()
	themeMu.Lock()
	before := theme
	theme = themeState{}
	themeMu.Unlock()
	t.Cleanup(func() {
		themeMu.Lock()
		theme = before
		themeMu.Unlock()
	})
}

// standInForSettingsHelper puts a path (or no path) where the settings tool
// is looked up, so the questions about asking it do not depend on the tool
// being installed on the machine the tests run on.
func standInForSettingsHelper(t *testing.T, path string) {
	t.Helper()
	before := settingsHelper
	settingsHelper = func() string { return path }
	t.Cleanup(func() { settingsHelper = before })
}

// TestThemeInNameIsWhatTheDesktopCalledItsTheme covers the reading the whole
// chain leans on: a desktop says dark by naming a theme that is, and every
// other name it gives is a light one — which is what the dark marker exists
// to turn on. No name is no answer at all, since "light" would be a guess
// about a desktop that had not spoken.
func TestThemeInNameIsWhatTheDesktopCalledItsTheme(t *testing.T) {
	for _, want := range []struct {
		value string
		dark  bool
		ok    bool
	}{
		{"", false, false},
		{"   ", false, false},
		{"Adwaita", false, true},
		{"Adwaita-dark", true, true},
		{"ADWAITA-DARK", true, true},
		{"Adwaita:dark", true, true},
		{"Arc-Dark", true, true},
		{"Breeze", false, true},
		{"Mint-Y-Aqua", false, true},
	} {
		dark, ok := themeInName(want.value)
		if dark != want.dark || ok != want.ok {
			t.Errorf("themeInName(%q) = %v, %v; want %v, %v",
				want.value, dark, ok, want.dark, want.ok)
		}
	}
}

// TestThemeInColorFGBGIsTheFaintestHint checks the one source that describes
// a window rather than the desktop: a terminal's background as a palette
// index. Black is dark, white is light, and the indices in between are a
// pair someone chose to read against.
func TestThemeInColorFGBGIsTheFaintestHint(t *testing.T) {
	for _, want := range []struct {
		value string
		dark  bool
		ok    bool
	}{
		{"15;0", true, true},  // white on black
		{"2;1", true, true},   // green on maroon
		{"0;15", false, true}, // black on white
		{"0;7", false, true},  // black on light grey
		{"0;9", false, false}, // a colour in between, which is a choice
		{"0;8", false, false},
		{"15;13", false, false},
		{"0", false, false}, // no background at all
		{"", false, false},
		{"0;white", false, false},
	} {
		dark, ok := themeInColorFGBG(want.value)
		if dark != want.dark || ok != want.ok {
			t.Errorf("themeInColorFGBG(%q) = %v, %v; want %v, %v",
				want.value, dark, ok, want.dark, want.ok)
		}
	}
}

// TestThemeInConfigReadsTheDesktopsFiles walks the files a desktop writes at
// login: the scheme KDE names, and the name or flag GTK keeps. A file that
// is not there is skipped, and a desktop with no file has said nothing.
func TestThemeInConfigReadsTheDesktopsFiles(t *testing.T) {
	for _, want := range []struct {
		name  string
		files map[string]string
		dark  bool
		ok    bool
	}{
		{
			name:  "KDE's dark scheme",
			files: map[string]string{"kdeglobals": "[General]\nColorScheme=BreezeDark\n"},
			dark:  true, ok: true,
		},
		{
			name:  "KDE's light scheme",
			files: map[string]string{"kdeglobals": "[General]\nColorScheme=Breeze\n"},
			dark:  false, ok: true,
		},
		{
			name:  "KDE's global theme, which writes the other key",
			files: map[string]string{"kdeglobals": "[KDE]\nWidgetStyle=qt6ct\n\n[General]\nLookAndFeelPackage=Sweet-Dark\n"},
			dark:  true, ok: true,
		},
		{
			name:  "a GTK session asked for dark whatever it is called",
			files: map[string]string{"gtk-3.0/settings.ini": "[Settings]\ngtk-application-prefer-dark-theme=1\ngtk-theme-name=Adwaita\n"},
			dark:  true, ok: true,
		},
		{
			name:  "a dark theme named, with the flag left off",
			files: map[string]string{"gtk-3.0/settings.ini": "[Settings]\ngtk-application-prefer-dark-theme=0\ngtk-theme-name=Arc-Dark\n"},
			dark:  true, ok: true,
		},
		{
			name:  "a light theme named",
			files: map[string]string{"gtk-3.0/settings.ini": "[Settings]\ngtk-theme-name=Adwaita\n"},
			dark:  false, ok: true,
		},
		{
			name:  "only the flag, left off, which is an explicit not-dark",
			files: map[string]string{"gtk-3.0/settings.ini": "[Settings]\ngtk-application-prefer-dark-theme=0\n"},
			dark:  false, ok: true,
		},
		{
			name:  "GTK 4 keeps its settings in its own file",
			files: map[string]string{"gtk-4.0/settings.ini": "[Settings]\ngtk-theme-name=Yaru-dark\n"},
			dark:  true, ok: true,
		},
		{
			name:  "KDE says light and GTK says dark: the file read first wins",
			files: map[string]string{"kdeglobals": "[General]\nColorScheme=Breeze\n", "gtk-3.0/settings.ini": "[Settings]\ngtk-theme-name=Arc-Dark\n"},
			dark:  false, ok: true,
		},
		{
			name:  "a file with no theme in it",
			files: map[string]string{"gtk-3.0/settings.ini": "[Settings]\ngtk-font-name=Sans 12\n"},
			dark:  false, ok: false,
		},
		{name: "a desktop that writes no file we read"},
	} {
		dir := t.TempDir()
		for name, body := range want.files {
			path := filepath.Join(dir, name)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatalf("%s: %v", want.name, err)
			}
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				t.Fatalf("%s: %v", want.name, err)
			}
		}
		dark, ok := themeInConfig(dir)
		if dark != want.dark || ok != want.ok {
			t.Errorf("%s: themeInConfig = %v, %v; want %v, %v", want.name, dark, ok, want.dark, want.ok)
		}
	}

	// No directory at all is a session with nowhere to look, not a light
	// desktop.
	if dark, ok := themeInConfig(""); dark || ok {
		t.Errorf("themeInConfig(\"\") = %v, %v; want false, false", dark, ok)
	}
}

// TestThemeFromAsksTheBestSourceFirst is the precedence the answer rests on:
// the environment says it for this one program, the desktop's settings for
// everything the user runs, the files for a desktop with no tool behind it,
// and a terminal's colours last of all. First word wins, and no word at all
// is (false, false) — never light.
func TestThemeFromAsksTheBestSourceFirst(t *testing.T) {
	darkDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(darkDir, "kdeglobals"),
		[]byte("[General]\nColorScheme=BreezeDark\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lightDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(lightDir, "kdeglobals"),
		[]byte("[General]\nColorScheme=Breeze\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	emptyDir := t.TempDir()

	for _, want := range []struct {
		name    string
		env     map[string]string
		dir     string
		setting darkAnswer
		dark    bool
		ok      bool
	}{
		{
			name:    "the environment over the desktop's settings",
			env:     map[string]string{"GTK_THEME": "Adwaita:dark"},
			dir:     lightDir,
			setting: darkAnswer{ok: true},
			dark:    true, ok: true,
		},
		{
			name:    "the desktop's settings over the files",
			env:     map[string]string{},
			dir:     lightDir,
			setting: darkAnswer{dark: true, ok: true},
			dark:    true, ok: true,
		},
		{
			name: "the files over a terminal's colours",
			env:  map[string]string{"COLORFGBG": "0;15"},
			dir:  darkDir,
			dark: true, ok: true,
		},
		{
			name: "a terminal's colours when nothing else spoke",
			env:  map[string]string{"COLORFGBG": "0;15"},
			dir:  emptyDir,
			dark: false, ok: true,
		},
		{
			name: "a setting that had nothing to say is not an answer",
			env:  map[string]string{}, dir: emptyDir,
			dark: false, ok: false,
		},
		{
			name: "nothing anywhere",
			env:  map[string]string{"GTK_THEME": "", "COLORFGBG": ""},
			dir:  emptyDir,
			dark: false, ok: false,
		},
	} {
		dark, ok := themeFrom(session(want.env), want.dir, want.setting)
		if dark != want.dark || ok != want.ok {
			t.Errorf("%s: themeFrom = %v, %v; want %v, %v", want.name, dark, ok, want.dark, want.ok)
		}
	}
}

// TestThemeInSettingsIsTheDesktopsOwnWord is the tool behind the setting:
// the scheme where there is one, the theme's name where there is not, and
// nothing at all when the tool answers neither question.
func TestThemeInSettingsIsTheDesktopsOwnWord(t *testing.T) {
	for _, want := range []struct {
		name   string
		scheme string
		theme  string
		dark   bool
		ok     bool
	}{
		{name: "dark asked for outright", scheme: "'prefer-dark'", theme: "'Adwaita'", dark: true, ok: true},
		{name: "light asked for outright", scheme: "'prefer-light'", theme: "'Adwaita-dark'", dark: false, ok: true},
		{name: "the theme's name when the scheme is at its default", scheme: "'default'", theme: "'Mint-Y-Dark'", dark: true, ok: true},
		{name: "a light theme at the default", scheme: "'default'", theme: "'Adwaita'", dark: false, ok: true},
		{name: "neither question answered", scheme: "", theme: "", dark: false, ok: false},
	} {
		ask := func(args ...string) string {
			switch strings.Join(args, " ") {
			case "get org.gnome.desktop.interface color-scheme":
				return want.scheme
			case "get org.gnome.desktop.interface gtk-theme":
				return want.theme
			}
			return ""
		}
		dark, ok := themeInSettings(ask)
		if dark != want.dark || ok != want.ok {
			t.Errorf("%s: themeInSettings = %v, %v; want %v, %v", want.name, dark, ok, want.dark, want.ok)
		}
	}

	// The tool prints its strings quoted and padded; none of that is part
	// of the answer.
	if dark, ok := themeInSettings(func(args ...string) string {
		return "  \"prefer-dark\"  "
	}); !dark || !ok {
		t.Errorf("themeInSettings on a quoted, padded answer = %v, %v; want true, true", dark, ok)
	}
}

// TestThemeNowHoldsItsReadingForASecond is the cache: the desktop is read
// once and the answer stands for the rest of the second, because reading it
// again would cost more than a theme change that arrives a second late is
// worth — and at the end of the second it is read again, which is how a
// theme the user changes reaches a frame.
func TestThemeNowHoldsItsReadingForASecond(t *testing.T) {
	freshTheme(t)
	base := time.Unix(1700000000, 0)
	dir := t.TempDir()
	dark := session(map[string]string{"GTK_THEME": "Adwaita:dark"})
	light := session(map[string]string{"GTK_THEME": "Adwaita"})

	if got, ok := themeNow(base, dark, dir, nil); !got || !ok {
		t.Fatalf("first reading = %v, %v; want true, true", got, ok)
	}
	if got, ok := themeNow(base.Add(themeTTL-time.Nanosecond), light, dir, nil); !got || !ok {
		t.Errorf("inside the second = %v, %v; want the reading to stand at true, true", got, ok)
	}
	if got, ok := themeNow(base.Add(themeTTL), light, dir, nil); got || !ok {
		t.Errorf("after the second = %v, %v; want it read again, as false, true", got, ok)
	}
}

// TestTheSettingsToolRunsOnlyWhenTheDesktopWrote is the gate on the one
// source that costs a process: it is asked when its file has changed, and
// not again until it changes — a desktop that has not written a new answer
// has not changed what it would say. Its answer then stands over the files.
func TestTheSettingsToolRunsOnlyWhenTheDesktopWrote(t *testing.T) {
	freshTheme(t)
	base := time.Unix(1700000000, 0)
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "dconf"), 0o755); err != nil {
		t.Fatal(err)
	}
	conf := filepath.Join(dir, "dconf", "user")
	write := func(body string) {
		t.Helper()
		if err := os.WriteFile(conf, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	answer := "'prefer-dark'"
	asks := 0
	ask := func(args ...string) string {
		asks++
		if strings.HasSuffix(strings.Join(args, " "), "color-scheme") {
			return answer
		}
		return "'Adwaita'"
	}
	plain := session(map[string]string{})

	write("the desktop's first answer")
	if dark, ok := themeNow(base, plain, dir, ask); !dark || !ok {
		t.Fatalf("with the tool saying dark = %v, %v; want true, true", dark, ok)
	}
	if asks == 0 {
		t.Fatal("the tool was never asked, and there was an answer to take")
	}

	// Nothing written since: the same answer stands, at no cost.
	asks = 0
	if dark, ok := themeNow(base.Add(themeTTL), plain, dir, ask); !dark || !ok {
		t.Errorf("with nothing written = %v, %v; want the answer to stand at true, true", dark, ok)
	}
	if asks != 0 {
		t.Errorf("the tool ran %d more times with its file unchanged; want none", asks)
	}

	// The desktop changes its mind, which is a write.
	answer = "'prefer-light'"
	write("the desktop's second answer, longer than the first")
	if dark, ok := themeNow(base.Add(2*themeTTL), plain, dir, ask); dark || !ok {
		t.Errorf("after the desktop wrote again = %v, %v; want false, true", dark, ok)
	}
	if asks == 0 {
		t.Error("the tool was not asked again after its file changed")
	}

	// And a session with no tool to ask has nothing it said standing over
	// the files, which here have nothing either.
	if dark, ok := themeNow(base.Add(3*themeTTL), plain, dir, nil); dark || ok {
		t.Errorf("with no tool to ask = %v, %v; want false, false", dark, ok)
	}
	if dark, ok := themeNow(base.Add(4*themeTTL), plain, dir, ask); dark || !ok {
		t.Errorf("a tool back on the desk = %v, %v; want false, true", dark, ok)
	}
}

// TestCanAskSettingsFollowsTheDesktopName is the gate on which desktops the
// tool is worth asking: it reads GNOME's settings, so it is asked by
// GNOME and the desktops that grew out of it, and not by a desktop that
// keeps its answer in a file or a session that named no desktop at all.
func TestCanAskSettingsFollowsTheDesktopName(t *testing.T) {
	standInForSettingsHelper(t, "gsettings")
	for _, want := range []struct {
		desktop string
		ask     bool
	}{
		{"GNOME", true},
		{"ubuntu:GNOME", true},
		{"X-Cinnamon", true},
		{"Budgie:GNOME", true},
		{"KDE", false},
		{"XFCE", false},
		{"", false},
	} {
		env := session(map[string]string{"XDG_CURRENT_DESKTOP": want.desktop})
		if got := canAskSettings(env); got != want.ask {
			t.Errorf("canAskSettings on %q = %v; want %v", want.desktop, got, want.ask)
		}
	}

	// A desktop that names itself but a machine with no tool to ask: there
	// is nothing to ask, whatever is running.
	standInForSettingsHelper(t, "")
	env := session(map[string]string{"XDG_CURRENT_DESKTOP": "GNOME"})
	if canAskSettings(env) {
		t.Error("canAskSettings = true with no settings tool on the machine; want false")
	}
}

// TestSystemDarkIsTheSameTwiceOver is the exported entry, which reads this
// machine rather than a fixture: only its being asked twice within one
// reading is asserted here, since what the machine says is the machine's.
func TestSystemDarkIsTheSameTwiceOver(t *testing.T) {
	freshTheme(t)
	first, firstOK := SystemDark()
	second, secondOK := SystemDark()
	if first != second || firstOK != secondOK {
		t.Errorf("SystemDark = %v, %v then %v, %v; want the same reading while it is fresh",
			first, firstOK, second, secondOK)
	}
}

// TestSystemDarkCostsNothingBetweenReadings is the frame loop's side of the
// cache: a frame inside the second gets a lock and a comparison, and
// allocates nothing — which is what makes it all right to ask on every
// frame of a window that is up.
func TestSystemDarkCostsNothingBetweenReadings(t *testing.T) {
	freshTheme(t)
	SystemDark() // a reading, so the next one is not due for a second
	if allocated := testing.AllocsPerRun(1000, func() { SystemDark() }); allocated != 0 {
		t.Errorf("SystemDark allocated %v times a call between readings; want 0", allocated)
	}
}
