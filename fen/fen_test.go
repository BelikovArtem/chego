package fen

import "testing"

func TestToBitboardArray(t *testing.T) {
	testcases := []struct {
		name     string
		fenStr   string
		expected [12]uint64
	}{
		{
			"Initial position",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			[12]uint64{
				0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
				0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
				0x8100000000000000, 0x800000000000000, 0x1000000000000000,
			},
		},
		{
			"Two rooks, two pawns",
			"8/4p3/1PR5/8/4R3/8/4p3/8",
			[12]uint64{
				0x20000000000, 0x0, 0x0, 0x40010000000, 0x0, 0x0,
				0x10000000001000, 0x0, 0x0, 0x0, 0x0, 0x0,
			},
		},
	}

	for _, tc := range testcases {
		for pieceType, bitboard := range ToBitboardArray(tc.fenStr) {
			if tc.expected[pieceType] != bitboard {
				t.Fatalf("%s\nexpected:%x\ngot:%x", tc.name, tc.expected[pieceType], bitboard)
			}
		}
	}
}

func TestFromBitboardArray(t *testing.T) {
	testcases := []struct {
		name      string
		bitboards [12]uint64
		expected  string
	}{
		{
			"Initial position",
			[12]uint64{
				0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
				0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
				0x8100000000000000, 0x800000000000000, 0x1000000000000000,
			},
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
		},
		{
			"Two rooks, two pawns",
			[12]uint64{
				0x20000000000, 0x0, 0x0, 0x40010000000, 0x0, 0x0,
				0x10000000001000, 0x0, 0x0, 0x0, 0x0, 0x0,
			},
			"8/4p3/1PR5/8/4R3/8/4p3/8",
		},
	}

	for _, tc := range testcases {
		got := FromBitboardArray(tc.bitboards)
		if tc.expected != got {
			t.Fatalf("expected: %s, got: %s", tc.expected, got)
		}
	}
}

// Best result: ~97 ns/op, 0 B/op, 0 allocs/op.
func BenchmarkToBitboardArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToBitboardArray("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	}
}

// Best result: ~280 ns/op, 120 B/op, 4 allocs/op.
func BenchmarkFromBitboardArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromBitboardArray([12]uint64{
			0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
			0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
			0x8100000000000000, 0x800000000000000, 0x1000000000000000,
		})
	}
}
