//go:build android

package jni

// Ready reports whether there is a virtual machine to talk to.
func Ready() bool {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return vm != nil
}
