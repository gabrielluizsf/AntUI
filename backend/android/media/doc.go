// Package media plays a video file into a [canvas.Canvas], and decodes its
// sound into samples.
//
// It draws nothing and plays nothing itself. It decodes, keeps time, and
// hands over the frame whose moment has come — which is what lets a video be
// a texture on something, or half the screen, or paused while the rest of a
// game carries on.
//
// It goes through the NDK's own extractor and codec, so the decoding is the
// device's hardware and nothing crosses into Java. The frames arrive through
// an image reader, which is API 24: below that, video cannot be got at from
// native code at all.
package media
