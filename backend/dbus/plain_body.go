package dbus

import "strings"

// markupEscaper writes the three characters the notification spec's body
// markup gives a meaning to as the entities that stand for them.
var markupEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

// plainBody is a notification's body as text to be shown as written. A server
// that announces body-markup reads `<a href>`, `<img src>` and `<b>` in a body
// as what they say, so text a program passes on — a chat message, a name
// somebody else chose — could put a link or a picture of its own into the
// notification; Windows and macOS show the same text as plain letters, and
// escaping it is what the spec asks of a client sending plain text.
func plainBody(s string) string {
	return markupEscaper.Replace(wireString(s))
}
