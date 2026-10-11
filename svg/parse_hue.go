package svg

// parseHue reads the hue of a colour, which is an angle: parseAngle is where it
// is read, and a hue is nothing more than one in degrees.
func parseHue(s string) (float64, error) { return parseAngle(s) }
