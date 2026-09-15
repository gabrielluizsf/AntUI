//go:build android

package jni

/*
#include "jni_wrap.h"
*/
import "C"

import (
	"errors"
	"math"
	"runtime"
	"unsafe"
)

// Value is one argument to a Java method: the eight bytes of a jvalue, which
// is a union of every type Java has.
//
// It is built with [Int], [Ref] and the rest rather than converted, because
// which eight bytes a number occupies depends on what it is. Passing an int
// where a long is expected is not a type error anywhere in JNI — it is a
// method that reads the wrong half of the stack.
type Value uint64

// A jvalue is a union containing a jlong and a jdouble, so it is eight bytes
// on every architecture Android has. Value being a uint64 depends on that,
// and these two lines are what would notice if it ever stopped being true.
var (
	_ [unsafe.Sizeof(C.jvalue{}) - 8]byte
	_ [8 - unsafe.Sizeof(C.jvalue{})]byte
)

// The argument builders, one per Java type. A Value is one slot of a JNI
// argument list, and which builder made it is what says how the eight bytes
// are read on the other side — an int and a float of the same width are not
// the same argument.
func Bool(v bool) Value {
	if v {
		return 1
	}
	return 0
}

// Byte builds a Java byte argument.
func Byte(v int8) Value { return Value(uint8(v)) }

// Char builds a Java char argument, which is one UTF-16 unit and not a rune.
func Char(v uint16) Value { return Value(v) }

// Short builds a Java short argument.
func Short(v int16) Value { return Value(uint16(v)) }

// Int builds a Java int argument.
func Int(v int32) Value { return Value(uint32(v)) }

// Long builds a Java long argument.
func Long(v int64) Value { return Value(uint64(v)) }

// Float builds a Java float argument.
func Float(v float32) Value { return Value(math.Float32bits(v)) }

// Double builds a Java double argument.
func Double(v float64) Value {
	return Value(math.Float64bits(v))
}

// Ref passes an object, or Java's null when the object is nil.
func Ref(o Object) Value { return Value(o.ref) }

// jvalues lays the arguments out the way the machine expects to read them.
func jvalues(vals []Value) []C.jvalue {
	if len(vals) == 0 {
		return nil
	}
	out := make([]C.jvalue, len(vals))
	for i, v := range vals {
		*(*uint64)(unsafe.Pointer(&out[i])) = uint64(v)
	}
	return out
}

func first(a []C.jvalue) *C.jvalue {
	if len(a) == 0 {
		return nil
	}
	return &a[0]
}

// keep stops the collector taking the argument array away while the call
// that is reading it is still running.
func keep(a []C.jvalue) { runtime.KeepAlive(a) }

var errNilObject = errors.New("jni: calling a method on a null object")

// The instance calls. Each checks for an exception afterwards and hands it
// back as an error, so a Java method that throws is a Go call that fails
// rather than a process that dies two calls later.

