//go:build linux && !android

package antui

import (
	"fmt"
)

// ask sends a command and waits for the answer to that one, letting anything
// else the server has to say go past. A server talks before it is spoken to
// — a stream can be asked for data before the command that made it has been
// replied to — so this cannot be a plain read.
func (p *pulse) ask(command uint32, build func(*tags)) (*read, error) {
	serial := p.next()

	var t tags
	t.u32(command)
	t.u32(serial)
	if build != nil {
		build(&t)
	}
	if err := p.send(noPacket, t.b); err != nil {
		return nil, err
	}

	for {
		channel, payload, err := p.receive()
		if err != nil {
			return nil, err
		}
		if channel != noPacket {
			continue // sound coming back at us, which cannot happen here
		}
		r := &read{b: payload}
		got, err := r.u32()
		if err != nil {
			return nil, err
		}
		tag, err := r.u32()
		if err != nil {
			return nil, err
		}
		if tag != serial {
			continue
		}
		switch got {
		case cmdReply:
			return r, nil
		case cmdError:
			code, _ := r.u32()
			return nil, fmt.Errorf("antui: the sound server refused command %d (error %d)",
				command, code)
		default:
			return nil, fmt.Errorf("antui: the sound server answered with command %d", got)
		}
	}
}
