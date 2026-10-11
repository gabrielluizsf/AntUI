//go:build linux && !android

package antui

// writeCreateStream writes the create-playback-stream request.
func writeCreateStream(t *tags, s pulsePlayback) {
	t.sampleSpec(sampleFloat32LE, byte(s.audio.channels), uint32(s.audio.rate))
	t.channelMap(frontLeft, frontRight)
	t.u32(noIndex) // any sink: the one the player chose
	t.noStr()      // and so no sink by name either
	t.u32(s.target * 4)
	t.boolean(false) // not corked: play as soon as there is anything
	t.u32(s.target)  // tlength, which is the latency being asked for
	t.u32(s.target)  // prebuf: start once there is a bufferful
	t.u32(s.minimum) // minreq
	t.u32(noIndex)   // not synchronised with another stream
	t.volume(s.audio.channels, volumeNorm)

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
	t.props([2]string{"media.name", s.name},
		[2]string{"media.role", "game"})
}
