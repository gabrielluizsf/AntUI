package apk

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// keepSymbols splits a freshly built library in two: a symbol file that can
// read a crash, and the library that ships.
//
// The symbol file keeps the symbol table and drops the DWARF, which is what
// the store's symbolication uses and a third of the size of keeping
// everything. Function names come back; line numbers do not.
func keepSymbols(tc *sdk.Toolchain, lib string) error {
	objcopy := tc.NDK.Tool("objcopy")
	strip := tc.NDK.Tool("strip")

	if err := copyFile(lib, symbolFile(lib)); err != nil {
		return fmt.Errorf("apk: keeping symbols: %w", err)
	}
	if _, err := run(objcopy, "--strip-debug", symbolFile(lib)); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	if _, err := run(strip, "--strip-all", lib); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return nil
}
