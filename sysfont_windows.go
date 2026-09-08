package antui

import (
	"os"
	"path/filepath"
	"strings"
)

// defaultUIPoints is what Windows draws its own interface at, and has since
// Vista.
const defaultUIPoints = 9

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

// fontDirs is where a font file lives when the registry names only the file:
// the system's own directory, and the per-user one Windows 10 added for
// fonts installed without administrator rights.
func fontDirs() []string {
	dirs := []string{`C:\Windows\Fonts`}
	if root := os.Getenv("SystemRoot"); root != "" {
		dirs[0] = filepath.Join(root, "Fonts")
	}
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		dirs = append(dirs, filepath.Join(local, "Microsoft", "Windows", "Fonts"))
	}
	return dirs
}

// systemFontPaths is the fallback, for a machine where reg is not on the
// path — a stripped container, or a program run somewhere odd.
func systemFontPaths() []string {
	var out []string
	for _, dir := range fontDirs() {
		for _, name := range []string{"segoeui.ttf", "tahoma.ttf", "arial.ttf", "verdana.ttf"} {
			out = append(out, filepath.Join(dir, name))
		}
	}
	return out
}