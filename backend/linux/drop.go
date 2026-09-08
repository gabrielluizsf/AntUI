//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"net/url"
	"strings"
	"time"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

// Files dragged onto the window, and files copied in a file manager.
//
// Both are the same question asked twice — "what files does the rest of the
// desktop want to give me?" — and on X11 both are answered through
// selections, which is the protocol's way of saying that one client holds
// something another wants and the two have to talk about it.
//
// Dragging is XDND: the source tells us it is over us, we say we will take
// it, it drops, we ask for the drop's contents as `text/uri-list`, it puts
// them on a property of our window and tells us, and we read the property.
// Five messages for one drag, which is what makes it worth writing once.
//
// Pasting is the same conversation without the drag: ask the CLIPBOARD
// selection for `text/uri-list`, and for `UTF8_STRING` when it has no files.

// XDND, and the version we speak. Version 5 is what every toolkit has spoken
// for fifteen years.
const xdndVersion = 5

// dropAtoms are looked up once, the first time they are wanted: a window that
// nobody ever drags anything onto should not pay for six round trips.
func (x *Driver) dropAtoms() {
	if x.atomXdndAware != 0 {
		return
	}
	x.atomXdndAware = x.internAtom("XdndAware")
	x.atomXdndEnter = x.internAtom("XdndEnter")
	x.atomXdndPosition = x.internAtom("XdndPosition")
	x.atomXdndStatus = x.internAtom("XdndStatus")
	x.atomXdndLeave = x.internAtom("XdndLeave")
	x.atomXdndDrop = x.internAtom("XdndDrop")
	x.atomXdndFinished = x.internAtom("XdndFinished")
	x.atomXdndSelection = x.internAtom("XdndSelection")
	x.atomXdndActionCopy = x.internAtom("XdndActionCopy")
	x.atomUriList = x.internAtom("text/uri-list")
	x.atomClipboard = x.internAtom("CLIPBOARD")
	x.atomDropProperty = x.internAtom("ANTUI_DROP")
	x.atomTargets = x.internAtom("TARGETS")
	x.atomText = x.internAtom("TEXT")
}

// announceDrops says the window will take a drag. Without this property no
// source will offer it anything at all.
func (x *Driver) announceDrops() {
	x.dropAtoms()
	if x.atomXdndAware == 0 {
		return
	}
	x.changeProperty(x.atomXdndAware, atomAtom, 32, le32(xdndVersion), 1)
}

// handleDrop answers one of the drag messages. It reports whether the message
// was one of them.
func (x *Driver) handleDrop(packet []byte) bool {
	x.dropAtoms()
	kind := binary.LittleEndian.Uint32(packet[8:])
	data := packet[12:]

	switch kind {
	case x.atomXdndEnter:
		x.dragFrom = binary.LittleEndian.Uint32(data[0:])
		return true

	case x.atomXdndPosition:
		// Where the pointer is, so that a window with more than one place to
		// drop into knows which it is over. The coordinates are on the root,
		// so they are turned into window ones the same way the pointer's are.
		at := binary.LittleEndian.Uint32(data[8:])
		x.mouseX = int(int16(at>>16)) - x.windowX
		x.mouseY = int(int16(at&0xFFFF)) - x.windowY
		x.win.SetMouse(x.mouseX, x.mouseY)
		x.win.Push(backend.Event{Type: backend.EventMouseMove, X: x.mouseX, Y: x.mouseY})

		x.dragFrom = binary.LittleEndian.Uint32(data[0:])
		x.sendDragStatus(true)
		return true

	case x.atomXdndLeave:
		x.dragFrom = 0
		return true

	case x.atomXdndDrop:
		x.dragFrom = binary.LittleEndian.Uint32(data[0:])
		when := binary.LittleEndian.Uint32(data[8:])
		files := x.askSelection(x.atomXdndSelection, x.atomUriList, when)
		x.finishDrag()
		if paths := parseURIList(files); len(paths) > 0 {
			x.win.Push(backend.Event{Type: backend.EventDropFiles, Files: paths,
				X: x.mouseX, Y: x.mouseY})
		}
		return true
	}
	return false
}

