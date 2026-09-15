// Package share hands a file to another app.
//
// Since Android 7 an app may not give another a path — an intent carrying a
// file:// address throws — so a file has to be offered as a content://
// address from a provider, and the receiving app is granted the right to
// read that one address and nothing else.
//
// The provider is in the library's Java shim and serves one directory, which
// is what [Dir] is. Put a copy of what is to be shared there, hand the
// address to [antui/backend/android/intent.ShareFile], and empty it afterwards.
package share
