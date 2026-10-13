//go:build android

package audio

import "errors"

// ErrRefused is the system saying no — a phone call in progress, or another
// app holding the sound in a way that cannot be interrupted.
var ErrRefused = errors.New("audio: the system refused audio focus")
