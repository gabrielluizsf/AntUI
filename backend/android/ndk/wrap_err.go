//go:build android

package ndk

import (
	"fmt"
	"syscall"
)

// wrapErr turns the NDK's convention — zero for success, a negative errno
// for anything else — into an error that says which call failed.
func wrapErr(call string, rc int32) error {
	if rc == 0 {
		return nil
	}
	if rc < 0 {
		return fmt.Errorf("ndk: %s: %w", call, syscall.Errno(-rc))
	}
	return fmt.Errorf("ndk: %s failed (%d)", call, rc)
}
