package apk

// SignAPKParams holds the parameters for signing a standalone APK.
type SignAPKParams struct {
	Keystore Keystore
	MinSDK   int
	In       string
	Out      string
}
