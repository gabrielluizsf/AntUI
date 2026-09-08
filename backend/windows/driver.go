//go:build windows

// Package windows drives windows on Win32 (user32 + gdi32), driven through syscall.
//
// The window core of AntUI lives elsewhere; this package is the platform
// half. It receives the window to drive as a [backend.Face], which is the
// half of Window a platform is allowed to touch, and drives the display with
// the exported calls below.
//
// The Windows backend talks to user32 and gdi32 through syscall, so there is
// no cgo here and no C toolchain needed to build. The frame is handed to GDI
// as a top-down 32-bit DIB, which is exactly the layout a Canvas already
// holds — no conversion per frame.
package windows

import (
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	gdi32  = syscall.NewLazyDLL("gdi32.dll")
	kernel = syscall.NewLazyDLL("kernel32.dll")

	procRegisterClassExW     = user32.NewProc("RegisterClassExW")
	procCreateWindowExW      = user32.NewProc("CreateWindowExW")
	procDestroyWindow        = user32.NewProc("DestroyWindow")
	procDefWindowProcW       = user32.NewProc("DefWindowProcW")
	procShowWindow           = user32.NewProc("ShowWindow")
	procUpdateWindow         = user32.NewProc("UpdateWindow")
	procSetForegroundWindow  = user32.NewProc("SetForegroundWindow")
	procPeekMessageW         = user32.NewProc("PeekMessageW")
	procTranslateMessage     = user32.NewProc("TranslateMessage")
	procDispatchMessageW     = user32.NewProc("DispatchMessageW")
	procSetWindowTextW       = user32.NewProc("SetWindowTextW")
	procBeginPaint           = user32.NewProc("BeginPaint")
	procEndPaint             = user32.NewProc("EndPaint")
	procGetDC                = user32.NewProc("GetDC")
	procReleaseDC            = user32.NewProc("ReleaseDC")
	procAdjustWindowRect     = user32.NewProc("AdjustWindowRect")
	procSetWindowPos         = user32.NewProc("SetWindowPos")
	procGetWindowLongW       = user32.NewProc("GetWindowLongW")
	procSetWindowLongW       = user32.NewProc("SetWindowLongW")
	procMonitorFromWindow    = user32.NewProc("MonitorFromWindow")
	procGetMonitorInfoW      = user32.NewProc("GetMonitorInfoW")
	procGetSystemMetrics     = user32.NewProc("GetSystemMetrics")
	procEnumDisplaySettingsW = user32.NewProc("EnumDisplaySettingsW")
	procTrackMouseEvent      = user32.NewProc("TrackMouseEvent")
	procSetCapture           = user32.NewProc("SetCapture")
	procReleaseCapture       = user32.NewProc("ReleaseCapture")
	procGetKeyState          = user32.NewProc("GetKeyState")
	procLoadCursorW          = user32.NewProc("LoadCursorW")
	procSetProcessDPIAware   = user32.NewProc("SetProcessDPIAware")
	procGetDpiForWindow      = user32.NewProc("GetDpiForWindow")
	procStretchDIBits        = gdi32.NewProc("StretchDIBits")
	procGetDeviceCaps        = gdi32.NewProc("GetDeviceCaps")
	procGetModuleHandleW     = kernel.NewProc("GetModuleHandleW")
)

// 64-bit Windows renames the window-long calls; 32-bit only has the short
// ones. Resolved once, so the rest of the file needs to know neither.
var (
	procGetWindowLongPtr = pickWindowLong("GetWindowLongPtrW", procGetWindowLongW)
	procSetWindowLongPtr = pickWindowLong("SetWindowLongPtrW", procSetWindowLongW)
)

func pickWindowLong(name string, fallback *syscall.LazyProc) *syscall.LazyProc {
	proc := user32.NewProc(name)
	if proc.Find() == nil {
		return proc
	}
	return fallback
}

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
	wmGetMinMaxInfo = 0x0024
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
)

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

	gwlStyle = ^uintptr(15) // GWL_STYLE, which is -16

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

