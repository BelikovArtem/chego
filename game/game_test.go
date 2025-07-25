package game

import (
	"os"
	"testing"

	"github.com/BelikovArtem/chego/fen"
	"github.com/BelikovArtem/chego/movegen"
	"github.com/BelikovArtem/chego/types"
)

func TestMain(m *testing.M) {
	// Setup.
	movegen.InitAttackTables()
	// Tests and benchmarks execution.
	os.Exit(m.Run())
}

func TestPushMove(t *testing.T) {
	testcases := []struct {
		name                    string
		move                    types.Move
		expectedEnPassantTarget int
		expectedCastlingRights  types.CastlingRights
		expectedActiveColor     types.Color
	}{
		{
			"h4", types.NewMove(types.SH4, types.SH2, types.MoveNormal),
			types.SH3, 0xF, types.ColorBlack,
		},
		{
			"e5", types.NewMove(types.SE5, types.SE7, types.MoveNormal),
			types.SE6, 0xF, types.ColorWhite,
		},
		{
			"c4", types.NewMove(types.SC4, types.SC2, types.MoveNormal),
			types.SC3, 0xF, types.ColorBlack,
		},
		{
			"nf6", types.NewMove(types.SF6, types.SG8, types.MoveNormal),
			0, 0xF, types.ColorWhite,
		},
		{
			"e3", types.NewMove(types.SE3, types.SE2, types.MoveNormal),
			0, 0xF, types.ColorBlack,
		},
		{
			"c6", types.NewMove(types.SC6, types.SC7, types.MoveNormal),
			0, 0xF, types.ColorWhite,
		},
		{
			"g4", types.NewMove(types.SG4, types.SG2, types.MoveNormal),
			types.SG3, 0xF, types.ColorBlack,
		},
		{
			"g6", types.NewMove(types.SG6, types.SG7, types.MoveNormal),
			0, 0xF, types.ColorWhite,
		},
		{
			"d4", types.NewMove(types.SD4, types.SD2, types.MoveNormal),
			types.SD3, 0xF, types.ColorBlack,
		},
		{
			"d6", types.NewMove(types.SD6, types.SD7, types.MoveNormal),
			0, 0xF, types.ColorWhite,
		},
		{
			"g5", types.NewMove(types.SG5, types.SG4, types.MoveNormal),
			0, 0xF, types.ColorBlack,
		},
		{
			"nh5", types.NewMove(types.SH5, types.SF6, types.MoveNormal),
			0, 0xF, types.ColorWhite,
		},
		{
			"dxe5", types.NewMove(types.SE5, types.SD4, types.MoveNormal),
			0, 0xF, types.ColorBlack,
		},
		{
			"dxe5", types.NewMove(types.SE5, types.SD6, types.MoveNormal),
			0, 0xF, types.ColorWhite,
		},
		{
			"Qcd8+", types.NewMove(types.SD8, types.SD1, types.MoveNormal),
			0, 0xF, types.ColorBlack,
		},
		{
			"kxd8", types.NewMove(types.SD8, types.SE8, types.MoveNormal),
			0, 0x3, types.ColorWhite,
		},
		{
			"Nf3", types.NewMove(types.SF3, types.SG1, types.MoveNormal),
			0, 0x3, types.ColorBlack,
		},
		{
			"bg7", types.NewMove(types.SG7, types.SF8, types.MoveNormal),
			0, 0x3, types.ColorWhite,
		},
		{
			"Nc3", types.NewMove(types.SC3, types.SB1, types.MoveNormal),
			0, 0x3, types.ColorBlack,
		},
		{
			"bg4", types.NewMove(types.SG4, types.SC8, types.MoveNormal),
			0, 0x3, types.ColorWhite,
		},
		{
			"Be2", types.NewMove(types.SE2, types.SF1, types.MoveNormal),
			0, 0x3, types.ColorBlack,
		},
		{
			"nd7", types.NewMove(types.SD7, types.SB8, types.MoveNormal),
			0, 0x3, types.ColorWhite,
		},
		{
			"Nd2", types.NewMove(types.SD2, types.SF3, types.MoveNormal),
			0, 0x3, types.ColorBlack,
		},
		{
			"bxe2", types.NewMove(types.SE2, types.SG4, types.MoveNormal),
			0, 0x3, types.ColorWhite,
		},
		{
			"Kxe2", types.NewMove(types.SE2, types.SE1, types.MoveNormal),
			0, 0x0, types.ColorBlack,
		},
		{
			"h6", types.NewMove(types.SH6, types.SH7, types.MoveNormal),
			0, 0x0, types.ColorWhite,
		},
		{
			"Kde4", types.NewMove(types.SE4, types.SD2, types.MoveNormal),
			0, 0x0, types.ColorBlack,
		},
		{
			"hxg5", types.NewMove(types.SG5, types.SH6, types.MoveNormal),
			0, 0x0, types.ColorWhite,
		},
		{
			"Nxg5", types.NewMove(types.SG5, types.SE4, types.MoveNormal),
			0, 0x0, types.ColorBlack,
		},
		{
			"ke7", types.NewMove(types.SE7, types.SD8, types.MoveNormal),
			0, 0x0, types.ColorWhite,
		},
		{
			"Rb1", types.NewMove(types.SB1, types.SA1, types.MoveNormal),
			0, 0x0, types.ColorBlack,
		},
		{
			"rg8", types.NewMove(types.SG8, types.SH8, types.MoveNormal),
			0, 0x0, types.ColorWhite,
		},
		{
			"rg1", types.NewMove(types.SG1, types.SH1, types.MoveNormal),
			0, 0x0, types.ColorBlack,
		},
		{
			"rb8", types.NewMove(types.SB8, types.SA8, types.MoveNormal),
			0, 0x0, types.ColorWhite,
		},
	}

	// Standard initial position.
	game := NewGame()

	for _, tc := range testcases {
		game.PushMove(tc.move)

		if game.Position.EPTarget != tc.expectedEnPassantTarget {
			t.Fatalf("test \"%s\" failed: expected EP square %d, got %d", tc.name,
				tc.expectedEnPassantTarget, game.Position.EPTarget)
		}
		if game.Position.CastlingRights != tc.expectedCastlingRights {
			t.Fatalf("test \"%s\" failed: expected castling rights %b, got %b", tc.name,
				tc.expectedCastlingRights, game.Position.CastlingRights)
		}
		if game.Position.ActiveColor != tc.expectedActiveColor {
			t.Fatalf("test \"%s\" failed: expected active color %b, got %b", tc.name,
				tc.expectedActiveColor, game.Position.ActiveColor)
		}
	}
}

