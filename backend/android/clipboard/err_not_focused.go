//go:build android

package clipboard

import "errors"

// ErrNotFocused is a read refused because the app is not the one the user is
// looking at.
//
// The platform does not say so: it hands back an empty clipboard, exactly as
// if there were nothing on it, and writes the refusal to the system log
// where the app cannot see it. That silence is a day of somebody's life, so
// this package looks at whether it has focus and says which of the two
// happened.
var ErrNotFocused = errors.New("clipboard: the app is not in focus, and " +
	"Android only lets the app in focus read the clipboard")
