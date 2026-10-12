//go:build android

package gallery

const (
	// Image is the user's photos, which is where a screenshot or a picture
	// the app made belongs.
	Image Kind = iota
	// Video is their films.
	Video
	// Audio is their music and recordings.
	Audio
)
