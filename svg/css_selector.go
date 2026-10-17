package svg

// complexSelector is a chain of compound selectors joined by combinators, read
// left to right, the rightmost standing on the element a rule is asked about.
type complexSelector struct {
	compounds []compoundSelector
	combs     []combinator
}

// compoundSelector is the simple selectors written with nothing between them,
// every one of which must hold for the compound to: `rect.big#hero[fill]`.
type compoundSelector struct {
	simples []simpleSelector
}

// combinator joins two compounds: a descendant (a space), `>`, `+` or `~`.
type combinator uint8

const (
	descendant combinator = iota
	child
	adjacent
	sibling
)

// simpleSelector is one condition on an element: its type, a class, an id, the
// universal `*`, an attribute condition, or a pseudo-class about where it stands.
type simpleSelector struct {
	kind simpleKind
	name string
	op   attrOp
	val  string
}

type simpleKind uint8

const (
	kindType simpleKind = iota
	kindClass
	kindID
	kindUniversal
	kindAttribute
	kindRoot
	kindFirstChild
	kindLastChild
	kindOnlyChild
)

// attrOp is how an attribute selector compares the value it names.
type attrOp uint8

const (
	attrExists attrOp = iota
	attrEquals
	attrIncludes
	attrDash
	attrPrefix
	attrSuffix
	attrSubstr
)
