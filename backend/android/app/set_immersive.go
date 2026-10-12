//go:build android

package app

// SetImmersive hides the status and navigation bars, giving the app the
// whole screen. A swipe from an edge brings them back for a few seconds and
// then they go again, which is what "sticky" means and what a game wants.
func SetImmersive(on bool) error {
	flags := 0
	if on {
		flags = uiFullscreen | uiHideNavigation | uiImmersiveSticky |
			uiLayoutStable | uiLayoutHideNav | uiLayoutFullscreen
	}
	return setSystemUI(flags)
}
