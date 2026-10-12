package event

// A SoundBank turns events into sound. It is the hearable half of a template:
// the visual half reports what happened, and the bank decides what it sounds
// like. antui/template/audio ships two of them (Tech and Simple), and any
// type implementing On can play anything at all — or nothing, which is a
// sound bank too.
type SoundBank interface {
	On(e Event)
}
