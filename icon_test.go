package antui

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// a picture with something in it, at a size worth scaling.
func aLogo(t *testing.T, size int) *canvas.Canvas {
	t.Helper()
	cv, err := canvas.NewCanvas(size, size)
	if err != nil {
		t.Fatal(err)
	}
	for y := range size {
		for x := range size {
			// A quarter of it transparent, so the scaling has alpha to carry.
			if x < size/2 && y < size/2 {
				cv.Put(x, y, canvas.RGBA(0, 0, 0, 0))
				continue
			}
			cv.Put(x, y, canvas.RGBA(uint8(x*255/size), uint8(y*255/size), 200, 255))
		}
	}
	return cv
}

func TestScalingAnIcon(t *testing.T) {
	big := aLogo(t, 256)
	small := canvas.IconScaled(big, 32)
	if small == nil || small.Width != 32 || small.Height != 32 {
		t.Fatalf("the icon came out %v", small)
	}
	// The transparent quarter is still transparent, and the rest is not.
	if got := small.At(2, 2); got.A() != 0 {
		t.Errorf("a transparent corner came out %v", got)
	}
	if got := small.At(30, 30); got.A() != 255 {
		t.Errorf("an opaque corner came out %v", got)
	}
	// And the colour is roughly what was there, rather than a smear towards
	// the transparent black it was averaged with.
	want := big.At(30*8, 30*8)
	got := small.At(30, 30)
	if awayFar(int(got.R()), int(want.R())) > 12 ||
		awayFar(int(got.G()), int(want.G())) > 12 {
		t.Errorf("the colour came out %v, want about %v", got, want)
	}

	if canvas.IconScaled(nil, 32) != nil || canvas.IconScaled(big, 0) != nil {
		t.Error("nothing came back as something")
	}
}

func awayFar(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func TestASetOfIconSizes(t *testing.T) {
	set := canvas.IconSet(aLogo(t, 256))
	if len(set) != len(canvas.IconSizes) {
		t.Fatalf("%d sizes, want %d", len(set), len(canvas.IconSizes))
	}
	// Largest first, which is the order every window system here wants.
	for i := 1; i < len(set); i++ {
		if set[i].Width >= set[i-1].Width {
			t.Errorf("size %d is %d, after %d", i, set[i].Width, set[i-1].Width)
		}
	}
	// Nothing is scaled up past what it was made from.
	small := canvas.IconSet(aLogo(t, 24))
	for _, one := range small {
		if one.Width > 24 && one.Width != 24 {
			if one.Width > 32 {
				t.Errorf("a 24-pixel picture was blown up to %d", one.Width)
			}
		}
	}
}

func TestWritingAnICO(t *testing.T) {
	var out bytes.Buffer
	if err := WriteICO(&out, canvas.IconSet(aLogo(t, 64))...); err != nil {
		t.Fatal(err)
	}
	data := out.Bytes()

	if binary.LittleEndian.Uint16(data[0:]) != 0 ||
		binary.LittleEndian.Uint16(data[2:]) != 1 {
		t.Fatal("it does not begin as an icon file")
	}
	count := int(binary.LittleEndian.Uint16(data[4:]))
	if count == 0 {
		t.Fatal("it holds no pictures")
	}

	// Every entry points at a PNG inside the file, and the sizes come down.
	last := 999
	for i := range count {
		entry := data[6+i*16:]
		width := int(entry[0])
		if width == 0 {
			width = 256
		}
		if width > last {
			t.Errorf("picture %d is %d wide, after %d", i, width, last)
		}
		last = width

		size := binary.LittleEndian.Uint32(entry[8:])
		at := binary.LittleEndian.Uint32(entry[12:])
		if int(at+size) > len(data) {
			t.Fatalf("picture %d runs past the end of the file", i)
		}
		if !bytes.HasPrefix(data[at:], []byte("\x89PNG\r\n\x1a\n")) {
			t.Errorf("picture %d is not a PNG", i)
		}
	}

	// A picture bigger than an icon can be is left out rather than written
	// with a size that means something else.
	var only bytes.Buffer
	if err := WriteICO(&only, aLogo(t, 512)); err == nil {
		t.Error("a 512-pixel icon was written")
	}
}

// The resource file has to be an object file the toolchain accepts, which is
// something the debug/pe reader can say better than any assertion about
// bytes.
func TestWritingAWindowsResource(t *testing.T) {
	for _, arch := range WindowsMachines() {
		var out bytes.Buffer
		if err := WriteWindowsResourceFor(&out, arch, canvas.IconSet(aLogo(t, 64))...); err != nil {
			t.Fatalf("%s: %v", arch, err)
		}

		file, err := pe.NewFile(bytes.NewReader(out.Bytes()))
		if err != nil {
			t.Fatalf("%s: it is not an object file: %v", arch, err)
		}
		defer file.Close()

		if len(file.Sections) != 1 || file.Sections[0].Name != ".rsrc" {
			t.Fatalf("%s: the sections are %v", arch, file.Sections)
		}
		section := file.Sections[0]
		if section.Size == 0 {
			t.Errorf("%s: the section is empty", arch)
		}
		if section.NumberOfRelocations == 0 {
			t.Errorf("%s: nothing is relocated, so the addresses in it are wrong", arch)
		}

		// Every relocation points at the section's own symbol, which is the
		// only one there is.
		for _, one := range section.Relocs {
			if one.SymbolTableIndex != 0 {
				t.Errorf("%s: a relocation points at symbol %d",
					arch, one.SymbolTableIndex)
			}
			if int(one.VirtualAddress) >= int(section.Size) {
				t.Errorf("%s: a relocation is outside the section", arch)
			}
		}

		// The directory begins with two types: icons, and the group that
		// names them.
		data, err := section.Data()
		if err != nil {
			t.Fatal(err)
		}
		if got := binary.LittleEndian.Uint16(data[14:]); got != 2 {
			t.Errorf("%s: the root holds %d types, want two", arch, got)
		}
		if got := binary.LittleEndian.Uint32(data[16:]); got != rtIcon {
			t.Errorf("%s: the first type is %d, want RT_ICON", arch, got)
		}
		if got := binary.LittleEndian.Uint32(data[24:]); got != rtGroupIcon {
			t.Errorf("%s: the second type is %d, want RT_GROUP_ICON", arch, got)
		}
	}

	if err := WriteWindowsResourceFor(nil, "sparc", aLogo(t, 32)); err == nil {
		t.Error("a resource was written for a machine Windows does not run on")
	}
}

// The window takes an icon, and refuses nothing at all.
func TestSettingAWindowIcon(t *testing.T) {
	stub := &stubBackend{}
	win := &Window{native: stub}

	if win.SetIcon() {
		t.Error("a window took no icon at all")
	}
	if win.SetIcon(nil, nil) {
		t.Error("a window took two nothings")
	}
	if !win.SetIcon(canvas.IconSet(aLogo(t, 64))...) {
		t.Fatal("the window would not take an icon")
	}
	if len(stub.icons) == 0 || stub.icons[0].Width < stub.icons[len(stub.icons)-1].Width {
		t.Error("the icons did not arrive largest first")
	}
}