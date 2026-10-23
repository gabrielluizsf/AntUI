package dbus

import "encoding/binary"

// Message types.
const (
	typeMethodCall   = 1
	typeMethodReturn = 2
	typeError        = 3
)

// Header field codes, and the one type each is sent as.
const (
	fieldPath        = 1 // o
	fieldInterface   = 2 // s
	fieldMember      = 3 // s
	fieldErrorName   = 4 // s
	fieldReplySerial = 5 // u
	fieldDestination = 6 // s
	fieldSignature   = 8 // g
)

// headerField is one entry of a message's header: a code and its value,
// which is a string for every field sent here but REPLY_SERIAL.
type headerField struct {
	code byte
	sig  byte // 's', 'o', 'g' or 'u'
	str  string
	num  uint32
}

// encodeMessage lays out a whole message: the fixed header, the header
// fields, padding to 8 and the body, whose signature becomes a field of its
// own when there is a body.
func encodeMessage(typ byte, serial uint32, fields []headerField, sig string, body []byte) []byte {
	if sig != "" {
		fields = append(fields, headerField{code: fieldSignature, sig: 'g', str: sig})
	}
	e := encoder{buf: make([]byte, 0, 128+len(body))}
	e.byte('l')
	e.byte(typ)
	e.byte(0) // flags
	e.byte(1) // protocol version
	e.uint32(uint32(len(body)))
	e.uint32(serial)

	lengthAt := len(e.buf)
	e.uint32(0) // the fields' length, filled in below
	start := len(e.buf)
	for _, f := range fields {
		e.align(8)
		e.byte(f.code)
		e.signature(string(f.sig))
		switch f.sig {
		case 'u':
			e.uint32(f.num)
		case 'g':
			e.signature(f.str)
		default:
			e.string(f.str)
		}
	}
	binary.LittleEndian.PutUint32(e.buf[lengthAt:], uint32(len(e.buf)-start))
	e.align(8)
	return append(e.buf, body...)
}
