//go:build android

package location

import (
	"errors"
)

// ErrNoGeocoder is a device with no service behind the geocoder. It is not
// rare: the geocoder is a network service and a device without Play Services
// or without a connection has none.
var ErrNoGeocoder = errors.New("location: this device has no geocoder")
