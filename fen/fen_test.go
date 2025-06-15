package fen

import (
	"chego/enum"
	"testing"
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
		expectedActiveColor     enum.Color
		expectedCastlingRights  enum.CastlingFlag
		expectedEnPassantTarget int
		expectedHalfmoveCnt     int
		expectedFullmoveCnt     int
	}{
		{
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			enum.ColorWhite, 0xF, 0, 0, 1,
		},
		{
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			enum.ColorBlack, 0xF, enum.SE3, 0, 1,
		},
	}

	for _, tc := range testcases {
		_, a, c, e, h, f := Parse(tc.fenStr)

		if a != tc.expectedActiveColor {
			t.Fatalf("test \"%s\" failed: expected color %b, got %b", tc.fenStr, tc.expectedActiveColor, a)
		}
		if c != tc.expectedCastlingRights {
			t.Fatalf("test \"%s\" failed: expected castling rights %b, got %b", tc.fenStr,
				tc.expectedCastlingRights, c)
		}
		if e != tc.expectedEnPassantTarget {
			t.Fatalf("test \"%s\" failed: expected en passant %d, got %d", tc.fenStr,
				tc.expectedEnPassantTarget, e)
		}
		if h != tc.expectedHalfmoveCnt {
			t.Fatalf("test \"%s\" failed: expected halfmove %d, got %d", tc.fenStr,
				tc.expectedHalfmoveCnt, h)
		}
		if f != tc.expectedFullmoveCnt {
			t.Fatalf("test \"%s\" failed: expected fullmove %d, got %d", tc.fenStr,
				tc.expectedFullmoveCnt, f)
		}
	}
}

func TestSerialize(t *testing.T) {
	testcases := []struct {
		bitboards       [12]uint64
		activeColor     enum.Color
		castlingRights  enum.CastlingFlag
		enPassantTarget int
		halfmoveCnt     int
		fullmoveCnt     int
		expected        string
	}{
		{
			ToBitboardArray("1r3r2/4bpkp/1qb1p1p1/3pP1P1/p1pP1Q2/PpP2N1R/1Pn1B2P/3RB2K"),
			enum.ColorWhite, 0x0, 0x0, 0, 1,
			"1r3r2/4bpkp/1qb1p1p1/3pP1P1/p1pP1Q2/PpP2N1R/1Pn1B2P/3RB2K w - - 0 1",
		},
		{
			ToBitboardArray("rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR"),
			enum.ColorBlack, 0xF, enum.SF3, 0, 1,
			"rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR b KQkq f3 0 1",
		},
		{
			ToBitboardArray("4k3/8/8/8/8/3P4/2K5/8"),
			enum.ColorWhite, 0x0, 0x0, 0, 64,
			"4k3/8/8/8/8/3P4/2K5/8 w - - 0 64",
		},
	}

	for _, tc := range testcases {
		got := Serialize(tc.bitboards, tc.activeColor, tc.castlingRights, tc.enPassantTarget,
			tc.halfmoveCnt, tc.fullmoveCnt)

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
		Serialize([12]uint64{
			0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
			0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
			0x8100000000000000, 0x800000000000000, 0x1000000000000000,
		}, enum.ColorWhite, 0xF, 0, 0, 1)
	}
}
