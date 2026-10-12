//go:build android

package sensor

// Value is the first number, for a sensor that measures a quantity: lux,
// hectopascal, degrees, or a count of steps.
func (r Reading) Value() float32 { return r.V[0] }
