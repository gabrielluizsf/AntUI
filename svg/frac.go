package svg

// frac is a colour's channel as the number an operation works on: between
// zero and one, where a byte's two hundred and fifty-five steps are.
func frac(c uint8) float64 { return float64(c) / 255 }
