//go:build windows

package windows

// 64-bit Windows renames the window-long calls; 32-bit only has the short
// ones. Resolved once, so the rest of the file needs to know neither.
var (
	procGetWindowLongPtr = pickWindowLong("GetWindowLongPtrW", procGetWindowLongW)
	procSetWindowLongPtr = pickWindowLong("SetWindowLongPtrW", procSetWindowLongW)
)
