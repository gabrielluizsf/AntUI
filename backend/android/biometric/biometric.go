//go:build android

package biometric

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Availability is whether the device can ask for a fingerprint or a face,
// and if not, why not. The distinction matters: hardware that is missing is
// permanent, and a user who has enrolled nothing can be sent to the settings
// screen to fix it.
type Availability int

const (
	// Available means the prompt will work.
	Available Availability = iota
	// NoHardware means this device has no sensor.
	NoHardware
	// NoneEnrolled means it has one and the user has not registered a
	// fingerprint or a face on it.
	NoneEnrolled
	// Unavailable means the hardware is there and busy or broken.
	Unavailable
	// TooOld means the platform predates the prompt: below API 28 there is
	// no BiometricPrompt, only the fingerprint API it replaced.
	TooOld
)

// String says why biometrics can or cannot be asked for, in a line that
// can be shown to a person.
func (a Availability) String() string {
	switch a {
	case Available:
		return "available"
	case NoHardware:
		return "no hardware"
	case NoneEnrolled:
		return "nothing enrolled"
	case Unavailable:
		return "unavailable"
	}
	return "too old"
}

// The errors a prompt can end in. Cancelled and Refused are the two a caller
// has to handle; the rest are worth telling apart when reporting.
var (
	// ErrCancelled is the user backing out, or the app cancelling.
	ErrCancelled = errors.New("biometric: cancelled")
	// ErrRefused is the negative button — the user choosing another way in.
	ErrRefused = errors.New("biometric: the user chose the other button")
	// ErrLockedOut is too many failed attempts.
	ErrLockedOut = errors.New("biometric: too many attempts; locked out for now")
	// ErrNotDeclared is the manifest missing USE_BIOMETRIC.
	//
	// It is a *normal* permission: granted at install by being declared, with
	// nothing to ask the user. Without it every call here throws a
	// SecurityException, which names the permission and reads like a bug in
	// the app rather than a line missing from its manifest — so it is caught
	// and said plainly.
	ErrNotDeclared = errors.New("biometric: this app's manifest does not declare " +
		"android.permission.USE_BIOMETRIC; add it and rebuild — it is granted at " +
		"install and the user is never asked")
)

// plainly turns the one failure everybody hits into a sentence.
func plainly(err error) error {
	var je *jni.Error
	if errors.As(err, &je) &&
		je.Class == "java.lang.SecurityException" &&
		strings.Contains(je.Message, "USE_BIOMETRIC") {
		return ErrNotDeclared
	}
	return err
}

// Can reports whether a prompt would work.
func Can() (Availability, error) {
	a := app.Current()
	if a == nil {
		return TooOld, app.ErrNoUIThread
	}
	sdk := a.Config().SDK
	if sdk < 28 {
		return TooOld, nil
	}
	if sdk < 29 {
		// BiometricManager arrived one version after the prompt did, so on
		// exactly this version the only way to ask is the older fingerprint
		// API — which is what "there is a sensor and something enrolled on
		// it" meant before faces existed.
		return fingerprintAvailability()
	}
	var out Availability
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, "biometric")
		if err != nil {
			return err
		}
		code, err := e.InvokeInt(manager, "canAuthenticate", jni.Sig(jni.TInt))
		if err != nil {
			return err
		}
		switch code {
		case 0:
			out = Available
		case 11:
			out = NoneEnrolled
		case 12:
			out = NoHardware
		default:
			out = Unavailable
		}
		return nil
	})
	return out, plainly(err)
}

func fingerprintAvailability() (Availability, error) {
	var out Availability
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, "fingerprint")
		if err != nil {
			return err
		}
		hardware, err := e.InvokeBool(manager, "isHardwareDetected", jni.Sig(jni.TBool))
		if err != nil {
			return err
		}
		if !hardware {
			out = NoHardware
			return nil
		}
		enrolled, err := e.InvokeBool(manager, "hasEnrolledFingerprints", jni.Sig(jni.TBool))
		if err != nil {
			return err
		}
		if !enrolled {
			out = NoneEnrolled
			return nil
		}
		out = Available
		return nil
	})
	return out, err
}

