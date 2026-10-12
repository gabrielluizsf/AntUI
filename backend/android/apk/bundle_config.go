package apk

import (
	"bytes"
)

// BundleConfig.pb, encoded by hand.
//
// It is a protocol buffer and there is no generated code here to write one
// with, but the message needed is small enough to spell out: which tool
// built it, and which files the store should leave uncompressed in the APKs
// it generates. Everything else is a default, and a default is what an
// ordinary app wants.
func bundleConfig() []byte {
	var out bytes.Buffer

	// Bundletool { version = 2 }. The field number is **2** and not 1, which
	// is the sort of thing that cannot be reasoned out: a wrong number
	// produces a perfectly well-formed protocol buffer, and bundletool
	// answers "Version must match the format '<major>.<minor>.<revision>',
	// but found ''" — an empty string, because it read a field that was not
	// there. It was settled by having bundletool write one and comparing the
	// bytes, which is the only way to settle it.
	out.Write(protoMessage(1, protoString(2, bundleFormat)))

	// Compression { uncompressed_glob = ... } — field 3, inner field 1,
	// checked the same way. The native libraries have to be stored rather
	// than deflated in the generated APKs, because the manifest says the
	// loader maps them out of the package and a compressed entry cannot be
	// mapped.
	var compression bytes.Buffer
	for _, glob := range []string{"**.so"} {
		compression.Write(protoString(1, glob))
	}
	out.Write(protoMessage(3, compression.Bytes()))
	return out.Bytes()
}
