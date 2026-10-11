//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
)

// changePropertyOn is the same on somebody else's window, which is how a
// selection is handed over: the asker names a property on their own window
// and the owner writes the answer into it.
func (x *Driver) changePropertyOn(window, property, typ uint32, format byte,
	value []byte, count uint32) {

	body := make([]byte, 20, 20+len(value)+3)
	binary.LittleEndian.PutUint32(body[0:], window)
	binary.LittleEndian.PutUint32(body[4:], property)
	binary.LittleEndian.PutUint32(body[8:], typ)
	body[12] = format
	binary.LittleEndian.PutUint32(body[16:], count)
	body = append(body, value...)
	x.send(xChangeProperty, 0, body) // detail 0 = replace
}
