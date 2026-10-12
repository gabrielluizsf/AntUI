//go:build android

package assets

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func listNames(dir string) ([]string, error) {
	var out []string
	err := jni.Do(func(e *jni.Env) error {
		listMu.Lock()
		held := javaAssets
		listMu.Unlock()
		if held.IsNil() {
			got, err := e.Invoke(app.Context(), "getAssets",
				jni.Sig(jni.TClass("android/content/res/AssetManager")))
			if err != nil {
				return err
			}
			held = e.Global(got)
			listMu.Lock()
			if javaAssets.IsNil() {
				javaAssets = held
			} else {
				e.DeleteGlobal(held)
				held = javaAssets
			}
			listMu.Unlock()
		}
		jdir, err := e.String(inner(dir))
		if err != nil {
			return err
		}
		arr, err := e.Invoke(held, "list",
			jni.Sig(jni.TArray(jni.TString), jni.TString), jni.Ref(jdir))
		if err != nil || arr.IsNil() {
			return err
		}
		out, err = e.GoStrings(arr)
		return err
	})
	return out, err
}
