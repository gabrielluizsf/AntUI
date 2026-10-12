//go:build android

package media

import (
	"sync"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// Player plays a file: the picture into a [canvas.Canvas], the sound through
// whatever the caller hands the samples to.
//
// It does not draw and it does not play. It decodes, keeps time, and hands
// over the frame whose moment has come — everything else is the app's, which
// is what lets a video be a texture on something, or half the screen, or
// paused while the rest of the game carries on.
type Player struct {
	info Info

	video *ndk.Decoder
	audio *ndk.Decoder

	mu sync.Mutex
	// cur is the frame being shown and next the one after it, decoded ahead
	// so that the moment it is due there is no wait. Two is enough: the
	// point is to be one ahead, not to buffer.
	cur, next       *canvas.Canvas
	curPts, nextPts time.Duration
	haveNext        bool

	// The clock. When there is sound it is the sound that keeps time —
	// counting frames actually handed to the device — because an ear notices
	// a hundredth of a second and an eye does not. With no sound it is the
	// wall clock.
	played  int64
	started time.Time
	base    time.Duration
	playing bool

	pcm     []float32
	pcmHead int
	pcmTail int

	ended  bool
	closed bool
}
