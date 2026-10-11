//go:build windows

package antui

import "unsafe"

// activate returns the audio client for a device. IMMDevice::Activate is the
// fourth method.
func (w *wasapi) activate(device uintptr) (uintptr, error) {
	var client uintptr
	if err := failed(comCall(device, 3,
		uintptr(unsafe.Pointer(&iidIAudioClient)), clsctxAll, 0,
		uintptr(unsafe.Pointer(&client))), "Activate"); err != nil {
		return 0, err
	}
	return client, nil
}
