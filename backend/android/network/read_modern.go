//go:build android

package network

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// readModern reads the connection from the active network's capabilities,
// which is the API 23 way.
func readModern(e *jni.Env, manager jni.Object, out *State) error {
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
}
