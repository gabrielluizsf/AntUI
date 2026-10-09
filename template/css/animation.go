package css

// Animation is one entry of the animation shorthand: the @keyframes block it
// plays (by name), its duration, easing and delay, how many times it runs
// (Infinite) and which direction and fill-mode govern it. Iterations without
// Infinite is the written number, defaulting to 1.
type Animation struct {
	Name       string
	Duration   Time
	Timing     Timing
	Delay      Time
	Iterations float64
	Infinite   bool
	Direction  uint8 // one of the Anim* constants
	Fill       uint8 // one of the Fill* constants
}
