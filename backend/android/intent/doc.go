// Package intent is how an app asks the rest of the device to do something:
// open a link, share some text, view a file — and how it reads what another
// app asked of it.
//
// Sharing a file needs a content:// address rather than a path, because no
// app may hand another a path any more; [antui/backend/android/share] is what makes
// one, and [ShareFile] is what sends it.
package intent
