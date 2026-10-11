//go:build windows

package antui

const (
	coinitMultithreaded = 0
	clsctxAll           = 23

	eRender  = 0 // the direction: out of the machine
	eConsole = 0 // and what for: games and everything else that is not a call

	shareModeShared     = 0
	streamEventCallback = 0x00040000

	waitObject0 = 0
	waitTimeout = 0x00000102

	formatFloat      = 3      // WAVE_FORMAT_IEEE_FLOAT
	formatPCM        = 1      // WAVE_FORMAT_PCM
	formatExtensible = 0xFFFE // and the tag that says look at the subformat

	// The one error worth telling apart: it is what a machine with no sound
	// card at all answers with.
	errNotFound          = 0x80070490 // ERROR_NOT_FOUND, as an HRESULT
	errDeviceInvalidated = 0x88890004
)
