//go:build linux

package antui

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// The PulseAudio native protocol, spoken straight over its unix socket.
//
// Not ALSA: a desktop runs PipeWire or PulseAudio, and whichever it is holds
// the card — a program that opens the hardware device behind the server's
// back either fails outright or takes the sound away from everything else.
// Both servers answer on this socket, and PipeWire's answer is the reason
// this works on a machine with no PulseAudio installed at all.
//
// Everything here is written for **protocol 13**, which is what this
// announces. A newer server speaks it happily; the versions above it add
// fields to the same commands and buy a playback stream nothing at all.

const (
	pulseVersion = 13

	// The tags a value carries. The protocol is big-endian throughout, and
	// the samples inside a data packet are not — they are whatever the sample
	// spec said, which here is little-endian floats.
	tagString     = 't'
	tagStringNil  = 'N'
	tagU32        = 'L'
	tagU8         = 'B'
	tagArbitrary  = 'x'
	tagTrue       = '1'
	tagFalse      = '0'
	tagSampleSpec = 'a'
	tagChannelMap = 'm'
	tagCvolume    = 'v'
	tagProplist   = 'P'

	// The commands this needs. A server says a great deal more than this and
	// all of it is ignored.
	cmdError                = 0
	cmdReply                = 2
	cmdCreatePlaybackStream = 3
	cmdDeletePlaybackStream = 4
	cmdAuth                 = 8
	cmdSetClientName        = 9
	cmdRequest              = 61

	sampleFloat32LE = 5 // PA_SAMPLE_FLOAT32LE
	frontLeft       = 1 // PA_CHANNEL_POSITION_FRONT_LEFT
	frontRight      = 2
	volumeNorm      = 0x10000 // PA_VOLUME_NORM, which is unattenuated

	noIndex    = 0xFFFFFFFF // PA_INVALID_INDEX
	noPacket   = 0xFFFFFFFF // a descriptor's channel, for a command
	maxPayload = 64 << 10   // the biggest packet this will read
)

// --- the tagstruct -----------------------------------------------------

// A tags is a tagstruct being built: values, each behind the byte that says
// what it is.
type tags struct{ b []byte }

func (t *tags) u32(v uint32) {
	t.b = append(t.b, tagU32)
	t.b = binary.BigEndian.AppendUint32(t.b, v)
}

func (t *tags) u8(v byte) { t.b = append(t.b, tagU8, v) }

func (t *tags) boolean(v bool) {
	if v {
		t.b = append(t.b, tagTrue)
		return
	}
	t.b = append(t.b, tagFalse)
}

// str writes a string, which carries its terminator: the protocol has no
// length in front of one.
func (t *tags) str(s string) {
	t.b = append(t.b, tagString)
	t.b = append(t.b, s...)
	t.b = append(t.b, 0)
}

func (t *tags) noStr() { t.b = append(t.b, tagStringNil) }

func (t *tags) bytes(b []byte) {
	t.b = append(t.b, tagArbitrary)
	t.b = binary.BigEndian.AppendUint32(t.b, uint32(len(b)))
	t.b = append(t.b, b...)
}

func (t *tags) sampleSpec(format, channels byte, rate uint32) {
	t.b = append(t.b, tagSampleSpec, format, channels)
	t.b = binary.BigEndian.AppendUint32(t.b, rate)
}

func (t *tags) channelMap(positions ...byte) {
	t.b = append(t.b, tagChannelMap, byte(len(positions)))
	t.b = append(t.b, positions...)
}

func (t *tags) volume(channels int, level uint32) {
	t.b = append(t.b, tagCvolume, byte(channels))
	for range channels {
		t.b = binary.BigEndian.AppendUint32(t.b, level)
	}
}

// props writes a property list, which is what a stream is named through from
// protocol 13 on. The keys are what a mixer shows the player, so they are
// worth filling in: an application called "audio stream" in the volume
// control is one nobody can turn down.
func (t *tags) props(pairs ...[2]string) {
	t.b = append(t.b, tagProplist)
	for _, pair := range pairs {
		t.str(pair[0])
		t.u32(uint32(len(pair[1]) + 1))
		t.bytes(append([]byte(pair[1]), 0))
	}
	t.b = append(t.b, tagStringNil) // the end of the list
}

