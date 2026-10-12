//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
)

// New makes an object. The method is a constructor from [Env.Constructor].
func (e *Env) New(c Class, m Method, args ...Value) (Object, error) {
	a := jvalues(args)
	r := C.mw_NewObject(e.ptr, c.ref, m.id, first(a))
	keep(a)
	if err := e.check(); err != nil {
		return Object{}, err
	}
	if r == 0 {
		return Object{}, errors.New("jni: the constructor returned nothing")
	}
	return Object{ref: r}, nil
}
