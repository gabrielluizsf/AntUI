package antui

import "math"

// The thresholds a gesture is recognised against, in points rather than
// pixels so that they mean the same thing on a dense screen as on a coarse
// one. They are the platform's own numbers: a gesture that behaves
// differently from every other app on the device is a gesture the user gets
// wrong, however well it is implemented.
const (
	// TapSlop is how far a finger may wander and still be a tap rather than
	// a drag.
	TapSlop = 8.0
	// LongPressTime is how long a finger must stay put to be a long press.
	LongPressTime = 0.5
	// DoubleTapTime and DoubleTapSlop bound the gap between two taps for
	// them to count as one double tap.
	DoubleTapTime = 0.3
	DoubleTapSlop = 32.0
	// FlingSpeed is the slowest lift that still counts as a throw, in points
	// a second.
	FlingSpeed = 50.0
)

// Gestures is what the fingers did this frame. Every field is about this
// frame alone: a Tap is true on the frame the finger lifted and false on the
// next one, the same way MousePressed is.
//
// It is touch only. A mouse has no gestures — it has buttons and a wheel,
// which are reported as themselves — so on a desktop this is always empty.
type Gestures struct {
	// X and Y are where the gesture is: the finger, or the point between two
	// of them for a pinch.
	X, Y int

	// Tap is a finger that landed and lifted without wandering. Count is 1
	// for a single tap, 2 for a double, and keeps going for as long as the
	// taps keep coming quickly enough — so a triple tap is Count 3, not
	// three separate taps.
	Tap   bool
	Count int

	// LongPress fires once, while the finger is still down, as soon as it
	// has been still for long enough. It does not wait for the lift, because
	// by then the user has been waiting for half a second with no answer.
	LongPress bool

	// Drag is a finger that has moved past the slop. DX and DY are this
	// frame's movement.
	Drag   bool
	DX, DY int

	// Fling is a drag that was let go while still moving. VX and VY are how
	// fast, in pixels a second.
	Fling  bool
	VX, VY float64

	// Pinch and Rotate need two fingers. Scale is the change this frame —
	// above 1 for fingers moving apart — and Angle is the change in radians,
	// positive clockwise.
	Pinch  bool
	Scale  float64
	Rotate bool
	Angle  float64
}

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

// Gestures is what the fingers did this frame.
func (win *Window) Gestures() Gestures { return win.gesture.out }

// slop is a threshold in points converted to pixels. A window whose backend
// does not know the display scale gets the number unchanged, which is right
// for a desktop and merely small on a phone that will not say.
func (win *Window) slop(points float64) float64 {
	if win.scale > 0 {
		return points * win.scale
	}
	return points
}

// recognize turns this frame's touches into gestures. It runs from Begin,
// after the backend has pushed everything, so it sees a complete frame.
//
// The time is passed in rather than read here, because half of what a
// gesture is is a duration — a tap is a press that was short and a long
// press is one that was not — and a recogniser that reads the clock itself
// can only be tested by waiting.
func (win *Window) recognize(now float64) {
	g := &win.gesture
	g.out = Gestures{}

	win.recognizeOne(now)
	win.recognizeTwo()
}