// A read is a tagstruct being taken apart.
type read struct {
	b  []byte
	at int
}

var errShort = errors.New("antui: the server's answer stopped in the middle")

func (r *read) tag(want byte) error {
	if r.at >= len(r.b) {
		return errShort
	}
	if got := r.b[r.at]; got != want {
		return fmt.Errorf("antui: the server sent a %q where a %q belongs", got, want)
	}
	r.at++
	return nil
}

func (r *read) u32() (uint32, error) {
	if err := r.tag(tagU32); err != nil {
		return 0, err
	}
	if r.at+4 > len(r.b) {
		return 0, errShort
	}
	v := binary.BigEndian.Uint32(r.b[r.at:])
	r.at += 4
	return v, nil
}

// --- the connection ----------------------------------------------------

// A pulse is one playback stream and the goroutine feeding it.
type pulse struct {
	conn    net.Conn
	audio   *Audio
	channel uint32 // the stream, as the server numbers it
	serial  uint32
	delay   time.Duration

	block []float32 // handed to the game to fill
	out   []byte    // and the same thing as bytes for the socket
}

func (p *pulse) latency() time.Duration { return p.delay }

func (p *pulse) close() error {
	if p.conn == nil {
		return nil
	}
	// The stream is deleted first so the server frees it now rather than
	// when it notices the socket has gone.
	var t tags
	t.u32(cmdDeletePlaybackStream)
	t.u32(p.next())
	t.u32(p.channel)
	_ = p.send(noPacket, t.b)
	return p.conn.Close()
}

func (p *pulse) next() uint32 {
	p.serial++
	return p.serial
}

// send writes one packet: the twenty-byte descriptor, then the payload.
func (p *pulse) send(channel uint32, payload []byte) error {
	var head [20]byte
	binary.BigEndian.PutUint32(head[0:], uint32(len(payload)))
	binary.BigEndian.PutUint32(head[4:], channel)
	binary.BigEndian.PutUint32(head[8:], 0)  // offset, high half
	binary.BigEndian.PutUint32(head[12:], 0) // and low
	binary.BigEndian.PutUint32(head[16:], 0) // flags
	if _, err := p.conn.Write(head[:]); err != nil {
		return err
	}
	_, err := p.conn.Write(payload)
	return err
}

// receive reads one packet.
func (p *pulse) receive() (channel uint32, payload []byte, err error) {
	var head [20]byte
	if _, err := io.ReadFull(p.conn, head[:]); err != nil {
		return 0, nil, err
	}
	length := binary.BigEndian.Uint32(head[0:])
	channel = binary.BigEndian.Uint32(head[4:])
	if length > maxPayload {
		return 0, nil, fmt.Errorf("antui: the server sent %d bytes at once", length)
	}
	payload = make([]byte, length)
	if _, err := io.ReadFull(p.conn, payload); err != nil {
		return 0, nil, err
	}
	return channel, payload, nil
}

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

// --- opening -----------------------------------------------------------

// pulseSocket is where the server is listening, in the order the ordinary
// tools look.
func pulseSocket() string {
	if server := os.Getenv("PULSE_SERVER"); server != "" {
		// A server address can name a transport and can list several; the
		// only one that matters here is a local socket.
		for part := range strings.SplitSeq(server, " ") {
			if path, ok := strings.CutPrefix(part, "unix:"); ok {
				return path
			}
			if strings.HasPrefix(part, "/") {
				return part
			}
		}
	}
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return filepath.Join(dir, "pulse", "native")
	}
	return fmt.Sprintf("/run/user/%d/pulse/native", os.Getuid())
}