func TestPopMove(t *testing.T) {
	testcases := []struct {
		CompletedMoves []types.Move
		expectedFen    string
	}{
		{[]types.Move{
			types.NewMove(types.SE4, types.SE2, types.MoveNormal),
		},
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
		{[]types.Move{
			types.NewMove(types.SE4, types.SE2, types.MoveNormal),
			types.NewMove(types.SB5, types.SB7, types.MoveNormal),
		},
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"},
	}

	for _, tc := range testcases {
		g := NewGame()

		for _, move := range tc.CompletedMoves {
			g.PushMove(move)
		}

		g.PopMove()
		got := fen.Serialize(g.Position)
		if got != tc.expectedFen {
			t.Fatalf("expected fen \"%s\", got \"%s\"", tc.expectedFen, got)
		}
	}
}

func TestIsThreefoldRepetition(t *testing.T) {
	testcases := []struct {
		moveStack []types.Move
		expected  bool
	}{
		{[]types.Move{
			types.NewMove(types.SH3, types.SH2, types.MoveNormal),
			types.NewMove(types.SH6, types.SH7, types.MoveNormal),
			types.NewMove(types.SH2, types.SH1, types.MoveNormal),
			types.NewMove(types.SH7, types.SH8, types.MoveNormal),
			types.NewMove(types.SH1, types.SH2, types.MoveNormal),
			types.NewMove(types.SH8, types.SH7, types.MoveNormal),
			types.NewMove(types.SH2, types.SH1, types.MoveNormal),
			types.NewMove(types.SH7, types.SH8, types.MoveNormal),
			types.NewMove(types.SH1, types.SH2, types.MoveNormal),
			types.NewMove(types.SH8, types.SH7, types.MoveNormal),
			types.NewMove(types.SH2, types.SH1, types.MoveNormal),
			types.NewMove(types.SH7, types.SH8, types.MoveNormal),
		}, true},
		{[]types.Move{
			types.NewMove(types.SH3, types.SH2, types.MoveNormal),
			types.NewMove(types.SH6, types.SH7, types.MoveNormal),
			types.NewMove(types.SH2, types.SH1, types.MoveNormal),
			types.NewMove(types.SH7, types.SH8, types.MoveNormal),
			types.NewMove(types.SH1, types.SH2, types.MoveNormal),
			types.NewMove(types.SH8, types.SH7, types.MoveNormal),
			types.NewMove(types.SH2, types.SH1, types.MoveNormal),
			types.NewMove(types.SH7, types.SH8, types.MoveNormal),
			types.NewMove(types.SH1, types.SH2, types.MoveNormal),
			types.NewMove(types.SH8, types.SH7, types.MoveNormal),
			types.NewMove(types.SH2, types.SH1, types.MoveNormal),
			types.NewMove(types.SD6, types.SD7, types.MoveNormal),
			types.NewMove(types.SD3, types.SD2, types.MoveNormal),
			types.NewMove(types.SD7, types.SD8, types.MoveNormal),
			types.NewMove(types.SD2, types.SD1, types.MoveNormal),
			types.NewMove(types.SD8, types.SD7, types.MoveNormal),
			types.NewMove(types.SD1, types.SD2, types.MoveNormal),
		}, false},
		{[]types.Move{
			types.NewMove(types.SE4, types.SE2, types.MoveNormal),
			types.NewMove(types.SE5, types.SE7, types.MoveNormal),
			types.NewMove(types.SF3, types.SG1, types.MoveNormal),
			types.NewMove(types.SC6, types.SB8, types.MoveNormal),
			types.NewMove(types.SG1, types.SF3, types.MoveNormal),
			types.NewMove(types.SB8, types.SC6, types.MoveNormal),
			types.NewMove(types.SF3, types.SG1, types.MoveNormal),
			types.NewMove(types.SC6, types.SB8, types.MoveNormal),
			types.NewMove(types.SG1, types.SF3, types.MoveNormal),
			types.NewMove(types.SB8, types.SC6, types.MoveNormal),
		}, true},
		{[]types.Move{
			types.NewMove(types.SD4, types.SD2, types.MoveNormal),
			types.NewMove(types.SD6, types.SD7, types.MoveNormal),
			types.NewMove(types.SD5, types.SD4, types.MoveNormal),
			types.NewMove(types.SA6, types.SA7, types.MoveNormal),
			types.NewMove(types.SA3, types.SA2, types.MoveNormal),
			types.NewMove(types.SC5, types.SC7, types.MoveNormal),
			types.NewMove(types.SA2, types.SA1, types.MoveNormal),
			types.NewMove(types.SA7, types.SA8, types.MoveNormal),
			types.NewMove(types.SA1, types.SA2, types.MoveNormal),
			types.NewMove(types.SA8, types.SA7, types.MoveNormal),
			types.NewMove(types.SA2, types.SA1, types.MoveNormal),
			types.NewMove(types.SA7, types.SA8, types.MoveNormal),
			types.NewMove(types.SA1, types.SA2, types.MoveNormal),
			types.NewMove(types.SA8, types.SA7, types.MoveNormal),
		}, false},
		{[]types.Move{
			types.NewMove(types.SD4, types.SD2, types.MoveNormal),
			types.NewMove(types.SD6, types.SD7, types.MoveNormal),
			types.NewMove(types.SD5, types.SD4, types.MoveNormal),
			types.NewMove(types.SA6, types.SA7, types.MoveNormal),
			types.NewMove(types.SA3, types.SA2, types.MoveNormal),
			types.NewMove(types.SC5, types.SC7, types.MoveNormal),
			types.NewMove(types.SA2, types.SA1, types.MoveNormal),
			types.NewMove(types.SA7, types.SA8, types.MoveNormal),
			types.NewMove(types.SA1, types.SA2, types.MoveNormal),
			types.NewMove(types.SA8, types.SA7, types.MoveNormal),
			types.NewMove(types.SA2, types.SA1, types.MoveNormal),
			types.NewMove(types.SA7, types.SA8, types.MoveNormal),
			types.NewMove(types.SA1, types.SA2, types.MoveNormal),
			types.NewMove(types.SA8, types.SA7, types.MoveNormal),
			types.NewMove(types.SA2, types.SA1, types.MoveNormal),
			types.NewMove(types.SA7, types.SA8, types.MoveNormal),
		}, true},
	}

	for _, tc := range testcases {
		game := NewGame()

		for _, move := range tc.moveStack {
			game.PushMove(move)
		}

		got := game.IsThreefoldRepetition()
		if tc.expected != got {
			t.Fatalf("expected %t, got %t", tc.expected, got)
		}
	}
}

func TestIsInsufficientMaterial(t *testing.T) {
	testcases := []struct {
		fenString string
		expected  bool
	}{
		{"3k1n2/8/8/8/8/5B2/4K3/8", false},
		{"3k4/8/8/8/8/8/4K3/8", true},
		{"3k4/8/8/8/8/5P2/4K3/8", false},
		{"3k4/2b5/8/8/8/8/4K3/8", true},
		{"3k4/8/8/8/8/8/3NK3/8", true},
		{"3k4/2b5/8/8/8/4B3/4K3/8", true},
		{"3k4/2b5/8/8/8/3B4/4K3/8", false},
		{"8/8/8/8/8/8/1n6/KN6", true},
	}

	game := NewGame()
	for _, tc := range testcases {
		game.Position.Bitboards = fen.ToBitboardArray(tc.fenString)

		got := game.IsInsufficientMaterial()
		if got != tc.expected {
			t.Fatalf("expected: %t, got: %t", tc.expected, got)
		}
	}
}

func TestIsCheckmate(t *testing.T) {
	testcases := []struct {
		fenString string
		expected  bool
	}{
		{"rnb1kbnr/pppp1ppp/4p3/8/6Pq/3P1P2/PPP1P2P/RNBQKBNR w KQkq - 0 1", false},
		{"rnb1kbnr/pppp1ppp/4p3/8/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 0 1", true},
		{"rnb1kbnr/pppp1ppp/4p3/8/6Pq/3P1P2/PPP1PN1P/R1BQKBNR w KQkq - 0 1", false},
	}

	game := NewGame()
	for _, tc := range testcases {
		game.Position = fen.Parse(tc.fenString)
		movegen.GenLegalMoves(game.Position, &game.LegalMoves)

		got := game.IsCheckmate()
		if got != tc.expected {
			t.Fatalf("expected: %t, got: %t", tc.expected, got)
		}
	}
}

func BenchmarkPushMove(b *testing.B) {
	game := NewGame()
	pos := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	for b.Loop() {
		game.Position = pos
		game.PushMove(types.NewMove(types.SE4, types.SE2, types.MoveNormal))
	}
}

// TODO: BenchmarkPopMove

func BenchmarkIsThreefoldRepetition(b *testing.B) {
	game := NewGame()
	moveStack := []types.Move{
		types.NewMove(types.SH3, types.SH2, types.MoveNormal),
		types.NewMove(types.SH6, types.SH7, types.MoveNormal),
		types.NewMove(types.SH2, types.SH1, types.MoveNormal),
		types.NewMove(types.SH7, types.SH8, types.MoveNormal),
		types.NewMove(types.SH1, types.SH2, types.MoveNormal),
		types.NewMove(types.SH8, types.SH7, types.MoveNormal),
		types.NewMove(types.SH2, types.SH1, types.MoveNormal),
		types.NewMove(types.SH7, types.SH8, types.MoveNormal),
		types.NewMove(types.SH1, types.SH2, types.MoveNormal),
		types.NewMove(types.SH8, types.SH7, types.MoveNormal),
		types.NewMove(types.SH2, types.SH1, types.MoveNormal),
		types.NewMove(types.SD6, types.SD7, types.MoveNormal),
		types.NewMove(types.SD3, types.SD2, types.MoveNormal),
		types.NewMove(types.SD7, types.SD8, types.MoveNormal),
		types.NewMove(types.SD2, types.SD1, types.MoveNormal),
		types.NewMove(types.SD8, types.SD7, types.MoveNormal),
		types.NewMove(types.SD1, types.SD2, types.MoveNormal),
	}

	for i, move := range moveStack {
		game.PushMove(move)
		movegen.GenLegalMoves(game.Position, &game.LegalMoves)
		if i < len(moveStack)-1 {
			key := repetitionKey(game.Position, game.LegalMoves)
			game.Repetitions[key]++
		}
	}

	for b.Loop() {
		game.IsThreefoldRepetition()
	}
}

func BenchmarkIsInsufficientMaterial(b *testing.B) {
	game := NewGame()

	for b.Loop() {
		game.IsInsufficientMaterial()
	}
}

func BenchmarkIsCheckmate(b *testing.B) {
	game := NewGame()

	for b.Loop() {
		game.IsCheckmate()
	}
}
