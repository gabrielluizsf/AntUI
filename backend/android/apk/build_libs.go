package apk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

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
		cc, err := tc.NDK.Clang(abi, cfg.MinSDK)
		if err != nil {
			return nil, fmt.Errorf("apk: %w", err)
		}
		out := filepath.Join(dir, string(abi), "lib"+cfg.Lib+".so")
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return nil, fmt.Errorf("apk: %w", err)
		}

		env := append(os.Environ(),
			"GOOS=android",
			"GOARCH="+abi.GOARCH(),
			"CGO_ENABLED=1",
			"CC="+cc,
			"CXX="+strings.TrimSuffix(cc, "clang")+"clang++",
		)
		if arm := abi.GOARM(); arm != "" {
			env = append(env, "GOARM="+arm)
		}

		// Nothing is stripped at the link. The library is built whole and
		// then split into what ships and what reads a crash — see
		// [keepSymbols], which explains why that is not the same as -s -w
		// even though the shipped library comes out identical in size.
		args := []string{"build", "-buildmode=c-shared", "-trimpath",
			"-overlay", overlay, "-o", out}
		if len(cfg.Tags) > 0 {
			args = append(args, "-tags", strings.Join(cfg.Tags, ","))
		}
		if cfg.Ldflags != "" {
			args = append(args, "-ldflags", cfg.Ldflags)
		}
		args = append(args, ".")

		if _, err := runIn(cfg.Dir, env, goBin, args...); err != nil {
			return nil, fmt.Errorf("apk: building for %s: %w", abi, err)
		}
		if *cfg.Strip {
			if err := keepSymbols(tc, out); err != nil {
				return nil, err
			}
		}
		libs[abi] = out
	}
	return libs, nil
}
