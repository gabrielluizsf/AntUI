package css

import "github.com/gabrielluizsf/antui/canvas"

func blendShadow(a, b Shadow, t float64) Shadow {
	if a.Inset != b.Inset {
		if t < 0.5 {
			return a
		}
		return b
	}
	return Shadow{
		X:      roundLerp(a.X, b.X, t),
		Y:      roundLerp(a.Y, b.Y, t),
		Blur:   roundLerp(a.Blur, b.Blur, t),
		Spread: roundLerp(a.Spread, b.Spread, t),
		Color:  canvas.Mix(a.Color, b.Color, float32(t)),
		Inset:  a.Inset,
	}
}
