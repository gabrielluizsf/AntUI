package template

import "github.com/gabrielluizsf/antui/canvas"

// image is one registered background-image. A name the stylesheet asks for
// that was never registered paints nothing.
func (cs *cssStyle) image(name string) *canvas.Canvas { return cs.images[name] }
