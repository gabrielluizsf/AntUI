//go:build android

package jni

// Error is a Java exception, carried across as a Go error.
type Error struct {
	// Class is the exception's class, dotted: "java.lang.SecurityException".
	Class string
	// Message is what getMessage said, which is often empty.
	Message string
}
