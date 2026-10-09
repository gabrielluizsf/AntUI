package css

import (
	"os"
)

// ParseFile reads a CSS file, exactly as [Parse] reads its text, and gives
// the file's directory to the @imports in it so they are read from beside it.
func ParseFile(cssFile string) (*Sheet, error) {
	data, err := os.ReadFile(cssFile)
	if err != nil {
		return nil, err
	}
	return parseText(string(data), cssFile)
}
