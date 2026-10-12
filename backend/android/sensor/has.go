//go:build android

package sensor

// Has reports whether the device has a sensor of a kind.
func Has(k Kind) bool {
	m, err := sensors()
	if err != nil {
		return false
	}
	return !m.Default(int(k)).IsNil()
}
