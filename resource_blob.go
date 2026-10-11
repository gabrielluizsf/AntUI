package antui

// resourceBlob is one leaf of the resource section: what it is called, its
// type, and its bytes.
type resourceBlob struct {
	id     uint16
	typeID uint16
	data   []byte
}
