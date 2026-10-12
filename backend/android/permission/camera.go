//go:build android

package permission

// The permissions an app built with this library is most likely to want,
// spelled the way the manifest and the request both take them.
//
// Only the dangerous ones are here — the ones the user is asked about. A
// normal permission like INTERNET is granted at install by being in the
// manifest, and there is nothing to request.
const (
	Camera             = "android.permission.CAMERA"
	Microphone         = "android.permission.RECORD_AUDIO"
	FineLocation       = "android.permission.ACCESS_FINE_LOCATION"
	CoarseLocation     = "android.permission.ACCESS_COARSE_LOCATION"
	BackgroundLocation = "android.permission.ACCESS_BACKGROUND_LOCATION"
	ReadStorage        = "android.permission.READ_EXTERNAL_STORAGE"
	WriteStorage       = "android.permission.WRITE_EXTERNAL_STORAGE"
	ReadImages         = "android.permission.READ_MEDIA_IMAGES"
	ReadVideo          = "android.permission.READ_MEDIA_VIDEO"
	ReadAudio          = "android.permission.READ_MEDIA_AUDIO"
	Notifications      = "android.permission.POST_NOTIFICATIONS"
	Contacts           = "android.permission.READ_CONTACTS"
	Calendar           = "android.permission.READ_CALENDAR"
	BodySensors        = "android.permission.BODY_SENSORS"
	Activity           = "android.permission.ACTIVITY_RECOGNITION"
	NearbyDevices      = "android.permission.BLUETOOTH_CONNECT"
)
