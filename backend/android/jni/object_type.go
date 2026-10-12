//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

// Object is a Java object. It is a reference, not the object: which kind of
// reference decides how long it lasts.
//
//   - A **local** reference, which is what every call gives back, is valid
//     until the [Do] that made it returns. That is the whole point of Do:
//     the frame it opens is what frees them.
//   - A **global** reference, from [Env.Global], lasts until it is deleted,
//     and may be used from any thread. Everything cached between calls has
//     to be one.
//   - A **weak** reference does not keep the object alive, and may become
//     nil at any moment.
//
// Using a local reference after its frame has closed does not fail. It reads
// whatever now sits at that slot, which is another object, and the crash
// happens somewhere else entirely.
type Object struct {
	ref C.jobject
}
