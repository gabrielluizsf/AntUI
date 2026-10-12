//go:build android

package picker

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Folder asks the user for a whole directory, and gives back its address.
//
// It is how an app is given a place to work in rather than one file — a
// folder of photographs to process, a directory to save into every time. The
// address is a *tree*, and it is not something [Open] takes: reading what is
// in it means the document API, which this library does not wrap yet.
//
// Pair it with [Keep], or the choice is forgotten the moment the app closes.
func Folder() (string, error) {
	var in jni.Object
	err := jni.Do(func(e *jni.Env) error {
		action, err := e.String("android.intent.action.OPEN_DOCUMENT_TREE")
		if err != nil {
			return err
		}
		made, err := e.Make("android/content/Intent",
			jni.Sig(jni.TVoid, jni.TString), jni.Ref(action))
		if err != nil {
			return err
		}
		if _, err := e.Invoke(made, "addFlags",
			jni.Sig(jni.TClass("android/content/Intent"), jni.TInt),
			jni.Int(grantRead|grantWrite|grantPersist)); err != nil {
			return err
		}
		in = e.Global(made)
		return nil
	})
	if err != nil {
		return "", err
	}
	return single(in)
}
