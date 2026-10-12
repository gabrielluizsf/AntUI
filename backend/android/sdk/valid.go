package sdk

// Valid reports whether the ABI is one this package knows.
func (a ABI) Valid() bool { return a.GOARCH() != "" }
