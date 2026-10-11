//go:build windows

package windows

const (
	wsThickFrame  = 0x00040000 // the frame is a resize handle
	wsMaximizeBox = 0x00010000

	wmSizing = 0x0214

	// Which edge of the rectangle is being dragged, from WM_SIZING's wparam.
	wmszLeft        = 1
	wmszRight       = 2
	wmszTop         = 3
	wmszBottom      = 4
	wmszTopLeft     = 5
	wmszTopRight    = 6
	wmszBottomLeft  = 7
	wmszBottomRight = 8

	swpNoZOrder       = 0x0004
	swpNoActivate     = 0x0010
	logPixelsX        = 88
	defaultWindowsDPI = 96
)
