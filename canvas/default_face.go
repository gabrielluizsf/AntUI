package canvas

import "sync/atomic"

// defaultFace is what Text and TextWidth use when nothing says otherwise.
var defaultFace atomic.Pointer[Face]
