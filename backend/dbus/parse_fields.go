package dbus

// parseFields reads the header fields read keeps. fields starts at offset 16
// of the message, which is what every alignment here is counted from; a field
// of a type this package does not know ends the parse, since its size cannot
// be stepped over without knowing it.
func (m *message) parseFields(fields []byte) {
	const base = 16
	i := 0
	align := func(n int) { i = (base+i+n-1)/n*n - base }
	for i < len(fields) {
		align(8)
		if i+3 > len(fields) {
			return
		}
		code, sigLen := fields[i], int(fields[i+1])
		if sigLen != 1 || i+3+sigLen > len(fields) {
			return
		}
		sig := fields[i+2]
		i += 4 // code, length, the type, its NUL
		switch sig {
		case 'u':
			align(4)
			if i+4 > len(fields) {
				return
			}
			if code == fieldReplySerial {
				m.replySerial = m.order.Uint32(fields[i:])
			}
			i += 4
		case 's', 'o':
			align(4)
			if i+4 > len(fields) {
				return
			}
			n := int(m.order.Uint32(fields[i:]))
			i += 4
			if n < 0 || i+n+1 > len(fields) {
				return
			}
			m.setField(code, string(fields[i:i+n]))
			i += n + 1
		case 'g':
			n := int(fields[i])
			i++
			if i+n+1 > len(fields) {
				return
			}
			m.setField(code, string(fields[i:i+n]))
			i += n + 1
		default:
			return
		}
	}
}

// setField keeps a string field read cares about.
func (m *message) setField(code byte, value string) {
	switch code {
	case fieldMember:
		m.member = value
	case fieldErrorName:
		m.errorName = value
	case fieldSignature:
		m.signature = value
	}
}
