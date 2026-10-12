//go:build android

package jni

// Error is the Java exception as a Go error: the class it was, and what it
// said.
func (e *Error) Error() string {
	if e.Message == "" {
		return "java: " + e.Class
	}
	return "java: " + e.Class + ": " + e.Message
}
