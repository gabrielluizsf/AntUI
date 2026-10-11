//go:build linux && !android

package antui

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
