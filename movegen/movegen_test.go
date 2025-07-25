package movegen

import (
	"os"
	"testing"

	"github.com/BelikovArtem/chego/bitutil"
	"github.com/BelikovArtem/chego/fen"
	"github.com/BelikovArtem/chego/format"
	"github.com/BelikovArtem/chego/types"
)

// Used to avoid writing InitAttackTables() each time.
func TestMain(m *testing.M) {
	// Setup.
	InitAttackTables()
	// Tests and benchmarks execution.
	os.Exit(m.Run())
}

func TestGenPawnAttacks(t *testing.T) {
	testcases := []struct {
		name     string
		color    types.Color
		bitboard uint64
		expected uint64
	}{
		{"White pawn B4", types.ColorWhite, types.B4, types.A5 | types.C5},
		{"White pawn A4", types.ColorWhite, types.A4, types.B5},
		{"White pawn H4", types.ColorWhite, types.H4, types.G5},
		{"White pawn B8", types.ColorWhite, types.B8, 0x0},
		{"Black pawn B4", types.ColorBlack, types.B4, types.A3 | types.C3},
		{"Black pawn A4", types.ColorBlack, types.A4, types.B3},
		{"Black pawn H4", types.ColorBlack, types.H4, types.G3},
		{"Black pawn B1", types.ColorBlack, types.B1, 0x0},
	}

	for _, tc := range testcases {
		got := genPawnAttacks(tc.bitboard, tc.color)
		if got != tc.expected {
			t.Logf("test \"%s\" failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", format.Bitboard(tc.expected, types.PieceWPawn))
			t.Logf("got bitboard:\n\n%s\n\n", format.Bitboard(got, types.PieceWPawn))
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
		{"Knight D4", types.D4, types.C2 | types.E2 | types.B3 | types.F3 | types.B5 |
			types.F5 | types.C6 | types.E6},
		{"Knight A8", types.A8, types.B6 | types.C7},
		{"Knight H1", types.H1, types.F2 | types.G3},
	}

	for _, tc := range testcases {
		got := genKnightAttacks(tc.bitboard)
		if got != tc.expected {
			t.Logf("test \"%s\" failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", format.Bitboard(tc.expected, types.PieceWKnight))
			t.Logf("got bitboard:\n\n%s\n\n", format.Bitboard(got, types.PieceWKnight))
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
		{"King D5", types.D5, types.C4 | types.D4 | types.E4 | types.C5 | types.E5 |
			types.C6 | types.D6 | types.E6},
		{"King A8", types.A8, types.A7 | types.B7 | types.B8},
	}

	for _, tc := range testcases {
		got := genKingAttacks(tc.bitboard)
		if got != tc.expected {
			t.Logf("test \"%s\" failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", format.Bitboard(tc.expected, types.PieceWKing))
			t.Logf("got bitboard:\n\n%s\n\n", format.Bitboard(got, types.PieceWKing))
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
		{"Bishop D5 - Blocked B3", types.D5, types.B3, types.C4 | types.B3 | types.E4 | types.F3 |
			types.G2 | types.H1 | types.C6 | types.B7 | types.A8 | types.E6 | types.F7 | types.G8},
		{"Bishop E2 - Blocked F3", types.E2, types.F3 | types.A6, types.D1 | types.F1 | types.D3 |
			types.F3 | types.C4 | types.B5 | types.A6},
	}

	for _, tc := range testcases {
		got := genBishopAttacks(tc.bitboard, tc.occupancy)
		if got != tc.expected {
			t.Logf("test \"%s\" failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", format.Bitboard(tc.expected, types.PieceWBishop))
			t.Logf("got bitboard:\n\n%s\n\n", format.Bitboard(got, types.PieceWBishop))
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
		{"Rook A1 - No blockers", types.A1, 0x0, types.B1 | types.C1 | types.D1 | types.E1 |
			types.F1 | types.G1 | types.H1 | types.A2 | types.A3 | types.A4 | types.A5 | types.A6 |
			types.A7 | types.A8},
		{"Rook D5 - Bloocked D2, B5, D7", types.D5, types.D2 | types.B5 | types.D7,
			types.D4 | types.D3 | types.D2 | types.C5 | types.B5 | types.E5 | types.F5 |
				types.G5 | types.H5 | types.D6 | types.D7},
	}

	for _, tc := range testcases {
		got := genRookAttacks(tc.bitboard, tc.occupancy)
		if got != tc.expected {
			t.Logf("test \"%s\" failed\n", tc.name)
			t.Logf("expected bitboard:\n\n%s\n\n", format.Bitboard(tc.expected, types.PieceWRook))
			t.Logf("got bitboard:\n\n%s\n\n", format.Bitboard(got, types.PieceWRook))
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
		got := bitutil.BitScan(tc.bitboard)
		if got != tc.expected {
			t.Fatalf("Testcase %s failed: expected %d, got %d", tc.name, tc.expected, got)
		}
	}
}

func TestGenMagicNumber(t *testing.T) {
	t.Logf("\n\n")
	for square := 0; square < 64; square++ {
		t.Logf("%x,\n", genMagicNumber(square, true))
	}

	t.Logf("\n\n")
	for square := 0; square < 64; square++ {
		t.Logf("%x,\n", genMagicNumber(square, false))
	}
}

func TestLookupBishopAttacks(t *testing.T) {
	var occupancy uint64 = types.F2 | types.B3 | types.F4 | types.D5 | types.G7
	for square := uint64(1); square != 0; square <<= 1 {
		got := lookupBishopAttacks(bitutil.BitScan(square), occupancy)
		expected := genBishopAttacks(square, occupancy)

		if got != expected {
			t.Logf("expected:\n\n%s\n\n", format.Bitboard(expected, types.PieceWBishop))
			t.Logf("got:\n\n%s\n\n", format.Bitboard(got, types.PieceWBishop))
			t.FailNow()
		}
	}
}

func TestLookupRookAttacks(t *testing.T) {
	var occupancy uint64 = types.F2 | types.B3 | types.F4 | types.D5 | types.G7
	for square := uint64(1); square != 0; square <<= 1 {
		got := lookupRookAttacks(bitutil.BitScan(square), occupancy)
		expected := genRookAttacks(square, occupancy)

		if got != expected {
			t.Logf("got:\n\n%s\n\n", format.Bitboard(got, types.PieceWRook))
			t.Logf("expected:\n\n%s\n\n", format.Bitboard(expected, types.PieceWRook))
			t.FailNow()
		}
	}
}

func TestLookupQueenAttacks(t *testing.T) {
	var occupancy uint64 = types.F2 | types.B3 | types.F4 | types.D5 | types.G7
	for square := uint64(1); square != 0; square <<= 1 {
		got := lookupQueenAttacks(bitutil.BitScan(square), occupancy)
		expected := genBishopAttacks(square, occupancy) |
			genRookAttacks(square, occupancy)

		if got != expected {
			t.Logf("got:\n\n%s\n\n", format.Bitboard(got, types.PieceWQueen))
			t.Logf("expected:\n\n%s\n\n", format.Bitboard(expected, types.PieceWQueen))
			t.FailNow()
		}
	}
}

func TestIsSquareUnderAttack(t *testing.T) {
	testcases := []struct {
		name     string
		fenStr   string
		square   int
		color    types.Color
		expected bool
	}{
		{
			"square D4 is not attacked by white",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			types.SD4,
			types.ColorWhite,
			false,
		},
		{
			"square D4 is attacked by white queen",
			"8/8/8/8/3p4/8/1Q6/8",
			types.SD4,
			types.ColorWhite,
			true,
		},
		{
			"square C3 is attacked by black pawn",
			"8/8/8/8/3p4/2K5/8/8",
			types.SC3,
			types.ColorBlack,
			true,
		},
		{
			"square F7 is attacked by white knight",
			"rnbq1bnr/2pppkpp/1p6/p3N3/4P3/8/PPPP1PPP/RNB1KB1R",
			types.SF7,
			types.ColorWhite,
			true,
		},
	}

	for _, tc := range testcases {
		bitboards := fen.ToBitboardArray(tc.fenStr)

		got := IsSquareUnderAttack(bitboards, tc.square, tc.color)
		if got != tc.expected {
			t.Logf("\n%s\n", format.Position(types.Position{
				Bitboards: bitboards, ActiveColor: tc.color}))
			t.Fatalf("test \"%s\" failed: got %t, expected %t\n", tc.name, got, tc.expected)
		}
	}
}

func TestGenKingMoves(t *testing.T) {
	testcases := []struct {
		fenStr   string
		square   int
		expected types.MoveList
	}{
		{
			"8/8/8/8/8/8/8/R3K2R w KQ - 0 1",
			types.SE1, types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SD1, types.SE1, types.MoveNormal),
				types.NewMove(types.SF1, types.SE1, types.MoveNormal),
				types.NewMove(types.SD2, types.SE1, types.MoveNormal),
				types.NewMove(types.SE2, types.SE1, types.MoveNormal),
				types.NewMove(types.SF2, types.SE1, types.MoveNormal),
				types.NewMove(types.SG1, types.SE1, types.MoveCastling),
				types.NewMove(types.SC1, types.SE1, types.MoveCastling),
			}, LastMoveIndex: 7},
		},
		{
			"1q4q1/8/8/8/8/8/8/R3K2R w KQ - 0 1",
			types.SE1, types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SD1, types.SE1, types.MoveNormal),
				types.NewMove(types.SF1, types.SE1, types.MoveNormal),
				types.NewMove(types.SD2, types.SE1, types.MoveNormal),
				types.NewMove(types.SE2, types.SE1, types.MoveNormal),
				types.NewMove(types.SF2, types.SE1, types.MoveNormal),
				types.NewMove(types.SC1, types.SE1, types.MoveCastling),
			}, LastMoveIndex: 6},
		},
		{
			"r3k2r/8/8/8/8/8/8/8 b kq - 0 1",
			types.SE8, types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SD7, types.SE8, types.MoveNormal),
				types.NewMove(types.SE7, types.SE8, types.MoveNormal),
				types.NewMove(types.SF7, types.SE8, types.MoveNormal),
				types.NewMove(types.SD8, types.SE8, types.MoveNormal),
				types.NewMove(types.SF8, types.SE8, types.MoveNormal),
				types.NewMove(types.SG8, types.SE8, types.MoveCastling),
				types.NewMove(types.SC8, types.SE8, types.MoveCastling),
			}, LastMoveIndex: 7},
		},
		{
			"r3k2r/8/8/8/8/8/8/2Q3Q1 b kq - 0 1",
			types.SE8, types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SD7, types.SE8, types.MoveNormal),
				types.NewMove(types.SE7, types.SE8, types.MoveNormal),
				types.NewMove(types.SF7, types.SE8, types.MoveNormal),
				types.NewMove(types.SD8, types.SE8, types.MoveNormal),
				types.NewMove(types.SF8, types.SE8, types.MoveNormal),
			}, LastMoveIndex: 5},
		},
	}

	for _, tc := range testcases {
		pos := fen.Parse(tc.fenStr)
		l := types.MoveList{}
		genKingMoves(pos, tc.square, &l)

		for i, move := range l.Moves {
			if tc.expected.Moves[i] != move {
				t.Fatalf("\n%s\nExpected %v\nGot %v", format.Position(pos), tc.expected, l)
			}
		}
	}
}

