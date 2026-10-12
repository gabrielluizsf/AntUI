//go:build android

package jni

import (
	"errors"
)

var errNilObject = errors.New("jni: calling a method on a null object")
