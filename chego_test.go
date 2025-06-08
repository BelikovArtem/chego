package chego

import (
	"chego/cli"
	"chego/enum"
	"chego/fen"
	"os"
	"testing"
)

// Used to avoid writing InitAttackTables() each time.
func TestMain(m *testing.M) {
	// Setup.
	InitAttackTables()
	// Tests execution.
	os.Exit(m.Run())
}

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
			t.Logf("test \"%s\" failed\n", tc.name)
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
			t.Logf("test \"%s\" failed\n", tc.name)
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
			t.Logf("test \"%s\" failed\n", tc.name)
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
			t.Logf("test \"%s\" failed\n", tc.name)
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
			t.Logf("test \"%s\" failed\n", tc.name)
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
	var occupancy uint64 = enum.F2 | enum.B3 | enum.F4 | enum.D5 | enum.G7
	for square := uint64(1); square != 0; square <<= 1 {
		got := LookupBishopAttacks(BitScan(square), occupancy)
		expected := GenBishopAttacks(square, occupancy)

		if got != expected {
			t.Logf("expected:\n\n%s\n\n", cli.FormatBitboard(expected, enum.PieceWBishop))
			t.Logf("got:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWBishop))
			t.FailNow()
		}
	}
}

func TestLookupRookAttacks(t *testing.T) {
	var occupancy uint64 = enum.F2 | enum.B3 | enum.F4 | enum.D5 | enum.G7
	for square := uint64(1); square != 0; square <<= 1 {
		got := LookupRookAttacks(BitScan(square), occupancy)
		expected := GenRookAttacks(square, occupancy)

		if got != expected {
			t.Logf("got:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWRook))
			t.Logf("expected:\n\n%s\n\n", cli.FormatBitboard(expected, enum.PieceWRook))
			t.FailNow()
		}
	}
}

func TestLookupQueenAttacks(t *testing.T) {
	var occupancy uint64 = enum.F2 | enum.B3 | enum.F4 | enum.D5 | enum.G7
	for square := uint64(1); square != 0; square <<= 1 {
		got := LookupQueenAttacks(BitScan(square), occupancy)
		expected := GenBishopAttacks(square, occupancy) |
			GenRookAttacks(square, occupancy)

		if got != expected {
			t.Logf("got:\n\n%s\n\n", cli.FormatBitboard(got, enum.PieceWQueen))
			t.Logf("expected:\n\n%s\n\n", cli.FormatBitboard(expected, enum.PieceWQueen))
			t.FailNow()
		}
	}
}

func TestIsSquareUnderAttack(t *testing.T) {
	testcases := []struct {
		name     string
		fenStr   string
		square   int
		color    enum.Color
		expected bool
	}{
		{
			"square D4 is not attacked by white",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			enum.SD4,
			enum.ColorWhite,
			false,
		},
		{
			"square D4 is attacked by white queen",
			"8/8/8/8/3p4/8/1Q6/8",
			enum.SD4,
			enum.ColorWhite,
			true,
		},
		{
			"square C3 is attacked by black pawn",
			"8/8/8/8/3p4/2K5/8/8",
			enum.SC3,
			enum.ColorBlack,
			true,
		},
	}

	for _, tc := range testcases {
		bitboards := fen.ToBitboardArray(tc.fenStr)
		var occupancy uint64
		for pieceType := enum.PieceWPawn; pieceType <= enum.PieceBKing; pieceType++ {
			occupancy |= bitboards[pieceType]
		}

		got := IsSquareUnderAttack(bitboards, occupancy, tc.square, tc.color)
		if got != tc.expected {
			t.Fatalf("test \"%s\" failed: got %t, expected %t\n", tc.name, got, tc.expected)
		}
	}
}

func TestGenPawnsPseudoLegalMoves(t *testing.T) {
	testcases := []struct {
		name     string
		bitboard uint64
		allies   uint64
		enemies  uint64
		color    enum.Color
		expected MoveList
	}{
		{
			"8/8/8/8/p1p1p1p1/8/PPPPPPPP/8 white pawns",
			0xFF00, 0x0, 0x55000000, enum.ColorWhite,
			MoveList{
				[218]Move{
					NewMove(enum.SA3, enum.SA2, 0, enum.MoveNormal),
					NewMove(enum.SB3, enum.SB2, 0, enum.MoveNormal),
					NewMove(enum.SB4, enum.SB2, 0, enum.MoveNormal),
					NewMove(enum.SC3, enum.SC2, 0, enum.MoveNormal),
					NewMove(enum.SD3, enum.SD2, 0, enum.MoveNormal),
					NewMove(enum.SD4, enum.SD2, 0, enum.MoveNormal),
					NewMove(enum.SE3, enum.SE2, 0, enum.MoveNormal),
					NewMove(enum.SF3, enum.SF2, 0, enum.MoveNormal),
					NewMove(enum.SF4, enum.SF2, 0, enum.MoveNormal),
					NewMove(enum.SG3, enum.SG2, 0, enum.MoveNormal),
					NewMove(enum.SH3, enum.SH2, 0, enum.MoveNormal),
					NewMove(enum.SH4, enum.SH2, 0, enum.MoveNormal),
				}, 12,
			},
		},
		{
			"8/8/8/8/p1p1p1p1/8/PPPPPPPP/8 black pawns",
			0x55000000, 0x0, 0xFF00, enum.ColorBlack,
			MoveList{
				[218]Move{
					NewMove(enum.SA3, enum.SA4, 0, enum.MoveNormal),
					NewMove(enum.SC3, enum.SC4, 0, enum.MoveNormal),
					NewMove(enum.SE3, enum.SE4, 0, enum.MoveNormal),
					NewMove(enum.SG3, enum.SG4, 0, enum.MoveNormal),
				}, 4,
			},
		},
		{
			"8/4P3/8/8/8/8/8/8 white quiet promotion",
			enum.E7, 0x0, 0x0, enum.ColorWhite,
			MoveList{
				[218]Move{
					NewMove(enum.SE8, enum.SE7, 0, enum.MovePromotion),
				}, 1,
			},
		},
		{
			"8/8/8/8/8/8/1p6/2P5 black quiet and capture promotions",
			enum.B2, 0x0, enum.C1, enum.ColorBlack,
			MoveList{
				[218]Move{
					NewMove(enum.SB1, enum.SB2, 0, enum.MovePromotion),
					NewMove(enum.SC1, enum.SB2, 0, enum.MovePromotion),
				}, 2,
			},
		},
	}

	for _, tc := range testcases {
		ml := MoveList{}
		genPawnsPseudoLegalMoves(tc.bitboard, tc.allies, tc.enemies, tc.color, &ml)

		for i, move := range ml.Moves {
			if tc.expected.Moves[i] != move {
				t.Fatalf("expected %v, got %v\n", tc.expected, ml)
			}
		}
	}
}

// TODO: TestGenNormalPseudoLegalMoves

func TestGenKingPsuedoLegalMoves(t *testing.T) {
	testcases := []struct {
		name           string
		allies         uint64
		enemies        uint64
		attacked       uint64
		square         int
		castlingRights enum.CastlingFlag
		color          enum.Color
		expected       MoveList
	}{
		{
			"8/8/8/8/8/8/8/R3K2R white king can castle both sides",
			0x81, 0x0, 0x0, enum.SE1, enum.CastlingWhiteKing | enum.CastlingWhiteQueen,
			enum.ColorWhite, MoveList{
				[218]Move{
					NewMove(enum.SD1, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SF1, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SD2, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SE2, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SF2, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SG1, enum.SE1, 0, enum.MoveCastling),
					NewMove(enum.SC1, enum.SE1, 0, enum.MoveCastling),
				}, 7,
			},
		},
		{
			"1q4q1/8/8/8/8/8/8/R3K2R white king cannot castle under attack",
			0x81, 0x42, 0x4242424242424242, enum.SE1, enum.CastlingWhiteKing | enum.CastlingWhiteQueen,
			enum.ColorWhite, MoveList{
				[218]Move{
					NewMove(enum.SD1, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SF1, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SD2, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SE2, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SF2, enum.SE1, 0, enum.MoveNormal),
					NewMove(enum.SC1, enum.SE1, 0, enum.MoveCastling),
				}, 6,
			},
		},
		{
			"r3k2r/8/8/8/8/8/8/8 black king can castle both sides",
			0x8100000000000000, 0x0, 0x0, enum.SE8, enum.CastlingBlackKing | enum.CastlingBlackQueen,
			enum.ColorBlack, MoveList{
				[218]Move{
					NewMove(enum.SD7, enum.SE8, 0, enum.MoveNormal),
					NewMove(enum.SE7, enum.SE8, 0, enum.MoveNormal),
					NewMove(enum.SF7, enum.SE8, 0, enum.MoveNormal),
					NewMove(enum.SD8, enum.SE8, 0, enum.MoveNormal),
					NewMove(enum.SF8, enum.SE8, 0, enum.MoveNormal),
					NewMove(enum.SG8, enum.SE8, 0, enum.MoveCastling),
					NewMove(enum.SC8, enum.SE8, 0, enum.MoveCastling),
				}, 7,
			},
		},
		{
			"r3k2r/8/8/8/8/8/8/2Q3Q1 black king cannot castle under attack",
			0x8100000000000000, 0x44, 0x4444444444444444, enum.SE8, enum.CastlingBlackKing |
				enum.CastlingBlackQueen, enum.ColorBlack, MoveList{
				[218]Move{
					NewMove(enum.SD7, enum.SE8, 0, enum.MoveNormal),
					NewMove(enum.SE7, enum.SE8, 0, enum.MoveNormal),
					NewMove(enum.SF7, enum.SE8, 0, enum.MoveNormal),
					NewMove(enum.SD8, enum.SE8, 0, enum.MoveNormal),
					NewMove(enum.SF8, enum.SE8, 0, enum.MoveNormal),
				}, 5,
			},
		},
	}

	for _, tc := range testcases {
		ml := MoveList{}
		genKingPseudoLegalMoves(tc.square, tc.allies, tc.enemies, tc.attacked,
			tc.castlingRights, &ml, tc.color)

		for i, move := range ml.Moves {
			if tc.expected.Moves[i] != move {
				t.Fatalf("testcase \"%s\" failed: expected %v, got %v\n", tc.name, tc.expected, ml)
			}
		}
	}
}

func TestMakeMove(t *testing.T) {
	testcases := []struct {
		name        string
		fenBefore   string
		fenExpected string
		move        Move
	}{
		{
			"pawn capture",
			"rnbqkbnr/ppp1pppp/8/3p4/2B1P3/5N2/PPPP1PPP/RNBQK2R",
			"rnbqkbnr/ppp1pppp/8/3P4/2B5/5N2/PPPP1PPP/RNBQK2R",
			NewMove(enum.SD5, enum.SE4, 0, enum.MoveNormal),
		},
		{
			"white en passant",
			"rnbqkbnr/ppp1pppp/8/8/1Pp5/5N2/P1PP1PPP/RNBQK2R",
			"rnbqkbnr/ppp1pppp/8/2P5/8/5N2/P1PP1PPP/RNBQK2R",
			NewMove(enum.SC5, enum.SB4, 0, enum.MoveEnPassant),
		},
		{
			"black en passant",
			"2bqkbnr/4p1pp/8/5pP1/8/3N1N2/P1PP1P1P/RqBQK2R",
			"2bqkbnr/4p1pp/8/8/6p1/3N1N2/P1PP1P1P/RqBQK2R",
			NewMove(enum.SG4, enum.SF5, 0, enum.MoveEnPassant),
		},
		{
			"capture promotion",
			"rnbqkbnr/ppP1pppp/8/8/8/5N2/P1PP1PPP/RNBQK2R",
			"rRbqkbnr/pp2pppp/8/8/8/5N2/P1PP1PPP/RNBQK2R",
			NewMove(enum.SB8, enum.SC7, enum.PromotionRook, enum.MovePromotion),
		},
		{
			"promotion",
			"2bqkbnr/4pppp/8/8/8/3N1N2/PpPP1PPP/R1BQK2R",
			"2bqkbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQK2R",
			NewMove(enum.SB1, enum.SB2, enum.PromotionQueen, enum.MovePromotion),
		},
		{
			"white O-O",
			"2bqkbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQK2R",
			"2bqkbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQ1RK1",
			NewMove(enum.SG1, enum.SE1, 0, enum.MoveCastling),
		},
		{
			"black O-O-O",
			"r3kbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQ1RK1",
			"2kr1bnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQ1RK1",
			NewMove(enum.SC8, enum.SE8, 0, enum.MoveCastling),
		},
	}

	for _, tc := range testcases {
		bitboards := fen.ToBitboardArray(tc.fenBefore)
		MakeMove(&bitboards, tc.move)

		got := fen.FromBitboardArray(bitboards)
		if got != tc.fenExpected {
			t.Fatalf("testcase \"%s\" failed: expected %s got %s", tc.name, tc.fenExpected, got)
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

// Using b.Loop here since it is the recommeded approach when the benchmark must have a setup.
func BenchmarkLookupBishopAttacks(b *testing.B) {
	InitAttackTables()

	for b.Loop() {
		LookupBishopAttacks(35, 0x0)
	}
}

func BenchmarkLookupRookAttacks(b *testing.B) {
	InitAttackTables()

	for b.Loop() {
		LookupRookAttacks(35, 0x0)
	}
}

func BenchmarkLookupQueenAttacks(b *testing.B) {
	InitAttackTables()

	for b.Loop() {
		LookupQueenAttacks(35, 0x0)
	}
}

func BenchmarkGenMagicNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenMagicNumber(23, false)
	}
}

func BenchmarkIsSquareUnderAttack(b *testing.B) {
	InitAttackTables()
	bitboards := fen.ToBitboardArray("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")

	for b.Loop() {
		IsSquareUnderAttack(bitboards, 0xFFFF00000000FFFF, enum.SD4, enum.ColorWhite)
	}
}

func BenchmarkGenPawnsPseudoLegalMoves(b *testing.B) {
	InitAttackTables()

	for b.Loop() {
		genPawnsPseudoLegalMoves(0xFF00, 0x0, 0x55000000, enum.ColorWhite, &MoveList{})
	}
}

func BenchmarkGenKingPseudoLegalMoves(b *testing.B) {
	InitAttackTables()

	for b.Loop() {
		genKingPseudoLegalMoves(enum.SE1, 0x81, 0x0, 0x0, enum.CastlingWhiteKing|
			enum.CastlingWhiteQueen, &MoveList{}, enum.ColorWhite)
	}
}

func BenchmarkInitAttackTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitAttackTables()
	}
}

func BenchmarkMakeMove(b *testing.B) {
	InitAttackTables()
	before := fen.ToBitboardArray("rnbqkbnr/pppppppp/8/8/2B1P3/5N2/PPPP1PPP/RNBQK2R")

	for b.Loop() {
		bitboard := before
		MakeMove(&bitboard, NewMove(enum.SG1, enum.SE1, 0, enum.MoveCastling))
	}
}
