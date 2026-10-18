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
// universal `*`, an attribute condition, or a pseudo-class about where it
// stands. A pseudo-class with an argument keeps it in the field its kind
// names.
type simpleSelector struct {
	kind     simpleKind
	name     string
	op       attrOp
	val      string
	caseless bool
	// sub is the selector list of `:not(…)`, `:is(…)` and `:where(…)`, whose
	// arguments stand on the element the compound is matched against.
	sub []complexSelector
	// has is the relative-selector list of `:has(…)`, each read from the
	// element outward.
	has []relativeSelector
	// nth is the `An+B` of an `:nth-…`, with the `of S` list some carry.
	nth nthExpr
	// lang is the list of ranges `:lang(…)` names, each of which may match the
	// language an element stands in.
	lang []string
	// uri is the namespace an attribute selector asks for: empty for one of no
	// namespace, which is what an attribute selector written without a prefix
	// means, and the namespace a prefix is bound to where one is worn. Any says
	// the selector was written `*|`, which names every namespace at once.
	uri string
	any bool
	// pe marks a pseudo-element, `::before` and the like. A still drawing has
	// no generated box for one, so it never matches; it is read all the same,
	// because the selector around it is valid and the rule is a rule.
	pe bool
}

type simpleKind uint8

const (
	kindType simpleKind = iota
	kindClass
	kindID
	kindUniversal
	kindAttribute
	kindRoot
	kindScope
	kindEmpty
	kindFirstChild
	kindLastChild
	kindOnlyChild
	kindFirstOfType
	kindLastOfType
	kindOnlyOfType
	kindNthChild
	kindNthLastChild
	kindNthOfType
	kindNthLastOfType
	kindNot
	kindIs
	kindWhere
	kindHas
	// kindLang and kindDir ask about the language and the direction an element
	// is written in, which a drawing carries on `xml:lang`/`lang` and `dir`
	// and which travel down from the element that gave one.
	kindLang
	kindDir
	// kindUnmatchable is a valid selector a still drawing cannot answer: an
	// interaction pseudo-class, a functional one about a state that is not
	// there, or a pseudo-element. It is kept so the rule around it is read,
	// and it never matches.
	kindUnmatchable
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
