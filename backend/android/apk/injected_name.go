package apk

// injectedName is what the file is called inside the package being built.
// Nothing is ever written there — the overlay only makes the go command
// believe it is — but the name is chosen to be one nobody would use.
const injectedName = "zz_antuiapk_android_main.go"
