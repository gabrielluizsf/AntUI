package svg

// maxPicturePixels is how many pixels a picture an `<image>` points at may
// say it is, every frame of an animated one counted. The decoders of the
// standard library allocate the size a header claims before reading the
// pixels it promises, so seventy bytes of PNG that say they are sixty
// thousand pixels on a side ask for fourteen gigabytes — the decompression
// bomb of Cloudreve's CVE-2026-55497. A hundred and thirty million is an
// eleven-thousand-pixel photograph, or five hundred frames of a GIF five
// hundred pixels square.
const maxPicturePixels = 1 << 27
