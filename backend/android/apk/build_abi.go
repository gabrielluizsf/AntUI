package apk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// buildAbi compiles the Go package for one ABI and returns the library it
// wrote. The compiler is the NDK's clang at the app's minimum SDK, not the
// newest available: a library built against newer headers can reference
// symbols the oldest device it claims to support does not have, and that
// failure happens at load time on a user's phone.
func buildAbi(tc *sdk.Toolchain, cfg *Config, dir string, abi sdk.ABI, goBin, overlay string) (string, error) {
	cc, err := tc.NDK.Clang(abi, cfg.MinSDK)
	if err != nil {
		return "", fmt.Errorf("apk: %w", err)
	}
	out := filepath.Join(dir, string(abi), "lib"+cfg.Lib+".so")
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return "", fmt.Errorf("apk: %w", err)
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

	// Nothing is stripped at the link: the library is built whole and then
	// split into what ships and what reads a crash — see [keepSymbols].
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
		return "", fmt.Errorf("apk: building for %s: %w", abi, err)
	}
	if *cfg.Strip {
		if err := keepSymbols(tc, out); err != nil {
			return "", err
		}
	}
	return out, nil
}
