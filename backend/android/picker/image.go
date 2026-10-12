//go:build android

package picker

// Image asks the user to choose a picture, and gives back its address.
//
// **No permission is needed, and none should be asked for.** The user
// choosing a file *is* the permission: what comes back is a content://
// address this app may read and nothing else. Asking for READ_EXTERNAL_STORAGE
// to open a picker is asking for the whole library in order to be handed one
// photo, and the store now refuses apps that do it.
//
// The address is not a path. It is opened with
// [antui/backend/android/app.Context]'s content resolver; there is no file behind it
// on many devices, because the picture may live in a cloud account.
//
// **Not from the frame loop.** The picker is another app, in front, for as
// long as the user wants it there.
func Image() (string, error) { return one("image/*") }
