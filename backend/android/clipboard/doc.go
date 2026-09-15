// Package clipboard is the system clipboard: text in and text out.
//
// It is the phone's answer to what antui.Window.ClipboardText does on a
// desktop, and it is a separate package because the rules are different
// enough to need saying — reading is only allowed while the app has focus,
// and reading is visible to the user.
package clipboard
