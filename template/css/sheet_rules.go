package css

// Rules returns every rule in the sheet, in stylesheet order.
func (sh *Sheet) Rules() []*Rule { return sh.rules }
