package apk

// Package apk builds an Android package out of a Go program.
//
// It runs on the machine doing the building, not on the phone: no cgo, no
// build tag, no device. What it needs is an SDK — see
// [antui/backend/android/sdk] — and a Go package that registers an app with
// [antui/backend/android/app].
//
// The pipeline is four steps, and each one is a separate function so that a
// caller who wants only part of it can have it:
//
//	BuildLibs   compile the Go package once per ABI, with the NDK's clang
//	Link        turn a manifest into a base APK, with aapt2
//	Assemble    put the libraries in the APK, uncompressed
//	Sign        align it and sign it
//
// [Build] is all four.
