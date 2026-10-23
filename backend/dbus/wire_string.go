package dbus

import "strings"

// wireString is text as the bus will carry it: valid UTF-8 with no NUL in
// it. The bus drops a connection that sends anything else, so a stray byte in
// a file name would lose the whole notification rather than show one wrong
// letter of it.
func wireString(s string) string {
	return strings.ReplaceAll(strings.ToValidUTF8(s, "�"), "\x00", "")
}
