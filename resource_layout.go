package antui

// resourceLayout is where every directory, leaf and byte lands once the
// section has been laid out.
type resourceLayout struct {
	out          []byte
	typeDirs     map[uint16]uint32
	nameDirs     []uint32
	groupNameDir uint32
	leaves       []uint32
	starts       []uint32
	icons        int
}
