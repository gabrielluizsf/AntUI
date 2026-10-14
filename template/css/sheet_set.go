package css

// SetProperty overrides a property for every element with the given classes,
// as a stylesheet appended after all the others — so it wins the cascade
// against anything the sheet said, while media queries still test it. Set
// classes to nil to target the element by tag alone.
func (sh *Sheet) SetProperty(q PropertyQuery, raw string) {
	s := Selector{Tag: q.Tag, Classes: append([]string(nil), q.Classes...)}
	sh.rules = append(sh.rules, &Rule{
		Selectors: []Selector{s},
		Decls:     []Declaration{{Prop: q.Prop, Raw: raw}},
		Order:     sh.order,
	})
	sh.order++
}
