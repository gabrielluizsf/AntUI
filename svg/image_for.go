package svg

import "slices"

// imageFor prepares the drawing parseText is about to read: where it came from,
// the memory of the gradients and pictures it will name, and its own `<style>`
// sheets, which are matched against every element before anything is built —
// because a rule reaches the shapes a clip, a mask, a marker or a pattern holds
// as much as the ones painted where they stand.
func imageFor(root *element, path string, chain []string) *Image {
	img := &Image{base: path, chain: append(slices.Clone(chain), path), Root: &Node{Name: "svg"}, grads: map[string]*gradient{}, pictures: map[string]*pictureFrames{}}
	img.readStyles(root)
	return img
}
