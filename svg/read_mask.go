package svg

import "strings"

// readMask reads one `<mask>`'s own sayings — the way it measures and where it
// reaches — into the place its id already put, before any picture is built.
// See [Image.readMasks] for why it is a separate step.
func (img *Image) readMask(id string, e *element) {
	warn := func(format string, args ...any) { img.warn(e, format, args...) }
	def := img.masks[id]
	def.measure = maskTypeIn(e, warn)
	units := strings.ToLower(strings.TrimSpace(e.attr("maskContentUnits")))
	def.contentBoxUnits = units != "" && units != "userspaceonuse"
	def.region = img.readMaskRegion(e, warn)
}
