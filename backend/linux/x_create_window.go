//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// Request opcodes.
const (
	xCreateWindow           = 1
	xChangeWindowAttributes = 2
	xDestroyWindow          = 4
	xMapWindow              = 8
	xConfigureWindow        = 12
	xInternAtom             = 16
	xChangeProperty         = 18
	xGetProperty            = 20
	xSetSelectionOwner      = 22
	xConvertSelection       = 24
	xSendEvent              = 25
	xCreateGC               = 55
	xFreeGC                 = 60
	xPutImage               = 72
	xQueryExtension         = 98
	xGetKeyboardMap         = 101
)
