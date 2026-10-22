//go:build windows

package windows

// Window messages.
const (
	wmDestroy       = 0x0002
	wmSize          = 0x0005
	wmSetFocus      = 0x0007
	wmKillFocus     = 0x0008
	wmPaint         = 0x000F
	wmClose         = 0x0010
	wmQuit          = 0x0012
	wmEraseBkgnd    = 0x0014
	wmSettingChange = 0x001A
	wmGetMinMaxInfo = 0x0024
	wmSetCursor     = 0x0020
	wmNCCreate      = 0x0081
	wmKeyDown       = 0x0100
	wmKeyUp         = 0x0101
	wmChar          = 0x0102
	wmSysKeyDown    = 0x0104
	wmSysKeyUp      = 0x0105
	wmSysChar       = 0x0106
	wmMouseMove     = 0x0200
	wmLButtonDown   = 0x0201
	wmLButtonUp     = 0x0202
	wmRButtonDown   = 0x0204
	wmRButtonUp     = 0x0205
	wmMButtonDown   = 0x0207
	wmMButtonUp     = 0x0208
	wmMouseWheel    = 0x020A
	wmMouseLeave    = 0x02A3
	wmThemeChanged  = 0x031A
)
