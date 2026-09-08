package antui

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/gabrielluizsf/antui/canvas"
)

// The icon a built program carries inside it.
//
// On Windows the file browser shows an executable's icon without running it,
// and it reads it out of a resource directory in the file. Go's linker puts
// any `.syso` file that is sitting in the main package into the link, so a
// program gets its icon by having one of those beside its main.go — no build
// tag, no external tool, and `go run` and `go build` both do it.
//
// So this writes one: a COFF object file with a `.rsrc` section holding the
// pictures as RT_ICON and one RT_GROUP_ICON that names them. That is all a
// `.syso` for an icon is, and writing it here means the engine does not send
// anybody off to install a resource compiler.
//
// The other two platforms have no such thing. Linux reads the icon off the
// running window (SetIcon), and macOS reads it out of a bundle's Info.plist,
// which is a folder rather than something inside the binary — a program run
// from a terminal there gets its icon from SetIcon too.

// The resource types Windows fixes.
const (
	rtIcon      = 3
	rtGroupIcon = 14
)

// WriteICO writes the pictures as a Windows .ico file, largest first.
//
// It is worth having on its own — an .ico is what a file browser or an
// installer wants — and it is half of what the resource below is made of.
func WriteICO(w io.Writer, images ...*canvas.Canvas) error {
	pictures, err := iconPictures(images)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	binary.Write(&out, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(&out, binary.LittleEndian, uint16(1)) // an icon, not a cursor
	binary.Write(&out, binary.LittleEndian, uint16(len(pictures)))

	at := uint32(6 + 16*len(pictures))
	for _, picture := range pictures {
		writeIconEntry(&out, picture, func(size uint32) {
			binary.Write(&out, binary.LittleEndian, size)
			binary.Write(&out, binary.LittleEndian, at)
		})
		at += uint32(len(picture.data))
	}
	for _, picture := range pictures {
		out.Write(picture.data)
	}
	_, err = w.Write(out.Bytes())
	return err
}

// WriteWindowsResource writes a .syso: an object file the Go linker puts into
// the program, holding the icon the file browser shows.
//
// Put what it writes in the main package as `icon_windows.syso` — the name
// ends in `_windows` so that the toolchain leaves it out of every other
// platform's build, which is what makes a checked-in resource harmless
// everywhere else.
func WriteWindowsResource(w io.Writer, images ...*canvas.Canvas) error {
	pictures, err := iconPictures(images)
	if err != nil {
		return err
	}
	section := buildResourceSection(pictures)
	return writeCOFF(w, section)
}

// --- the pictures --------------------------------------------------------------

// an iconPicture is one size, as the bytes that go in the file.
type iconPicture struct {
	width, height int
	data          []byte // a PNG, which every Windows since Vista reads
}

// iconPictures turns canvases into what the formats above want, largest
// first. A size over 256 is left out: the entry that names it has one byte
// for the width, and 256 is written as a nought.
func iconPictures(images []*canvas.Canvas) ([]iconPicture, error) {
	var out []iconPicture
	for _, image := range images {
		if image == nil || image.Width <= 0 || image.Height <= 0 {
			continue
		}
		if image.Width > 256 || image.Height > 256 {
			continue
		}
		data, err := encodeIconPNG(image)
		if err != nil {
			return nil, err
		}
		out = append(out, iconPicture{image.Width, image.Height, data})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("antui: there is no picture to make an icon from")
	}
	// Largest first, which is what every reader here expects.
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].width > out[j-1].width; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
}

// writeIconEntry writes the sixteen-byte description of one picture that both
// an .ico and a group resource begin with, and calls `tail` for the four
// bytes at the end, which is a file offset in one and a resource number in
// the other.
func writeIconEntry(out *bytes.Buffer, picture iconPicture, tail func(size uint32)) {
	// 256 is written as nought: the field is one byte.
	out.WriteByte(byte(picture.width % 256))
	out.WriteByte(byte(picture.height % 256))
	out.WriteByte(0)                                   // colours in the palette: none, it is a true-colour icon
	out.WriteByte(0)                                   // reserved
	binary.Write(out, binary.LittleEndian, uint16(1))  // planes
	binary.Write(out, binary.LittleEndian, uint16(32)) // bits a pixel
	tail(uint32(len(picture.data)))
}

// groupIcon is the little directory that names the pictures. It is what
// Windows looks up first, and it is why the icons need numbers.
func groupIcon(pictures []iconPicture) []byte {
	var out bytes.Buffer
	binary.Write(&out, binary.LittleEndian, uint16(0))
	binary.Write(&out, binary.LittleEndian, uint16(1))
	binary.Write(&out, binary.LittleEndian, uint16(len(pictures)))
	for i, picture := range pictures {
		writeIconEntry(&out, picture, func(size uint32) {
			binary.Write(&out, binary.LittleEndian, size)
			// The number of the RT_ICON that holds it. They are numbered from
			// one, because nought is not a resource name.
			binary.Write(&out, binary.LittleEndian, uint16(i+1))
		})
	}
	return out.Bytes()
}

