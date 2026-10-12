package apk

import (
	"image"
	_ "image/jpeg"
)

func (c canvasImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.cv.Width, c.cv.Height)
}
