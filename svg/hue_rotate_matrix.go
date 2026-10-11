package svg

import "math"

func hueRotateMatrix(deg float64) [20]float64 {
	a := deg * math.Pi / 180
	c, s := math.Cos(a), math.Sin(a)
	at := func(base, cosc, sinc float64) float64 { return base + cosc*c + sinc*s }
	return [20]float64{
		at(0.213, 0.787, -0.213), at(0.715, -0.715, -0.715), at(0.072, -0.072, 0.928), 0, 0,
		at(0.213, -0.213, 0.143), at(0.715, 0.285, 0.140), at(0.072, -0.072, -0.283), 0, 0,
		at(0.213, -0.213, -0.787), at(0.715, -0.715, 0.715), at(0.072, 0.928, 0.072), 0, 0,
		0, 0, 0, 1, 0,
	}
}