// sendDragStatus tells the source whether we will take what it is dragging.
func (x *Driver) sendDragStatus(accept bool) {
	if x.dragFrom == 0 {
		return
	}
	accepted := uint32(0)
	action := uint32(0)
	if accept {
		accepted = 1
		action = x.atomXdndActionCopy
	}
	x.sendDragMessage(x.atomXdndStatus, [5]uint32{
		x.window, accepted, 0, 0, action,
	})
}

// finishDrag tells the source we are done with the drop, which is what
// releases it to carry on.
func (x *Driver) finishDrag() {
	if x.dragFrom == 0 {
		return
	}
	x.sendDragMessage(x.atomXdndFinished, [5]uint32{
		x.window, 1, x.atomXdndActionCopy, 0, 0,
	})
	x.dragFrom = 0
}

func (x *Driver) sendDragMessage(kind uint32, data [5]uint32) {
	body := make([]byte, 40)
	binary.LittleEndian.PutUint32(body[0:], x.dragFrom)
	binary.LittleEndian.PutUint32(body[4:], 0) // straight to the client

	event := body[8:]
	event[0] = xClientMessage
	event[1] = 32 // format
	binary.LittleEndian.PutUint32(event[4:], x.dragFrom)
	binary.LittleEndian.PutUint32(event[8:], kind)
	for i, value := range data {
		binary.LittleEndian.PutUint32(event[12+i*4:], value)
	}
	x.send(xSendEvent, 0, body)
}

// askSelection asks whoever holds a selection to write it onto a property of
// our window, waits for them to say they have, and reads it.
//
// The wait is bounded, and generously: a selection is held by another process,
// and a process that has stopped answering must not stop this one — but a
// busy machine can take a while to schedule it, and a drop that quietly
// brings in nothing because the answer arrived a moment late is worse than a
// frame that took a second. A second is long past what any live source needs
// and far short of a hang.
func (x *Driver) askSelection(selection, target, when uint32) []byte {
	if !x.alive || selection == 0 || target == 0 {
		return nil
	}
	x.send(xConvertSelection, 0, le32(x.window, selection, target,
		x.atomDropProperty, when))

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		packet, _, ok := x.takePacket()
		if !ok {
			x.fill(20 * time.Millisecond)
			continue
		}
		if packet[0]&0x7F != xSelectionNotify {
			if packet[0] != xReply && x.win != nil {
				x.handleEvent(packet)
			}
			continue
		}
		// The property is None when the holder could not give us what was
		// asked for — a file manager with text on the clipboard and nothing
		// that looks like a file, say.
		if binary.LittleEndian.Uint32(packet[20:]) == 0 {
			return nil
		}
		return x.readProperty(x.atomDropProperty)
	}
	return nil
}

// readProperty reads a property off our own window and deletes it, which is
// what a selection's protocol asks of whoever receives one.
func (x *Driver) readProperty(property uint32) []byte {
	var out []byte
	for offset := uint32(0); ; {
		// Delete on the last read only, so a property that needs two reads is
		// not thrown away halfway through.
		seq := x.send(xGetProperty, 0,
			le32(x.window, property, 0, offset, 4096))
		reply, extra, ok := x.awaitReply(seq)
		if !ok {
			break
		}
		// The reply pads its data to a four-byte boundary, and a property
		// whose length is not a multiple of four comes back with the padding
		// still attached — a 26-character copy as 28 bytes with two zeros on
		// the end. The reply says how many items are really there, which is
		// what is cut to. (The original port of this function carried the
		// padding; the integration test caught it on a read-back.)
		items := int(binary.LittleEndian.Uint32(reply[16:]))
		length := items * int(reply[1]) / 8
		if length < 0 || length > len(extra) {
			break
		}
		out = append(out, extra[:length]...)
		if binary.LittleEndian.Uint32(reply[12:]) == 0 { // bytes-after
			break
		}
		offset += uint32(length / 4)
	}
	// Deleting says "I have taken it", which is what the holder waits for.
	x.send(xGetProperty, 1, le32(x.window, property, 0, 0, 0))
	return out
}

