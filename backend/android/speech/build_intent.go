//go:build android

package speech

import "github.com/gabrielluizsf/antui/backend/android/jni"

func buildIntent(e *jni.Env, opt Options) (jni.Object, error) {
	action, err := e.ConstantString("android/speech/RecognizerIntent", "ACTION_RECOGNIZE_SPEECH")
	if err != nil {
		return jni.Object{}, err
	}
	jaction, err := e.String(action)
	if err != nil {
		return jni.Object{}, err
	}
	in, err := e.Make("android/content/Intent", jni.Sig(jni.TVoid, jni.TString),
		jni.Ref(jaction))
	if err != nil {
		return jni.Object{}, err
	}
	put := func(key, value string) error {
		k, err := e.String(key)
		if err != nil {
			return err
		}
		v, err := e.String(value)
		if err != nil {
			return err
		}
		_, err = e.Invoke(in, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TString),
			jni.Ref(k), jni.Ref(v))
		return err
	}
	// The free-form model rather than the web-search one: the second is
	// tuned for short queries and gets a sentence wrong.
	if err := put("android.speech.extra.LANGUAGE_MODEL", "free_form"); err != nil {
		return jni.Object{}, err
	}
	if opt.Language != "" {
		if err := put("android.speech.extra.LANGUAGE", opt.Language); err != nil {
			return jni.Object{}, err
		}
	}
	if opt.Partial {
		k, err := e.String("android.speech.extra.PARTIAL_RESULTS")
		if err != nil {
			return jni.Object{}, err
		}
		if _, err := e.Invoke(in, "putExtra",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TBool),
			jni.Ref(k), jni.Bool(true)); err != nil {
			return jni.Object{}, err
		}
	}
	return in, nil
}
