package chego

import (
	"testing"
)

func TestParseBitboards(t *testing.T) {
	testcases := []struct {
		name     string
		fen      string
		expected [15]uint64
	}{
		{
			"Initial position",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			[15]uint64{
				0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
				0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
				0x8100000000000000, 0x800000000000000, 0x1000000000000000,
				0xFFFF, 0xFFFF000000000000, 0xFFFF00000000FFFF,
			},
		},
		{
			"Two rooks, two pawns",
			"8/4p3/1PR5/8/4R3/8/4p3/8",
			[15]uint64{
				0x20000000000, 0x0, 0x0, 0x40010000000, 0x0, 0x0,
				0x10000000001000, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x60010000000, 0x10000000001000, 0x10060010001000,
			},
		},
	}

	for _, tc := range testcases {
		for piece, bitboard := range ParseBitboards(tc.fen) {
			if tc.expected[piece] != bitboard {
				t.Fatalf("expected %v\ngot %v", tc.expected, bitboard)
			}
		}
	}
}

func TestSerializeBitboards(t *testing.T) {
	testcases := []struct {
		name      string
		bitboards [15]uint64
		expected  string
	}{
		{
			"Initial position",
			[15]uint64{
				0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
				0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
				0x8100000000000000, 0x800000000000000, 0x1000000000000000,
			}, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
		},
		{
			"Two rooks, two pawns",
			[15]uint64{
				0x20000000000, 0x0, 0x0, 0x40010000000, 0x0, 0x0,
				0x10000000001000, 0x0, 0x0, 0x0, 0x0, 0x0,
			}, "8/4p3/1PR5/8/4R3/8/4p3/8",
		},
	}

	for _, tc := range testcases {
		got := SerializeBitboards(tc.bitboards)
		if tc.expected != got {
			t.Fatalf("expected %s\ngot %s", tc.expected, got)
		}
	}
}

// TestParseFEN does not check the parsed bitboards, since that is the job
// of TestParseBitboards
func TestParseFEN(t *testing.T) {
	testcases := []struct {
		fen      string
		expected Position
	}{
		{
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			Position{
				ActiveColor:    ColorWhite,
				CastlingRights: 0xF,
				EPTarget:       SA1,
				HalfmoveCnt:    0,
				FullmoveCnt:    1,
			},
		},
		{
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			Position{
				ActiveColor:    ColorBlack,
				CastlingRights: 0xF,
				EPTarget:       SE3,
				HalfmoveCnt:    0,
				FullmoveCnt:    1,
			},
		},
	}

	for _, tc := range testcases {
		p := ParseFEN(tc.fen)
		tc.expected.Bitboards = p.Bitboards

		if p != tc.expected {
			t.Fatalf("expected %v\ngot %v", tc.expected, p)
		}
	}
}

// TestSerializeFEN does not check the serialized bitboards, since that is the job
// of TestSerializeBitboards.
func TestSerializeFEN(t *testing.T) {
	testcases := []struct {
		position Position
		expected string
	}{
		{Position{
			Bitboards:   ParseBitboards("1r3r2/4bpkp/1qb1p1p1/3pP1P1/p1pP1Q2/PpP2N1R/1Pn1B2P/3RB2K"),
			ActiveColor: ColorWhite, CastlingRights: 0x0, EPTarget: 0x0,
			HalfmoveCnt: 0, FullmoveCnt: 1,
		}, "1r3r2/4bpkp/1qb1p1p1/3pP1P1/p1pP1Q2/PpP2N1R/1Pn1B2P/3RB2K w - - 0 1"},
		{Position{
			Bitboards:   ParseBitboards("rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR"),
			ActiveColor: ColorBlack, CastlingRights: 0xF, EPTarget: SF3,
			HalfmoveCnt: 0, FullmoveCnt: 1,
		}, "rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR b KQkq f3 0 1"},
		{Position{
			Bitboards:   ParseBitboards("4k3/8/8/8/8/3P4/2K5/8"),
			ActiveColor: ColorWhite, CastlingRights: 0x0,
			EPTarget: 0x0, HalfmoveCnt: 0, FullmoveCnt: 64,
		}, "4k3/8/8/8/8/3P4/2K5/8 w - - 0 64"},
	}

	for _, tc := range testcases {
		got := SerializeFEN(tc.position)

		if got != tc.expected {
			t.Fatalf("expected \"%s\", got \"%s\"", tc.expected, got)
		}
	}
}

func BenchmarkParseBitboards(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseBitboards("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	}
}

func BenchmarkSerializeBitboards(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SerializeBitboards([15]uint64{
			0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
			0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
			0x8100000000000000, 0x800000000000000, 0x1000000000000000,
		})
	}
}

func BenchmarkParseFEN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	}
}

func BenchmarkSerializeFEN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SerializeFEN(Position{
			Bitboards: [15]uint64{
				0xFF00, 0x42, 0x24, 0x81, 0x8, 0x10,
				0xFF000000000000, 0x4200000000000000, 0x2400000000000000,
				0x8100000000000000, 0x800000000000000, 0x1000000000000000,
			},
			ActiveColor:    ColorWhite,
			CastlingRights: 0xF,
			FullmoveCnt:    1,
		})
	}
}
