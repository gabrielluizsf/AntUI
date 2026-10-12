//go:build android

package app

// The system-UI flags, as android.view.View numbers them. They are
// deprecated in favour of WindowInsetsController from API 30 and still work,
// which is why they are what is used: one path that works everywhere beats
// two that each work half the time.
const (
	uiFullscreen       = 0x00000004 // hide the status bar
	uiHideNavigation   = 0x00000002 // hide the navigation bar
	uiImmersiveSticky  = 0x00001000 // let a swipe bring them back briefly
	uiLayoutStable     = 0x00000100
	uiLayoutHideNav    = 0x00000200
	uiLayoutFullscreen = 0x00000400
)
