//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Field identifies a field, on the same terms as [Method].
type Field struct{ id C.jfieldID }
