//go:build android

package jni

import "errors"

var errNoFrame = errors.New("jni: cannot open a local reference frame")
