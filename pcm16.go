//go:build windows

package antui

import (
	"math"
)

// pcm16 is a sample as the older devices want it, clipped rather than wrapped
// — a sample that overflows into the opposite sign is a click, and a click is
// the loudest thing a speaker can make.
func pcm16(sample float32) int16 {
	v := float64(sample) * 32767
	if v > 32767 {
		v = 32767
	}
	if v < -32768 {
		v = -32768
	}
	return int16(math.Round(v))
}