// encodeIconPNG writes a picture as a PNG, which is what an icon over 48
// pixels is these days and what every Windows since Vista reads. The old
// format is a bitmap with the rows upside down and a mask nobody uses; there
// is no reason to write one.
func encodeIconPNG(image_ *canvas.Canvas) ([]byte, error) {
	out := image.NewNRGBA(image.Rect(0, 0, image_.Width, image_.Height))
	for y := range image_.Height {
		for x := range image_.Width {
			c := image_.At(x, y)
			at := out.PixOffset(x, y)
			out.Pix[at+0] = c.R()
			out.Pix[at+1] = c.G()
			out.Pix[at+2] = c.B()
			out.Pix[at+3] = c.A()
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, out); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// --- the resource section ------------------------------------------------------

// A resource section is three levels of directory — by type, then by name,
// then by language — with a sixteen-byte leaf at the bottom saying where the
// bytes are and how many. The bytes themselves sit after the directories.
//
// The leaf's address is where the linker comes in: it is an address in the
// program that does not exist yet, so the field is left as an offset from the
// start of the section and a relocation tells the linker to add the section's
// own address to it. That is the whole of why this is an object file rather
// than a blob.

const (
	resourceDirSize   = 16
	resourceEntrySize = 8
	resourceLeafSize  = 16
	// The high bit of an entry's offset says it points at another directory
	// rather than at a leaf.
	resourceSubdir = 0x80000000
)

// resourceSection is the bytes of the section and where the addresses in it
// are, so that relocations can be written for them.
type resourceSection struct {
	data []byte
	// fixups are offsets into data holding a section-relative address the
	// linker has to turn into a real one.
	fixups []uint32
}

func buildResourceSection(pictures []iconPicture) resourceSection {
	// What is in it: one leaf per icon, then one for the group.
	type blob struct {
		id     uint16
		typeID uint16
		data   []byte
	}
	var blobs []blob
	for i, picture := range pictures {
		blobs = append(blobs, blob{uint16(i + 1), rtIcon, picture.data})
	}
	blobs = append(blobs, blob{1, rtGroupIcon, groupIcon(pictures)})

	// The two types, in the order Windows wants them: by number.
	icons := len(pictures)

	// Where everything lands. The directories first, then the leaves, then
	// the bytes.
	at := uint32(resourceDirSize + 2*resourceEntrySize) // the root: two types
	typeDirs := map[uint16]uint32{}

	typeDirs[rtIcon] = at
	at += uint32(resourceDirSize + icons*resourceEntrySize)
	nameDirs := make([]uint32, 0, icons+1)
	for range icons {
		nameDirs = append(nameDirs, at)
		at += resourceDirSize + resourceEntrySize // one language
	}

	typeDirs[rtGroupIcon] = at
	at += resourceDirSize + resourceEntrySize
	groupNameDir := at
	at += resourceDirSize + resourceEntrySize

	leaves := make([]uint32, len(blobs))
	for i := range blobs {
		leaves[i] = at
		at += resourceLeafSize
	}

	// The bytes, each starting on an eight-byte boundary — which is what
	// every resource compiler does and what keeps a reader from having to
	// care.
	starts := make([]uint32, len(blobs))
	for i := range blobs {
		at = (at + 7) &^ 7
		starts[i] = at
		at += uint32(len(blobs[i].data))
	}

	out := make([]byte, at)
	put32 := func(offset, value uint32) {
		binary.LittleEndian.PutUint32(out[offset:], value)
	}
	put16 := func(offset uint32, value uint16) {
		binary.LittleEndian.PutUint16(out[offset:], value)
	}

	// A directory header: no name entries, and however many by number.
	dir := func(offset uint32, count int) {
		put16(offset+12, 0)             // named entries
		put16(offset+14, uint16(count)) // entries by number
	}
	// One entry of a directory: what it is called, and where it points.
	entry := func(offset uint32, id uint16, target uint32, isDir bool) {
		put32(offset, uint32(id))
		if isDir {
			target |= resourceSubdir
		}
		put32(offset+4, target)
	}

	// The root, by type.
	dir(0, 2)
	entry(resourceDirSize, rtIcon, typeDirs[rtIcon], true)
	entry(resourceDirSize+resourceEntrySize, rtGroupIcon, typeDirs[rtGroupIcon], true)

	// The icons, by number.
	dir(typeDirs[rtIcon], icons)
	for i := range icons {
		entry(typeDirs[rtIcon]+resourceDirSize+uint32(i)*resourceEntrySize,
			uint16(i+1), nameDirs[i], true)
		// And its one language.
		dir(nameDirs[i], 1)
		entry(nameDirs[i]+resourceDirSize, langNeutral, leaves[i], false)
	}

	// The group, which is the one Windows looks up to find the rest.
	dir(typeDirs[rtGroupIcon], 1)
	entry(typeDirs[rtGroupIcon]+resourceDirSize, 1, groupNameDir, true)
	dir(groupNameDir, 1)
	entry(groupNameDir+resourceDirSize, langNeutral, leaves[len(blobs)-1], false)

	// The leaves, and the bytes they point at.
	var fixups []uint32
	for i, one := range blobs {
		put32(leaves[i], starts[i]) // the linker adds the section's address
		put32(leaves[i]+4, uint32(len(one.data)))
		put32(leaves[i]+8, 0) // code page: none, the data is not text
		put32(leaves[i]+12, 0)
		fixups = append(fixups, leaves[i])
		copy(out[starts[i]:], one.data)
	}
	return resourceSection{data: out, fixups: fixups}
}

// langNeutral is the language a resource has when it is the same in every
// language, which a picture is.
const langNeutral = 0

// --- the object file -----------------------------------------------------------

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

// WindowsMachines is the architectures a resource can be written for: the
// three Go builds Windows programs for.
func WindowsMachines() []string { return []string{"amd64", "386", "arm64"} }

// writeCOFF writes the object file: one section, its relocations, and the two
// symbols a relocation needs to point at something.
func writeCOFF(w io.Writer, section resourceSection) error {
	return writeCOFFFor(w, "amd64", section)
}

func writeCOFFFor(w io.Writer, arch string, section resourceSection) error {
	machine, ok := coffMachines[arch]
	if !ok {
		return fmt.Errorf("antui: no Windows resource for %q", arch)
	}

	const headerSize = 20
	const sectionHeaderSize = 40
	rawAt := uint32(headerSize + sectionHeaderSize)
	relocationsAt := rawAt + uint32(len(section.data))
	symbolsAt := relocationsAt + uint32(10*len(section.fixups))

	var out bytes.Buffer
	// The file header.
	binary.Write(&out, binary.LittleEndian, machine.machine)
	binary.Write(&out, binary.LittleEndian, uint16(1)) // one section
	binary.Write(&out, binary.LittleEndian, uint32(0)) // no timestamp: same bytes every time
	binary.Write(&out, binary.LittleEndian, symbolsAt)
	binary.Write(&out, binary.LittleEndian, uint32(2)) // the section symbol and its aux
	binary.Write(&out, binary.LittleEndian, uint16(0)) // no optional header
	binary.Write(&out, binary.LittleEndian, uint16(0))

	// The section header.
	out.Write([]byte{'.', 'r', 's', 'r', 'c', 0, 0, 0})
	binary.Write(&out, binary.LittleEndian, uint32(0)) // virtual size: an object has none
	binary.Write(&out, binary.LittleEndian, uint32(0)) // and no address yet
	binary.Write(&out, binary.LittleEndian, uint32(len(section.data)))
	binary.Write(&out, binary.LittleEndian, rawAt)
	binary.Write(&out, binary.LittleEndian, relocationsAt)
	binary.Write(&out, binary.LittleEndian, uint32(0)) // no line numbers
	binary.Write(&out, binary.LittleEndian, uint16(len(section.fixups)))
	binary.Write(&out, binary.LittleEndian, uint16(0))
	// Initialised data, readable.
	binary.Write(&out, binary.LittleEndian, uint32(0x40000040))

	out.Write(section.data)

	// The relocations, all of them against the section's own symbol: the
	// value already in the field is the offset from the start of the section,
	// and the linker adds where the section ended up.
	for _, at := range section.fixups {
		binary.Write(&out, binary.LittleEndian, at)
		binary.Write(&out, binary.LittleEndian, uint32(0)) // the first symbol
		binary.Write(&out, binary.LittleEndian, machine.relocation)
	}

	// The symbol table: the section, and the auxiliary record that says how
	// big it is. Eight characters is exactly what ".rsrc" fits in, so there
	// is no string table to write.
	out.Write([]byte{'.', 'r', 's', 'r', 'c', 0, 0, 0})
	binary.Write(&out, binary.LittleEndian, uint32(0)) // value: the start of the section
	binary.Write(&out, binary.LittleEndian, uint16(1)) // section number, from one
	binary.Write(&out, binary.LittleEndian, uint16(0)) // type: none
	out.WriteByte(3)                                   // storage class: static
	out.WriteByte(1)                                   // one auxiliary record

	// The auxiliary section record.
	binary.Write(&out, binary.LittleEndian, uint32(len(section.data)))
	binary.Write(&out, binary.LittleEndian, uint16(len(section.fixups)))
	binary.Write(&out, binary.LittleEndian, uint16(0)) // line numbers
	binary.Write(&out, binary.LittleEndian, uint32(0)) // checksum
	binary.Write(&out, binary.LittleEndian, uint16(0)) // number of the section this is a copy of
	// The selection byte and three unused ones. An auxiliary record is
	// eighteen bytes exactly, the same as a symbol, and one byte short of
	// that makes every reader lose the string table.
	out.Write([]byte{0, 0, 0, 0})

	// An empty string table, which is four bytes saying its own size. A COFF
	// file without one is a COFF file some readers refuse.
	binary.Write(&out, binary.LittleEndian, uint32(4))

	_, err := w.Write(out.Bytes())
	return err
}

// WriteWindowsResourceFor is WriteWindowsResource for an architecture other
// than the usual one.
func WriteWindowsResourceFor(w io.Writer, arch string, images ...*canvas.Canvas) error {
	pictures, err := iconPictures(images)
	if err != nil {
		return err
	}
	return writeCOFFFor(w, arch, buildResourceSection(pictures))
}