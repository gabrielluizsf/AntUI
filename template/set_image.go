package template

import "github.com/gabrielluizsf/antui/canvas"

// SetImage registers a canvas a stylesheet's background-image url() can point
// at. The stylesheet names images by the literal text of the url(), so it is
// lowercase, exactly as the declaration wrote it. A registered image paints at
// its own pixel size until background-size scales it.
func (c *CSS) SetImage(name string, cv *canvas.Canvas) { c.style.images[name] = cv }
