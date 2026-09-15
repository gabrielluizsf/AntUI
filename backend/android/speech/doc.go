// Package speech turns text into sound and sound into text.
//
// Both halves are services rather than parts of Android: a device without
// them has neither, and [Available] and [ErrNoEngine] are how that is found
// out rather than guessed at. Recognition in particular usually sends what
// was said to a server, which is worth telling the user about in the app and
// not only in a policy nobody reads.
//
// # The manifest
//
// Both halves are other apps, and since Android 11 an app cannot see another
// unless its manifest says which it is looking for. Without these, a device
// with an engine installed behaves exactly like one without:
//
//	antuiapk build -queries android.intent.action.TTS_SERVICE,android.speech.RecognitionService
//
// Listening also needs android.permission.RECORD_AUDIO.
package speech
