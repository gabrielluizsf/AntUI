//go:build android

package camera

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// Preview is a running camera, producing frames.
type Preview struct {
	info    Info
	manager *ndk.CameraManager
	device  *ndk.CameraDevice

	mu     sync.Mutex
	canvas *canvas.Canvas
	closed bool
	// asked is the size the image reader was made at, which is the size
	// every frame will be.
	askedW, askedH int

	// upright is how far a frame has to be turned to look right, worked out
	// once from the sensor's mounting and the screen's rotation.
	upright int
}
