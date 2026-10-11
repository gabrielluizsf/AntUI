//go:build windows

package antui

import "unsafe"

// defaultDevice returns the default render endpoint.
func (w *wasapi) defaultDevice() (uintptr, error) {
	var enumerator uintptr
	hr, _, _ := procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&clsidMMDeviceEnumerator)), 0, clsctxAll,
		uintptr(unsafe.Pointer(&iidIMMDeviceEnumerator)),
		uintptr(unsafe.Pointer(&enumerator)))
	if err := failed(hr, "CoCreateInstance"); err != nil {
		return 0, err
	}
	defer comRelease(enumerator)

	// IMMDeviceEnumerator::GetDefaultAudioEndpoint is the fourth method.
	var device uintptr
	if err := failed(comCall(enumerator, 4, eRender, eConsole,
		uintptr(unsafe.Pointer(&device))), "GetDefaultAudioEndpoint"); err != nil {
		return 0, err
	}
	return device, nil
}
