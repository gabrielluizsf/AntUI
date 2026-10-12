//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Engines is every text-to-speech engine installed, by package name. The
// list is empty on a device with none.
func Engines() ([]string, error) {
	var out []string
	err := jni.Do(func(e *jni.Env) error {
		// A throwaway engine, only to ask it what else is there. It is the
		// only way the platform offers: the list is an instance method.
		tts, err := e.Make("android/speech/tts/TextToSpeech",
			jni.Sig(jni.TVoid, jni.TContext,
				jni.TClass("android/speech/tts/TextToSpeech$OnInitListener")),
			jni.Ref(app.Context()), jni.Ref(jni.Object{}))
		if err != nil {
			return err
		}
		defer e.InvokeVoid(tts, "shutdown", jni.Sig(jni.TVoid))

		list, err := e.Invoke(tts, "getEngines", jni.Sig(jni.TClass("java/util/List")))
		if err != nil || list.IsNil() {
			return err
		}
		n, err := e.InvokeInt(list, "size", jni.Sig(jni.TInt))
		if err != nil {
			return err
		}
		return e.Frame(int(n)*4+16, func() error {
			for i := range int(n) {
				info, err := e.Invoke(list, "get", jni.Sig(jni.TObject, jni.TInt),
					jni.Int(int32(i)))
				if err != nil || info.IsNil() {
					continue
				}
				cls, err := e.Class("android/speech/tts/TextToSpeech$EngineInfo")
				if err != nil {
					return err
				}
				f, err := e.Field(cls, "name", jni.TString)
				if err != nil {
					return err
				}
				name, err := e.GetString(info, f)
				if err == nil && name != "" {
					out = append(out, name)
				}
			}
			return nil
		})
	})
	return out, err
}
