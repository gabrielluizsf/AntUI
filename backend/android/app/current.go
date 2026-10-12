//go:build android

package app

// Current is the running activity, or nil when this process is not one —
// which is every desktop build, and an Android process before
// ANativeActivity_onCreate has been called.
func Current() *App { return current() }
