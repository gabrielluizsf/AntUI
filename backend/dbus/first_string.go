package dbus

// firstString is the body's first argument when that is a string — the
// human-readable half of an error — or "" when it is not.
func (m message) firstString() string {
	if len(m.signature) == 0 || m.signature[0] != 's' || len(m.body) < 4 {
		return ""
	}
	n := int(m.order.Uint32(m.body))
	if n < 0 || 4+n > len(m.body) {
		return ""
	}
	return string(m.body[4 : 4+n])
}
