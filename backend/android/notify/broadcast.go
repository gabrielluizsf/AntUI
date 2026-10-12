//go:build android

package notify

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// broadcast builds the PendingIntent one action sends.
//
// The mutable flag is the trap. A PendingIntent has had to be immutable since
// API 31 — except one carrying a reply, because the system's whole job there
// is to *add* the typed text to it. Marking a reply immutable produces a
// button that works and a reply that is always empty.
func broadcast(e *jni.Env, id, index int32, name string, token int64, mutable bool) (jni.Object, error) {
	cls, err := e.Class("dev/antui/AntuiReceiver")
	if err != nil {
		return jni.Object{}, err
	}
	intent, err := e.Make("android/content/Intent",
		jni.Sig(jni.TVoid, jni.TContext, jni.TClass("java/lang/Class")),
		jni.Ref(app.Context()), jni.Ref(cls.Object()))
	if err != nil {
		return jni.Object{}, err
	}
	// A distinct action string per button, so that two PendingIntents for
	// the same notification are not treated as the same one and merged —
	// which they are if only their extras differ.
	if err := setAction(e, intent, fmt.Sprintf("dev.antui.%d.%d", id, index)); err != nil {
		return jni.Object{}, err
	}
	if err := putString(e, intent, extraAction, name); err != nil {
		return jni.Object{}, err
	}
	if err := putInt(e, intent, extraID, id); err != nil {
		return jni.Object{}, err
	}
	if err := putLong(e, intent, app.ReceiverToken, token); err != nil {
		return jni.Object{}, err
	}

	flags := flagUpdateCurrent | flagImmutable
	if mutable {
		flags = flagUpdateCurrent | flagMutable
	}
	// The request code has to differ per button for the same reason the
	// action string does.
	return e.Static("android/app/PendingIntent", "getBroadcast",
		jni.Sig(jni.TClass("android/app/PendingIntent"),
			jni.TContext, jni.TInt, jni.TClass("android/content/Intent"), jni.TInt),
		jni.Ref(app.Context()), jni.Int(id*100+index), jni.Ref(intent), jni.Int(flags))
}
