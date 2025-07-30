package chego

import "testing"

func TestMakeMove(t *testing.T) {
	testcases := []struct {
		name     string
		fenStr   string
		expected string
		move     Move
	}{
		{
			"pawn capture",
			"rnbqkbnr/ppp1pppp/8/3p4/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1",
			"rnbqkbnr/ppp1pppp/8/3P4/2B5/5N2/PPPP1PPP/RNBQK2R b KQkq - 0 1",
			NewMove(SD5, SE4, MoveNormal),
		},
		{
			"white en passant",
			"rnbqkbnr/ppp1pppp/8/8/1Pp5/5N2/P1PP1PPP/RNBQK2R w KQkq b3 0 1",
			"rnbqkbnr/ppp1pppp/8/2P5/8/5N2/P1PP1PPP/RNBQK2R b KQkq - 0 1",
			NewMove(SC5, SB4, MoveEnPassant),
		},
		{
			"black en passant",
			"2bqkbnr/4p1pp/8/5pP1/8/3N1N2/P1PP1P1P/RqBQK2R b KQkq g4 0 1",
			"2bqkbnr/4p1pp/8/8/6p1/3N1N2/P1PP1P1P/RqBQK2R w KQkq - 0 2",
			NewMove(SG4, SF5, MoveEnPassant),
		},
		{
			"capture promotion",
			"rnbqkbnr/ppP1pppp/8/8/8/5N2/P1PP1PPP/RNBQK2R w KQkq - 0 1",
			"rRbqkbnr/pp2pppp/8/8/8/5N2/P1PP1PPP/RNBQK2R b KQkq - 0 1",
			NewPromotionMove(SB8, SC7, PromotionRook),
		},
		{
			"promotion",
			"2bqkbnr/4pppp/8/8/8/3N1N2/PpPP1PPP/R1BQK2R b KQkq - 0 1",
			"2bqkbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQK2R w KQkq - 0 2",
			NewPromotionMove(SB1, SB2, PromotionQueen),
		},
		{
			"white O-O",
			"2bqkbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQK2R w KQkq - 0 1",
			"2bqkbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQ1RK1 b kq - 1 1",
			NewMove(SG1, SE1, MoveCastling),
		},
		{
			"black O-O-O",
			"r3kbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQ1RK1 b KQkq - 0 1",
			"2kr1bnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQ1RK1 w KQ - 1 2",
			NewMove(SC8, SE8, MoveCastling),
		},
		{
			"white rook",
			"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			"r3k2r/8/8/8/8/8/8/1R2K2R b Kkq - 1 1",
			NewMove(SB1, SA1, MoveNormal),
		},
		{
			"black rook",
			"r3k2r/8/8/8/8/8/8/1R2K2R b Kkq - 1 1",
			"r3k1r1/8/8/8/8/8/8/1R2K2R w Kq - 2 2",
			NewMove(SG8, SH8, MoveNormal),
		},
		{
			"white double pawn push",
			"4k3/4p3/8/8/8/8/4P3/4K3 w - - 0 1",
			"4k3/4p3/8/8/4P3/8/8/4K3 b - e3 0 1",
			NewMove(SE4, SE2, MoveNormal),
		},
		{
			"black double pawn push",
			"4k3/4p3/8/8/4P3/8/8/4K3 b - e3 0 1",
			"4k3/8/8/4p3/4P3/8/8/4K3 w - e6 0 2",
			NewMove(SE5, SE7, MoveNormal),
		},
	}

	for _, tc := range testcases {
		pos := ParseFEN(tc.fenStr)
		pos.MakeMove(tc.move)

		got := SerializeFEN(pos)
		if got != tc.expected {
			t.Fatalf("test \"%s\" failed: expected %s got %s", tc.name, tc.expected, got)
		}
	}
}

func BenchmarkMakeMove(b *testing.B) {
	before := ParseFEN("rnbqkbnr/pppppppp/8/8/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1")

	for b.Loop() {
		pos := before
		pos.MakeMove(NewMove(SG1, SE1, MoveCastling))
	}
}
