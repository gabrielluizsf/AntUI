//go:build android

package location

// The two providers a device without Play Services has. There is no third:
// the "fused" provider everyone knows is part of Play Services, not of
// Android, and an app that must work without them uses these.
const (
	// GPS is the satellites: accurate to a few metres, slow to get a first
	// answer, useless indoors, and expensive in battery.
	GPS = "gps"
	// Network is the towers and the wifi around: accurate to tens or
	// hundreds of metres, answers immediately, works indoors, costs almost
	// nothing.
	Network = "network"
)
