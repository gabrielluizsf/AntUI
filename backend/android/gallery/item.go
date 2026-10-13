//go:build android

package gallery

import "time"

// Item is one thing in the device's library.
type Item struct {
	// URI is its address, and is what [antui/backend/android/picker.Open] takes.
	URI string
	// Name is what the user sees.
	Name string
	// Size is in bytes, and Mime is what it holds.
	Size int64
	Mime string
	// Added is when it appeared in the library.
	Added time.Time
	// Duration is how long it plays, and is zero for a picture.
	Duration time.Duration
	// Album is the folder inside the library, on Android 10 and up.
	Album string
}
