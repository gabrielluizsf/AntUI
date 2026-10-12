package sdk

import (
	"strings"
)

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
