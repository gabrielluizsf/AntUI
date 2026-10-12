//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// build makes the intent.
//
// ACTION_OPEN_DOCUMENT rather than ACTION_GET_CONTENT: the first goes to the
// system's own picker and gives an address that keeps working, and the second
// goes to whichever app claims the type and gives one that may not outlive
// the call. The system picker is also the one that shows cloud accounts.
func build(mime string, many bool) (jni.Object, error) {
	var out jni.Object
	err := jni.Do(func(e *jni.Env) error {
		action, err := e.String("android.intent.action.OPEN_DOCUMENT")
		if err != nil {
			return err
		}
		in, err := e.Make("android/content/Intent",
			jni.Sig(jni.TVoid, jni.TString), jni.Ref(action))
		if err != nil {
			return err
		}
		jmime, err := e.String(mime)
		if err != nil {
			return err
		}
		if _, err := e.Invoke(in, "setType",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString),
			jni.Ref(jmime)); err != nil {
			return err
		}
		// CATEGORY_OPENABLE keeps out the things that have no bytes behind
		// them, which the picker will otherwise happily offer.
		cat, err := e.String("android.intent.category.OPENABLE")
		if err != nil {
			return err
		}
		if _, err := e.Invoke(in, "addCategory",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TString),
			jni.Ref(cat)); err != nil {
			return err
		}
		// Read, and persistable so that Keep can make the choice survive a
		// restart. A picker started without the persist flag gives an
		// address that stops working when the process does, and takes the
		// permission with it.
		if _, err := e.Invoke(in, "addFlags",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
			jni.Int(grantRead|grantPersist)); err != nil {
			return err
		}
		if many {
			key, err := e.String("android.intent.extra.ALLOW_MULTIPLE")
			if err != nil {
				return err
			}
			if _, err := e.Invoke(in, "putExtra",
				jni.Sig(jni.TClass("android/content/Intent"), jni.TString, jni.TBool),
				jni.Ref(key), jni.Bool(true)); err != nil {
				return err
			}
		}
		out = e.Global(in)
		return nil
	})
	return out, err
}
