package css

// BackSize is a background-size: explicit widths in each axis, or the cover
// and contain keywords. The zero value is "auto auto", the picture at its own
// size.
type BackSize struct {
	W, H    Length
	Cover   bool
	Contain bool
}
