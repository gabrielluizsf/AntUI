//go:build android

package jni

const (
	kindMethod uint8 = iota
	kindStaticMethod
	kindField
	kindStaticField
)
