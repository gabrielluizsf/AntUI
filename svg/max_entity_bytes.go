package svg

// maxEntityBytes is how much text the names a doctype defines may be read
// into, all of them together, before the drawing is refused. A name may hold
// ten of another and that one ten of the next: the "billion laughs" of
// CVE-2003-1564 is a file of half a kilobyte that reads into gigabytes, and
// stopping only at a name that comes back to itself does nothing about it. A
// drawing that leans on its doctype for a namespace or a repeated colour reads
// a few hundred bytes out of it, never megabytes.
const maxEntityBytes = 1 << 20
