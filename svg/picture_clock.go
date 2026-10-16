package svg

import "github.com/gabrielluizsf/antui"

// pictureClock is the clock an animated picture picks its frame with: the
// seconds an `<image>` of a moving GIF reads to know which frame is showing,
// so two drawings of the same GIF move together. It is a variable so a test
// can move time instead of waiting for it.
var pictureClock = antui.Now