// Virtual key codes.
const (
	vkBack    = 0x08
	vkTab     = 0x09
	vkReturn  = 0x0D
	vkShift   = 0x10
	vkControl = 0x11
	vkMenu    = 0x12
	vkCapital = 0x14
	vkEscape  = 0x1B
	vkSpace   = 0x20
	vkPrior   = 0x21
	vkNext    = 0x22
	vkEnd     = 0x23
	vkHome    = 0x24
	vkLeft    = 0x25
	vkUp      = 0x26
	vkRight   = 0x27
	vkDown    = 0x28
	vkInsert  = 0x2D
	vkDelete  = 0x2E
	vkLWin    = 0x5B
	vkRWin    = 0x5C
	vkNumpad0 = 0x60
	vkNumpad9 = 0x69
	vkF1      = 0x70
	vkF4      = 0x73
	vkF12     = 0x7B
	vkLShift  = 0xA0
	vkRShift  = 0xA1
	vkLCtrl   = 0xA2
	vkRCtrl   = 0xA3
	vkLMenu   = 0xA4
	vkRMenu   = 0xA5

	vkOEMPlus   = 0xBB
	vkOEMComma  = 0xBC
	vkOEMMinus  = 0xBD
	vkOEMPeriod = 0xBE
)

type point struct{ X, Y int32 }

type rect struct{ Left, Top, Right, Bottom int32 }

type msg struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
}

type wndClassExW struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}

type paintStruct struct {
	HDC         uintptr
	Erase       int32
	RcPaint     rect
	Restore     int32
	IncUpdate   int32
	RgbReserved [32]byte
}

type bitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type monitorInfo struct {
	Size      uint32
	RcMonitor rect
	RcWork    rect
	Flags     uint32
}

type trackMouseEvent struct {
	Size      uint32
	Flags     uint32
	HWndTrack uintptr
	HoverTime uint32
}

type minMaxInfo struct {
	Reserved     point
	MaxSize      point
	MaxPosition  point
	MinTrackSize point
	MaxTrackSize point
}

// windowsFromHandle maps a window handle back to the Go window it belongs to.
// Keeping it here rather than in GWLP_USERDATA sidesteps the 32/64-bit split
// in the window-long calls, and keeps a Go pointer out of a C-visible slot.
var (
	windowsMu sync.RWMutex
	windows   = map[uintptr]*Driver{}
)

// Driver is one window on Windows. It owns the platform connections the way
// the window core owns the input state.
type Driver struct {
	win  backend.Face
	hwnd uintptr
	// The icons the window is showing, kept so they can be let go of when a
	// new one takes their place.
	iconBig, iconSmall uintptr
	instance           uintptr

	limits        backend.Limits
	info          bitmapInfoHeader
	highSurrogate uint16 // the high half of a pending UTF-16 pair
	trackingMouse bool
	classAtom     uintptr

	fullscreen                    bool
	windowedWidth, windowedHeight int // the size to come back out to
	mouseX, mouseY                int // the last position the pointer reported
}

var (
	classOnce sync.Once
	classErr  error
	wndProcCB = syscall.NewCallback(wndProc)
)

const className = "AntUIWindow"

func utf16Ptr(s string) *uint16 {
	p, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		// A string with an interior NUL cannot be a title; an empty one can.
		p, _ = syscall.UTF16PtrFromString("")
	}
	return p
}

// registerClass registers the window class once for the whole process. Two
// windows in one program share it, which is what the class is for.
func registerClass(instance uintptr) error {
	classOnce.Do(func() {
		cursor, _, _ := procLoadCursorW.Call(0, 32512) // IDC_ARROW
		wc := wndClassExW{
			Size:      uint32(unsafe.Sizeof(wndClassExW{})),
			Style:     csHRedraw | csVRedraw | csOwnDC,
			WndProc:   wndProcCB,
			Instance:  instance,
			Cursor:    cursor,
			ClassName: utf16Ptr(className),
		}
		atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
		if atom == 0 {
			if errno, ok := err.(syscall.Errno); !ok || uintptr(errno) != errClassAlreadyExists {
				classErr = fmt.Errorf("antui: could not register the window class: %w", err)
			}
		}
	})
	return classErr
}

