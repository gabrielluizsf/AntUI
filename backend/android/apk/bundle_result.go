package apk

import (
	"fmt"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// bundleResult describes the signed bundle. It already carries the symbols
// inside it, so a copy beside it would be the same three megabytes twice —
// one is written only when the caller asked for one by name.
func bundleResult(cfg Config, libs map[sdk.ABI]string) (*Result, error) {
	st, err := os.Stat(cfg.Out)
	if err != nil {
		return nil, fmt.Errorf("apk: %w", err)
	}
	res := &Result{APK: cfg.Out, Libs: libs, Size: st.Size()}
	if cfg.Symbols != "" {
		if wrote, err := WriteSymbols(libs, cfg.Lib, cfg.Symbols); err != nil {
			return nil, err
		} else if wrote {
			res.Symbols = cfg.Symbols
		}
	}
	return res, nil
}
