//go:build windows

package windows

// Styles, flags and the handful of constants the calls above want.
const (
	wsOverlappedWindow = 0x00CF0000
	wsPopup            = 0x80000000

	csHRedraw = 0x0002
	csVRedraw = 0x0001
	csOwnDC   = 0x0020

	swShowNormal = 1
	pmRemove     = 0x0001

	cwUseDefault = ^uintptr(0) - 0x7FFFFFFF // 0x80000000 as a signed int

	gwlStyle   = ^uintptr(15) // GWL_STYLE, which is -16
	gwlExStyle = ^uintptr(19) // GWL_EXSTYLE, which is -20

	wsExLayered = 0x00080000
	lwaAlpha    = 0x00000002

	swpNoOwnerZOrder = 0x0200
	swpFrameChanged  = 0x0020
	swpNoMove        = 0x0002

	hwndTop        = 0
	hwndNoTopMost  = ^uintptr(1) // (HWND)-2
	monitorNearest = 0x00000002

	smCXScreen = 0
	smCYScreen = 1

	enumCurrentSettings = ^uintptr(0) // (DWORD)-1

	tmeLeave    = 0x00000002
	wheelDelta  = 120
	biRGB       = 0
	dibRGBColor = 0
	srcCopy     = 0x00CC0020

	sizeMinimized = 1

	errClassAlreadyExists = 1410
)
