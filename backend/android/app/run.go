//go:build android

package app

// Run registers a function that is the whole app and takes no argument. It
// is for a program that reaches Android through antui.Window rather than
// through this package.
func Run(f func()) { Main(func(*App) { f() }) }
