package dbus

// notifySignature is the Notify call's arguments: app_name, replaces_id,
// app_icon, summary, body, actions, hints and expire_timeout.
const notifySignature = "susssasa{sv}i"

// notifyBody encodes a Notify call's arguments: a new notification (no id to
// replace), no icon, no actions, no hints and the server's own timeout.
func notifyBody(appName, summary, body string) []byte {
	var e encoder
	e.string(appName)
	e.uint32(0)
	e.string("")
	e.string(summary)
	e.string(body)
	e.emptyArray(4) // as
	e.emptyArray(8) // a{sv}: a dict entry aligns like a struct
	e.int32(-1)
	return e.buf
}
