//go:build windows

package windows

// hkeyCurrentUser is HKEY_CURRENT_USER: 0x80000001 with the rest of the
// machine word filled in, which is how the windows.h header gets there — it
// casts a LONG — and what the registry expects to be handed. keyRead is the
// access the theme value is read with, and regDword the type it is kept in.
const (
	hkeyCurrentUser = ^uintptr(0x7FFFFFFE) // 0x80000001 sign-extended
	keyRead         = 0x20019
	regDword        = 4
)
