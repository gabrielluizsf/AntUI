//go:build android

package antui

func newBackend() platform { return &androidWindow{primary: -1} }
