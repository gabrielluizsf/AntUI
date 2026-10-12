//go:build android

package ndk

// Rect is a rectangle in window pixels, half-open on the right and bottom.
type Rect struct{ Left, Top, Right, Bottom int32 }