// Clipboard is what the X selections hold: the files they name, and their
// text.
func (x *Driver) Clipboard() (string, []string) {
	if !x.alive {
		return "", nil
	}
	x.dropAtoms()

	// Files first: a file manager puts both on the clipboard, and the files
	// are the more useful answer when there are any.
	if data := x.askSelection(x.atomClipboard, x.atomUriList, 0); len(data) > 0 {
		if paths := parseURIList(data); len(paths) > 0 {
			return string(data), paths
		}
	}
	if x.atomUTF8String != 0 {
		if data := x.askSelection(x.atomClipboard, x.atomUTF8String, 0); len(data) > 0 {
			return string(data), parseURIList(data)
		}
	}
	data := x.askSelection(x.atomClipboard, atomString, 0)
	return string(data), parseURIList(data)
}

// parseURIList turns what a drag or a clipboard holds into paths on this
// machine.
//
// The format is one URI a line, with comment lines beginning with a hash, and
// the lines end with a carriage return as well as a newline. Anything that is
// not a local file — a web address, a file on another machine — is left out
// rather than handed on as a path that does not exist.
func parseURIList(data []byte) []string {
	var out []string
	for line := range strings.Lines(string(data)) {
		line = strings.TrimRight(line, "\r\n")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if path, ok := pathFromURI(line); ok {
			out = append(out, path)
		}
	}
	return out
}

// pathFromURI is the local path a file:// URI names, and false for anything
// else.
func pathFromURI(uri string) (string, bool) {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return "", false
	}
	if !strings.HasPrefix(uri, "file://") {
		// A bare path is what some sources send, and what a person pasting a
		// path into the window means.
		if strings.HasPrefix(uri, "/") {
			return uri, true
		}
		return "", false
	}
	rest := uri[len("file://"):]
	// file://host/path is a file on another machine unless the host is this
	// one, which is written as an empty host or `localhost`.
	if at := strings.IndexByte(rest, '/'); at > 0 {
		if host := rest[:at]; host != "localhost" {
			return "", false
		}
		rest = rest[at:]
	}
	unescaped, err := url.PathUnescape(rest)
	if err != nil {
		return "", false
	}
	return unescaped, true
}

// --- the window's icon ---------------------------------------------------------

// SetIcon gives the window an icon for the task bar, by writing _NET_WM_ICON,
// which is what every window manager from this century reads.
//
// The property is one long list of 32-bit numbers: for each size, its width,
// its height, and then its pixels as ARGB, a row at a time. Several sizes go
// one after another in the same property and the window manager picks the one
// it wants — which is why this takes a list rather than a picture.
//
// The values are CARDINAL, and a CARDINAL is 32 bits on the wire whatever the
// machine is, so a pixel is one number and not two.
func (x *Driver) SetIcon(images []*canvas.Canvas) bool {
	if !x.alive || len(images) == 0 {
		return false
	}
	if x.atomNetWMIcon == 0 {
		x.atomNetWMIcon = x.internAtom("_NET_WM_ICON")
		if x.atomNetWMIcon == 0 {
			return false
		}
	}

	count := 0
	for _, image := range images {
		count += 2 + image.Width*image.Height
	}
	// A property is sent in one request, so it has to fit in one. The biggest
	// sizes go in first — they are what a task switcher wants — and anything
	// past what will fit is left out rather than making a request the server
	// refuses whole.
	most := int(x.maxRequestBytes/4) - 8
	if most <= 0 {
		most = 16384
	}

	values := make([]uint32, 0, min(count, most))
	for _, image := range images {
		need := 2 + image.Width*image.Height
		if len(values)+need > most {
			continue
		}
		values = append(values, uint32(image.Width), uint32(image.Height))
		for y := range image.Height {
			for x := range image.Width {
				values = append(values, uint32(image.At(x, y)))
			}
		}
	}
	if len(values) == 0 {
		return false
	}

	x.changeProperty(x.atomNetWMIcon, atomCardinal, 32,
		le32(values...), uint32(len(values)))
	return x.alive
}

