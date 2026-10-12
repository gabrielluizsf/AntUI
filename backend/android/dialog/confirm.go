//go:build android

package dialog

// Confirm is a yes-or-no dialog, and reports whether the answer was yes.
func Confirm(title, message, yes, no string) (bool, error) {
	b, err := Ask(title, message, yes, no, "")
	return b == Positive, err
}
