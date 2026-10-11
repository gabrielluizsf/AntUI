//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// WM_SIZE_HINTS flag bits, from ICCCM. Only four of the nine are set here:
// the window's position is the window manager's business.
const (
	xSizeHintMinSize  = 16
	xSizeHintMaxSize  = 32
	xSizeHintAspect   = 128
	xSizeHintBaseSize = 256
	xWMNormalHints    = 40 // XA_WM_NORMAL_HINTS
	xWMSizeHints      = 41 // XA_WM_SIZE_HINTS
	xResourceManager  = 23 // XA_RESOURCE_MANAGER, where Xft.dpi lives
)
