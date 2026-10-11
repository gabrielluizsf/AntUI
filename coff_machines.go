package antui

// The machine numbers a COFF header can carry, and the relocation each of
// them uses for "the address of this, without the image base".
var coffMachines = map[string]struct {
	machine    uint16
	relocation uint16
}{
	"amd64": {0x8664, 3}, // IMAGE_REL_AMD64_ADDR32NB
	"386":   {0x014c, 7}, // IMAGE_REL_I386_DIR32NB
	"arm64": {0xAA64, 2}, // IMAGE_REL_ARM64_ADDR32NB
	"arm":   {0x01C0, 2}, // IMAGE_REL_ARM_ADDR32NB
}
