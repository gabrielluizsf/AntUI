//go:build android

package audio

const (
	// Gained means the app has it, or has it back.
	Gained Change = 1
	// Lost means something else took it for good. Stop, and do not expect
	// it back.
	Lost Change = -1
	// LostBriefly means something else took it for a moment — a call, a
	// spoken direction. Pause, and wait for Gained.
	LostBriefly Change = -2
	// Duck means something else is playing over the app. Turn down, do not
	// stop; the other sound is short and stopping would be worse.
	Duck Change = -3
)
