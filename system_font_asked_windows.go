package antui

import (
	"path/filepath"
	"strings"
)

// systemFontAsked reads the registry's font table, which is where Windows
// records what file each installed font is in.
//
// The name is asked for and the file is looked up, rather than a path being
// assumed: a font's file is not always named after it, a user may have
// installed their own copy, and a machine with Windows somewhere other than
// C:\Windows is unusual rather than impossible.
//
// What is *not* asked here is which font the interface is set to, because
// that comes from SystemParametersInfo and not from the registry, and
// calling it needs a syscall layer this library does not have. Segoe UI is
// what every Windows since Vista uses and what changing the setting changes
// away from, so it is asked for by name.
func systemFontAsked() []string {
	const key = `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`
	var out []string
	for _, name := range []string{
		"Segoe UI (TrueType)",
		"Tahoma (TrueType)",
		"Arial (TrueType)",
	} {
		answer := ask("reg", "query", key, "/v", name)
		if answer == "" {
			continue
		}
		// reg prints "    Segoe UI (TrueType)    REG_SZ    segoeui.ttf",
		// and the file is the last field. A bare name means the font is in
		// the system's own directory.
		fields := strings.Fields(answer)
		if len(fields) == 0 {
			continue
		}
		file := fields[len(fields)-1]
		if !strings.HasSuffix(strings.ToLower(file), ".ttf") &&
			!strings.HasSuffix(strings.ToLower(file), ".ttc") {
			continue
		}
		if filepath.IsAbs(file) {
			out = append(out, file)
			continue
		}
		for _, dir := range fontDirs() {
			out = append(out, filepath.Join(dir, file))
		}
	}
	return out
}
