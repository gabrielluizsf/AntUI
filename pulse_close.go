//go:build linux && !android

package antui

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
