//go:build android

package jni

import "sync"

// The feature packages above this one all do the same thing: find a class,
// find a method on it, call it, read the result. Written out that is five
// lines and four error checks for one call, and a package that wraps twenty
// platform calls is four hundred lines of it.
//
// The Invoke calls below collapse that to one line each. They are for calls
// made occasionally — asking the battery level, putting something on the
// clipboard — and not for a frame loop, though the cache below means even
// that would only cost a map lookup.

// memberKey identifies one method or field of one class.
type memberKey struct {
	class uintptr
	name  string
	sig   string
	// kind separates the four namespaces, which overlap: a class may have a
	// static and an instance member of the same name and signature.
	kind uint8
}

const (
	kindMethod uint8 = iota
	kindStaticMethod
	kindField
	kindStaticField
)

var (
	memberMu sync.RWMutex
	methods  = map[memberKey]Method{}
	fields   = map[memberKey]Field{}
)

// cachedMethod is Method or StaticMethod with the lookup remembered.
//
// A method id stays valid for as long as its class is loaded, and every
// class this package hands out is a global reference, so they never go. The
// cache is therefore permanent and that is correct rather than a leak: there
// are only so many methods in the platform.
func (e *Env) cachedMethod(c Class, name, sig string, kind uint8) (Method, error) {
	key := memberKey{uintptr(c.ref), name, sig, kind}
	memberMu.RLock()
	m, ok := methods[key]
	memberMu.RUnlock()
	if ok {
		return m, nil
	}
	var err error
	if kind == kindStaticMethod {
		m, err = e.StaticMethod(c, name, sig)
	} else {
		m, err = e.Method(c, name, sig)
	}
	if err != nil {
		return Method{}, err
	}
	memberMu.Lock()
	methods[key] = m
	memberMu.Unlock()
	return m, nil
}

func (e *Env) cachedField(c Class, name, sig string, kind uint8) (Field, error) {
	key := memberKey{uintptr(c.ref), name, sig, kind}
	memberMu.RLock()
	f, ok := fields[key]
	memberMu.RUnlock()
	if ok {
		return f, nil
	}
	var err error
	if kind == kindStaticField {
		f, err = e.StaticField(c, name, sig)
	} else {
		f, err = e.Field(c, name, sig)
	}
	if err != nil {
		return Field{}, err
	}
	memberMu.Lock()
	fields[key] = f
	memberMu.Unlock()
	return f, nil
}

// method finds a method on an object's own class.
//
// **It does not cache, and it lets the class reference go before it
// returns.** Both halves of that were wrong before and both were invisible.
//
// [Env.ClassOf] hands back a *local* reference, and a local reference is a
// slot in a table that is reused as soon as the frame that made it ends — so
// the same number means one class now and a different class a minute later.
// Keying a cache on it is how `close()V`, resolved once for an
// SQLiteCursor, came back for an SQLiteDatabase. Nothing notices: the ids
// are opaque, the call goes through, and the VM does whatever that method's
// code does to an object of the wrong type. The only thing that says so is
// the VM's own CheckJNI —
//
//	adb shell setprop debug.checkjni 1 && adb shell stop && adb shell start
//
// which aborts with "can't call void android.database.sqlite.
// SQLiteCursor.close() on instance of android.database.sqlite.
// SQLiteDatabase". It is worth turning on now and then for exactly this.
//
// Not deleting the reference was the same mistake's other half: a loop of
// Invoke calls filled the local reference table with class references nobody
// would look at again, and the table is 512 entries deep.
//
// The named-class calls — [Env.Static], [Env.Constant], [Env.Make] — do
// cache, and safely: [Env.Class] keeps one **global** reference per name,
// and a global is unique and permanent. A caller who wants a cached lookup
// on a hot path should name the class.
func (e *Env) method(o Object, name, sig string) (Method, error) {
	c, err := e.ClassOf(o)
	if err != nil {
		return Method{}, err
	}
	defer e.Delete(c.Object())
	return e.Method(c, name, sig)
}

// Invoke calls a method on an object, finding the method by name.
func (e *Env) Invoke(o Object, name, sig string, args ...Value) (Object, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return Object{}, err
	}
	return e.CallObject(o, m, args...)
}