// Open creates the window that Options describes.
//
// Windows delivers messages to the thread that created the window, so the
// frame loop has to stay on it. Callers must drive the window from the
// goroutine that opened it.
func Open(win backend.Face, opts backend.Options) (*Driver, error) {
	if win == nil {
		return nil, fmt.Errorf("antui: windows backend with no window to drive")
	}
	runtime.LockOSThread()

	d := &Driver{
		win:            win,
		windowedWidth:  opts.Width,
		windowedHeight: opts.Height,
	}
	d.instance, _, _ = procGetModuleHandleW.Call(0)

	// Keeps the window sharp on scaled displays, with no manifest needed.
	if procSetProcessDPIAware.Find() == nil {
		procSetProcessDPIAware.Call()
	}

	if err := registerClass(d.instance); err != nil {
		runtime.UnlockOSThread()
		return nil, err
	}

	// The size asked for is the drawable area, so the frame is added on top.
	r := rect{0, 0, int32(opts.Width), int32(opts.Height)}
	procAdjustWindowRect.Call(uintptr(unsafe.Pointer(&r)), wsOverlappedWindow, 0)

	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr(className))),
		uintptr(unsafe.Pointer(utf16Ptr(opts.Title))),
		wsOverlappedWindow,
		cwUseDefault, cwUseDefault,
		uintptr(r.Right-r.Left), uintptr(r.Bottom-r.Top),
		0, 0, d.instance, 0,
	)
	if hwnd == 0 {
		runtime.UnlockOSThread()
		return nil, fmt.Errorf("antui: could not create the window: %w", err)
	}
	d.hwnd = hwnd

	windowsMu.Lock()
	windows[hwnd] = d
	windowsMu.Unlock()

	// Say the window takes a drag. Without it the pointer shows the no-entry
	// sign over the window and the message never comes.
	d.acceptDrops()

	procShowWindow.Call(hwnd, swShowNormal)
	procUpdateWindow.Call(hwnd)
	procSetForegroundWindow.Call(hwnd)
	return d, nil
}

// Close destroys the window and lets go of everything the platform held.
func (d *Driver) Close() {
	if d.hwnd == 0 {
		return
	}
	windowsMu.Lock()
	delete(windows, d.hwnd)
	windowsMu.Unlock()

	procDestroyWindow.Call(d.hwnd)
	d.hwnd = 0
	runtime.UnlockOSThread()
}

