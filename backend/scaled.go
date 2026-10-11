package backend

// scaled is the limits with every length in points turned into pixels.
func (l Limits) Scaled(scale float64) Limits {
	l.MinWidth, l.MinHeight = scaleTo(l.MinWidth, scale), scaleTo(l.MinHeight, scale)
	l.MaxWidth, l.MaxHeight = scaleTo(l.MaxWidth, scale), scaleTo(l.MaxHeight, scale)
	return l
}
