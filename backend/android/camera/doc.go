// Package camera opens the device's cameras and hands each frame over as a
// [canvas.Canvas].
//
// It goes through the NDK's Camera2 API, so nothing crosses into Java, and
// the conversion from the camera's YUV into the canvas's pixels happens in C
// — two million pixels thirty times a second is not work to do in a loop
// that the collector may pause.
//
// The frames arrive turned. A camera sensor is mounted on its side on almost
// every phone, and nothing rotates its image; [Options.Upright] does, during
// the conversion, which costs nothing.
package camera
