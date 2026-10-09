package css

func blendTransformFunc(fa, fb TransformFunc, t float64) TransformFunc {
	out := fa
	switch fa.Kind {
	case TransformTranslate:
		out.Dx = blendLength(fa.Dx, fb.Dx, t, Units{})
		out.Dy = blendLength(fa.Dy, fb.Dy, t, Units{})
	case TransformScale:
		out.Sx = fa.Sx + (fb.Sx-fa.Sx)*t
		out.Sy = fa.Sy + (fb.Sy-fa.Sy)*t
	case TransformRotate:
		out.Ax = Rad(fa.Ax.Rad() + (fb.Ax.Rad()-fa.Ax.Rad())*t)
	case TransformSkew:
		out.Ax = Rad(fa.Ax.Rad() + (fb.Ax.Rad()-fa.Ax.Rad())*t)
		out.Ay = Rad(fa.Ay.Rad() + (fb.Ay.Rad()-fa.Ay.Rad())*t)
	case TransformMatrix:
		out.M = [6]float64{}
		for j := 0; j < 6; j++ {
			out.M[j] = fa.M[j] + (fb.M[j]-fa.M[j])*t
		}
	}
	return out
}
