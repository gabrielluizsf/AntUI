package template

import "github.com/gabrielluizsf/antui/template/audio"

// Cyberpunk is the futuristic template: neon edges and glowing controls, the
// level of finish a modern website skins itself in, voiced by
// [github.com/gabrielluizsf/antui/template/audio.Tech].
func Cyberpunk() Template {
	return New(cyberStyle{}, audio.Tech)
}
