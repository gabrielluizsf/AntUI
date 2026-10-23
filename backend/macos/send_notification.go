//go:build darwin

package macos

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo LDFLAGS: -framework Cocoa -framework UserNotifications
#include <stdlib.h>
#include "backend_darwin.h"
*/
import "C"
import "unsafe"

// SendNotification shows a notification in the notification centre,
// reporting whether it could be requested. A bundled app goes through
// UserNotifications and is asked for permission the first time; a bare
// binary — go run, a test — has no bundle for that framework to accept, and
// goes through osascript instead.
func (d *Driver) SendNotification(title, body string) bool {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	cBody := C.CString(body)
	defer C.free(unsafe.Pointer(cBody))
	switch C.antui_d_send_notification(cTitle, cBody) {
	case 1:
		return true
	case -1:
		return osascriptNotify(title, body) == nil
	}
	return false
}
