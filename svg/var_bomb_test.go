package svg

import (
	"fmt"
	"strings"
	"testing"
)

// TestResolveVarsRefusesAnExponentialValue is the billion laughs written in
// custom properties: each names the one before ten times, seven deep, which is
// ten million copies of a colour out of a few hundred bytes. The CSS Variables
// spec has a value past the limit be invalid, which is a value left with its
// var() in it — the one [Style.withVars] drops.
func TestResolveVarsRefusesAnExponentialValue(t *testing.T) {
	customs := map[string]string{"--v0": "red"}
	for i := 1; i <= 7; i++ {
		customs[fmt.Sprintf("--v%d", i)] = strings.Repeat(fmt.Sprintf("var(--v%d)", i-1), 10)
	}
	got := resolveVars("var(--v7)", customs)
	if !strings.Contains(got, "var(") || len(got) > maxVarBytes {
		t.Fatalf("the value resolved to %d bytes, want it left unresolved", len(got))
	}
}

// TestResolveVarsStillReadsNestedValues keeps the limit from costing what it
// protects: a value naming one that names another reads through both.
func TestResolveVarsStillReadsNestedValues(t *testing.T) {
	customs := map[string]string{"--base": "#f00", "--accent": "var(--base)"}
	if got := resolveVars("var(--accent)", customs); got != "#f00" {
		t.Fatalf("the value resolved to %q, want #f00", got)
	}
}