// CallVoid calls a method that returns nothing.
func (e *Env) CallVoid(o Object, m Method, args ...Value) error {
	if o.IsNil() {
		return errNilObject
	}
	a := jvalues(args)
	C.mw_CallVoid(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return e.check()
}

// CallBool calls a method that returns a Java boolean.
func (e *Env) CallBool(o Object, m Method, args ...Value) (bool, error) {
	if o.IsNil() {
		return false, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallBoolean(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return r != C.JNI_FALSE, e.check()
}

// CallInt calls a method that returns a Java int.
func (e *Env) CallInt(o Object, m Method, args ...Value) (int32, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallInt(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return int32(r), e.check()
}

// CallLong calls a method that returns a Java long.
func (e *Env) CallLong(o Object, m Method, args ...Value) (int64, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallLong(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return int64(r), e.check()
}

// CallFloat calls a method that returns a Java float.
func (e *Env) CallFloat(o Object, m Method, args ...Value) (float32, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallFloat(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return float32(r), e.check()
}

// CallDouble calls a method that returns a Java double.
func (e *Env) CallDouble(o Object, m Method, args ...Value) (float64, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallDouble(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return float64(r), e.check()
}

// CallObject calls a method that returns an object. The result is
// a local reference and lives until the frame it was made in ends — see
// [Env.NewGlobal] for one that has to outlive the call.
func (e *Env) CallObject(o Object, m Method, args ...Value) (Object, error) {
	if o.IsNil() {
		return Object{}, errNilObject
	}
	a := jvalues(args)
	r := C.mw_CallObject(e.ptr, o.ref, m.id, first(a))
	keep(a)
	return Object{ref: r}, e.check()
}

// CallString is CallObject and then reading the String, which is common
// enough to be worth one call instead of three lines.
func (e *Env) CallString(o Object, m Method, args ...Value) (string, error) {
	r, err := e.CallObject(o, m, args...)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}

// The static calls.

// CallStaticVoid calls a class method that returns nothing.
func (e *Env) CallStaticVoid(c Class, m Method, args ...Value) error {
	a := jvalues(args)
	C.mw_CallStaticVoid(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return e.check()
}

// CallStaticBool calls a class method that returns a Java boolean.
func (e *Env) CallStaticBool(c Class, m Method, args ...Value) (bool, error) {
	a := jvalues(args)
	r := C.mw_CallStaticBoolean(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return r != C.JNI_FALSE, e.check()
}

// CallStaticInt calls a class method that returns a Java int.
func (e *Env) CallStaticInt(c Class, m Method, args ...Value) (int32, error) {
	a := jvalues(args)
	r := C.mw_CallStaticInt(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return int32(r), e.check()
}

// CallStaticLong calls a class method that returns a Java long.
func (e *Env) CallStaticLong(c Class, m Method, args ...Value) (int64, error) {
	a := jvalues(args)
	r := C.mw_CallStaticLong(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return int64(r), e.check()
}

// CallStaticFloat calls a class method that returns a Java float.
func (e *Env) CallStaticFloat(c Class, m Method, args ...Value) (float32, error) {
	a := jvalues(args)
	r := C.mw_CallStaticFloat(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return float32(r), e.check()
}

// CallStaticDouble calls a class method that returns a Java double.
func (e *Env) CallStaticDouble(c Class, m Method, args ...Value) (float64, error) {
	a := jvalues(args)
	r := C.mw_CallStaticDouble(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return float64(r), e.check()
}

// CallStaticObject calls a class method that returns an object.
func (e *Env) CallStaticObject(c Class, m Method, args ...Value) (Object, error) {
	a := jvalues(args)
	r := C.mw_CallStaticObject(e.ptr, c.ref, m.id, first(a))
	keep(a)
	return Object{ref: r}, e.check()
}

// CallStaticString calls a class method that returns a String, and
// reads it out as Go text.
func (e *Env) CallStaticString(c Class, m Method, args ...Value) (string, error) {
	r, err := e.CallStaticObject(c, m, args...)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}

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

// Instance fields.

// GetInt reads an int field of an object.
func (e *Env) GetInt(o Object, f Field) (int32, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	r := C.mw_GetIntField(e.ptr, o.ref, f.id)
	return int32(r), e.check()
}

// GetLong reads a long field of an object.
func (e *Env) GetLong(o Object, f Field) (int64, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	r := C.mw_GetLongField(e.ptr, o.ref, f.id)
	return int64(r), e.check()
}

// GetBool reads a boolean field of an object.
func (e *Env) GetBool(o Object, f Field) (bool, error) {
	if o.IsNil() {
		return false, errNilObject
	}
	r := C.mw_GetBooleanField(e.ptr, o.ref, f.id)
	return r != C.JNI_FALSE, e.check()
}

// GetFloat reads a float field of an object.
func (e *Env) GetFloat(o Object, f Field) (float32, error) {
	if o.IsNil() {
		return 0, errNilObject
	}
	r := C.mw_GetFloatField(e.ptr, o.ref, f.id)
	return float32(r), e.check()
}

// GetObject reads an object field, as a local reference.
func (e *Env) GetObject(o Object, f Field) (Object, error) {
	if o.IsNil() {
		return Object{}, errNilObject
	}
	r := C.mw_GetObjectField(e.ptr, o.ref, f.id)
	return Object{ref: r}, e.check()
}

// GetString reads a String field and hands it back as Go text.
func (e *Env) GetString(o Object, f Field) (string, error) {
	r, err := e.GetObject(o, f)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}

// SetInt writes an int field of an object.
func (e *Env) SetInt(o Object, f Field, v int32) error {
	if o.IsNil() {
		return errNilObject
	}
	C.mw_SetIntField(e.ptr, o.ref, f.id, C.jint(v))
	return e.check()
}

// SetFloat writes a float field of an object.
func (e *Env) SetFloat(o Object, f Field, v float32) error {
	if o.IsNil() {
		return errNilObject
	}
	C.mw_SetFloatField(e.ptr, o.ref, f.id, C.jfloat(v))
	return e.check()
}

// SetObject writes an object field.
func (e *Env) SetObject(o Object, f Field, v Object) error {
	if o.IsNil() {
		return errNilObject
	}
	C.mw_SetObjectField(e.ptr, o.ref, f.id, v.ref)
	return e.check()
}

// Static fields, which is where the platform keeps its constants —
// Build.VERSION.SDK_INT, Context.VIBRATOR_SERVICE, and a thousand others.

// StaticInt reads an int field of a class rather than of an object.
func (e *Env) StaticInt(c Class, f Field) (int32, error) {
	r := C.mw_GetStaticIntField(e.ptr, c.ref, f.id)
	return int32(r), e.check()
}

// StaticLong reads a long field of a class.
func (e *Env) StaticLong(c Class, f Field) (int64, error) {
	r := C.mw_GetStaticLongField(e.ptr, c.ref, f.id)
	return int64(r), e.check()
}

// StaticObject reads an object field of a class, as a local reference.
func (e *Env) StaticObject(c Class, f Field) (Object, error) {
	r := C.mw_GetStaticObjectField(e.ptr, c.ref, f.id)
	return Object{ref: r}, e.check()
}

// StaticString reads a String field of a class, as Go text.
func (e *Env) StaticString(c Class, f Field) (string, error) {
	r, err := e.StaticObject(c, f)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}
