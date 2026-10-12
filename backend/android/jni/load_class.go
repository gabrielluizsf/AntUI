//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"strings"
)

// loadClass asks the app's own class loader. Its argument is a dotted name,
// not a slashed one — the two spellings are used in different halves of the
// same API.
func (e *Env) loadClass(name string) (C.jclass, error) {
	l, m, err := e.appLoader()
	if err != nil {
		return 0, err
	}
	dotted, err := e.String(strings.ReplaceAll(name, "/", "."))
	if err != nil {
		return 0, err
	}
	a := jvalues([]Value{Ref(dotted)})
	got := C.mw_CallObject(e.ptr, l, m, first(a))
	keep(a)
	if err := e.check(); err != nil {
		return 0, errors.New("jni: no class " + name + ": " + err.Error())
	}
	if got == 0 {
		return 0, errors.New("jni: no class " + name)
	}
	return C.jclass(got), nil
}
