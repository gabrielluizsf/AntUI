package svg

// fetchPictureLater reads an address off the goroutine that is reading the
// drawing, so a server that is slow or gone cannot hold up a frame, and hands
// the answer to every element that is waiting on it. The picture is decoded
// once, here, so the elements sharing the address share the canvas too. The
// entry is settled before they are told, and a read that did not answer is
// dropped from the table on the way out, so the next drawing that names the
// address asks the server again.
func fetchPictureLater(target string, e *pictureBytes) {
	data, fetchErr := fetchPicture(target)
	read := pictureRead{data: data, err: fetchErr}
	pic, err := read.decoded()

	e.mu.Lock()
	e.read, e.done = read, true
	waiting := e.waiting
	e.waiting = nil
	e.mu.Unlock()

	if fetchErr != nil {
		forgetPictureBytes(target)
	}
	for _, tell := range waiting {
		tell(pic, err)
	}
}
