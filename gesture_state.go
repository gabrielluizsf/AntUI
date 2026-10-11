package antui

// gestureState is what has to be remembered between frames to recognise any
// of it.
type gestureState struct {
	// The finger being watched for a tap, a long press and a fling.
	tracking  bool
	id        int
	startX    int
	startY    int
	startTime float64
	lastX     int
	lastY     int
	lastTime  float64
	vx, vy    float64
	moved     bool
	longFired bool

	// The tap before this one, for counting doubles.
	tapTime  float64
	tapX     int
	tapY     int
	tapCount int

	// The two fingers being watched for a pinch.
	twoFingers bool
	prevDist   float64
	prevAngle  float64

	out Gestures
}
