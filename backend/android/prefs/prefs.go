//go:build android

package prefs

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// modePrivate is Context.MODE_PRIVATE: readable by this app and no other.
// It is the only mode still allowed — the world-readable ones were removed
// in Android 7 and throw if asked for.
const modePrivate = 0

// Prefs is one named store of small values, kept as XML in the app's own
// directory and read back on the next run.
//
// It is for settings and for the scraps of state a program wants to
// remember: which level was reached, whether the sound is on, the last
// window it had. It is not a database and not a file — a value larger than a
// few kilobytes belongs somewhere else, and the whole store is loaded into
// memory the first time it is opened.
type Prefs struct {
	name string
}

// Open names a store. Nothing happens until something is read or written, so
// this cannot fail.
func Open(name string) *Prefs { return &Prefs{name: name} }

// Default is the store Android gives an app that never asks for one by name.
func Default() *Prefs { return &Prefs{name: ""} }

// with runs f with the SharedPreferences object.
func (p *Prefs) with(f func(*jni.Env, jni.Object) error) error {
	return jni.Do(func(e *jni.Env) error {
		var obj jni.Object
		var err error
		if p.name == "" {
			obj, err = e.Static("android/preference/PreferenceManager",
				"getDefaultSharedPreferences",
				jni.Sig(jni.TClass("android/content/SharedPreferences"), jni.TContext),
				jni.Ref(app.Context()))
		} else {
			var jname jni.Object
			if jname, err = e.String(p.name); err != nil {
				return err
			}
			obj, err = e.Invoke(app.Context(), "getSharedPreferences",
				jni.Sig(jni.TClass("android/content/SharedPreferences"), jni.TString, jni.TInt),
				jni.Ref(jname), jni.Int(modePrivate))
		}
		if err != nil {
			return err
		}
		return f(e, obj)
	})
}

// String reads a string, or fallback when the key is not there.
func (p *Prefs) String(key, fallback string) (string, error) {
	var out string
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		d, err := e.String(fallback)
		if err != nil {
			return err
		}
		out, err = e.InvokeString(o, "getString",
			jni.Sig(jni.TString, jni.TString, jni.TString), jni.Ref(k), jni.Ref(d))
		return err
	})
	return out, err
}

// Int, Bool, Float and Long are String for the other types the store holds.
func (p *Prefs) Int(key string, fallback int32) (int32, error) {
	var out int32
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeInt(o, "getInt",
			jni.Sig(jni.TInt, jni.TString, jni.TInt), jni.Ref(k), jni.Int(fallback))
		return err
	})
	return out, err
}

// Bool reads a true-or-false setting, or the fallback when it is not
// there.
func (p *Prefs) Bool(key string, fallback bool) (bool, error) {
	var out bool
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(o, "getBoolean",
			jni.Sig(jni.TBool, jni.TString, jni.TBool), jni.Ref(k), jni.Bool(fallback))
		return err
	})
	return out, err
}

// Float reads a float setting, or the fallback when it is not there.
func (p *Prefs) Float(key string, fallback float32) (float32, error) {
	var out float32
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeFloat(o, "getFloat",
			jni.Sig(jni.TFloat, jni.TString, jni.TFloat), jni.Ref(k), jni.Float(fallback))
		return err
	})
	return out, err
}

// Long reads a 64-bit setting, or the fallback when it is not there.
func (p *Prefs) Long(key string, fallback int64) (int64, error) {
	var out int64
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeLong(o, "getLong",
			jni.Sig(jni.TLong, jni.TString, jni.TLong), jni.Ref(k), jni.Long(fallback))
		return err
	})
	return out, err
}

// Has reports whether a key is in the store at all, which is how a missing
// value is told from one that happens to equal the fallback.
func (p *Prefs) Has(key string) (bool, error) {
	var out bool
	err := p.with(func(e *jni.Env, o jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(o, "contains", jni.Sig(jni.TBool, jni.TString), jni.Ref(k))
		return err
	})
	return out, err
}

// edit opens an editor, runs f, and commits with apply — which writes to
// memory now and to disk on a background thread. commit would write to disk
// before returning, and doing that from a frame loop is how an app stutters.
func (p *Prefs) edit(f func(*jni.Env, jni.Object) error) error {
	return p.with(func(e *jni.Env, o jni.Object) error {
		ed, err := e.Invoke(o, "edit",
			jni.Sig(jni.TClass("android/content/SharedPreferences$Editor")))
		if err != nil {
			return err
		}
		if err := f(e, ed); err != nil {
			return err
		}
		return e.InvokeVoid(ed, "apply", jni.Sig(jni.TVoid))
	})
}

// SetString, SetInt, SetBool, SetFloat and SetLong write one value.
func (p *Prefs) SetString(key, value string) error {
	return p.edit(func(e *jni.Env, ed jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		v, err := e.String(value)
		if err != nil {
			return err
		}
		_, err = e.Invoke(ed, "putString",
			jni.Sig(jni.TClass("android/content/SharedPreferences$Editor"),
				jni.TString, jni.TString), jni.Ref(k), jni.Ref(v))
		return err
	})
}

// SetInt writes an int setting.
func (p *Prefs) SetInt(key string, value int32) error {
	return p.putScalar(key, "putInt", jni.TInt, jni.Int(value))
}

// SetBool writes a true-or-false setting.
func (p *Prefs) SetBool(key string, value bool) error {
	return p.putScalar(key, "putBoolean", jni.TBool, jni.Bool(value))
}

// SetFloat writes a float setting.
func (p *Prefs) SetFloat(key string, value float32) error {
	return p.putScalar(key, "putFloat", jni.TFloat, jni.Float(value))
}

// SetLong writes a 64-bit setting.
func (p *Prefs) SetLong(key string, value int64) error {
	return p.putScalar(key, "putLong", jni.TLong, jni.Long(value))
}

func (p *Prefs) putScalar(key, method, typ string, value jni.Value) error {
	return p.edit(func(e *jni.Env, ed jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		_, err = e.Invoke(ed, method,
			jni.Sig(jni.TClass("android/content/SharedPreferences$Editor"), jni.TString, typ),
			jni.Ref(k), value)
		return err
	})
}

// Remove takes a key out.
func (p *Prefs) Remove(key string) error {
	return p.edit(func(e *jni.Env, ed jni.Object) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		_, err = e.Invoke(ed, "remove",
			jni.Sig(jni.TClass("android/content/SharedPreferences$Editor"), jni.TString),
			jni.Ref(k))
		return err
	})
}

// Clear empties the store.
func (p *Prefs) Clear() error {
	return p.edit(func(e *jni.Env, ed jni.Object) error {
		_, err := e.Invoke(ed, "clear",
			jni.Sig(jni.TClass("android/content/SharedPreferences$Editor")))
		return err
	})
}