// --- being the one who holds the clipboard -------------------------------------

// The rest of the selection protocol: answering, rather than asking.
//
// X11 has no clipboard of its own. Whoever copied something *is* the
// clipboard: they own the selection, they keep the bytes, and every paste is
// a question put to them. So a window that wants Copy to work has to stay
// and answer, which is what this is.

// SetClipboard puts text on the X selections, taking the clipboard over and
// keeping the text to hand out.
func (x *Driver) SetClipboard(text string) bool {
	if !x.alive || x.window == 0 {
		return false
	}
	x.dropAtoms()
	if x.atomClipboard == 0 {
		return false
	}
	if x.atomTargets == 0 {
		x.atomTargets = x.internAtom("TARGETS")
		x.atomText = x.internAtom("TEXT")
	}

	x.copied = text
	// SetSelectionOwner: window, selection, time. CurrentTime is nought,
	// which the protocol allows and every toolkit uses for a copy that
	// happens now.
	x.send(xSetSelectionOwner, 0, le32(x.window, x.atomClipboard, 0))
	return x.alive
}

// answerSelection hands over what was copied. It reports whether the event
// was one of the selection ones.
func (x *Driver) answerSelection(packet []byte) bool {
	switch packet[0] & 0x7F {
	case xSelectionClear:
		// Somebody else copied something: what this window holds is no longer
		// the clipboard, and keeping it would be keeping a lie.
		x.copied = ""
		return true

	case xSelectionRequest:
		// The fields, in the order the protocol puts them: the time, who is
		// asking, which selection, what they want it as, and where to put it.
		when := binary.LittleEndian.Uint32(packet[4:])
		requestor := binary.LittleEndian.Uint32(packet[12:])
		target := binary.LittleEndian.Uint32(packet[20:])
		property := binary.LittleEndian.Uint32(packet[24:])
		if property == 0 {
			// An old-style requestor: the target is where the answer goes.
			property = target
		}

		switch {
		case target == x.atomTargets:
			// What this window can give: the list itself, and the three ways
			// of asking for text that are still in use.
			kinds := le32(x.atomTargets, x.atomUTF8String, atomString, x.atomText)
			x.changePropertyOn(requestor, property, atomAtom, 32, kinds, 4)

		case target == x.atomUTF8String || target == atomString || target == x.atomText:
			x.changePropertyOn(requestor, property, target, 8,
				[]byte(x.copied), uint32(len(x.copied)))

		default:
			property = 0 // nothing this window can turn it into
		}
		x.sendSelectionNotify(requestor, property, target, when)
		return true
	}
	return false
}

// sendSelectionNotify tells the asker the answer is on their window, or that
// there is no answer.
func (x *Driver) sendSelectionNotify(requestor, property, target, when uint32) {
	body := make([]byte, 40)
	binary.LittleEndian.PutUint32(body[0:], requestor)
	binary.LittleEndian.PutUint32(body[4:], 0) // straight to the client

	event := body[8:]
	event[0] = xSelectionNotify
	binary.LittleEndian.PutUint32(event[4:], when)
	binary.LittleEndian.PutUint32(event[8:], requestor)
	binary.LittleEndian.PutUint32(event[12:], x.atomClipboard)
	binary.LittleEndian.PutUint32(event[16:], target)
	binary.LittleEndian.PutUint32(event[20:], property)

	x.send(xSendEvent, 0, body)
}