package apk

import "github.com/gabrielluizsf/antui/backend/android/sdk"

// BuildLibs compiles the Go package once per ABI, writing each shared
// library into dir. The compiler is the NDK's clang for that ABI at the
// app's minimum SDK — not the newest available, because a library built
// against a newer platform's headers can reference symbols the oldest device
// it claims to support does not have, and that failure happens at load time
// on a user's phone rather than here.
func BuildLibs(tc *sdk.Toolchain, cfg *Config, dir string) (map[sdk.ABI]string, error) {
	if err := cfg.check(); err != nil {
		return nil, err
	}
	goBin, err := goCommand()
	if err != nil {
		return nil, err
	}
	overlay, err := overlayFor(cfg, dir)
	if err != nil {
		return nil, err
	}

	libs := make(map[sdk.ABI]string, len(cfg.ABIs))
	for _, abi := range cfg.ABIs {
		out, err := buildAbi(tc, cfg, dir, abi, goBin, overlay)
		if err != nil {
			return nil, err
		}
		libs[abi] = out
	}
	return libs, nil
}
