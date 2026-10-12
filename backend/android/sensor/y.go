//go:build android

package sensor

// Y is the second number. See [Reading.X] for which way it points.
func (r Reading) Y() float32 { return r.V[1] }
