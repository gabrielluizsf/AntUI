package sdk

import (
	"os"
	"strings"
)

// Missing is the error a partial install produces. It is an error and also a
// report: it names each thing that is wrong and, where one exists, the
// sdkmanager line that puts it right.
type Missing struct {
	Root     string
	Problems []Problem
}

// Problem is one thing wrong with the install.
type Problem struct {
	What string // "build-tools", "ndk", "a JDK"
	// Have and Want are set when the thing is there but too old, and empty
	// when it is not there at all.
	Have, Want string
	// Install is the sdkmanager package that fixes it, empty when nothing
	// sdkmanager has would.
	Install string
}

func (m *Missing) add(what, install string) {
	m.Problems = append(m.Problems, Problem{What: what, Install: install})
}

func (m *Missing) old(what, have, want string) {
	m.Problems = append(m.Problems, Problem{What: what, Have: have, Want: want})
}

// Error says what is not installed and what to run to install it.
func (m *Missing) Error() string {
	var b strings.Builder
	b.WriteString("sdk: the Android SDK at " + m.Root + " is not ready: ")
	for i, p := range m.Problems {
		if i > 0 {
			b.WriteString("; ")
		}
		if p.Have != "" {
			b.WriteString(p.What + " is " + p.Have + ", want " + p.Want + " or newer")
		} else {
			b.WriteString(p.What + " is not installed")
		}
	}
	if fix := m.Fix(); fix != "" {
		b.WriteString(". Run: " + fix)
	}
	return b.String()
}

// Fix is the one sdkmanager command that installs everything missing, or the
// empty string when nothing missing can be installed that way.
func (m *Missing) Fix() string {
	var pkgs []string
	for _, p := range m.Problems {
		if p.Install != "" {
			pkgs = append(pkgs, `"`+p.Install+`"`)
		}
	}
	if len(pkgs) == 0 {
		return ""
	}
	return "sdkmanager " + strings.Join(pkgs, " ")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
