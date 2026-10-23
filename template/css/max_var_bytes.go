package css

// maxVarBytes is how many bytes the var()s of one value may read in, all of
// them together. A real value — a colour, a length, a list of a few gradients
// — is tens of bytes; this is room for any of them many times over, and far
// short of the gigabytes a value naming itself through ten others reads into.
const maxVarBytes = 256 << 10
