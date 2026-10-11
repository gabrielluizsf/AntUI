//go:build linux && !android

package antui

func (t *tags) channelMap(positions ...byte) {
	t.b = append(t.b, tagChannelMap, byte(len(positions)))
	t.b = append(t.b, positions...)
}
