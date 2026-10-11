package svg

import "strings"

// attr reads one attribute, answering the empty string for one that was not
// there, so a caller can ask without checking twice.
//
// The name is matched without regard to case, because the names a drawing uses
// are written the way XML has them — `viewBox` and `preserveAspectRatio` are
// camel case, and everything else is not — and a drawing that spelled one of them
// in a case this package did not expect would otherwise be read as though it
// had said nothing at all.
func (e *element) attr(name string) string { return e.Attr[strings.ToLower(name)] }
