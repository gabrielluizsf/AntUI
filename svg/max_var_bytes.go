package svg

// maxVarBytes is how long one value may grow while its var()s are read into
// it. A real value — a colour, a length, a list of a few gradients — is tens
// of bytes; this is room for any of them many times over, and far short of the
// gigabytes a value naming itself through ten others can read into.
const maxVarBytes = 256 << 10
