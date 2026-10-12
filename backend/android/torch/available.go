//go:build android

package torch

// Available reports whether this device has a lamp to light.
func Available() bool {
	_, err := cameraWithFlash()
	return err == nil
}
