package template

// multiPieces bounds the fragments one box may be cut into. Each fragment is
// one pass over the container's callback, and a box that would need more than
// this — a very tall box in a very short column — is drawn whole in the piece
// that overflows instead, so a stylesheet cannot cost a frame a thousand
// passes.
const multiPieces = 64
