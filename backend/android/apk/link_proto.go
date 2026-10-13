package apk

import "github.com/gabrielluizsf/antui/backend/android/sdk"

// LinkProto is [Link] asking aapt2 for the protocol-buffer forms, which is
// what a bundle holds instead of binary XML and an arsc.
func LinkProto(tc *sdk.Toolchain, cfg *Config, out string) error {
	return link(tc, cfg, out, true)
}
