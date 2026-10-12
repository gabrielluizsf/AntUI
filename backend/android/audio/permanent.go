//go:build android

package audio

const (
	// Permanent is for something that plays until the user stops it — a
	// game, a music player. Whatever was playing before stops.
	Permanent Kind = 1
	// Brief is for something short that the other app should pause for and
	// resume after: a spoken direction, a video clip.
	Brief Kind = 2
	// BriefDucking is for something short that can play *over* the other
	// app, which turns itself down rather than stopping. It is what almost
	// every notification sound should ask for.
	BriefDucking Kind = 3
)
