package template

import "github.com/gabrielluizsf/antui/canvas"

// fillClip fills a clip-path polygon by fanning triangles from its first
// vertex — exact for a convex polygon, which is the shape every widget's
// chamfer is.
func fillClip(cv *canvas.Canvas, pts []canvas.Point, c canvas.Color) {
	for i := 1; i+1 < len(pts); i++ {
		cv.FillTriangle([3]canvas.Point{pts[0], pts[i], pts[i+1]}, c)
	}
}
