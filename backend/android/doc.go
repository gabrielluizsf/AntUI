// Package android is AntUI on a phone: the platform underneath an Android
// app, and the tools that turn a Go program into one the Play Store takes.
//
// An Android app written with this library is an ordinary Go program. It is
// compiled as a shared library and loaded into an app process by the system,
// which is the only shape Android allows — so unlike the rest of AntUI this
// corner needs cgo and the Android NDK. Nothing else about writing the
// program changes: a [antui.Window] opens, a [canvas.Canvas] is drawn into,
// and the same code runs on a desktop.
//
// # The layers
//
// Each one may use the ones below it and none of the ones above.
//
//	android/ndk      the C APIs the NDK exposes: the native window, the
//	                 looper, the input queue, assets, sensors, logging
//	android/jni      calling Java, for the far larger part of Android that
//	                 the NDK does not expose
//	android/app      the running activity — its lifecycle, its surface, and
//	                 the context every feature package needs
//	android/<name>   one package per feature: permission, notify, sensor,
//	                 location, camera, storage, share, prefs and the rest
//	android/apk      the builder. It runs on the developer's machine, not on
//	                 the phone, and needs neither cgo nor a device.
//
// The window backend itself is not here. It is the antui platform's Android
// implementation and lives in the parent package, because that interface is
// unexported; it is a thin thing that calls into android/app.
//
// # Building
//
// Everything except android/apk is behind a build tag and compiles only for
// GOOS=android, so an ordinary build on a desktop ignores it entirely. The
// packages still exist off Android — they just have nothing in them, which
// is what keeps "go build ./..." honest.
//
// # From nothing to an app on a phone
//
// Four things have to be installed, and only the first is Go's:
//
//	go            1.27 or newer
//	Android SDK   the platform, the build-tools and the platform-tools
//	Android NDK   the compiler; a Go app is native code
//	a JDK         only for apksigner and jarsigner, which are Java programs
//
// The SDK lives wherever ANDROID_HOME says, or in the usual place for the
// system. Rather than checking that by hand, ask:
//
//	go run ./cmd/antuiapk doctor
//
// which says what is there, what is missing, and the exact sdkmanager line
// that installs each missing piece. Then:
//
//	antuiapk init myapp     # a program, an icon and a config
//	cd myapp
//	go run .                # on this machine, with no device and no SDK
//	antuiapk run .          # on a phone or an emulator
//	antuiapk bundle .       # the .aab a store takes
//
// The program `init` writes is cross-platform on purpose: it opens a
// [antui.Window], draws, and reads gestures, and none of that mentions
// Android. Being able to run it with "go run" while writing it — no device,
// no install, no two-minute build — is most of what makes writing a phone
// app bearable.
//
// # What a real app needs beyond that
//
//   - **Permissions** are declared in the config and asked for at run time.
//     [antui/backend/android/permission.Ensure] does both halves — but not from the
//     frame loop, for the reason its documentation gives.
//   - **Assets** go in with `-assets`, and are read with
//     [antui/backend/android/assets]. They are not files: they stay inside the
//     package, which is never unpacked, so nothing on disk corresponds to
//     one. A game should use multg.Files, which is already the right thing on
//     both a desktop and a phone.
//   - **The safe area** is [antui.Window.SafeArea]. An app owns every pixel
//     of the screen and the system paints its status bar, its navigation bar
//     and the notch over some of them. Anything to read or press belongs
//     inside that rectangle.
//   - **The font**. [antui.Window.UseSystemFont] draws the program in
//     whatever the phone writes its own interface in, which is the
//     difference between an app and a program that looks like a terminal.
//
// examples/android/tracker is all of that in one app: a database, settings,
// a share sheet, a notification, an icon and a bundle.
//
// # What is bound, and what is not
//
// Bound, with a package each:
//
//	app          the activity: lifecycle, surface, intents, results, paths
//	assets       what was packed into the package
//	audio        recording, and audio focus
//	biometric    whether a fingerprint can be asked for, and asking
//	camera       a preview, upright, as a canvas
//	clipboard    the clipboard, and why it refuses when not in front
//	device       what phone this is, and what Android it runs
//	dialog       the system's own alert
//	display      size, density, rotation, refresh rate
//	gallery      the user's photos and films, through MediaStore
//	intent       opening a URL, the share sheet, another app's activity
//	location     where the device is
//	media        decoding video and audio
//	network      what the device is connected by, and whether it is online
//	notify       channels, notifications, and buttons that come back
//	permission   asking, and telling "refused" from "refused for good"
//	picker       the document picker
//	power        battery, charging, saving, doze
//	prefs        settings that survive being moved to a new phone
//	sensor       the accelerometer and everything beside it
//	share        a file another app can be handed
//	speech       text to speech, and speech to text
//	sqlite       the database every Android has
//	toast        the small message
//	torch        the flash
//
// **Not bound, and each for a reason worth knowing:**
//
//   - **WebView.** It is a Java view, and this library draws into a surface
//     it owns. Putting a view over that canvas means a view hierarchy, which
//     is a different kind of app from the one this builds. It is possible
//     and it is not small.
//   - **The hardware keystore and encrypted storage.** A key that lives in
//     the phone's secure element, unlocked by a fingerprint. It is the right
//     way to keep a credential and it is a large surface to get exactly
//     right, which is a bad thing to do carelessly.
//   - **NFC and Bluetooth.** Both need the app to be in the foreground in
//     particular ways, and Bluetooth needs a permission story of its own.
//   - **Vibration and haptics.**
//   - **The predictive-back gesture**, which is Android 13's and needs the
//     activity to say in advance what going back would do.
//   - **The keyboard's insets** — the layout has to move when it comes up,
//     and [antui/backend/android/app.WindowInsets] reports them, but nothing here
//     moves anything yet.
//   - **CFF fonts, hinting, kerning and variable-font axes.** See
//     [canvas.Face] for what is read and what is not.
//
// # Publishing
//
// Building the file is perhaps a third of a release. The rest is a listing,
// a set of declarations and a review, none of which any build produces:
//
//	antuiapk publish
//
// prints all of it — the account, the 14-day wait a personal account has
// before production, the listing assets, the declarations that block a
// release, and the thing that surprises people, which is that with Play App
// Signing the key from "antuiapk keygen" is the *upload* key and Google
// re-signs with one it holds.
//
//	antuiapk check <apk|aab>
//
// answers the half a machine can: 64-bit, the target API level, whether it
// is debuggable, whether it is signed with the key every developer has.
//
// See the antuiapk command for the rest of the flags.
package android
