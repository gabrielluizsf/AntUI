package apk

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/canvas"
)

// iconKey is what the archive depends on: the picture itself, the colour
// behind it, and the tool that compiled it.
func iconKey(tc *sdk.Toolchain, icon string, bg canvas.Color) (string, error) {
	body, err := os.ReadFile(icon)
	if err != nil {
		return "", fmt.Errorf("apk: %w", err)
	}
	h := sha256.New()
	fmt.Fprintf(h, "%s\n%s\n%s\n", cacheVersion, bg, tc.BuildTools.Version)
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))[:32], nil
}
