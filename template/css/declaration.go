package css

// Declaration is a single property/value pair.
type Declaration struct {
	Prop      string
	Raw       string
	Important bool // the declaration carried !important
}
