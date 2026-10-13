package apk

import "encoding/binary"

func protoBytes(field int, body []byte) []byte {
	var out []byte
	out = binary.AppendUvarint(out, uint64(field)<<3|2)
	out = binary.AppendUvarint(out, uint64(len(body)))
	return append(out, body...)
}
