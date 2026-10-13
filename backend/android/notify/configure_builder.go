//go:build android

package notify

import "github.com/gabrielluizsf/antui/backend/android/jni"

// configureBuilder fills the builder in, in the order the properties are
// read.
func configureBuilder(e *jni.Env, builder jni.Object, n Notification, sdk int) error {
	if err := builderText(e, builder, n); err != nil {
		return err
	}
	if err := builderIcon(e, builder, n); err != nil {
		return err
	}
	if err := builderFlags(e, builder, n); err != nil {
		return err
	}
	if err := builderProgress(e, builder, n, sdk); err != nil {
		return err
	}
	if err := addActions(e, builder, n); err != nil {
		return err
	}
	return builderIntent(e, builder, n)
}
