package antui

import (
	"bytes"
	"fmt"
	"io"
)

func writeCOFFFor(w io.Writer, arch string, section resourceSection) error {
	machine, ok := coffMachines[arch]
	if !ok {
		return fmt.Errorf("antui: no Windows resource for %q", arch)
	}

	var out bytes.Buffer
	writeCOFFHeader(&out, section, machine.machine)
	out.Write(section.data)
	writeCOFFRelocations(&out, section, machine.relocation)
	writeCOFFSymbols(&out, section)

	_, err := w.Write(out.Bytes())
	return err
}
