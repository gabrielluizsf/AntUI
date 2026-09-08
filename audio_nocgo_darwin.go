//go:build darwin && !cgo

package antui

import "fmt"

// CoreAudio is a C framework, so the macOS backend needs cgo — the same as
// the window backend beside it. Without it there is no sound, and saying so
// plainly beats an undefined symbol at link time.
func openAudio(name string, a *Audio) (driver, error) {
	return nil, fmt.Errorf("%w: macOS needs cgo for CoreAudio", ErrNoAudio)
}