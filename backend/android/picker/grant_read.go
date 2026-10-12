//go:build android

package picker

// The intent flags that decide what may be done with an address and for how
// long.
const (
	grantRead  int32 = 0x00000001
	grantWrite int32 = 0x00000002
	// grantPersist is what makes a permission survive a restart. Without it
	// the address the user chose works until the app is closed and then
	// stops, which looks like the app losing the file.
	grantPersist int32 = 0x00000040
)