// recognizeOne watches the first finger down for a tap, a long press, a drag
// and a fling.
func (win *Window) recognizeOne(now float64) {
	g := &win.gesture
	touches := win.touches
	if len(touches) == 0 {
		g.tracking = false
		return
	}

	// The finger being tracked, or the first one if there is no tracking yet
	// or the tracked one has gone.
	i := -1
	if g.tracking {
		i = win.touchIndex(g.id)
	}
	if i < 0 {
		t := touches[0]
		if !t.Began {
			// A finger that was already down when tracking was lost — after
			// a second finger came and went, say. Picking it up mid-gesture
			// would report a tap the user never made.
			g.tracking = false
			return
		}
		*g = gestureState{
			tracking:  true,
			id:        t.ID,
			startX:    t.X,
			startY:    t.Y,
			startTime: now,
			lastX:     t.X,
			lastY:     t.Y,
			lastTime:  now,
			tapTime:   g.tapTime,
			tapX:      g.tapX,
			tapY:      g.tapY,
			tapCount:  g.tapCount,
		}
		i = 0
	}

	t := touches[i]
	g.out.X, g.out.Y = t.X, t.Y

	// Velocity, smoothed. One frame's movement is far too noisy to throw
	// anything with — a finger that stalls for a single frame before lifting
	// would otherwise register as a dead stop.
	if dt := now - g.lastTime; dt > 0 {
		const smoothing = 0.6
		vx := float64(t.X-g.lastX) / dt
		vy := float64(t.Y-g.lastY) / dt
		g.vx = g.vx*(1-smoothing) + vx*smoothing
		g.vy = g.vy*(1-smoothing) + vy*smoothing
		g.lastX, g.lastY, g.lastTime = t.X, t.Y, now
	}

	dx := float64(t.X - g.startX)
	dy := float64(t.Y - g.startY)
	if !g.moved && math.Hypot(dx, dy) > win.slop(TapSlop) {
		g.moved = true
	}

	// A long press fires while the finger is still down, once.
	if !g.moved && !g.longFired && now-g.startTime >= LongPressTime {
		g.longFired = true
		g.out.LongPress = true
	}

	if g.moved && (t.DX != 0 || t.DY != 0) {
		g.out.Drag = true
		g.out.DX, g.out.DY = t.DX, t.DY
	}

	if !t.Ended {
		return
	}
	g.tracking = false

	// A cancelled touch is the system taking the gesture away — a
	// notification pulled down over the app. It is not a tap and it is not a
	// throw.
	if t.Cancelled {
		return
	}

	speed := math.Hypot(g.vx, g.vy)
	switch {
	case g.moved:
		if speed >= win.slop(FlingSpeed) {
			g.out.Fling = true
			g.out.VX, g.out.VY = g.vx, g.vy
		}
	case now-g.startTime < LongPressTime:
		g.out.Tap = true
		near := math.Hypot(float64(t.X-g.tapX), float64(t.Y-g.tapY)) <= win.slop(DoubleTapSlop)
		if near && now-g.tapTime <= DoubleTapTime {
			g.tapCount++
		} else {
			g.tapCount = 1
		}
		g.out.Count = g.tapCount
		g.tapTime, g.tapX, g.tapY = now, t.X, t.Y
	}
}

// recognizeTwo watches exactly two fingers for a pinch and a twist. Three or
// more is nobody's gesture and is left alone.
func (win *Window) recognizeTwo() {
	g := &win.gesture
	live := win.touches[:0:0] // a view, not a copy of the backing array
	for _, t := range win.touches {
		if !t.Ended {
			live = append(live, t)
		}
	}
	if len(live) != 2 {
		g.twoFingers = false
		return
	}

	a, b := live[0], live[1]
	dx := float64(b.X - a.X)
	dy := float64(b.Y - a.Y)
	dist := math.Hypot(dx, dy)
	angle := math.Atan2(dy, dx)
	g.out.X, g.out.Y = (a.X+b.X)/2, (a.Y+b.Y)/2

	if !g.twoFingers {
		g.twoFingers = true
		g.prevDist, g.prevAngle = dist, angle
		return
	}
	// A pinch is meaningless when the fingers are on top of each other, and
	// dividing by that distance is worse than meaningless.
	if g.prevDist > 1 && dist > 1 && dist != g.prevDist {
		g.out.Pinch = true
		g.out.Scale = dist / g.prevDist
	}
	// Wrap the angle into (-pi, pi] so that crossing the top of the circle
	// is a small turn and not a full one backwards.
	turn := math.Mod(angle-g.prevAngle+3*math.Pi, 2*math.Pi) - math.Pi
	if turn != 0 {
		g.out.Rotate = true
		g.out.Angle = turn
	}
	g.prevDist, g.prevAngle = dist, angle
}
