package css

// parseImport reads another stylesheet in where the @import stands, so its
// rules land in this sheet at that point, carrying the media the statement is
// inside and the media query the statement itself names. Every reason the file
// cannot be read is a warning rather than an error: a stylesheet that imports
// something this canvas cannot open still styles everything else in it.
func (p *parser) parseImport(sh *Sheet, outer Media) error {
	prelude, semi, err := p.readHeader(true)
	if err != nil {
		// The statement never ended, and it is the last thing left to read:
		// either the file was cut off before its ';', or what follows is
		// inside a string that never closes. Both end the sheet there.
		sh.Warn = append(sh.Warn, fmtErrf("ignoring an unfinished @import").Error())
		p.i = len(p.src)
		return nil
	}
	if !semi {
		// An @import has no block, but the one written all the same is
		// swallowed so the rest of the sheet still is.
		if p.peek() == '{' {
			p.next()
			if err := p.skipBlock(); err != nil {
				return err
			}
		}
		sh.Warn = append(sh.Warn, fmtErrf("ignoring an @import with a block").Error())
		return nil
	}
	if p.peek() == ';' {
		p.next()
	}
	ref, mediaText, warns := splitImport(prelude)
	sh.Warn = append(sh.Warn, warns...)
	if ref == "" {
		return nil
	}
	return p.importSheet(sh, outer, importRef{name: ref, media: mediaText})
}
