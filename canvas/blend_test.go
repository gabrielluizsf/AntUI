package canvas

import (
	"math/rand/v2"
	"testing"
)

// The assembly exists only if it agrees with the Go loop on every pixel. A
// blend that is nearly right is a renderer whose output depends on which
// machine drew it, which is the one thing a software renderer is for.

func TestBlendRowMatchesTheReference(t *testing.T) {
	if !hasSIMD {
		t.Skip("no assembly on this architecture")
	}

	// Every alpha, against a spread of destinations, at every length that
	// exercises the tail.
	for alpha := range 256 {
		for _, n := range []int{1, 2, 3, 4, 5, 7, 8, 9, 15, 16, 17, 33, 64} {
			src := make([]Color, n)
			for i := range src {
				src[i] = RGBA(uint8(i*37), uint8(i*11+5), uint8(i*97), uint8(alpha))
			}
			want := make([]Color, n)
			got := make([]Color, n)
			for i := range want {
				c := RGBA(uint8(i*53+7), uint8(i*29), uint8(i*3+1), uint8(i*17))
				want[i], got[i] = c, c
			}

			blendRowGo(want, src)
			BlendRow(got, src)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("alpha %d, n %d, pixel %d: got %08X, want %08X",
						alpha, n, i, uint32(got[i]), uint32(want[i]))
				}
			}
		}
	}
}

func TestBlendRowSolidMatchesTheReference(t *testing.T) {
	if !hasSIMD {
		t.Skip("no assembly on this architecture")
	}
	for alpha := range 256 {
		for _, n := range []int{1, 3, 4, 5, 8, 17, 64} {
			c := RGBA(200, 30, 90, uint8(alpha))
			want := make([]Color, n)
			got := make([]Color, n)
			for i := range want {
				p := RGBA(uint8(i*13), uint8(255-i*5), uint8(i*61), uint8(i*23))
				want[i], got[i] = p, p
			}

			// The reference for a solid run is the same as blending a run of
			// one repeated colour, which is what the caller means by it.
			srcRun := make([]Color, n)
			for i := range srcRun {
				srcRun[i] = c
			}
			blendRowGo(want, srcRun)
			BlendRowSolid(got, c)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("alpha %d, n %d, pixel %d: got %08X, want %08X",
						alpha, n, i, uint32(got[i]), uint32(want[i]))
				}
			}
		}
	}
}

// Random pixels catch anything the patterned ones above happen to line up on.
func TestBlendRowFuzz(t *testing.T) {
	if !hasSIMD {
		t.Skip("no assembly on this architecture")
	}
	rng := rand.New(rand.NewPCG(1, 2))
	for range 2000 {
		n := 1 + rng.IntN(40)
		src := make([]Color, n)
		want := make([]Color, n)
		got := make([]Color, n)
		for i := range src {
			src[i] = Color(rng.Uint32())
			p := Color(rng.Uint32())
			want[i], got[i] = p, p
		}
		blendRowGo(want, src)
		BlendRow(got, src)
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("n %d, pixel %d: src %08X, got %08X, want %08X",
					n, i, uint32(src[i]), uint32(got[i]), uint32(want[i]))
			}
		}
	}
}

func TestBlendRowHandlesEmptyAndMismatchedRuns(t *testing.T) {
	BlendRow(nil, nil)
	BlendRow([]Color{1, 2, 3}, nil)
	BlendRow(nil, []Color{1, 2, 3})
	BlendRowSolid(nil, Red)

	// A shorter destination bounds the work, rather than reading past it.
	dst := []Color{Black, Black}
	BlendRow(dst, []Color{White, White, White, White, White})
	if dst[0] != White || dst[1] != White {
		t.Errorf("dst = %v, want both blended", dst)
	}
}

func BenchmarkBlendRow(b *testing.B) {
	src := make([]Color, 800)
	dst := make([]Color, 800)
	for i := range src {
		src[i] = RGBA(uint8(i), uint8(i*3), uint8(i*7), 160)
	}
	b.Run("simd", func(b *testing.B) {
		for b.Loop() {
			BlendRow(dst, src)
		}
	})
	b.Run("go", func(b *testing.B) {
		for b.Loop() {
			blendRowGo(dst, src)
		}
	})
}

func BenchmarkBlendRowSolid(b *testing.B) {
	dst := make([]Color, 800)
	c := RGBA(200, 30, 90, 160)
	b.Run("simd", func(b *testing.B) {
		for b.Loop() {
			BlendRowSolid(dst, c)
		}
	})
	b.Run("go", func(b *testing.B) {
		for b.Loop() {
			blendRowSolidGo(dst, c)
		}
	})
}
