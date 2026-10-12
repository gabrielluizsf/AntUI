package apk

import (
	"image/color"
	_ "image/jpeg"
)

func (c canvasImage) ColorModel() color.Model { return color.NRGBAModel }
