package apk

// DefaultStore is what [Config.Store] leaves uncompressed: everything that
// is already compressed, so that the package does not compress it twice.
var DefaultStore = []string{
	"png", "jpg", "jpeg", "webp", "gif",
	"ogg", "mp3", "m4a", "opus", "wav",
	"mp4", "m4v", "webm", "mkv",
	"zip", "gz", "br", "woff", "woff2",
}
