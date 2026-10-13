package template

// buttonRadius is how rounded a control's corners are. It is capped at eight
// scale units, so a button never turns into a pill: just enough rounding to
// sit like the modern dashboard, whatever its height.
func buttonRadius(h, u int) int { return min(h/2, radiusMax*u) }
