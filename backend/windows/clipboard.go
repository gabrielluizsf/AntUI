//go:build windows

package windows

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
