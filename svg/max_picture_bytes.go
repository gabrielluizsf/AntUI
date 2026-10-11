package svg

// maxPictureBytes is how big a picture fetched over the network may be before
// it is refused rather than read: a server that answers with something the
// size of a disk should not be able to fill memory before the decode ever
// gets a say, and a truncated picture fails to decode with a message of its
// own anyway, so the refusal at the door is the clearer of the two.
const maxPictureBytes = 64 << 20