// Ask puts the platform's prompt on screen and waits.
//
// It returns nil when the user proved who they are, and one of the errors
// above otherwise. A finger that was not recognised is **not** a failure —
// the prompt stays up and the user tries again — so this only returns when
// the matter is settled one way or the other.
//
// cancel is the label on the way out, and cannot be empty: the platform
// insists that the user always have one.
//
// **Not from the frame loop**, for the same reason as every other thing that
// waits for a person.
func Ask(title, subtitle, cancel string) error {
	if cancel == "" {
		cancel = "Cancel"
	}
	can, err := Can()
	if err != nil {
		return err
	}
	if can != Available {
		return fmt.Errorf("biometric: %s", can)
	}

	done := make(chan error, 1)
	var once sync.Once
	finish := func(err error) { once.Do(func() { done <- err }) }

	token, release := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		switch method {
		case "onAuthenticationSucceeded":
			finish(nil)
		case "onAuthenticationFailed":
			// Not an answer: a finger the sensor did not recognise. The
			// prompt is still up.
		case "onAuthenticationError":
			finish(promptError(e, args))
		case "onClick":
			// The negative button, which arrives here as well as through an
			// error on most versions.
			finish(ErrRefused)
		}
	})
	defer release()

	if e := app.RunOnUISync(func() { err = authenticate(token, title, subtitle, cancel) }); e != nil {
		return e
	}
	if err != nil {
		return plainly(err)
	}
	return <-done
}

// promptError turns the platform's numbered reason into one of ours.
func promptError(e *jni.Env, args []jni.Object) error {
	if len(args) == 0 {
		return ErrCancelled
	}
	code, err := e.InvokeInt(args[0], "intValue", jni.Sig(jni.TInt))
	if err != nil {
		return ErrCancelled
	}
	switch code {
	case 5, 10: // ERROR_CANCELED, ERROR_USER_CANCELED
		return ErrCancelled
	case 13: // ERROR_NEGATIVE_BUTTON
		return ErrRefused
	case 7, 9: // ERROR_LOCKOUT, ERROR_LOCKOUT_PERMANENT
		return ErrLockedOut
	}
	message := ""
	if len(args) > 1 {
		message, _ = e.GoString(args[1])
	}
	if message == "" {
		message = fmt.Sprintf("error %d", code)
	}
	return fmt.Errorf("biometric: %s", message)
}

func authenticate(token int64, title, subtitle, cancel string) error {
	return jni.Do(func(e *jni.Env) error {
		builder, err := e.Make("android/hardware/biometrics/BiometricPrompt$Builder",
			jni.Sig(jni.TVoid, jni.TContext), jni.Ref(app.Context()))
		if err != nil {
			return err
		}
		bt := jni.TClass("android/hardware/biometrics/BiometricPrompt$Builder")
		cs := jni.TClass("java/lang/CharSequence")
		text := func(method, value string) error {
			if value == "" {
				return nil
			}
			js, err := e.String(value)
			if err != nil {
				return err
			}
			_, err = e.Invoke(builder, method, jni.Sig(bt, cs), jni.Ref(js))
			return err
		}
		if err := text("setTitle", title); err != nil {
			return err
		}
		if err := text("setSubtitle", subtitle); err != nil {
			return err
		}

		// The executor the platform runs its callbacks on. The main one will
		// do: everything it calls hands straight over to a channel.
		executor, err := e.Invoke(app.Context(), "getMainExecutor",
			jni.Sig(jni.TClass("java/util/concurrent/Executor")))
		if err != nil {
			return err
		}
		listener, err := app.Proxy(e, token,
			"android.content.DialogInterface$OnClickListener")
		if err != nil {
			return err
		}
		jcancel, err := e.String(cancel)
		if err != nil {
			return err
		}
		if _, err := e.Invoke(builder, "setNegativeButton",
			jni.Sig(bt, cs, jni.TClass("java/util/concurrent/Executor"),
				jni.TClass("android/content/DialogInterface$OnClickListener")),
			jni.Ref(jcancel), jni.Ref(executor), jni.Ref(listener)); err != nil {
			return err
		}
		prompt, err := e.Invoke(builder, "build",
			jni.Sig(jni.TClass("android/hardware/biometrics/BiometricPrompt")))
		if err != nil {
			return err
		}

		signal, err := e.Make("android/os/CancellationSignal", jni.Sig(jni.TVoid))
		if err != nil {
			return err
		}
		callback, err := app.Listener(e, app.AuthCallbackClass, token)
		if err != nil {
			return err
		}
		return e.InvokeVoid(prompt, "authenticate",
			jni.Sig(jni.TVoid, jni.TClass("android/os/CancellationSignal"),
				jni.TClass("java/util/concurrent/Executor"),
				jni.TClass("android/hardware/biometrics/BiometricPrompt$AuthenticationCallback")),
			jni.Ref(signal), jni.Ref(executor), jni.Ref(callback))
	})
}
