//go:build android

package prefs

// Default is the store Android gives an app that never asks for one by name.
func Default() *Prefs { return &Prefs{name: ""} }
