package apk

import (
	"fmt"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// buildResult describes the APK that was just signed and writes the symbols
// beside it: the moment they are made is the only moment they can be kept.
func buildResult(cfg Config, libs map[sdk.ABI]string) (*Result, error) {
	st, err := os.Stat(cfg.Out)
	if err != nil {
		return nil, fmt.Errorf("apk: %w", err)
	}
	res := &Result{APK: cfg.Out, Libs: libs, Size: st.Size()}
	syms := symbolsPath(&cfg)
	if wrote, err := WriteSymbols(libs, cfg.Lib, syms); err != nil {
		return nil, err
	} else if wrote {
		res.Symbols = syms
	}
	return res, nil
}
