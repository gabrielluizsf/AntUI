package svg

import (
	"sync"
)

// pictureBytes is one network address as the whole process knows it: read once
// and kept, so every drawing that names the address shares the read instead of
// fetching again, and one goroutine fetches it however many drawings ask at
// once. A failed read is dropped when its waiters have been told, so the next
// drawing that names the address asks the server again rather than being
// poisoned by one moment of it being unreachable.
type pictureBytes struct {
	mu      sync.Mutex
	read    pictureRead
	done    bool
	waiting []func(*pictureFrames, error)
}

// pictureByteCache is the entries by resolved address. It is one table for the
// process rather than one per drawing, which is what makes a picture read for
// one drawing available to the next.
var pictureByteCache = struct {
	mu      sync.Mutex
	entries map[string]*pictureBytes
}{entries: map[string]*pictureBytes{}}

// pictureBytesFor is the entry for an address, made and set fetching by whoever
// asks first: the fetch is started with the table held, so one address is read
// once even when two drawings name it at the same moment.
func pictureBytesFor(target string) *pictureBytes {
	pictureByteCache.mu.Lock()
	defer pictureByteCache.mu.Unlock()
	if e := pictureByteCache.entries[target]; e != nil {
		return e
	}
	e := &pictureBytes{}
	pictureByteCache.entries[target] = e
	go fetchPictureLater(target, e)
	return e
}

// forgetPictureBytes drops a failed read from the table.
func forgetPictureBytes(target string) {
	pictureByteCache.mu.Lock()
	delete(pictureByteCache.entries, target)
	pictureByteCache.mu.Unlock()
}
