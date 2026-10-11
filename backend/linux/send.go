//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// send writes a request and returns the sequence number the server will
// answer it under.
func (x *Driver) send(opcode, detail byte, body []byte) uint16 {
	x.seq++
	x.write(x.request(opcode, detail, body))
	return x.seq
}
