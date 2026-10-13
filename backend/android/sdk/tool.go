package sdk

import "path/filepath"

// Tool is one of the llvm programs beside clang — llvm-strip, llvm-objcopy,
// llvm-readelf — named without its "llvm-" prefix.
func (n NDK) Tool(name string) string { return filepath.Join(n.Bin(), exe("llvm-"+name)) }
