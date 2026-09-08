//go:build !linux && !windows && !darwin

package antui

import "fmt"

// openAudio is what OpenAudio needs on a platform this library has not been
// taught about.
//
// Nothing is written here yet, and a program on one of these platforms hears
// silence rather than failing to start: the game runs, mixed, and nothing is
// heard, which is what [ErrNoAudio] is for.
func openAudio(name string, a *Audio) (driver, error) {
	return nil, fmt.Errorf("%w: nothing is written for this platform", ErrNoAudio)
}