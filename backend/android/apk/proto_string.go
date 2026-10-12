package apk

// The two pieces of protocol buffer encoding this needs. Both fields used
// above are length-delimited, which is wire type 2.
func protoString(field int, s string) []byte {
	return protoBytes(field, []byte(s))
}
