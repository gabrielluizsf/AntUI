//go:build android

package notify

// The PendingIntent flags. IMMUTABLE has been required since API 31 and is
// harmless before it; without it the call throws on a modern device.
const (
	flagUpdateCurrent int32 = 1 << 27
	flagImmutable     int32 = 1 << 26
	// Bring the activity to the front rather than starting a second copy.
	flagSingleTop int32 = 0x20000000
	flagClearTop  int32 = 0x04000000
)
