//go:build android

package dialog

// Confirm is a yes-or-no dialog, and reports whether the answer was yes.
func Confirm(title, message string, labels Labels) (bool, error) {
	b, err := Ask(title, message, labels)
	return b == Positive, err
}
