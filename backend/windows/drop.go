//go:build windows

package windows

import (
	"syscall"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// Files dragged onto the window, and files copied in Explorer.
//
// Windows says both in the same currency: a drop is an HDROP handle in
// WM_DROPFILES, and Explorer puts an HDROP on the clipboard as CF_HDROP when
// somebody copies files. So one function reads the handle and both paths use
// it.

var (
	shell32 = syscall.NewLazyDLL("shell32.dll")

	procDragAcceptFiles = shell32.NewProc("DragAcceptFiles")
	procDragQueryFileW  = shell32.NewProc("DragQueryFileW")
	procDragQueryPoint  = shell32.NewProc("DragQueryPoint")
	procDragFinish      = shell32.NewProc("DragFinish")

	procOpenClipboard    = user32.NewProc("OpenClipboard")
	procCloseClipboard   = user32.NewProc("CloseClipboard")
	procGetClipboardData = user32.NewProc("GetClipboardData")
	procGlobalLock       = kernel.NewProc("GlobalLock")
	procGlobalUnlock     = kernel.NewProc("GlobalUnlock")
)

const (
	wmDropFiles = 0x0233

	cfUnicodeText = 13
	cfHDrop       = 15

	// What DragQueryFile answers with when it is asked for file 0xFFFFFFFF:
	// how many there are rather than the name of one.
	dragQueryCount = 0xFFFFFFFF
)

// acceptDrops tells the shell this window will take a drag. Without it the
// pointer shows the no-entry sign and WM_DROPFILES never arrives.
func (d *Driver) acceptDrops() {
	if d.hwnd != 0 {
		procDragAcceptFiles.Call(d.hwnd, 1)
	}
}

// dropped reads the files out of a WM_DROPFILES and lets the shell go.
func (d *Driver) dropped(hdrop uintptr) {
	if d.win == nil || hdrop == 0 {
		return
	}
	files := filesFromDrop(hdrop)

	// Where the pointer was when it let go, in window coordinates — which is
	// what everything else here counts in.
	var point struct{ X, Y int32 }
	procDragQueryPoint.Call(hdrop, uintptr(unsafe.Pointer(&point)))
	procDragFinish.Call(hdrop)

	if len(files) > 0 {
		d.win.Push(backend.Event{Type: backend.EventDropFiles, Files: files,
			X: int(point.X), Y: int(point.Y)})
	}
}

// filesFromDrop is the paths an HDROP holds, whether it came from a drag or
// off the clipboard.
func filesFromDrop(hdrop uintptr) []string {
	count, _, _ := procDragQueryFileW.Call(hdrop, dragQueryCount, 0, 0)
	if count == 0 || count > 4096 {
		return nil
	}
	out := make([]string, 0, count)
	for i := uintptr(0); i < count; i++ {
		// Asked for the length first, because a path may be long and the old
		// 260-character limit is not one any more.
		length, _, _ := procDragQueryFileW.Call(hdrop, i, 0, 0)
		if length == 0 {
			continue
		}
		buf := make([]uint16, length+1)
		got, _, _ := procDragQueryFileW.Call(hdrop, i,
			uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if got == 0 {
			continue
		}
		out = append(out, syscall.UTF16ToString(buf[:got]))
	}
	return out
}

// Clipboard is what the clipboard holds: the files it names, and its text.
//
// The clipboard is a shared thing that another process may have open, so a
// failure to open it is an ordinary answer — nothing — rather than an error
// worth stopping for.
func (d *Driver) Clipboard() (string, []string) {
	if ok, _, _ := procOpenClipboard.Call(d.hwnd); ok == 0 {
		return "", nil
	}
	defer procCloseClipboard.Call()

	var files []string
	if handle, _, _ := procGetClipboardData.Call(cfHDrop); handle != 0 {
		// A clipboard HDROP is read the same way a dropped one is, and is
		// not finished with: it belongs to the clipboard, not to us.
		files = filesFromDrop(handle)
	}

	text := ""
	if handle, _, _ := procGetClipboardData.Call(cfUnicodeText); handle != 0 {
		if at, _, _ := procGlobalLock.Call(handle); at != 0 {
			text = utf16PtrToString(at)
			procGlobalUnlock.Call(handle)
		}
	}
	return text, files
}

// utf16PtrToString reads a zero-terminated UTF-16 string out of memory the
// clipboard owns. It stops at a length no clipboard text reaches, so a
// handle that is not what it claims cannot walk off the end of the world.
//
// The pointer comes from GlobalLock rather than from Go's heap, so it is
// walked as one — `unsafe.Slice` over the locked block, which is what the
// checker wants to see rather than arithmetic on a uintptr.
func utf16PtrToString(at uintptr) string {
	const most = 1 << 20
	block := unsafe.Slice((*uint16)(unsafe.Pointer(at)), most)
	for i, c := range block {
		if c == 0 {
			return syscall.UTF16ToString(block[:i])
		}
	}
	return syscall.UTF16ToString(block)
}

// --- the window's icon ---------------------------------------------------------

var (
	procCreateIconIndirect = user32.NewProc("CreateIconIndirect")
	procDestroyIcon        = user32.NewProc("DestroyIcon")
	procSendMessageW       = user32.NewProc("SendMessageW")
	procCreateBitmap       = gdi32.NewProc("CreateBitmap")
	procCreateDIBSection   = gdi32.NewProc("CreateDIBSection")
	procDeleteObject       = gdi32.NewProc("DeleteObject")
)

const (
	wmSetIcon = 0x0080
	iconSmall = 0
	iconBig   = 1
)

type iconInfo struct {
	icon     int32 // 1 for an icon, 0 for a cursor
	xHotspot uint32
	yHotspot uint32
	mask     uintptr
	colour   uintptr
}

// SetIcon gives the window a big icon and a small one. Windows asks for both:
// the big one is the task bar and Alt-Tab, the small one the title bar.
func (d *Driver) SetIcon(images []*canvas.Canvas) bool {
	if d.hwnd == 0 || len(images) == 0 {
		return false
	}
	// Largest first is the order they arrive in; the small icon wants
	// something near 16 and the big one something near 32.
	big := iconFor(images, 32)
	small := iconFor(images, 16)
	if big == 0 && small == 0 {
		return false
	}

	if big != 0 {
		procSendMessageW.Call(d.hwnd, wmSetIcon, iconBig, big)
		if d.iconBig != 0 {
			procDestroyIcon.Call(d.iconBig)
		}
		d.iconBig = big
	}
	if small != 0 {
		procSendMessageW.Call(d.hwnd, wmSetIcon, iconSmall, small)
		if d.iconSmall != 0 {
			procDestroyIcon.Call(d.iconSmall)
		}
		d.iconSmall = small
	}
	return true
}

// iconFor makes an HICON out of whichever picture is nearest a size.
func iconFor(images []*canvas.Canvas, want int) uintptr {
	best := images[0]
	for _, image := range images {
		if awayFrom(image.Width, want) < awayFrom(best.Width, want) {
			best = image
		}
	}
	return hicon(best)
}

// awayFrom is how far a size is from the one that is wanted.
func awayFrom(size, want int) int {
	if size < want {
		return want - size
	}
	return size - want
}

// hicon builds an icon from a picture: a 32-bit colour bitmap with the alpha
// in it, and an empty mask, which is what a modern icon is. The mask has to
// exist even though nothing reads it.
func hicon(image *canvas.Canvas) uintptr {
	if image == nil || image.Width <= 0 || image.Height <= 0 {
		return 0
	}
	// Windows wants the rows the other way up and the channels as BGRA,
	// which is what a Color already is on a little-endian machine — so only
	// the rows are turned over here.
	pixels := make([]uint32, image.Width*image.Height)
	for y := range image.Height {
		row := (image.Height - 1 - y) * image.Width
		for x := range image.Width {
			pixels[row+x] = uint32(image.At(x, y))
		}
	}

	colour, _, _ := procCreateBitmap.Call(uintptr(image.Width), uintptr(image.Height),
		1, 32, uintptr(unsafe.Pointer(&pixels[0])))
	if colour == 0 {
		return 0
	}
	defer procDeleteObject.Call(colour)

	mask, _, _ := procCreateBitmap.Call(uintptr(image.Width), uintptr(image.Height),
		1, 1, 0)
	if mask == 0 {
		return 0
	}
	defer procDeleteObject.Call(mask)

	info := iconInfo{icon: 1, mask: mask, colour: colour}
	handle, _, _ := procCreateIconIndirect.Call(uintptr(unsafe.Pointer(&info)))
	return handle
}

// --- putting something on the clipboard ----------------------------------------

var (
	procEmptyClipboard   = user32.NewProc("EmptyClipboard")
	procSetClipboardData = user32.NewProc("SetClipboardData")
	procGlobalAlloc      = kernel.NewProc("GlobalAlloc")
	procGlobalFree       = kernel.NewProc("GlobalFree")
)

// GMEM_MOVEABLE. The clipboard owns what it is given and moves it about, so
// it will not take a fixed block.
const gmemMoveable = 0x0002

// SetClipboard puts text on the clipboard as UTF-16, which is what every
// Windows program reads.
func (d *Driver) SetClipboard(text string) bool {
	if ok, _, _ := procOpenClipboard.Call(d.hwnd); ok == 0 {
		return false
	}
	defer procCloseClipboard.Call()
	procEmptyClipboard.Call()

	chars, err := syscall.UTF16FromString(text)
	if err != nil {
		return false
	}
	size := uintptr(len(chars) * 2)
	handle, _, _ := procGlobalAlloc.Call(gmemMoveable, size)
	if handle == 0 {
		return false
	}
	at, _, _ := procGlobalLock.Call(handle)
	if at == 0 {
		procGlobalFree.Call(handle)
		return false
	}
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(at)), len(chars)), chars)
	procGlobalUnlock.Call(handle)

	// Once it is handed over the clipboard owns the block; freeing it here
	// would be freeing somebody else's memory.
	if ok, _, _ := procSetClipboardData.Call(cfUnicodeText, handle); ok == 0 {
		procGlobalFree.Call(handle)
		return false
	}
	return true
}