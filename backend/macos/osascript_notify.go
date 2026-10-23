//go:build darwin

package macos

import "os/exec"

// osascriptNotify shows a notification through AppleScript's display
// notification, which needs no bundle and no cgo. It does not wait for
// osascript to finish — the notification is up as soon as it starts, and a
// frame is no place to wait out an interpreter — but it does reap it.
func osascriptNotify(title, body string) error {
	cmd := exec.Command("osascript", osascriptArgs(title, body)...)
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}
