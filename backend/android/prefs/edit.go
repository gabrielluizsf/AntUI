//go:build android

package prefs

import "github.com/gabrielluizsf/antui/backend/android/jni"

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
