//go:build windows

package antui

import "runtime"

func (w *wasapi) run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer close(w.done)

	hr, _, _ := procCoInitializeEx.Call(0, coinitMultithreaded)
	if err := failed(hr, "CoInitializeEx"); err != nil {
		w.ready <- err
		return
	}
	defer procCoUninitialize.Call()

	client, render, frames, format, err := w.start()
	if err != nil {
		w.ready <- err
		return
	}
	defer comRelease(render)
	defer comRelease(client)
	defer func() {
		if w.event != 0 {
			procCloseHandle.Call(w.event)
		}
	}()
	defer comCall(client, 11) // IAudioClient::Stop

	w.ready <- nil
	w.feed(client, render, frames, format)
}
