//go:build android

package intent

// grantRead is Intent.FLAG_GRANT_READ_URI_PERMISSION: the receiving app may
// read the address this intent carries, and only that address, and only
// until the app that sent it goes away.
const grantRead int32 = 0x00000001
