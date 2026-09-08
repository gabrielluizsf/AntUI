package antui

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/font"
)

// Finding the font the platform writes its own interface in.
//
// **The system is asked, not guessed.** Every desktop keeps a register of
// what its interface font is and where the file lives, and every one of them
// answers a different way:
//
//	Linux, BSD  fontconfig, through fc-match, which is the same resolver
//	            every toolkit on the machine uses
//	Android     /system/fonts/fonts.xml, the platform's own font
//	            configuration, read in the order it lists families
//	Windows     the registry's font table, which maps a font's name to its
//	            file, read with reg query
//	macOS       fontconfig when it is installed; CoreText is the right
//	            answer and is not here — see systemFontAsked
//
// Only when asking fails does a list of well-known paths get tried, and that
// list is a fallback and not the method. A machine with a font nobody here
// has heard of should still get its own font, and the only way that happens
// is by asking.
//
// The answer is cached: it runs a program, and a face at three sizes should
// not run it three times.

// ErrNoSystemFont is what [SystemFace] gives back when the system has no
// answer and nothing on the fallback list is there.
//
// It is worth telling apart from a broken font file, because the response is
// the same either way — carry on with the built-in — and a program should
// not stop for it.
var ErrNoSystemFont = errors.New("antui: no system font found")

var (
	sysOnce sync.Once
	sysData []byte
	sysPath string
	sysErr  error
)

// SystemFont is the file this platform writes its interface in, and the
// bytes of it.
//
// The path is there so a program can say which font it chose — a settings
// screen showing "Noto Sans" is a great deal more use than one showing
// nothing — and so that a bug report can name it.
func SystemFont() (path string, data []byte, err error) {
	sysOnce.Do(func() {
		// Asked for first, guessed at second.
		for _, candidate := range append(systemFontAsked(), systemFontPaths()...) {
			body, err := os.ReadFile(candidate)
			if err != nil {
				continue
			}
			// Read rather than merely found: a path that exists and is not a
			// TrueType file — a CFF .otf, a link to nothing, a stub — is no
			// more use than a missing one, and the next candidate is better
			// than failing here.
			if _, err := font.ParseTTF(body); err != nil {
				continue
			}
			sysPath, sysData = candidate, body
			return
		}
		sysErr = ErrNoSystemFont
	})
	return sysPath, sysData, sysErr
}

// SystemFace is the platform's own interface font at a size in pixels.
//
//	face, err := antui.SystemFace(15)
//	if err == nil {
//	    antui.SetDefaultFace(face)
//	}
//
// A program that wants to look like the machine it is on asks for this and
// carries on with the built-in when it is not there, which is what the two
// lines above do. Refusing to start because a font is missing would be a
// strange thing for a window library to do.
func SystemFace(pixels float64) (*canvas.Face, error) {
	_, data, err := SystemFont()
	if err != nil {
		return nil, err
	}
	return canvas.ParseFace(data, pixels)
}

// UIFace is [SystemFace] at the size the platform draws its own interface
// at, scaled for the display.
//
// Points, not pixels: the size an interface is specified in is a physical
// one, and the number of pixels it comes to depends on the screen. A program
// with a window passes win.ContentScale(); zero points means the platform's
// own default, which differs by platform and is why it is not a constant
// here.
func UIFace(points, scale float64) (*canvas.Face, error) {
	if points <= 0 {
		points = defaultUIPoints
	}
	if scale <= 0 {
		scale = 1
	}
	return SystemFace(points * scale)
}

// SetSystemFace makes the program draw in the platform's own font, and says
// whether it managed.
//
// It is the whole of what most programs want:
//
//	antui.SetSystemFace(0, win.ContentScale())
//
// and it is deliberately not done automatically. The built-in font is exact
// and identical on every machine, which is what a test card, a screenshot
// comparison and a game with a fixed layout all depend on; changing that
// underneath every program that already exists would break them for a
// nicety.
func SetSystemFace(points, scale float64) bool {
	f, err := UIFace(points, scale)
	if err != nil {
		return false
	}
	canvas.SetDefaultFace(f)
	return true
}

// ask runs a program that answers with a path, and returns nothing at all if
// it is not there, fails, or takes too long.
//
// The timeout is not caution for its own sake: fc-match on a machine with a
// cold font cache can take seconds while it builds one, and a library that
// blocks a program's start-up on that has made things worse rather than
// better. Two seconds is far more than a warm cache needs.
func ask(name string, args ...string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	done := make(chan string, 1)
	cmd := exec.Command(path, args...)
	go func() {
		out, err := cmd.Output()
		if err != nil {
			done <- ""
			return
		}
		done <- strings.TrimSpace(string(out))
	}()
	select {
	case answer := <-done:
		return answer
	case <-time.After(2 * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return ""
	}
}