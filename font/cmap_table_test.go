package font

import (
	"encoding/binary"
	"testing"

	"github.com/gabrielluizsf/antui/assert"
)

func TestCmapTable_Glyph(t *testing.T) {
	// Format 0 Data Setup: minimum length 262 (6 + 256)
	f0Data := make([]byte, 262)
	f0Data[6+65] = 99 // Rune 65 ('A') maps to glyph 99

	// Format 4 Data Setup: Segments = 1
	f4Data := make([]byte, 24)
	binary.BigEndian.PutUint16(f4Data[14:16], 60)  // end
	binary.BigEndian.PutUint16(f4Data[18:20], 40)  // start
	binary.BigEndian.PutUint16(f4Data[20:22], 100) // delta
	binary.BigEndian.PutUint16(f4Data[22:24], 0)   // offset range

	// Format 6 Data Setup
	f6Data := make([]byte, 14)
	binary.BigEndian.PutUint16(f6Data[6:8], 10)   // first
	binary.BigEndian.PutUint16(f6Data[8:10], 2)   // count
	binary.BigEndian.PutUint16(f6Data[10:12], 42) // glyph for rune 10
	binary.BigEndian.PutUint16(f6Data[12:14], 43) // glyph for rune 11

	// Format 12 Data Setup: Groups = 1
	f12Data := make([]byte, 28)
	binary.BigEndian.PutUint32(f12Data[16:20], 1000) // start
	binary.BigEndian.PutUint32(f12Data[20:24], 1005) // end
	binary.BigEndian.PutUint32(f12Data[24:28], 50)   // offset

	tests := []struct {
		name     string
		table    CmapTable
		rune     rune
		expected int
	}{
		{
			name: "Format 0 - Valid mapping",
			table: CmapTable{
				Format: 0,
				Data:   f0Data,
			},
			rune:     65,
			expected: 99,
		},
		{
			name: "Format 0 - Out of bounds rune",
			table: CmapTable{
				Format: 0,
				Data:   f0Data,
			},
			rune:     300,
			expected: 0,
		},
		{
			name: "Format 4 - Valid mapping",
			table: CmapTable{
				Format:   4,
				Segments: 1,
				Data:     f4Data,
			},
			rune:     50,
			expected: 150, // 50 + 100 (delta) = 150
		},
		{
			name: "Format 4 - Rune out of segment bounds",
			table: CmapTable{
				Format:   4,
				Segments: 1,
				Data:     f4Data,
			},
			rune:     70,
			expected: 0,
		},
		{
			name: "Format 6 - Valid mapping",
			table: CmapTable{
				Format: 6,
				Data:   f6Data,
			},
			rune:     11,
			expected: 43,
		},
		{
			name: "Format 6 - Rune out of bounds",
			table: CmapTable{
				Format: 6,
				Data:   f6Data,
			},
			rune:     20,
			expected: 0,
		},
		{
			name: "Format 12 - Valid mapping",
			table: CmapTable{
				Format: 12,
				Groups: 1,
				Data:   f12Data,
			},
			rune:     1002,
			expected: 52, // 50 + 1002 - 1000 = 52
		},
		{
			name: "Format 12 - Rune out of bounds",
			table: CmapTable{
				Format: 12,
				Groups: 1,
				Data:   f12Data,
			},
			rune:     2000,
			expected: 0,
		},
		{
			name: "Unknown format defaults to zero",
			table: CmapTable{
				Format: 99,
				Data:   []byte{},
			},
			rune:     10,
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.table.Glyph(tc.rune)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestBeHelpers(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}

	t.Run("be16 valid parsing", func(t *testing.T) {
		got := be16(data, 0)
		assert.Equal(t, uint16(0x0102), got)
	})

	t.Run("be16 out of bounds returns zero", func(t *testing.T) {
		got := be16(data, 3)
		assert.Equal(t, uint16(0), got)
	})

	t.Run("be32 valid parsing", func(t *testing.T) {
		got := be32(data, 0)
		assert.Equal(t, uint32(0x01020304), got)
	})

	t.Run("be32 out of bounds returns zero", func(t *testing.T) {
		got := be32(data, 1)
		assert.Equal(t, uint32(0), got)
	})
}
