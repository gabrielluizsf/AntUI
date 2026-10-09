package css

// SetProperty overrides a property for every element with the given classes,
// as a stylesheet appended after all the others — so it wins the cascade
// against anything the sheet said, while media queries still test it. Set
// classes to nil to target the element by tag alone.
func (sh *Sheet) SetProperty(tag string, classes []string, prop, raw string) {
	s := Selector{Tag: tag, Classes: append([]string(nil), classes...)}
	sh.rules = append(sh.rules, &Rule{
		Selectors: []Selector{s},
		Decls:     []Declaration{{Prop: prop, Raw: raw}},
		Order:     sh.order,
	})
	sh.order++
}
