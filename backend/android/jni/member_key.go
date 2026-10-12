//go:build android

package jni

// memberKey identifies one method or field of one class.
type memberKey struct {
	class uintptr
	name  string
	sig   string
	// kind separates the four namespaces, which overlap: a class may have a
	// static and an instance member of the same name and signature.
	kind uint8
}
