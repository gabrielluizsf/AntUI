package svg

import "testing"

// A sheet may name a value once as a custom property and read it back with
// var(), including through the elements a custom property inherits down and
// through another custom property it was written as. These check that the value
// a var() reads is the one the cascade settled on the element, that a fallback
// stands where nothing defined the name, and that a name nothing answers leaves
// the declaration out rather than painting with it.

func TestSheetResolvesAVariable(t *testing.T) {
	got := cssFill(t, `<style>:root { --c: red } rect { fill: var(--c) }</style>`+
		`<rect width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the red the custom property named", got)
	}
}

func TestAVariableInheritsDownTheDrawing(t *testing.T) {
	got := cssFill(t, `<style>g { --c: red }</style><g fill="var(--c)"><rect width="10" height="10"/></g>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the group's custom property inherited by the rect", got)
	}
}

func TestAVariableTakesItsFallback(t *testing.T) {
	got := cssFill(t, `<style>rect { fill: var(--missing, #00f) }</style>`+
		`<rect width="10" height="10"/>`)
	if got != cssBlue {
		t.Errorf("fill = %v, want the fallback the var() wrote", got)
	}
}

func TestAVariableWrittenAsAnotherVariable(t *testing.T) {
	got := cssFill(t, `<style>:root { --a: var(--b); --b: red } rect { fill: var(--a) }</style>`+
		`<rect width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the chain of custom properties followed to the red", got)
	}
}

func TestTheStyleAttributeReadsTheSheetsVariable(t *testing.T) {
	got := cssFill(t, `<style>rect { --c: red }</style>`+
		`<rect style="fill: var(--c)" width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the sheet's custom property read by the style attribute", got)
	}
}

func TestTheStyleAttributeDeclaresAVariable(t *testing.T) {
	got := cssFill(t, `<style>rect { --c: blue } rect { fill: var(--c) }</style>`+
		`<rect style="--c: red" width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the style attribute's custom property over the sheet's", got)
	}
}

func TestAnImportantVariableBeatsANormalOne(t *testing.T) {
	got := cssFill(t, `<style>rect { --c: red } rect { --c: blue !important } rect { fill: var(--c) }</style>`+
		`<rect width="10" height="10"/>`)
	if got != cssBlue {
		t.Errorf("fill = %v, want the important custom property to win", got)
	}
}

func TestAnUnresolvedVariableLeavesTheDeclarationOut(t *testing.T) {
	got := cssFill(t, `<style>rect { fill: var(--nope) }</style><rect width="10" height="10"/>`)
	if got != cssBlack {
		t.Errorf("fill = %v, want the rect left black where var() had nothing to read", got)
	}
}

func TestAVariableCycleIsStoppedRatherThanSpun(t *testing.T) {
	got := cssFill(t, `<style>:root { --a: var(--b); --b: var(--a) } rect { fill: var(--a) }</style>`+
		`<rect width="10" height="10"/>`)
	if got != cssBlack {
		t.Errorf("fill = %v, want a cycle to resolve to nothing rather than hang", got)
	}
}
