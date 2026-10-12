package apk

import (
	"bytes"
	"encoding/xml"
)

// attr writes an XML attribute value, quotes and all, with everything
// escaped. It goes through encoding/xml rather than a handful of Replaces so
// that a label with an apostrophe or a permission from a config file cannot
// produce a manifest that is not XML.
func attr(s string) string {
	var buf bytes.Buffer
	buf.WriteByte('"')
	if err := xml.EscapeText(&buf, []byte(s)); err != nil {
		// EscapeText only fails if the writer does, and a bytes.Buffer does
		// not.
		panic(err)
	}
	buf.WriteByte('"')
	return buf.String()
}
