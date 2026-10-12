//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// document builds an OPEN_DOCUMENT or CREATE_DOCUMENT intent.
func document(e *jni.Env, action, mime, title string, many bool) (jni.Object, error) {
	jaction, err := e.String(action)
	if err != nil {
		return jni.Object{}, err
	}
	in, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TString), jni.Ref(jaction))
	if err != nil {
		return jni.Object{}, err
	}
	jmime, err := e.String(mime)
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(in, "setType",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString), jni.Ref(jmime)); err != nil {
		return jni.Object{}, err
	}
	cat, err := e.String("android.intent.category.OPENABLE")
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(in, "addCategory",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString), jni.Ref(cat)); err != nil {
		return jni.Object{}, err
	}
	if title != "" {
		key, err := e.String("android.intent.extra.TITLE")
		if err != nil {
			return jni.Object{}, err
		}
		jtitle, err := e.String(title)
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(in, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TString),
			jni.Ref(key), jni.Ref(jtitle)); err != nil {
			return jni.Object{}, err
		}
	}
	if _, err := e.Invoke(in, "addFlags",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
		jni.Int(grantRead|grantWrite|grantPersist)); err != nil {
		return jni.Object{}, err
	}
	return e.Global(in), nil
}
