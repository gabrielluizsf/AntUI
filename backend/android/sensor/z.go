//go:build android

package sensor

// Z is the third number. See [Reading.X] for which way it points.
func (r Reading) Z() float32 { return r.V[2] }
