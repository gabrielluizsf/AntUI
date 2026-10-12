//go:build android

package notify

// Channel is a category of notification the user can silence on its own.
//
// Every notification from API 26 onwards belongs to one, and posting to a
// channel that was never created shows nothing at all — with no error, which
// is the single most common way notifications fail to appear.
type Channel struct {
	// ID is what notifications refer to. It never changes and never appears
	// on screen.
	ID string
	// Name is what the user sees in the settings screen, so it should read
	// like a category: "Messages", "Downloads".
	Name string
	// Description is the line under the name, and may be empty.
	Description string
	Importance  Importance
}
