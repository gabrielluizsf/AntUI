package apk

import (
	_ "image/jpeg"
)

// safeFraction is how much of an adaptive layer is certain to be visible.
// The layer is 108dp and the guaranteed circle inside it is 72dp, so
// anything that must be seen has to be within two thirds of the middle —
// the rest may be cropped, and on some launchers is.
const safeFraction = 72.0 / 108.0
