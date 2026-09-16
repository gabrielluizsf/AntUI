// Package audio is the hearable half of a template: a set of short sounds,
// one per interactive component, that plays through [antui.OpenAudio] when
// the component is used.
//
// Two sound banks come ready-made, matching the two ready-made looks. Tech is
// the crisp electronic voice of the Cyberpunk look. Simple is the bouncy,
// playful voice of the built-in look:
//
//	player := audio.NewPlayer(44100)
//	stream, err := antui.OpenAudio("my game", 44100, player.Fill)
//	audio.Use(player)            // the banks below play through it now
//
//	audio.Tech.Button.Play()             // a futuristic click
//	audio.Simple.TextInput.PlayStep(3)   // a cartoon keystroke
//
// A sound bank is also an [event.SoundBank]: hand one to a template and the
// template plays the right sound for each event all by itself, without the
// game touching audio at all. Nothing plays until a Player is in use, so a
// machine without speakers merely stays quiet.
//
// Licensing: the WAV files embedded in sfx/ are CC0 1.0 — public domain, free
// for commercial use, no attribution required. See sfx/THIRD_PARTY.md.
package audio

import (
	"sync/atomic"
)

// defaultDevice is the player the packages's sounds play through, or nil
// when nothing is open and everything should stay quiet.
var defaultPlayer atomic.Pointer[Player]

// Use names the Player the standalone sounds play through. Hand it a player
// whose Fill you passed to [antui.OpenAudio].
func Use(p *Player) {
	var next *Player
	if p != nil {
		next = p
	}
	defaultPlayer.Store(next)
}
