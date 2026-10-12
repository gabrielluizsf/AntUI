package apk

import (
	"fmt"
)

// ErrNoBundletool is a machine without the jar that takes a bundle apart.
//
// Nothing needs it to build one — this package writes the format itself —
// and it is not part of the SDK. It is one jar from Google's releases, and
// it is the only way to see what the store will actually make of an upload.
var ErrNoBundletool = fmt.Errorf("apk: bundletool is not on this machine; " +
	"get the jar from github.com/google/bundletool/releases and put it at " +
	"$ANDROID_HOME/bundletool.jar, or point $BUNDLETOOL at it")
