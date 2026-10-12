package apk

// Symbols are how a crash on someone else's phone stays readable.
//
// A native crash arrives from the store as a list of addresses. Turning
// those back into function names needs a copy of the library with its symbol
// table still in it, and the copy has to be the *same build* — a rebuild from
// the same source lands the code somewhere else and symbolises to nonsense.
// So the moment to keep it is the moment it is made, and there is no second
// chance: once the stripped library has shipped, a crash in it is a list of
// addresses forever.
//
// This is why stripping happens here rather than at the link. Go's -s and -w
// tell the linker never to write the symbols at all, which costs the same as
// stripping afterwards — the shipped library comes out the same size to the
// byte — and throws away the only copy. Building whole and stripping after
// leaves something to keep.
//
// Sizes, for one real arm64 app:
//
//	4.4 MB   as built, with DWARF
//	3.3 MB   the symbol file: DWARF dropped, symbol table kept
//	3.0 MB   what ships, stripped
//	3.0 MB   what -s -w used to ship, with nothing kept
//
// The archive is what Play asks for on the release page, and [Bundle] also
// puts it inside the bundle so that one upload carries both.
const symbolSuffix = ".sym"