func TestGenPseudoLegalMoves(t *testing.T) {
	testcases := []struct {
		fenStr   string
		expected types.MoveList
	}{
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SA3, types.SA2, types.MoveNormal),
				types.NewMove(types.SA4, types.SA2, types.MoveNormal),
				types.NewMove(types.SB3, types.SB2, types.MoveNormal),
				types.NewMove(types.SG3, types.SG2, types.MoveNormal),
				types.NewMove(types.SG4, types.SG2, types.MoveNormal),
				types.NewMove(types.SH3, types.SG2, types.MoveNormal),
				types.NewMove(types.SD6, types.SD5, types.MoveNormal),
				types.NewMove(types.SE6, types.SD5, types.MoveNormal),
				types.NewMove(types.SB1, types.SC3, types.MoveNormal),
				types.NewMove(types.SD1, types.SC3, types.MoveNormal),
				types.NewMove(types.SA4, types.SC3, types.MoveNormal),
				types.NewMove(types.SB5, types.SC3, types.MoveNormal),
				types.NewMove(types.SD3, types.SE5, types.MoveNormal),
				types.NewMove(types.SC4, types.SE5, types.MoveNormal),
				types.NewMove(types.SG4, types.SE5, types.MoveNormal),
				types.NewMove(types.SC6, types.SE5, types.MoveNormal),
				types.NewMove(types.SG6, types.SE5, types.MoveNormal),
				types.NewMove(types.SD7, types.SE5, types.MoveNormal),
				types.NewMove(types.SF7, types.SE5, types.MoveNormal),
				types.NewMove(types.SC1, types.SD2, types.MoveNormal),
				types.NewMove(types.SE3, types.SD2, types.MoveNormal),
				types.NewMove(types.SF4, types.SD2, types.MoveNormal),
				types.NewMove(types.SG5, types.SD2, types.MoveNormal),
				types.NewMove(types.SH6, types.SD2, types.MoveNormal),
				types.NewMove(types.SD1, types.SE2, types.MoveNormal),
				types.NewMove(types.SF1, types.SE2, types.MoveNormal),
				types.NewMove(types.SD3, types.SE2, types.MoveNormal),
				types.NewMove(types.SC4, types.SE2, types.MoveNormal),
				types.NewMove(types.SB5, types.SE2, types.MoveNormal),
				types.NewMove(types.SA6, types.SE2, types.MoveNormal),
				types.NewMove(types.SB1, types.SA1, types.MoveNormal),
				types.NewMove(types.SC1, types.SA1, types.MoveNormal),
				types.NewMove(types.SD1, types.SA1, types.MoveNormal),
				types.NewMove(types.SF1, types.SH1, types.MoveNormal),
				types.NewMove(types.SG1, types.SH1, types.MoveNormal),
				types.NewMove(types.SD3, types.SF3, types.MoveNormal),
				types.NewMove(types.SE3, types.SF3, types.MoveNormal),
				types.NewMove(types.SG3, types.SF3, types.MoveNormal),
				types.NewMove(types.SH3, types.SF3, types.MoveNormal),
				types.NewMove(types.SF4, types.SF3, types.MoveNormal),
				types.NewMove(types.SG4, types.SF3, types.MoveNormal),
				types.NewMove(types.SF5, types.SF3, types.MoveNormal),
				types.NewMove(types.SH5, types.SF3, types.MoveNormal),
				types.NewMove(types.SF6, types.SF3, types.MoveNormal),
				types.NewMove(types.SD1, types.SE1, types.MoveNormal),
				types.NewMove(types.SF1, types.SE1, types.MoveNormal),
				types.NewMove(types.SG1, types.SE1, types.MoveCastling),
				types.NewMove(types.SC1, types.SE1, types.MoveCastling),
			}, LastMoveIndex: 35},
		},
	}

	for _, tc := range testcases {
		l := types.MoveList{}
		pos := fen.Parse(tc.fenStr)

		genPseudoLegalMoves(pos, &l)

		for i, move := range l.Moves {
			if move != tc.expected.Moves[i] {
				t.Fatalf("\n%s\nExpected %v\nGot %v", format.Position(pos), tc.expected, l)
			}
		}
	}
}

