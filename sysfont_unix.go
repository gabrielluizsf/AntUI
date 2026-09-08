//go:build linux || freebsd || netbsd || openbsd || dragonfly || darwin

package antui

// systemFontAsked asks fontconfig, which is the resolver every toolkit on
// the machine uses. Its answer is the file the desktop is actually written
// in, whatever the distribution chose and whatever the user changed it to.
//
// Two patterns, in order. "system-ui" is what a desktop's own interface font
// is called, and fontconfig maps it to whatever that desktop set; not every
// configuration defines it, and "sans-serif" is the one that always
// resolves.
//
// fc-match always answers something — that is its job — so a machine with no
// fonts at all still gets a path, and the file simply fails to open.
func systemFontAsked() []string {
	var out []string
	for _, pattern := range []string{"system-ui", "sans-serif"} {
		if path := ask("fc-match", "--format=%{file}", pattern); path != "" {
			out = append(out, path)
		}
	}
	return out
}