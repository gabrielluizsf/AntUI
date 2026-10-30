//go:build android

package antui

// nativeHandle has nothing to give yet: the ANativeWindow comes and goes with
// the activity, and a handle that dies on the next WindowDown is not one a
// renderer can be handed once.
func (b *androidWindow) nativeHandle() (NativeHandle, bool) { return NativeHandle{}, false }