func TestGenLegalMoves(t *testing.T) {
	testcases := []struct {
		fenStr   string
		expected types.MoveList
	}{
		{
			"8/8/8/8/4P2q/2N5/PPPP1PPP/R1BQKBNR w KQkq - 0 1",
			types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SA3, types.SA2, types.MoveNormal),
				types.NewMove(types.SA4, types.SA2, types.MoveNormal),
				types.NewMove(types.SB3, types.SB2, types.MoveNormal),
				types.NewMove(types.SB4, types.SB2, types.MoveNormal),
				types.NewMove(types.SD3, types.SD2, types.MoveNormal),
				types.NewMove(types.SD4, types.SD2, types.MoveNormal),
				types.NewMove(types.SG3, types.SG2, types.MoveNormal),
				types.NewMove(types.SG4, types.SG2, types.MoveNormal),
				types.NewMove(types.SH3, types.SH2, types.MoveNormal),
				types.NewMove(types.SE5, types.SE4, types.MoveNormal),
				types.NewMove(types.SE2, types.SG1, types.MoveNormal),
				types.NewMove(types.SF3, types.SG1, types.MoveNormal),
				types.NewMove(types.SH3, types.SG1, types.MoveNormal),
				types.NewMove(types.SB1, types.SC3, types.MoveNormal),
				types.NewMove(types.SE2, types.SC3, types.MoveNormal),
				types.NewMove(types.SA4, types.SC3, types.MoveNormal),
				types.NewMove(types.SB5, types.SC3, types.MoveNormal),
				types.NewMove(types.SD5, types.SC3, types.MoveNormal),
				types.NewMove(types.SE2, types.SF1, types.MoveNormal),
				types.NewMove(types.SD3, types.SF1, types.MoveNormal),
				types.NewMove(types.SC4, types.SF1, types.MoveNormal),
				types.NewMove(types.SB5, types.SF1, types.MoveNormal),
				types.NewMove(types.SA6, types.SF1, types.MoveNormal),
				types.NewMove(types.SB1, types.SA1, types.MoveNormal),
				types.NewMove(types.SE2, types.SD1, types.MoveNormal),
				types.NewMove(types.SF3, types.SD1, types.MoveNormal),
				types.NewMove(types.SG4, types.SD1, types.MoveNormal),
				types.NewMove(types.SH5, types.SD1, types.MoveNormal),
				types.NewMove(types.SE2, types.SE1, types.MoveNormal),
			}, LastMoveIndex: 29},
		},
		{
			"3q4/8/8/8/8/8/3p1p2/r3K3 w - - 0 1",
			types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SE2, types.SE1, types.MoveNormal),
				types.NewMove(types.SF2, types.SE1, types.MoveNormal),
			}, LastMoveIndex: 2},
		},
		{
			"2q1k3/8/8/8/8/8/8/R3K2R w Kqkq - 0 1",
			types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SB1, types.SA1, types.MoveNormal),
				types.NewMove(types.SC1, types.SA1, types.MoveNormal),
				types.NewMove(types.SD1, types.SA1, types.MoveNormal),
				types.NewMove(types.SA2, types.SA1, types.MoveNormal),
				types.NewMove(types.SA3, types.SA1, types.MoveNormal),
				types.NewMove(types.SA4, types.SA1, types.MoveNormal),
				types.NewMove(types.SA5, types.SA1, types.MoveNormal),
				types.NewMove(types.SA6, types.SA1, types.MoveNormal),
				types.NewMove(types.SA7, types.SA1, types.MoveNormal),
				types.NewMove(types.SA8, types.SA1, types.MoveNormal),
				types.NewMove(types.SF1, types.SH1, types.MoveNormal),
				types.NewMove(types.SG1, types.SH1, types.MoveNormal),
				types.NewMove(types.SH2, types.SH1, types.MoveNormal),
				types.NewMove(types.SH3, types.SH1, types.MoveNormal),
				types.NewMove(types.SH4, types.SH1, types.MoveNormal),
				types.NewMove(types.SH5, types.SH1, types.MoveNormal),
				types.NewMove(types.SH6, types.SH1, types.MoveNormal),
				types.NewMove(types.SH7, types.SH1, types.MoveNormal),
				types.NewMove(types.SH8, types.SH1, types.MoveNormal),
				types.NewMove(types.SD1, types.SE1, types.MoveNormal),
				types.NewMove(types.SF1, types.SE1, types.MoveNormal),
				types.NewMove(types.SD2, types.SE1, types.MoveNormal),
				types.NewMove(types.SE2, types.SE1, types.MoveNormal),
				types.NewMove(types.SF2, types.SE1, types.MoveNormal),
				types.NewMove(types.SG1, types.SE1, types.MoveCastling),
			}, LastMoveIndex: 25},
		},
		{
			"4k3/4Q3/8/8/7B/8/8/8 b - - 0 1",
			types.MoveList{Moves: [218]types.Move{}, LastMoveIndex: 0},
		},
		{
			"8/K7/1p4p1/7p/1P4k1/8/8/8 w - - 0 1",
			types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SB5, types.SB4, types.MoveNormal),
				types.NewMove(types.SA6, types.SA7, types.MoveNormal),
				types.NewMove(types.SB6, types.SA7, types.MoveNormal),
				types.NewMove(types.SB7, types.SA7, types.MoveNormal),
				types.NewMove(types.SA8, types.SA7, types.MoveNormal),
				types.NewMove(types.SB8, types.SA7, types.MoveNormal),
			}, LastMoveIndex: 6},
		},
		{
			"4q3/8/8/8/4Pp2/8/8/4K3 w - f5 0 1",
			types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SE5, types.SE4, types.MoveNormal),
				types.NewMove(types.SD1, types.SE1, types.MoveNormal),
				types.NewMove(types.SF1, types.SE1, types.MoveNormal),
				types.NewMove(types.SD2, types.SE1, types.MoveNormal),
				types.NewMove(types.SE2, types.SE1, types.MoveNormal),
				types.NewMove(types.SF2, types.SE1, types.MoveNormal),
			}, LastMoveIndex: 6},
		},
		{
			"4q3/8/8/8/4Pp2/5K2/8/8 w - f5 0 1",
			types.MoveList{Moves: [218]types.Move{
				types.NewMove(types.SE5, types.SE4, types.MoveNormal),
				types.NewMove(types.SF5, types.SE4, types.MoveEnPassant),
				types.NewMove(types.SE2, types.SF3, types.MoveNormal),
				types.NewMove(types.SF2, types.SF3, types.MoveNormal),
				types.NewMove(types.SG2, types.SF3, types.MoveNormal),
				types.NewMove(types.SF4, types.SF3, types.MoveNormal),
				types.NewMove(types.SG4, types.SF3, types.MoveNormal),
			}, LastMoveIndex: 7},
		},
		{
			"rnbqkbnr/2pp1Qpp/1p6/p3N3/4P3/8/PPPP1PPP/RNB1KB1R b - - 0 1",
			types.MoveList{Moves: [218]types.Move{}, LastMoveIndex: 0},
		},
	}

	for _, tc := range testcases {
		pos := fen.Parse(tc.fenStr)
		l := types.MoveList{}
		GenLegalMoves(pos, &l)

		if l.LastMoveIndex != tc.expected.LastMoveIndex {
			t.Fatalf("\n%s\nExpected %v\nGot %v", format.Position(pos), tc.expected, l)
		}

		for i, move := range tc.expected.Moves {
			if l.Moves[i] != move {
				t.Fatalf("\n%s\nExpected %v\nGot %v", format.Position(pos), tc.expected, l)
			}
		}
	}
}

func BenchmarkGenPawnAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genPawnAttacks(types.B4, types.ColorWhite)
	}
}

func BenchmarkGenKnightAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genKnightAttacks(types.B4)
	}
}

func BenchmarkGenKingAttakcs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genKingAttacks(types.B4)
	}
}

func BenchmarkGenBishopAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genBishopAttacks(types.D5, types.B3)
	}
}

func BenchmarkGenRookAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genRookAttacks(types.D5, types.B3)
	}
}

func BenchmarkInitBishopReleventOccupancy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initBishopRelevantOccupancy()
	}
}

func BenchmarkInitRookReleventOccupancy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initRookRelevantOccupancy()
	}
}

func BenchmarkLookupBishopAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lookupBishopAttacks(35, 0x0)
	}
}

func BenchmarkLookupRookAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lookupRookAttacks(35, 0x0)
	}
}

func BenchmarkLookupQueenAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lookupQueenAttacks(35, 0x0)
	}
}

func BenchmarkGenMagicNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genMagicNumber(23, false)
	}
}

func BenchmarkIsSquareUnderAttack(b *testing.B) {
	bitboards := fen.ToBitboardArray("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")

	for b.Loop() {
		IsSquareUnderAttack(bitboards, types.SD4, types.ColorWhite)
	}
}

func BenchmarkGenPawnMoves(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genPawnMoves(types.SE4, 0x0, 0x0, 0, types.ColorWhite, &types.MoveList{})
	}
}

func BenchmarkGenKingMoves(b *testing.B) {
	pos := fen.Parse("8/8/8/8/8/8/8/R3K2R w - - 0 1")

	for b.Loop() {
		genKingMoves(pos, types.SE1, &types.MoveList{})
	}
}

func BenchmarkGenPseudoLegalMoves(b *testing.B) {
	pos := fen.Parse("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")

	for b.Loop() {
		genPseudoLegalMoves(pos, &types.MoveList{})
	}
}

func BenchmarkGenLegalMoves(b *testing.B) {
	pos := fen.Parse("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")

	for b.Loop() {
		GenLegalMoves(pos, &types.MoveList{})
	}
}

func BenchmarkInitAttackTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitAttackTables()
	}
}
