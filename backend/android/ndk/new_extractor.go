//go:build android

package ndk

/*
#cgo LDFLAGS: -lmediandk -landroid

#include "media_cgo.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
)

// newExtractor opens a file and hands the extractor its descriptor.
//
// Not setDataSource, which takes a *URI* rather than a path — a bare path is
// refused with AMEDIA_ERROR_UNSUPPORTED, which says nothing about why. A
// descriptor also means the same code will one day serve a content:// address
// from the picker, which arrives as a descriptor and never as a path.
//
// The file has to stay open for as long as the extractor reads from it, so
// it is handed back to be closed with it.
func newExtractor(path string) (*C.AMediaExtractor, *os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	st, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	ex := C.AMediaExtractor_new()
	if ex == nil {
		f.Close()
		return nil, nil, errors.New("ndk: cannot make a media extractor")
	}
	if r := C.AMediaExtractor_setDataSourceFd(ex, C.int(f.Fd()), 0,
		C.off64_t(st.Size())); r != C.AMEDIA_OK {
		C.AMediaExtractor_delete(ex)
		f.Close()
		return nil, nil, fmt.Errorf("ndk: nothing on this device can read %s (%d)",
			path, int(r))
	}
	return ex, f, nil
}
