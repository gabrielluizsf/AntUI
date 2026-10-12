//go:build android

package camera

// Info is which camera this is.
func (p *Preview) Info() Info { return p.info }
