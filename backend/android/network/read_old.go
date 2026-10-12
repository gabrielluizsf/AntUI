//go:build android

package network

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
