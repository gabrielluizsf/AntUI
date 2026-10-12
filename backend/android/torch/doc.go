//go:build android

package torch

// Package torch turns the camera's flash on and steady, as a lamp.
//
// It needs **no permission**, which is worth saying because everything else
// to do with the camera does. Since API 23 the platform treats the torch as
// its own thing: an app may light it without being allowed to see through
// the lens.
