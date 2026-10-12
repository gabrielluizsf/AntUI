//go:build android

package location

// Place is an address, and where it is.
type Place struct {
	// Address is the whole thing as one line, the way it would be written on
	// an envelope in that country.
	Address string
	// Locality is the town or city, Region the state or province.
	Locality string
	Region   string
	Country  string
	PostCode string
	// Latitude and Longitude are where the address is, which for a forward
	// lookup is the answer and for a reverse one is roughly what was asked.
	Latitude, Longitude float64
}
