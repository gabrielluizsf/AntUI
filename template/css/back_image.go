package css

// BackImage is one background layer: either a url() the template looks up in
// its image registry, or a gradient painted directly.
type BackImage struct {
	URL  string
	Grad *Gradient
}
