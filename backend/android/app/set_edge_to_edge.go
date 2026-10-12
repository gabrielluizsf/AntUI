//go:build android

package app

// SetEdgeToEdge draws the app behind the system bars without hiding them —
// the bars stay, over the app's own pixels, and [WindowInsets] is how the
// layout keeps out from under them.
func SetEdgeToEdge(on bool) error {
	flags := 0
	if on {
		flags = uiLayoutStable | uiLayoutHideNav | uiLayoutFullscreen
	}
	return setSystemUI(flags)
}
