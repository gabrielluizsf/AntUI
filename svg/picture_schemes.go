package svg

// pictureSchemes is what each scheme means to a picture: the ways this drawing
// knows how to get bytes out of an address. A scheme that is not here is said
// out loud rather than handed to the file system, so an `ftp:` or a `blob:`
// address is a resource that failed rather than a path read by accident.
var pictureSchemes = map[string]func(string) ([]byte, error){
	"data":  dataURI,
	"http":  getPicture,
	"https": getPicture,
	"file":  pictureFile,
}
