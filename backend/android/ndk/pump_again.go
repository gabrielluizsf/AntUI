//go:build android

package ndk

// The states a pump can be in.
const (
	pumpAgain = 0
	pumpFrame = 1
	pumpEnd   = 2
	pumpError = -1
)
