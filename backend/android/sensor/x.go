//go:build android

package sensor

// X, Y and Z are the first three numbers, for a sensor that measures a
// direction. On a phone held upright facing the user, X points right, Y
// points up and Z points out of the screen — and that is the **device's**
// frame, which does not turn when the screen does. See [Reading.ForDisplay].
func (r Reading) X() float32 { return r.V[0] }
