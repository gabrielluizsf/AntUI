package css

// Selector is one simple selector that must match for a rule to apply.
// `.a.b:hover` carries the classes "a" and "b" and the Hover state. Descendant
// and child combinators are parsed but can never match: AntUI widgets are
// flat, so the engine keeps the rule with a warning and the selector never
// fires. The same goes for parts of the CSS3 selector language the flat
// model cannot evaluate — attributes, ids, pseudo-elements and structural
// pseudo-classes — all recorded in [Selector.Unsupported].
type Selector struct {
	Tag     string   // element name; "" means any element with the classes
	Classes []string // every class must be present
	State   State    // every bit must be present
	All     bool     // the universal selector *

	// Unsupported names the selector feature the flat widget model cannot
	// evaluate (a combinator, an attribute, an id, a pseudo-element, a
	// structural pseudo-class). Non-empty means the rule parsed fine but must
	// never match — the engine drops it with a warning rather than an error.
	Unsupported string
}