// Pump handles whatever happened this frame, folding it into the window
// through win.Push.
func (d *Driver) Pump(win backend.Face) {
	var m msg
	for {
		got, _, _ := procPeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, pmRemove)
		if got == 0 {
			return
		}
		if m.Message == wmQuit {
			win.SetShouldClose()
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

// SetTitle names the window in the system's decoration.
func (d *Driver) SetTitle(title string) {
	if d.hwnd != 0 {
		procSetWindowTextW.Call(d.hwnd, uintptr(unsafe.Pointer(utf16Ptr(title))))
	}
}

// blit hands the whole canvas to GDI as a top-down 32-bit DIB. The negative
// height is what "top-down" means to Windows; without it the frame arrives
// upside down.
func (d *Driver) blit(hdc uintptr) {
	cv := d.win.Canvas()
	if cv.Pixels == nil || cv.Width <= 0 || cv.Height <= 0 {
		return
	}
	d.info = bitmapInfoHeader{
		Size:        uint32(unsafe.Sizeof(bitmapInfoHeader{})),
		Width:       int32(cv.Width),
		Height:      -int32(cv.Height),
		Planes:      1,
		BitCount:    32,
		Compression: biRGB,
	}
	procStretchDIBits.Call(
		hdc,
		0, 0, uintptr(cv.Width), uintptr(cv.Height),
		0, 0, uintptr(cv.Width), uintptr(cv.Height),
		uintptr(unsafe.Pointer(&cv.Pixels[0])),
		uintptr(unsafe.Pointer(&d.info)),
		dibRGBColor, srcCopy,
	)
}

// Present copies the window's dirty rectangle to the platform surface.
//
// GDI is not told the rectangle is dirty; it is handed the whole frame the
// way a game hands one over, and the dirty rectangle decides on the other
// side whether the frame is even copied. Which is this call as it is.
func (d *Driver) Present(win backend.Face, dirty canvas.Area) {
	if dirty.Width <= 0 || dirty.Height <= 0 || d.hwnd == 0 {
		return
	}
	hdc, _, _ := procGetDC.Call(d.hwnd)
	if hdc == 0 {
		return
	}
	d.blit(hdc)
	procReleaseDC.Call(d.hwnd, hdc)
}

// SetFullscreen takes the border off and covers the monitor.
//
// Windows has no full-screen state to ask for: full screen is a window with
// no border covering the monitor. So the border comes off, the window moves
// over the monitor it is mostly on — not the primary one, which is the bug
// every second implementation of this has — and coming back out puts the
// style and the placement back.
func (d *Driver) SetFullscreen(on bool) bool {
	if d.hwnd == 0 {
		return false
	}
	style, _, _ := procGetWindowLongPtr.Call(d.hwnd, gwlStyle)

	if on {
		// Remember the size to come back out to. The core has recorded it in
		// the canvas, which is about to change; this is the one time it is
		// still the windowed size.
		if cv := d.win.Canvas(); cv != nil && cv.Width > 0 {
			d.windowedWidth, d.windowedHeight = cv.Width, cv.Height
		}

		monitor, _, _ := procMonitorFromWindow.Call(d.hwnd, monitorNearest)
		info := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
		if ok, _, _ := procGetMonitorInfoW.Call(monitor, uintptr(unsafe.Pointer(&info))); ok == 0 {
			return false
		}
		procSetWindowLongPtr.Call(d.hwnd, gwlStyle,
			style&^uintptr(wsOverlappedWindow)|wsPopup)
		procSetWindowPos.Call(d.hwnd, hwndTop,
			uintptr(info.RcMonitor.Left), uintptr(info.RcMonitor.Top),
			uintptr(info.RcMonitor.Right-info.RcMonitor.Left),
			uintptr(info.RcMonitor.Bottom-info.RcMonitor.Top),
			swpNoOwnerZOrder|swpFrameChanged)
		d.fullscreen = true
		return true
	}

	width, height := d.windowedWidth, d.windowedHeight
	if width <= 0 {
		width = 640
	}
	if height <= 0 {
		height = 480
	}
	procSetWindowLongPtr.Call(d.hwnd, gwlStyle,
		style&^uintptr(wsPopup)|wsOverlappedWindow)
	r := rect{0, 0, int32(width), int32(height)}
	procAdjustWindowRect.Call(uintptr(unsafe.Pointer(&r)), wsOverlappedWindow, 0)
	procSetWindowPos.Call(d.hwnd, hwndNoTopMost, cwUseDefault, cwUseDefault,
		uintptr(r.Right-r.Left), uintptr(r.Bottom-r.Top),
		swpNoMove|swpNoOwnerZOrder|swpFrameChanged)
	d.fullscreen = false
	return true
}

// DisplaySize is the size of the display the window is on.
func (d *Driver) DisplaySize() (w, h int, ok bool) {
	if d.hwnd != 0 {
		monitor, _, _ := procMonitorFromWindow.Call(d.hwnd, monitorNearest)
		info := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
		if ok, _, _ := procGetMonitorInfoW.Call(monitor, uintptr(unsafe.Pointer(&info))); ok != 0 {
			return int(info.RcMonitor.Right - info.RcMonitor.Left),
				int(info.RcMonitor.Bottom - info.RcMonitor.Top), true
		}
	}
	width, _, _ := procGetSystemMetrics.Call(smCXScreen)
	height, _, _ := procGetSystemMetrics.Call(smCYScreen)
	if width == 0 || height == 0 {
		return 0, 0, false
	}
	return int(width), int(height), true
}

// DEVMODEW is large and its layout is fixed by the ABI rather than by any
// struct this package would want to write out. Only two of its fields matter
// here, so it is a byte buffer with the two offsets named.
const (
	devModeSize         = 220
	devModeSizeOffset   = 68
	devModeFrequencyOff = 184
)

// DisplayRefresh is how often the display repaints, in hertz.
func (d *Driver) DisplayRefresh() int {
	var mode [devModeSize]byte
	*(*uint16)(unsafe.Pointer(&mode[devModeSizeOffset])) = devModeSize

	ok, _, _ := procEnumDisplaySettingsW.Call(0, enumCurrentSettings,
		uintptr(unsafe.Pointer(&mode[0])))
	if ok == 0 {
		return 0
	}
	hz := *(*uint32)(unsafe.Pointer(&mode[devModeFrequencyOff]))
	// 0 and 1 both mean "the hardware default", which is not a number.
	if hz <= 1 {
		return 0
	}
	return int(hz)
}

// ---------------------------------------------------------------------------
// The window procedure
// ---------------------------------------------------------------------------

func winMods() backend.Mod {
	var mods backend.Mod
	down := func(vk uintptr) bool {
		state, _, _ := procGetKeyState.Call(vk)
		return state&0x8000 != 0
	}
	if down(vkShift) {
		mods |= backend.ModShift
	}
	if down(vkControl) {
		mods |= backend.ModControl
	}
	if down(vkMenu) {
		mods |= backend.ModAlt
	}
	if down(vkLWin) || down(vkRWin) {
		mods |= backend.ModSuper
	}
	return mods
}

// winKey maps a virtual key code to the key it stands for. The extended bit
// in lparam is what separates the left and right control and alt keys.
func winKey(vk, lparam uintptr) backend.Key {
	extended := lparam&(1<<24) != 0

	if vk >= 'A' && vk <= 'Z' || vk >= '0' && vk <= '9' {
		return backend.Key(vk)
	}
	switch vk {
	case vkSpace:
		return backend.KeySpace
	case vkEscape:
		return backend.KeyEscape
	case vkReturn:
		return backend.KeyEnter
	case vkTab:
		return backend.KeyTab
	case vkBack:
		return backend.KeyBackspace
	case vkInsert:
		return backend.KeyInsert
	case vkDelete:
		return backend.KeyDelete
	case vkRight:
		return backend.KeyRight
	case vkLeft:
		return backend.KeyLeft
	case vkDown:
		return backend.KeyDown
	case vkUp:
		return backend.KeyUp
	case vkPrior:
		return backend.KeyPageUp
	case vkNext:
		return backend.KeyPageDown
	case vkHome:
		return backend.KeyHome
	case vkEnd:
		return backend.KeyEnd
	case vkCapital:
		return backend.KeyCapsLock
	case vkShift, vkLShift:
		return backend.KeyLeftShift
	case vkRShift:
		return backend.KeyRightShift
	case vkControl:
		if extended {
			return backend.KeyRightControl
		}
		return backend.KeyLeftControl
	case vkMenu:
		if extended {
			return backend.KeyRightAlt
		}
		return backend.KeyLeftAlt
	case vkLCtrl:
		return backend.KeyLeftControl
	case vkRCtrl:
		return backend.KeyRightControl
	case vkLMenu:
		return backend.KeyLeftAlt
	case vkRMenu:
		return backend.KeyRightAlt
	case vkLWin:
		return backend.KeyLeftSuper
	case vkRWin:
		return backend.KeyRightSuper
	case vkOEMComma:
		return backend.KeyComma
	case vkOEMPeriod:
		return backend.KeyPeriod
	case vkOEMMinus:
		return backend.KeyMinus
	case vkOEMPlus:
		return backend.KeyEqual
	}
	switch {
	case vk >= vkF1 && vk <= vkF12:
		return backend.KeyF1 + backend.Key(vk-vkF1)
	case vk >= vkNumpad0 && vk <= vkNumpad9:
		return backend.Key0 + backend.Key(vk-vkNumpad0)
	}
	return backend.KeyUnknown
}

func loWord(v uintptr) int { return int(int16(uint16(v))) }
func hiWord(v uintptr) int { return int(int16(uint16(v >> 16))) }

func defWindowProc(hwnd, message, wparam, lparam uintptr) uintptr {
	ret, _, _ := procDefWindowProcW.Call(hwnd, message, wparam, lparam)
	return ret
}

func wndProc(hwnd, message, wparam, lparam uintptr) uintptr {
	windowsMu.RLock()
	d := windows[hwnd]
	windowsMu.RUnlock()

	// Messages arrive before CreateWindowExW has returned the handle we key
	// the map on, so there is a window of time with nothing to dispatch to.
	if d == nil || d.win == nil {
		return defWindowProc(hwnd, message, wparam, lparam)
	}
	win := d.win

	pushMouse := func(t backend.EventType, button backend.MouseButton) uintptr {
		win.Push(backend.Event{
			Type:   t,
			Button: button,
			X:      loWord(lparam),
			Y:      hiWord(lparam),
			Mods:   winMods(),
		})
		return 0
	}

	switch message {
	case wmDropFiles:
		d.dropped(wparam)
		return 0

	case wmClose:
		win.PushSimple(backend.EventClose)
		return 0

	case wmDestroy:
		win.SetShouldClose()
		return 0

	case wmEraseBkgnd:
		return 1 // the frame covers everything, so erasing only flickers

	case wmPaint:
		var ps paintStruct
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		d.blit(hdc)
		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		win.Push(backend.Event{Type: backend.EventExpose})
		return 0

	case wmSize:
		width, height := int(uint16(lparam)), int(uint16(lparam>>16))
		cv := win.Canvas()
		if wparam != sizeMinimized && width > 0 && height > 0 &&
			(width != cv.Width || height != cv.Height) {
			if win.ResizeCanvas(width, height) {
				win.Push(backend.Event{Type: backend.EventResize, Width: width, Height: height})
			}
		}
		return 0

	case wmMouseMove:
		d.mouseX, d.mouseY = loWord(lparam), hiWord(lparam)
		// Ask to be told when the pointer leaves, so the tracking flag can be
		// reset and the request made again next time it comes back.
		if !d.trackingMouse {
			track := trackMouseEvent{
				Size:      uint32(unsafe.Sizeof(trackMouseEvent{})),
				Flags:     tmeLeave,
				HWndTrack: hwnd,
			}
			procTrackMouseEvent.Call(uintptr(unsafe.Pointer(&track)))
			d.trackingMouse = true
		}
		return pushMouse(backend.EventMouseMove, 0)

	case wmMouseLeave:
		d.trackingMouse = false
		return 0

	case wmLButtonDown:
		// Capturing means a drag that leaves the window still reports its
		// release, so a button or slider cannot get stuck held down.
		procSetCapture.Call(hwnd)
		return pushMouse(backend.EventMouseDown, backend.MouseLeft)
	case wmLButtonUp:
		procReleaseCapture.Call()
		return pushMouse(backend.EventMouseUp, backend.MouseLeft)
	case wmRButtonDown:
		return pushMouse(backend.EventMouseDown, backend.MouseRight)
	case wmRButtonUp:
		return pushMouse(backend.EventMouseUp, backend.MouseRight)
	case wmMButtonDown:
		return pushMouse(backend.EventMouseDown, backend.MouseMiddle)
	case wmMButtonUp:
		return pushMouse(backend.EventMouseUp, backend.MouseMiddle)

	case wmMouseWheel:
		win.Push(backend.Event{
			Type:  backend.EventMouseWheel,
			Wheel: int(int16(uint16(wparam>>16))) / wheelDelta,
			X:     d.mouseX,
			Y:     d.mouseY,
			Mods:  winMods(),
		})
		return 0

	case wmKeyDown, wmSysKeyDown:
		mods := winMods()
		win.Push(backend.Event{
			Type:   backend.EventKeyDown,
			Key:    winKey(wparam, lparam),
			Mods:   mods,
			Repeat: lparam&0x40000000 != 0,
		})
		// Alt+F4 has to reach the system, or the window cannot be closed the
		// way every other window on the desktop closes.
		if message == wmSysKeyDown && wparam == vkF4 && mods&backend.ModAlt != 0 {
			break
		}
		return 0

	case wmKeyUp, wmSysKeyUp:
		win.Push(backend.Event{
			Type: backend.EventKeyUp,
			Key:  winKey(wparam, lparam),
			Mods: winMods(),
		})
		return 0

	case wmChar, wmSysChar:
		return d.handleChar(uint16(wparam))

	case wmSetFocus, wmKillFocus:
		win.Push(backend.Event{Type: backend.EventFocus, Focused: message == wmSetFocus})
		return 0

	case wmGetMinMaxInfo:
		// The lparam of this message is a pointer the window manager owns and
		// expects to be written through — that is the whole message. Turning
		// it back into one is the only way to answer, so `go vet` flags this
		// line as a possible misuse of unsafe.Pointer and is wrong to: the
		// value came from the OS as a pointer and is live for this call.
		mmi := (*minMaxInfo)(unsafe.Pointer(lparam)) //lint:ignore unsafeptr the OS handed us this pointer
		d.minMax(mmi)
		return 0

	case wmSizing:
		r := (*rect)(unsafe.Pointer(lparam)) //lint:ignore unsafeptr the OS handed us this pointer
		d.sizing(wparam, r)
		return 1
	}

	return defWindowProc(hwnd, message, wparam, lparam)
}

// handleChar turns a UTF-16 code unit into a text event, joining the two
// halves of a surrogate pair into the one character they stand for.
func (d *Driver) handleChar(unit uint16) uintptr {
	var r rune
	switch {
	case unit >= 0xD800 && unit <= 0xDBFF: // the high half
		d.highSurrogate = unit
		return 0
	case unit >= 0xDC00 && unit <= 0xDFFF: // the low half
		if d.highSurrogate != 0 {
			r = 0x10000 + (rune(d.highSurrogate)-0xD800)<<10 + (rune(unit) - 0xDC00)
			d.highSurrogate = 0
		}
	default:
		r = rune(unit)
	}

	// Below 32 is a control character, and 127 is delete: neither is text.
	if r >= 32 && r != 127 {
		mods := winMods()
		if mods&backend.ModControl == 0 {
			d.win.Push(backend.Event{Type: backend.EventText, Rune: r, Text: string(r), Mods: mods})
		}
	}
	return 0
}