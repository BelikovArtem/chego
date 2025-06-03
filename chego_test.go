package chego

import (
	"chego/cli"
	"chego/enum"
	"testing"
)

func TestGenPawnAttacks(t *testing.T) {
	testcases := []struct {
		name     string
		color    enum.Color
		bitboard uint64
		expected uint64
	}{
		{"White pawn B4", enum.ColorWhite, enum.B4, enum.A5 | enum.C5},
		{"White pawn A4", enum.ColorWhite, enum.A4, enum.B5},
		{"White pawn H4", enum.ColorWhite, enum.H4, enum.G5},
		{"White pawn B8", enum.ColorWhite, enum.B8, 0x0},
		{"Black pawn B4", enum.ColorBlack, enum.B4, enum.A3 | enum.C3},
		{"Black pawn A4", enum.ColorBlack, enum.A4, enum.B3},
		{"Black pawn H4", enum.ColorBlack, enum.H4, enum.G3},
		{"Black pawn B1", enum.ColorBlack, enum.B1, 0x0},
	}

	for _, tc := range testcases {
		got := GenPawnAttacks(tc.bitboard, tc.color)
		if got != tc.expected {
			t.Logf("test %s failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", cli.FormatBitboard(tc.expected, enum.PieceWPawn))
			t.Logf("got bitboard:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWPawn))
			t.FailNow()
		}
	}
}

func TestGenKnightAttacks(t *testing.T) {
	testcases := []struct {
		name     string
		bitboard uint64
		expected uint64
	}{
		{"Knight D4", enum.D4, enum.C2 | enum.E2 | enum.B3 | enum.F3 | enum.B5 |
			enum.F5 | enum.C6 | enum.E6},
		{"Knight A8", enum.A8, enum.B6 | enum.C7},
		{"Knight H1", enum.H1, enum.F2 | enum.G3},
	}

	for _, tc := range testcases {
		got := GenKnightAttacks(tc.bitboard)
		if got != tc.expected {
			t.Logf("test %s failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", cli.FormatBitboard(tc.expected, enum.PieceWKnight))
			t.Logf("got bitboard:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWKnight))
			t.FailNow()
		}
	}
}

func TestGenKingAttacks(t *testing.T) {
	testcases := []struct {
		name     string
		bitboard uint64
		expected uint64
	}{
		{"King D5", enum.D5, enum.C4 | enum.D4 | enum.E4 | enum.C5 | enum.E5 |
			enum.C6 | enum.D6 | enum.E6},
		{"King A8", enum.A8, enum.A7 | enum.B7 | enum.B8},
	}

	for _, tc := range testcases {
		got := GenKingAttacks(tc.bitboard)
		if got != tc.expected {
			t.Logf("test %s failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", cli.FormatBitboard(tc.expected, enum.PieceWKing))
			t.Logf("got bitboard:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWKing))
			t.FailNow()
		}
	}
}

func TestGenBishopAttacks(t *testing.T) {
	testcases := []struct {
		name      string
		bitboard  uint64
		occupancy uint64
		expected  uint64
	}{
		{"Bishop D5 - Blocked B3", enum.D5, enum.B3, enum.C4 | enum.B3 | enum.E4 | enum.F3 |
			enum.G2 | enum.H1 | enum.C6 | enum.B7 | enum.A8 | enum.E6 | enum.F7 | enum.G8},
	}

	for _, tc := range testcases {
		got := GenBishopAttacks(tc.bitboard, tc.occupancy)
		if got != tc.expected {
			t.Logf("test %s failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", cli.FormatBitboard(tc.expected, enum.PieceWBishop))
			t.Logf("got bitboard:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWBishop))
			t.FailNow()
		}
	}
}

func TestGenRookAttacks(t *testing.T) {
	testcases := []struct {
		name      string
		bitboard  uint64
		occupancy uint64
		expected  uint64
	}{
		{"Rook A1 - No blockers", enum.A1, 0x0, enum.B1 | enum.C1 | enum.D1 | enum.E1 |
			enum.F1 | enum.G1 | enum.H1 | enum.A2 | enum.A3 | enum.A4 | enum.A5 | enum.A6 |
			enum.A7 | enum.A8},
		{"Rook D5 - Bloocked D2, B5, D7", enum.D5, enum.D2 | enum.B5 | enum.D7,
			enum.D4 | enum.D3 | enum.D2 | enum.C5 | enum.B5 | enum.E5 | enum.F5 |
				enum.G5 | enum.H5 | enum.D6 | enum.D7},
	}

	for _, tc := range testcases {
		got := GenRookAttacks(tc.bitboard, tc.occupancy)
		if got != tc.expected {
			t.Logf("test %s failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", cli.FormatBitboard(tc.expected, enum.PieceWRook))
			t.Logf("got bitboard:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWRook))
			t.FailNow()
		}
	}
}

func TestBitScan(t *testing.T) {
	testcases := []struct {
		name     string
		bitboard uint64
		expected int
	}{
		{"F0000", 0xF0000, 16},
	}

	for _, tc := range testcases {
		got := BitScan(tc.bitboard)
		if got != tc.expected {
			t.Fatalf("Testcase %s failed: expected %d, got %d", tc.name, tc.expected, got)
		}
	}
}

func TestGenMagicNumber(t *testing.T) {
	InitBishopRelevantOccupancy()
	InitRookRelevantOccupancy()

	t.Logf("\n\n")
	for square := 0; square < 64; square++ {
		t.Logf("%x,\n", GenMagicNumber(square, true))
	}

	t.Logf("\n\n")
	for square := 0; square < 64; square++ {
		t.Logf("%x,\n", GenMagicNumber(square, false))
	}
}

func TestLookupBishopAttacks(t *testing.T) {
	InitAttackTables()

	var occupied uint64 = enum.F2 | enum.B3 | enum.F4 | enum.D5 | enum.G7
	for square := uint64(1); square != 0; square <<= 1 {
		got := LookupBishopAttacks(BitScan(square), occupied)
		expected := GenBishopAttacks(square, occupied)

		if got != expected {
			t.Logf("expected:\n\n%s\n\n", cli.FormatBitboard(expected, enum.PieceWBishop))
			t.Logf("got:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWBishop))
			t.FailNow()
		}
	}
}

func TestLookupRookAttacks(t *testing.T) {
	InitAttackTables()

	var occupied uint64 = enum.F2 | enum.B3 | enum.F4 | enum.D5 | enum.G7
	for square := uint64(1); square != 0; square <<= 1 {
		got := LookupRookAttacks(BitScan(square), occupied)
		expected := GenRookAttacks(square, occupied)

		if got != expected {
			t.Logf("got:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWRook))
			t.Logf("expected:\n\n%s\n\n", cli.FormatBitboard(expected, enum.PieceWRook))
			t.FailNow()
		}
	}
}

func BenchmarkGenPawnAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenPawnAttacks(enum.B4, enum.ColorWhite)
	}
}

func BenchmarkGenKnightAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenKnightAttacks(enum.B4)
	}
}

func BenchmarkGenKingAttakcs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenKingAttacks(enum.B4)
	}
}

func BenchmarkGenBishopAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenBishopAttacks(enum.D5, enum.B3)
	}
}

func BenchmarkGenRookAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenRookAttacks(enum.D5, enum.B3)
	}
}

func BenchmarkInitBishopReleventOccupancy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitBishopRelevantOccupancy()
	}
}

func BenchmarkInitRookReleventOccupancy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitRookRelevantOccupancy()
	}
}

func BenchmarkGenMagicNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenMagicNumber(23, false)
	}
}

func BenchmarkInitAttackTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitAttackTables()
	}
}
