package font

// f2dot14 is the fixed-point form composite transforms are written in.
func f2dot14(v uint16) float64 { return float64(int16(v)) / 16384 }
