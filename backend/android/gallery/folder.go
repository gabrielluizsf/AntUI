//go:build android

package gallery

func (k Kind) folder() string {
	switch k {
	case Video:
		return "Movies"
	case Audio:
		return "Music"
	}
	return "Pictures"
}
