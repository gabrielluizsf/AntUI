//go:build windows

package windows

import (
	"unsafe"

	"github.com/gabrielluizsf/antui/canvas"
)

// hcursor builds a cursor from a picture, hung from the point hotX,hotY —
// the pixel in it that the pointer actually points at. It is hicon's twin:
// the same colour-bitmap-plus-empty-mask trick, but marked as a cursor and
// given a hotspot instead of being centred.
func hcursor(image *canvas.Canvas, hotX, hotY int) uintptr {
	if image == nil || image.Width <= 0 || image.Height <= 0 {
		return 0
	}
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

	info := iconInfo{icon: 0, xHotspot: uint32(hotX), yHotspot: uint32(hotY), mask: mask, colour: colour}
	handle, _, _ := procCreateIconIndirect.Call(uintptr(unsafe.Pointer(&info)))
	return handle
}
