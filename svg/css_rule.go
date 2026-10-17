package svg

// cssRule is one rule of a `<style>` sheet: the selectors that turn it on and
// the declarations it makes when one of them matches an element.
type cssRule struct {
	selectors []complexSelector
	decls     []cssDecl
}

// cssDecl is one `name: value` pair of a rule or of a `style="…"` attribute,
// with the `!important` it was written with: an important declaration beats
// one written without it, whatever selector either came from.
type cssDecl struct {
	name      string
	value     string
	important bool
}
