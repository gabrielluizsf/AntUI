//go:build android

package picker

import "github.com/gabrielluizsf/antui/backend/android/jni"

// openDocumentIntent builds the ACTION_OPEN_DOCUMENT intent for one or many
// files of mime.
func openDocumentIntent(e *jni.Env, mime string, many bool) (jni.Object, error) {
	action, err := e.String("android.intent.action.OPEN_DOCUMENT")
	if err != nil {
		return jni.Object{}, err
	}
	in, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TString), jni.Ref(action))
	if err != nil {
		return jni.Object{}, err
	}
	jmime, err := e.String(mime)
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(in, "setType",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString),
		jni.Ref(jmime)); err != nil {
		return jni.Object{}, err
	}
	// CATEGORY_OPENABLE keeps out the things that have no bytes behind
	// them, which the picker will otherwise happily offer.
	cat, err := e.String("android.intent.category.OPENABLE")
	if err != nil {
		return jni.Object{}, err
	}
	if _, err := e.Invoke(in, "addCategory",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TString),
		jni.Ref(cat)); err != nil {
		return jni.Object{}, err
	}
	// Read, and persistable so that Keep can make the choice survive a
	// restart. A picker started without the persist flag gives an
	// address that stops working when the process does, and takes the
	// permission with it.
	if _, err := e.Invoke(in, "addFlags",
		jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
		jni.Int(grantRead|grantPersist)); err != nil {
		return jni.Object{}, err
	}
	if many {
		if err := allowMultiple(e, in); err != nil {
			return jni.Object{}, err
		}
	}
	return in, nil
}
