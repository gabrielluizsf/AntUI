//go:build android

package dialog

// Message is a dialog with one button, for saying something rather than
// asking. It waits until the button is pressed.
func Message(title, message, dismiss string) error {
	if dismiss == "" {
		dismiss = "OK"
	}
	_, err := Ask(title, message, dismiss, "", "")
	return err
}
