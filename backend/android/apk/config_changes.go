package apk

// configChanges is what the activity says it handles itself. Without it,
// Android destroys and recreates the activity on every rotation, keyboard
// change and theme switch — which for a native app means tearing the surface
// down and building it back up, losing everything the app had in memory.
// A NativeActivity is expected to cope, so it says so.
const configChanges = "orientation|screenSize|smallestScreenSize|screenLayout|" +
	"keyboard|keyboardHidden|navigation|uiMode|density|fontScale|locale|layoutDirection"
