package font

type CmapTable struct {
	Format uint16
	Data   []byte

	// format 4
	Segments int
	// format 12
	Groups int
}
