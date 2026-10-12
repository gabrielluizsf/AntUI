//go:build android

package clipboard

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// SetText puts text on the clipboard.
//
// label is what the system shows when it describes what was copied — some
// versions put it in a toast. It is not the text and is not pasted anywhere.
func SetText(label, text string) error {
	var err error
	if e := app.RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			manager, err := app.Service(e, app.ServiceClipboard)
			if err != nil {
				return err
			}
			jlabel, err := e.String(label)
			if err != nil {
				return err
			}
			jtext, err := e.String(text)
			if err != nil {
				return err
			}
			// A ClipData is a label plus one or more items, because the
			// clipboard can hold a URI or an intent as easily as a string.
			clip, err := e.Static("android/content/ClipData", "newPlainText",
				jni.Sig(jni.TClass("android/content/ClipData"),
					jni.TClass("java/lang/CharSequence"),
					jni.TClass("java/lang/CharSequence")),
				jni.Ref(jlabel), jni.Ref(jtext))
			if err != nil {
				return err
			}
			return e.InvokeVoid(manager, "setPrimaryClip",
				jni.Sig(jni.TVoid, jni.TClass("android/content/ClipData")), jni.Ref(clip))
		})
	}); e != nil {
		return e
	}
	return err
}
