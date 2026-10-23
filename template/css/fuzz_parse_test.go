package css

import "testing"

// FuzzParse feeds the style sheet reader text nobody vetted: whatever a sheet
// says must end in an error or a sheet, never a panic.
func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		`a{color:red;margin:1px 2px}`,
		`@media (min-width:100px){.x{--v:calc(1px + 2%);width:var(--v)}}`,
		`.g{display:grid;grid-template-columns:repeat(3,1fr);background:linear-gradient(red,blue)}`,
		`@keyframes k{from{opacity:0}to{opacity:1}} .a{animation:k 1s;transform:rotate(45deg)}`,
		`.c{clip-path:polygon(0 0,100% 0,50% 100%);font:12px/1.5 sans-serif}`,
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, src string) {
		if len(src) > 4096 {
			return
		}
		Parse(src)
	})
}
