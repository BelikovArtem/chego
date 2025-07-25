package fen

import (
	"testing"

	"github.com/BelikovArtem/chego/types"
)

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

func TestParse(t *testing.T) {
	testcases := []struct {
		fenStr                  string
		expectedActiveColor     types.Color
		expectedCastlingRights  types.CastlingRights
		expectedEnPassantTarget int
		expectedHalfmoveCnt     int
		expectedFullmoveCnt     int
	}{
		{
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			types.ColorWhite, 0xF, 0, 0, 1,
		},
		{
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			types.ColorBlack, 0xF, types.SE3, 0, 1,
		},
	}

	for _, tc := range testcases {
		p := Parse(tc.fenStr)

		if p.ActiveColor != tc.expectedActiveColor {
			t.Fatalf("test \"%s\" failed: expected color %b, got %b", tc.fenStr,
				tc.expectedActiveColor, p.ActiveColor)
		}
		if p.CastlingRights != tc.expectedCastlingRights {
			t.Fatalf("test \"%s\" failed: expected castling rights %b, got %b", tc.fenStr,
				tc.expectedCastlingRights, p.CastlingRights)
		}
		if p.EPTarget != tc.expectedEnPassantTarget {
			t.Fatalf("test \"%s\" failed: expected en passant %d, got %d", tc.fenStr,
				tc.expectedEnPassantTarget, p.EPTarget)
		}
		if p.HalfmoveCnt != tc.expectedHalfmoveCnt {
			t.Fatalf("test \"%s\" failed: expected halfmove %d, got %d", tc.fenStr,
				tc.expectedHalfmoveCnt, p.HalfmoveCnt)
		}
		if p.FullmoveCnt != tc.expectedFullmoveCnt {
			t.Fatalf("test \"%s\" failed: expected fullmove %d, got %d", tc.fenStr,
				tc.expectedFullmoveCnt, p.FullmoveCnt)
		}
	}
}

func TestSerialize(t *testing.T) {
	testcases := []struct {
		position types.Position
		expected string
	}{
		{types.Position{
			Bitboards:   ToBitboardArray("1r3r2/4bpkp/1qb1p1p1/3pP1P1/p1pP1Q2/PpP2N1R/1Pn1B2P/3RB2K"),
			ActiveColor: types.ColorWhite, CastlingRights: 0x0, EPTarget: 0x0,
			HalfmoveCnt: 0, FullmoveCnt: 1,
		}, "1r3r2/4bpkp/1qb1p1p1/3pP1P1/p1pP1Q2/PpP2N1R/1Pn1B2P/3RB2K w - - 0 1"},
		{types.Position{
			Bitboards:   ToBitboardArray("rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR"),
			ActiveColor: types.ColorBlack, CastlingRights: 0xF, EPTarget: types.SF3,
			HalfmoveCnt: 0, FullmoveCnt: 1,
		}, "rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR b KQkq f3 0 1"},
		{types.Position{
			Bitboards:   ToBitboardArray("4k3/8/8/8/8/3P4/2K5/8"),
			ActiveColor: types.ColorWhite, CastlingRights: 0x0,
			EPTarget: 0x0, HalfmoveCnt: 0, FullmoveCnt: 64,
		}, "4k3/8/8/8/8/3P4/2K5/8 w - - 0 64"},
	}

	for _, tc := range testcases {
		got := Serialize(tc.position)

		if got != tc.expected {
			t.Fatalf("expected \"%s\", got \"%s\"", tc.expected, got)
		}
	}
}

func BenchmarkToBitboardArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToBitboardArray("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	}
}

func BenchmarkFromBitboardArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromBitboardArray([12]uint64{
			0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
			0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
			0x8100000000000000, 0x800000000000000, 0x1000000000000000,
		})
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	}
}

func BenchmarkSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Serialize(types.Position{
			Bitboards: [12]uint64{
				0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
				0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
				0x8100000000000000, 0x800000000000000, 0x1000000000000000,
			},
			ActiveColor:    types.ColorWhite,
			CastlingRights: 0xF,
			EPTarget:       0,
			HalfmoveCnt:    0,
			FullmoveCnt:    1,
		})
	}
}
