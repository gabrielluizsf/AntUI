//go:build linux && !android

package antui

import ()

// run is the device's goroutine: it answers the server's requests for sound
// until the stream is closed or the socket breaks.
func (p *pulse) run(credit uint32) {
	for {
		for credit >= uint32(len(p.out)) {
			if err := p.write(len(p.block)); err != nil {
				p.stop(err)
				return
			}
			credit -= uint32(len(p.out))
		}
		// What is left is less than a block. Sending it now rather than
		// holding on to it is what keeps the buffer full: the server asks
		// for exactly what it has room for.
		if frame := p.audio.channels * 4; credit >= uint32(frame) {
			frames := int(credit) / frame
			if err := p.write(frames * p.audio.channels); err != nil {
				p.stop(err)
				return
			}
			credit -= uint32(frames * frame)
		}

		channel, payload, err := p.receive()
		if err != nil {
			p.stop(err)
			return
		}
		if channel != noPacket {
			continue
		}
		r := &read{b: payload}
		command, err := r.u32()
		if err != nil {
			p.stop(err)
			return
		}
		if _, err := r.u32(); err != nil { // the tag, which nothing answers
			p.stop(err)
			return
		}
		if command == cmdRequest {
			if _, err := r.u32(); err == nil { // the stream, which is ours
				if bytes, err := r.u32(); err == nil {
					credit += bytes
				}
			}
		}
		// Anything else — an underflow, a stream moved to another sink, a
		// volume changed — is the server telling us about the world, and
		// there is nothing here that would do anything differently.
	}
}