// pulseCookie is the shared secret a session's programs authenticate with.
// A server that does not care about it — PipeWire does not — is given one
// anyway, because the field is not optional.
func pulseCookie() []byte {
	cookie := make([]byte, 256)
	paths := []string{os.Getenv("PULSE_COOKIE")}
	if config := os.Getenv("XDG_CONFIG_HOME"); config != "" {
		paths = append(paths, filepath.Join(config, "pulse", "cookie"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths,
			filepath.Join(home, ".config", "pulse", "cookie"),
			filepath.Join(home, ".pulse-cookie"))
	}
	for _, path := range paths {
		if path == "" {
			continue
		}
		if data, err := os.ReadFile(path); err == nil && len(data) >= len(cookie) {
			copy(cookie, data)
			return cookie
		}
	}
	return cookie
}

// openAudio is the Linux half of OpenAudio.
func openAudio(name string, a *Audio) (driver, error) {
	conn, err := net.Dial("unix", pulseSocket())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoAudio, err)
	}
	p := &pulse{conn: conn, audio: a}

	if _, err := p.ask(cmdAuth, func(t *tags) {
		t.u32(pulseVersion)
		t.bytes(pulseCookie())
	}); err != nil {
		conn.Close()
		return nil, err
	}

	if _, err := p.ask(cmdSetClientName, func(t *tags) {
		t.props([2]string{"application.name", name},
			[2]string{"application.process.binary", filepath.Base(os.Args[0])},
			[2]string{"application.icon_name", "applications-games"})
	}); err != nil {
		conn.Close()
		return nil, err
	}

	// What the buffer is worth in bytes. A game wants this short — a sound
	// that arrives a tenth of a second after the thing that made it is a
	// sound that belongs to nothing — and not so short that a frame taking
	// longer than usual leaves a hole in it.
	const (
		wanted = 25 * time.Millisecond // between writing a sample and hearing it
		least  = 10 * time.Millisecond // the least the server should ask for
	)
	frame := uint32(a.channels * 4)
	target := uint32(a.rate) * frame * uint32(wanted) / uint32(time.Second)
	minimum := uint32(a.rate) * frame * uint32(least) / uint32(time.Second)

	reply, err := p.ask(cmdCreatePlaybackStream, func(t *tags) {
		t.sampleSpec(sampleFloat32LE, byte(a.channels), uint32(a.rate))
		t.channelMap(frontLeft, frontRight)
		t.u32(noIndex) // any sink: the one the player chose
		t.noStr()      // and so no sink by name either
		t.u32(target * 4)
		t.boolean(false) // not corked: play as soon as there is anything
		t.u32(target)    // tlength, which is the latency being asked for
		t.u32(target)    // prebuf: start once there is a bufferful
		t.u32(minimum)   // minreq
		t.u32(noIndex)   // not synchronised with another stream
		t.volume(a.channels, volumeNorm)

		// Protocol 12 added these seven, and every one of them is no.
		t.boolean(false) // no_remap
		t.boolean(false) // no_remix
		t.boolean(false) // fix_format
		t.boolean(false) // fix_rate
		t.boolean(false) // fix_channels
		t.boolean(false) // no_move
		t.boolean(false) // variable_rate

		// And 13 added these three.
		t.boolean(false) // muted
		t.boolean(true)  // adjust_latency: hold the buffer to what was asked
		t.props([2]string{"media.name", name},
			[2]string{"media.role", "game"})
	})
	if err != nil {
		conn.Close()
		return nil, err
	}

	channel, err := reply.u32()
	if err != nil {
		conn.Close()
		return nil, err
	}
	if _, err := reply.u32(); err != nil { // the sink input's own index
		conn.Close()
		return nil, err
	}
	missing, err := reply.u32()
	if err != nil {
		conn.Close()
		return nil, err
	}
	p.channel = channel
	p.delay = time.Duration(target) * time.Second / time.Duration(uint32(a.rate)*frame)

	// The block handed to the game each time. Big enough that the goroutine
	// is not woken for a handful of samples, small enough to fit inside the
	// latency asked for.
	frames := max(int(minimum/frame), 64)
	p.block = make([]float32, frames*a.channels)
	p.out = make([]byte, len(p.block)*4)

	go p.run(missing)
	return p, nil
}

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

// write pulls one block from the game and sends it.
func (p *pulse) write(samples int) error {
	block := p.block[:samples]
	p.audio.pull(block)
	for i, sample := range block {
		binary.LittleEndian.PutUint32(p.out[i*4:], math.Float32bits(sample))
	}
	return p.send(p.channel, p.out[:samples*4])
}

// stop records why the sound stopped, unless it stopped because it was told
// to.
func (p *pulse) stop(err error) {
	if p.audio.stopped() {
		return
	}
	p.audio.fail(fmt.Errorf("antui: the sound server went away: %w", err))
}