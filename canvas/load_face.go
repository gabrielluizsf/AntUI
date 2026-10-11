package canvas

import (
	"fmt"
	"os"
)

// LoadFace reads a TrueType file at a size in pixels per em.
func LoadFace(path string, pixels float64) (*Face, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("antui: %w", err)
	}
	return ParseFace(body, pixels)
}
