//go:build android

package app

/*
#cgo LDFLAGS: -landroid

#include <android/configuration.h>
#include <android/native_activity.h>
*/
import "C"

func twoLetters(b [2]C.char) string {
	if b[0] == 0 {
		return ""
	}
	return string([]byte{byte(b[0]), byte(b[1])})
}
