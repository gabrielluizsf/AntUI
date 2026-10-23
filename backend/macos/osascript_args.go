//go:build darwin

package macos

// osascriptArgs is the command line for a notification. The title and body
// travel as the script's own arguments rather than spliced into its source,
// so no quote or backslash in them can end a string early and run as
// AppleScript.
func osascriptArgs(title, body string) []string {
	return []string{
		"-e", "on run argv",
		"-e", "display notification (item 2 of argv) with title (item 1 of argv)",
		"-e", "end run",
		title, body,
	}
}