// InvokeVoid, InvokeBool and the rest are Invoke for the other return types.
func (e *Env) InvokeVoid(o Object, name, sig string, args ...Value) error {
	m, err := e.method(o, name, sig)
	if err != nil {
		return err
	}
	return e.CallVoid(o, m, args...)
}

// InvokeBool is [Env.Invoke] for a method that returns a boolean.
func (e *Env) InvokeBool(o Object, name, sig string, args ...Value) (bool, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return false, err
	}
	return e.CallBool(o, m, args...)
}

// InvokeInt is [Env.Invoke] for a method that returns an int.
func (e *Env) InvokeInt(o Object, name, sig string, args ...Value) (int32, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return 0, err
	}
	return e.CallInt(o, m, args...)
}

// InvokeLong is [Env.Invoke] for a method that returns a long.
func (e *Env) InvokeLong(o Object, name, sig string, args ...Value) (int64, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return 0, err
	}
	return e.CallLong(o, m, args...)
}

// InvokeFloat is [Env.Invoke] for a method that returns a float.
func (e *Env) InvokeFloat(o Object, name, sig string, args ...Value) (float32, error) {
	m, err := e.method(o, name, sig)
	if err != nil {
		return 0, err
	}
	return e.CallFloat(o, m, args...)
}

// InvokeString is [Env.Invoke] for a method that returns a String,
// read out as Go text.
func (e *Env) InvokeString(o Object, name, sig string, args ...Value) (string, error) {
	r, err := e.Invoke(o, name, sig, args...)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}

// The static equivalents, named by the class rather than by an instance.

// Static calls a class method and hands back what it returned, looking the
// class and the method up once and keeping them — see [Env.Invoke], which is
// the same thing for a method on an object.
func (e *Env) Static(class, name, sig string, args ...Value) (Object, error) {
	c, m, err := e.staticOf(class, name, sig)
	if err != nil {
		return Object{}, err
	}
	return e.CallStaticObject(c, m, args...)
}

// StaticVoid is [Env.Static] for a class method that returns nothing.
func (e *Env) StaticVoid(class, name, sig string, args ...Value) error {
	c, m, err := e.staticOf(class, name, sig)
	if err != nil {
		return err
	}
	return e.CallStaticVoid(c, m, args...)
}

// StaticIntOf is [Env.Static] for a class method that returns an int.
func (e *Env) StaticIntOf(class, name, sig string, args ...Value) (int32, error) {
	c, m, err := e.staticOf(class, name, sig)
	if err != nil {
		return 0, err
	}
	return e.CallStaticInt(c, m, args...)
}

func (e *Env) staticOf(class, name, sig string) (Class, Method, error) {
	c, err := e.Class(class)
	if err != nil {
		return Class{}, Method{}, err
	}
	m, err := e.cachedMethod(c, name, sig, kindStaticMethod)
	return c, m, err
}

// Constant reads a static field, which is where the platform keeps most of
// its named values — Context.VIBRATOR_SERVICE, Intent.ACTION_SEND, and a
// thousand others.
func (e *Env) Constant(class, name, sig string) (Object, error) {
	c, err := e.Class(class)
	if err != nil {
		return Object{}, err
	}
	f, err := e.cachedField(c, name, sig, kindStaticField)
	if err != nil {
		return Object{}, err
	}
	return e.StaticObject(c, f)
}

// ConstantString and ConstantInt are Constant for the two types most of them
// are.
func (e *Env) ConstantString(class, name string) (string, error) {
	o, err := e.Constant(class, name, TString)
	if err != nil {
		return "", err
	}
	return e.GoString(o)
}

// ConstantInt reads an int constant off a class — the API levels,
// the flags and the modes Android keeps as static final fields.
func (e *Env) ConstantInt(class, name string) (int32, error) {
	c, err := e.Class(class)
	if err != nil {
		return 0, err
	}
	f, err := e.cachedField(c, name, TInt, kindStaticField)
	if err != nil {
		return 0, err
	}
	return e.StaticInt(c, f)
}

// Make builds an object: find the class, find the constructor, call it.
func (e *Env) Make(class, sig string, args ...Value) (Object, error) {
	c, err := e.Class(class)
	if err != nil {
		return Object{}, err
	}
	ctor, err := e.cachedMethod(c, "<init>", sig, kindMethod)
	if err != nil {
		return Object{}, err
	}
	return e.New(c, ctor, args...)
}
