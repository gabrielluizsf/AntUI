//go:build android

package antui

// setExternal cannot hand the surface over: present locks the ANativeWindow
// for a CPU blit, which an EGL or Vulkan surface on the same window refuses.
func (b *androidWindow) setExternal(bool) bool { return false }
