//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/canvas"
	"unsafe"
)

// hicon builds an icon from a picture: a 32-bit colour bitmap with the alpha
// in it, and an empty mask, which is what a modern icon is. The mask has to
// exist even though nothing reads it.
func hicon(image *canvas.Canvas) uintptr {
	if image == nil || image.Width <= 0 || image.Height <= 0 {
		return 0
	}
	// Windows wants the rows the other way up and the channels as BGRA,
	// which is what a Color already is on a little-endian machine — so only
	// the rows are turned over here.
	pixels := make([]uint32, image.Width*image.Height)
	for y := range image.Height {
		row := (image.Height - 1 - y) * image.Width
		for x := range image.Width {
			pixels[row+x] = uint32(image.At(x, y))
		}
	}

	colour, _, _ := procCreateBitmap.Call(uintptr(image.Width), uintptr(image.Height),
		1, 32, uintptr(unsafe.Pointer(&pixels[0])))
	if colour == 0 {
		return 0
	}
	defer procDeleteObject.Call(colour)

	mask, _, _ := procCreateBitmap.Call(uintptr(image.Width), uintptr(image.Height),
		1, 1, 0)
	if mask == 0 {
		return 0
	}
	defer procDeleteObject.Call(mask)

	info := iconInfo{icon: 1, mask: mask, colour: colour}
	handle, _, _ := procCreateIconIndirect.Call(uintptr(unsafe.Pointer(&info)))
	return handle
}
