//go:build android

package network

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Transport is how the device is connected.
type Transport int

// The transports, numbered as NetworkCapabilities numbers them.
const (
	Cellular  Transport = 0
	WiFi      Transport = 1
	Bluetooth Transport = 2
	Ethernet  Transport = 3
	VPN       Transport = 4
	None      Transport = -1
)

// String names how the device is connected: wifi, mobile, ethernet.
func (t Transport) String() string {
	switch t {
	case Cellular:
		return "cellular"
	case WiFi:
		return "wifi"
	case Bluetooth:
		return "bluetooth"
	case Ethernet:
		return "ethernet"
	case VPN:
		return "vpn"
	}
	return "none"
}

// The capability numbers this package asks about.
const (
	capNotMetered = 11
	capInternet   = 12
	capValidated  = 16
)

// State is what the device's connection is like now.
type State struct {
	// Online is whether there is a network that reaches the internet.
	Online bool
	// Validated is stronger: the system has actually checked that the
	// internet is reachable, rather than merely being attached to something.
	// A captive portal in a hotel is Online and not Validated.
	Validated bool
	// Metered is a connection the user pays for by the byte. An app that
	// downloads anything large should ask, and should not do it here.
	Metered bool
	// Transport is what it is.
	Transport Transport
}

// Read asks about the current connection.
//
// It needs android.permission.ACCESS_NETWORK_STATE in the manifest, which is
// a normal permission — granted at install, with nothing to ask the user.
// Without it every call here reports a SecurityException.
func Read() (State, error) {
	out := State{Transport: None}
	a := app.Current()
	if a == nil {
		return out, app.ErrNoUIThread
	}
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceConnectivity)
		if err != nil {
			return err
		}
		if a.Config().SDK < 23 {
			return readOld(e, manager, &out)
		}
		network, err := e.Invoke(manager, "getActiveNetwork",
			jni.Sig(jni.TClass("android/net/Network")))
		if err != nil || network.IsNil() {
			return err
		}
		caps, err := e.Invoke(manager, "getNetworkCapabilities",
			jni.Sig(jni.TClass("android/net/NetworkCapabilities"),
				jni.TClass("android/net/Network")), jni.Ref(network))
		if err != nil || caps.IsNil() {
			return err
		}
		has := func(name string, n int32) (bool, error) {
			return e.InvokeBool(caps, name, jni.Sig(jni.TBool, jni.TInt), jni.Int(n))
		}
		if out.Online, err = has("hasCapability", capInternet); err != nil {
			return err
		}
		if out.Validated, err = has("hasCapability", capValidated); err != nil {
			return err
		}
		unmetered, err := has("hasCapability", capNotMetered)
		if err != nil {
			return err
		}
		out.Metered = !unmetered
		for _, t := range []Transport{WiFi, Cellular, Ethernet, VPN, Bluetooth} {
			on, err := has("hasTransport", int32(t))
			if err != nil {
				return err
			}
			if on {
				out.Transport = t
				break
			}
		}
		return nil
	})
	return out, err
}

// readOld is the path for API 21 and 22, where NetworkCapabilities does not
// exist and NetworkInfo says much less.
func readOld(e *jni.Env, manager jni.Object, out *State) error {
	info, err := e.Invoke(manager, "getActiveNetworkInfo",
		jni.Sig(jni.TClass("android/net/NetworkInfo")))
	if err != nil || info.IsNil() {
		return err
	}
	if out.Online, err = e.InvokeBool(info, "isConnected", jni.Sig(jni.TBool)); err != nil {
		return err
	}
	// No way to check, so the honest answer is the same as Online rather
	// than a claim the platform did not make.
	out.Validated = out.Online
	if out.Metered, err = e.InvokeBool(manager, "isActiveNetworkMetered", jni.Sig(jni.TBool)); err != nil {
		return err
	}
	kind, err := e.InvokeInt(info, "getType", jni.Sig(jni.TInt))
	if err != nil {
		return err
	}
	// ConnectivityManager's old type numbers are not NetworkCapabilities'.
	switch kind {
	case 0:
		out.Transport = Cellular
	case 1:
		out.Transport = WiFi
	case 7:
		out.Transport = Bluetooth
	case 9:
		out.Transport = Ethernet
	case 17:
		out.Transport = VPN
	}
	return nil
}

// Online is Read().Online, for the one question most callers have.
func Online() bool {
	s, err := Read()
	return err == nil && s.Online
}

// Watch calls f whenever the connection changes — a network appearing, the
// one in use going away, or what it can do changing, which is how a
// connection going from metered to unmetered is noticed.
//
// f runs on a thread of the system's choosing, not the app's, and is given
// the whole state rather than the difference, because the difference is
// almost never what a caller wants. Call stop when done; a callback left
// registered outlives the app's interest in it and goes on being called.
func Watch(f func(State)) (stop func(), err error) {
	token, release := app.Listen(func(e *jni.Env, method string, args []jni.Object) {
		// Every one of the three means the same thing to a caller: something
		// changed, here is how it is now. Reading it again costs a few calls
		// and avoids having to piece the state together from three
		// different callbacks that each carry part of it.
		if s, err := Read(); err == nil {
			f(s)
		}
	})

	var callback jni.Object
	err = jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServiceConnectivity)
		if err != nil {
			return err
		}
		cb, err := app.Listener(e, app.NetworkCallbackClass, token)
		if err != nil {
			return err
		}
		// A global reference, and not only so it can be unregistered: the
		// only thing holding this object is the platform, and a local
		// reference would be gone before the first callback.
		callback = e.Global(cb)

		builder, err := e.Make("android/net/NetworkRequest$Builder", jni.Sig(jni.TVoid))
		if err != nil {
			return err
		}
		if _, err := e.Invoke(builder, "addCapability",
			jni.Sig(jni.TClass("android/net/NetworkRequest$Builder"), jni.TInt),
			jni.Int(capInternet)); err != nil {
			return err
		}
		request, err := e.Invoke(builder, "build",
			jni.Sig(jni.TClass("android/net/NetworkRequest")))
		if err != nil {
			return err
		}
		return e.InvokeVoid(manager, "registerNetworkCallback",
			jni.Sig(jni.TVoid, jni.TClass("android/net/NetworkRequest"),
				jni.TClass("android/net/ConnectivityManager$NetworkCallback")),
			jni.Ref(request), jni.Ref(callback))
	})
	if err != nil {
		release()
		return func() {}, err
	}

	var once sync.Once
	return func() {
		once.Do(func() {
			jni.Do(func(e *jni.Env) error {
				manager, err := app.Service(e, app.ServiceConnectivity)
				if err == nil {
					e.InvokeVoid(manager, "unregisterNetworkCallback",
						jni.Sig(jni.TVoid,
							jni.TClass("android/net/ConnectivityManager$NetworkCallback")),
						jni.Ref(callback))
				}
				e.DeleteGlobal(callback)
				return nil
			})
			release()
		})
	}, nil
}
