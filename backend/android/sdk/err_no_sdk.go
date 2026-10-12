package sdk

import (
	"errors"
)

// ErrNoSDK is what [Find] reports when it cannot even locate the SDK, as
// distinct from finding one with pieces missing.
var ErrNoSDK = errors.New("sdk: no Android SDK found")
