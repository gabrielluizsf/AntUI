package sdk

// Package sdk finds the Android SDK, the NDK and a JDK on the machine doing
// the building, and says where each tool inside them is.
//
// It is host-side: it runs on Linux, macOS and Windows, needs no cgo and no
// device, and knows nothing about a running app. Everything that builds an
// APK goes through here so that "where is aapt2" is answered once.
