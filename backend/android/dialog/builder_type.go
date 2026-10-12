//go:build android

package dialog

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

var builderType = jni.TClass("android/app/AlertDialog$Builder")
