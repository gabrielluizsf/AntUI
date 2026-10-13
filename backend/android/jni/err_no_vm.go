//go:build android

package jni

import "errors"

// ErrNoVM is what every call reports before [Init] has run — which is every
// call in a process that Android did not start as an app.
var ErrNoVM = errors.New("jni: no Java virtual machine; this process is not an Android app")
