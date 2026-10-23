package font

// maxColorGlyphSide is how many pixels a colour glyph may be on a side. An
// emoji font cuts its pictures at a hundred and some pixels, and the largest
// strikes in use are a few hundred; a glyph that says it is thousands wide is
// a decompression bomb rather than a letter.
const maxColorGlyphSide = 2048
