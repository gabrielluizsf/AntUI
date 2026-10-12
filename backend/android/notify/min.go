//go:build android

package notify

// The importances, as NotificationManager numbers them.
const (
	// Min shows nothing in the status bar.
	Min Importance = 1
	// Low appears silently.
	Low Importance = 2
	// Default makes a sound.
	Default Importance = 3
	// High makes a sound and slides in over whatever is on screen. It is for
	// things the user has to see now, and using it for anything else is how
	// an app gets its notifications turned off.
	High Importance = 4
)
