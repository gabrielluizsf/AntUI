//go:build android

package app

// WindowInsets asks the system what it is covering.
//
// It runs on the UI thread, and waits — a view may not be read from anywhere
// else. Call it when something changes rather than every frame; the events
// that change it are ContentRect and ConfigChanged.
func WindowInsets() (Insets, error) {
	var out Insets
	var err error
	if e := RunOnUISync(func() { out, err = readInsets() }); e != nil {
		return out, e
	}
	return out, err
}
