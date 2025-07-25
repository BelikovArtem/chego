package types_test

import (
	"testing"

	"github.com/BelikovArtem/chego/fen"
	"github.com/BelikovArtem/chego/types"
)

func TestMakeMove(t *testing.T) {
	testcases := []struct {
		name     string
		fenStr   string
		expected string
		move     types.Move
	}{
		{
			"pawn capture",
			"rnbqkbnr/ppp1pppp/8/3p4/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1",
			"rnbqkbnr/ppp1pppp/8/3P4/2B5/5N2/PPPP1PPP/RNBQK2R b KQkq - 0 1",
			types.NewMove(types.SD5, types.SE4, types.MoveNormal),
		},
		{
			"white en passant",
			"rnbqkbnr/ppp1pppp/8/8/1Pp5/5N2/P1PP1PPP/RNBQK2R w KQkq b3 0 1",
			"rnbqkbnr/ppp1pppp/8/2P5/8/5N2/P1PP1PPP/RNBQK2R b KQkq - 0 1",
			types.NewMove(types.SC5, types.SB4, types.MoveEnPassant),
		},
		{
			"black en passant",
			"2bqkbnr/4p1pp/8/5pP1/8/3N1N2/P1PP1P1P/RqBQK2R b KQkq g4 0 1",
			"2bqkbnr/4p1pp/8/8/6p1/3N1N2/P1PP1P1P/RqBQK2R w KQkq - 0 2",
			types.NewMove(types.SG4, types.SF5, types.MoveEnPassant),
		},
		{
			"capture promotion",
			"rnbqkbnr/ppP1pppp/8/8/8/5N2/P1PP1PPP/RNBQK2R w KQkq - 0 1",
			"rRbqkbnr/pp2pppp/8/8/8/5N2/P1PP1PPP/RNBQK2R b KQkq - 0 1",
			types.NewPromotionMove(types.SB8, types.SC7, types.PromotionRook),
		},
		{
			"promotion",
			"2bqkbnr/4pppp/8/8/8/3N1N2/PpPP1PPP/R1BQK2R b KQkq - 0 1",
			"2bqkbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQK2R w KQkq - 0 2",
			types.NewPromotionMove(types.SB1, types.SB2, types.PromotionQueen),
		},
		{
			"white O-O",
			"2bqkbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQK2R w KQkq - 0 1",
			"2bqkbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQ1RK1 b kq - 0 1",
			types.NewMove(types.SG1, types.SE1, types.MoveCastling),
		},
		{
			"black O-O-O",
			"r3kbnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQ1RK1 b KQkq - 0 1",
			"2kr1bnr/4pppp/8/8/8/3N1N2/P1PP1PPP/RqBQ1RK1 w KQ - 0 2",
			types.NewMove(types.SC8, types.SE8, types.MoveCastling),
		},
		{
			"white rook",
			"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			"r3k2r/8/8/8/8/8/8/1R2K2R b Kkq - 1 1",
			types.NewMove(types.SB1, types.SA1, types.MoveNormal),
		},
		{
			"black rook",
			"r3k2r/8/8/8/8/8/8/1R2K2R b Kkq - 1 1",
			"r3k1r1/8/8/8/8/8/8/1R2K2R w Kq - 2 2",
			types.NewMove(types.SG8, types.SH8, types.MoveNormal),
		},
		{
			"white double pawn push",
			"4k3/4p3/8/8/8/8/4P3/4K3 w - - 0 1",
			"4k3/4p3/8/8/4P3/8/8/4K3 b - e3 0 1",
			types.NewMove(types.SE4, types.SE2, types.MoveNormal),
		},
		{
			"black double pawn push",
			"4k3/4p3/8/8/4P3/8/8/4K3 b - e3 0 1",
			"4k3/8/8/4p3/4P3/8/8/4K3 w - e6 0 2",
			types.NewMove(types.SE5, types.SE7, types.MoveNormal),
		},
	}

	for _, tc := range testcases {
		pos := fen.Parse(tc.fenStr)
		pos.MakeMove(tc.move)

		got := fen.Serialize(pos)
		if got != tc.expected {
			t.Fatalf("test \"%s\" failed: expected %s got %s", tc.name, tc.expected, got)
		}
	}
}

func BenchmarkMakeMove(b *testing.B) {
	before := fen.Parse("rnbqkbnr/pppppppp/8/8/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1")

	for b.Loop() {
		pos := before
		pos.MakeMove(types.NewMove(types.SG1, types.SE1, types.MoveCastling))
	}
}
