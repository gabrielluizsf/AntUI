//go:build android

package notify

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// builderIcon sets the small icon, falling back to the platform's own.
func builderIcon(e *jni.Env, builder jni.Object, n Notification) error {
	icon := n.Icon
	if icon == 0 {
		// Any notification without a small icon is refused outright, so
		// there has to be a default, and the only one every device has is
		// the platform's own.
		var err error
		if icon, err = e.ConstantInt("android/R$drawable", "ic_dialog_info"); err != nil {
			return err
		}
	}
	return setBuilder(e, builder, "setSmallIcon",
		jni.Sig(jni.TClass("android/app/Notification$Builder"), jni.TInt), jni.Int(icon))
}
