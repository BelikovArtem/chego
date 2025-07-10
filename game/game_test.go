package game

import (
	"chego/enum"
	"chego/fen"
	"chego/movegen"
	"os"
	"testing"
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
		move                    movegen.Move
		expectedEnPassantTarget int
		expectedCastlingRights  enum.CastlingFlag
		expectedActiveColor     enum.Color
	}{
		{
			"h4", movegen.NewMove(enum.SH4, enum.SH2, 0, enum.MoveNormal),
			enum.SH3, 0xF, enum.ColorBlack,
		},
		{
			"e5", movegen.NewMove(enum.SE5, enum.SE7, 0, enum.MoveNormal),
			enum.SE6, 0xF, enum.ColorWhite,
		},
		{
			"c4", movegen.NewMove(enum.SC4, enum.SC2, 0, enum.MoveNormal),
			enum.SC3, 0xF, enum.ColorBlack,
		},
		{
			"nf6", movegen.NewMove(enum.SF6, enum.SG8, 0, enum.MoveNormal),
			0, 0xF, enum.ColorWhite,
		},
		{
			"e3", movegen.NewMove(enum.SE3, enum.SE2, 0, enum.MoveNormal),
			0, 0xF, enum.ColorBlack,
		},
		{
			"c6", movegen.NewMove(enum.SC6, enum.SC7, 0, enum.MoveNormal),
			0, 0xF, enum.ColorWhite,
		},
		{
			"g4", movegen.NewMove(enum.SG4, enum.SG2, 0, enum.MoveNormal),
			enum.SG3, 0xF, enum.ColorBlack,
		},
		{
			"g6", movegen.NewMove(enum.SG6, enum.SG7, 0, enum.MoveNormal),
			0, 0xF, enum.ColorWhite,
		},
		{
			"d4", movegen.NewMove(enum.SD4, enum.SD2, 0, enum.MoveNormal),
			enum.SD3, 0xF, enum.ColorBlack,
		},
		{
			"d6", movegen.NewMove(enum.SD6, enum.SD7, 0, enum.MoveNormal),
			0, 0xF, enum.ColorWhite,
		},
		{
			"g5", movegen.NewMove(enum.SG5, enum.SG4, 0, enum.MoveNormal),
			0, 0xF, enum.ColorBlack,
		},
		{
			"nh5", movegen.NewMove(enum.SH5, enum.SF6, 0, enum.MoveNormal),
			0, 0xF, enum.ColorWhite,
		},
		{
			"dxe5", movegen.NewMove(enum.SE5, enum.SD4, 0, enum.MoveNormal),
			0, 0xF, enum.ColorBlack,
		},
		{
			"dxe5", movegen.NewMove(enum.SE5, enum.SD6, 0, enum.MoveNormal),
			0, 0xF, enum.ColorWhite,
		},
		{
			"Qcd8+", movegen.NewMove(enum.SD8, enum.SD1, 0, enum.MoveNormal),
			0, 0xF, enum.ColorBlack,
		},
		{
			"kxd8", movegen.NewMove(enum.SD8, enum.SE8, 0, enum.MoveNormal),
			0, 0x3, enum.ColorWhite,
		},
		{
			"Nf3", movegen.NewMove(enum.SF3, enum.SG1, 0, enum.MoveNormal),
			0, 0x3, enum.ColorBlack,
		},
		{
			"bg7", movegen.NewMove(enum.SG7, enum.SF8, 0, enum.MoveNormal),
			0, 0x3, enum.ColorWhite,
		},
		{
			"Nc3", movegen.NewMove(enum.SC3, enum.SB1, 0, enum.MoveNormal),
			0, 0x3, enum.ColorBlack,
		},
		{
			"bg4", movegen.NewMove(enum.SG4, enum.SC8, 0, enum.MoveNormal),
			0, 0x3, enum.ColorWhite,
		},
		{
			"Be2", movegen.NewMove(enum.SE2, enum.SF1, 0, enum.MoveNormal),
			0, 0x3, enum.ColorBlack,
		},
		{
			"nd7", movegen.NewMove(enum.SD7, enum.SB8, 0, enum.MoveNormal),
			0, 0x3, enum.ColorWhite,
		},
		{
			"Nd2", movegen.NewMove(enum.SD2, enum.SF3, 0, enum.MoveNormal),
			0, 0x3, enum.ColorBlack,
		},
		{
			"bxe2", movegen.NewMove(enum.SE2, enum.SG4, 0, enum.MoveNormal),
			0, 0x3, enum.ColorWhite,
		},
		{
			"Kxe2", movegen.NewMove(enum.SE2, enum.SE1, 0, enum.MoveNormal),
			0, 0x0, enum.ColorBlack,
		},
		{
			"h6", movegen.NewMove(enum.SH6, enum.SH7, 0, enum.MoveNormal),
			0, 0x0, enum.ColorWhite,
		},
		{
			"Kde4", movegen.NewMove(enum.SE4, enum.SD2, 0, enum.MoveNormal),
			0, 0x0, enum.ColorBlack,
		},
		{
			"hxg5", movegen.NewMove(enum.SG5, enum.SH6, 0, enum.MoveNormal),
			0, 0x0, enum.ColorWhite,
		},
		{
			"Nxg5", movegen.NewMove(enum.SG5, enum.SE4, 0, enum.MoveNormal),
			0, 0x0, enum.ColorBlack,
		},
		{
			"ke7", movegen.NewMove(enum.SE7, enum.SD8, 0, enum.MoveNormal),
			0, 0x0, enum.ColorWhite,
		},
		{
			"Rb1", movegen.NewMove(enum.SB1, enum.SA1, 0, enum.MoveNormal),
			0, 0x0, enum.ColorBlack,
		},
		{
			"rg8", movegen.NewMove(enum.SG8, enum.SH8, 0, enum.MoveNormal),
			0, 0x0, enum.ColorWhite,
		},
		{
			"rg1", movegen.NewMove(enum.SG1, enum.SH1, 0, enum.MoveNormal),
			0, 0x0, enum.ColorBlack,
		},
		{
			"rb8", movegen.NewMove(enum.SB8, enum.SA8, 0, enum.MoveNormal),
			0, 0x0, enum.ColorWhite,
		},
	}

	// Standard initial position.
	game := NewGame()

	for _, tc := range testcases {
		game.PushMove(tc.move)

		if game.EnPassantTarget != tc.expectedEnPassantTarget {
			t.Fatalf("test \"%s\" failed: expected EP square %d, got %d", tc.name,
				tc.expectedEnPassantTarget, game.EnPassantTarget)
		}
		if game.CastlingRights != tc.expectedCastlingRights {
			t.Fatalf("test \"%s\" failed: expected castling rights %b, got %b", tc.name,
				tc.expectedCastlingRights, game.CastlingRights)
		}
		if game.ActiveColor != tc.expectedActiveColor {
			t.Fatalf("test \"%s\" failed: expected active color %b, got %b", tc.name,
				tc.expectedActiveColor, game.ActiveColor)
		}
	}
}

func TestPopMove(t *testing.T) {
	testcases := []struct {
		CompletedMoves []movegen.Move
		expectedFen    string
	}{
		{[]movegen.Move{
			movegen.NewMove(enum.SE4, enum.SE2, 0, enum.MoveNormal),
		},
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
		{[]movegen.Move{
			movegen.NewMove(enum.SE4, enum.SE2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SB5, enum.SB7, 0, enum.MoveNormal),
		},
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"},
	}

	for _, tc := range testcases {
		g := NewGame()

		for _, move := range tc.CompletedMoves {
			g.PushMove(move)
		}

		g.PopMove()
		got := fen.Serialize(g.Bitboards, g.ActiveColor, g.CastlingRights, g.EnPassantTarget,
			g.HalfmoveCnt, g.FullmoveCnt)
		if got != tc.expectedFen {
			t.Fatalf("expected fen \"%s\", got \"%s\"", tc.expectedFen, got)
		}
	}
}

func TestIsThreefoldRepetition(t *testing.T) {
	testcases := []struct {
		moveStack []movegen.Move
		expected  bool
	}{
		{[]movegen.Move{
			movegen.NewMove(enum.SH3, enum.SH2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH6, enum.SH7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH2, enum.SH1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH7, enum.SH8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH1, enum.SH2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH8, enum.SH7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH2, enum.SH1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH7, enum.SH8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH1, enum.SH2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH8, enum.SH7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH2, enum.SH1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH7, enum.SH8, 0, enum.MoveNormal),
		}, true},
		{[]movegen.Move{
			movegen.NewMove(enum.SH3, enum.SH2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH6, enum.SH7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH2, enum.SH1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH7, enum.SH8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH1, enum.SH2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH8, enum.SH7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH2, enum.SH1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH7, enum.SH8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH1, enum.SH2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH8, enum.SH7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SH2, enum.SH1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD6, enum.SD7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD3, enum.SD2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD7, enum.SD8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD2, enum.SD1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD8, enum.SD7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD1, enum.SD2, 0, enum.MoveNormal),
		}, false},
		{[]movegen.Move{
			movegen.NewMove(enum.SE4, enum.SE2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SE5, enum.SE7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SF3, enum.SG1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SC6, enum.SB8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SG1, enum.SF3, 0, enum.MoveNormal),
			movegen.NewMove(enum.SB8, enum.SC6, 0, enum.MoveNormal),
			movegen.NewMove(enum.SF3, enum.SG1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SC6, enum.SB8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SG1, enum.SF3, 0, enum.MoveNormal),
			movegen.NewMove(enum.SB8, enum.SC6, 0, enum.MoveNormal),
		}, true},
		{[]movegen.Move{
			movegen.NewMove(enum.SD4, enum.SD2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD6, enum.SD7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD5, enum.SD4, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA6, enum.SA7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA3, enum.SA2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SC5, enum.SC7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA2, enum.SA1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA7, enum.SA8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA1, enum.SA2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA8, enum.SA7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA2, enum.SA1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA7, enum.SA8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA1, enum.SA2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA8, enum.SA7, 0, enum.MoveNormal),
		}, false},
		{[]movegen.Move{
			movegen.NewMove(enum.SD4, enum.SD2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD6, enum.SD7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SD5, enum.SD4, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA6, enum.SA7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA3, enum.SA2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SC5, enum.SC7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA2, enum.SA1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA7, enum.SA8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA1, enum.SA2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA8, enum.SA7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA2, enum.SA1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA7, enum.SA8, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA1, enum.SA2, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA8, enum.SA7, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA2, enum.SA1, 0, enum.MoveNormal),
			movegen.NewMove(enum.SA7, enum.SA8, 0, enum.MoveNormal),
		}, true},
	}

	for _, tc := range testcases {
		game := NewGame()

		for i, move := range tc.moveStack {
			game.PushMove(move)
			game.LegalMoves = movegen.GenLegalMoves(game.Bitboards, game.ActiveColor, game.CastlingRights,
				game.EnPassantTarget)
			if i < len(tc.moveStack)-1 {
				positionKey := position{game.LegalMoves, game.Bitboards, game.ActiveColor, game.CastlingRights}.repetitionKey()
				game.Repetitions[positionKey]++
			}
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
		game.Bitboards = fen.ToBitboardArray(tc.fenString)

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
		bb, ac, cr, ep, _, _ := fen.Parse(tc.fenString)
		game.Bitboards = bb
		game.ActiveColor = ac
		game.CastlingRights = cr
		game.EnPassantTarget = ep
		game.LegalMoves = movegen.GenLegalMoves(bb, ac, cr, ep)

		got := game.IsCheckmate()
		if got != tc.expected {
			t.Fatalf("expected: %t, got: %t", tc.expected, got)
		}
	}
}

func BenchmarkPushMove(b *testing.B) {
	game := NewGame()
	bitboards := fen.ToBitboardArray("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")

	for b.Loop() {
		game.PushMove(movegen.NewMove(enum.SE4, enum.SE2, 0, enum.MoveNormal))
		// Restore the game state.
		game.Bitboards = bitboards
	}
}

// TODO: BenchmarkPopMove

func BenchmarkIsThreefoldRepetition(b *testing.B) {
	game := NewGame()
	moveStack := []movegen.Move{
		movegen.NewMove(enum.SH3, enum.SH2, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH6, enum.SH7, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH2, enum.SH1, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH7, enum.SH8, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH1, enum.SH2, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH8, enum.SH7, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH2, enum.SH1, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH7, enum.SH8, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH1, enum.SH2, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH8, enum.SH7, 0, enum.MoveNormal),
		movegen.NewMove(enum.SH2, enum.SH1, 0, enum.MoveNormal),
		movegen.NewMove(enum.SD6, enum.SD7, 0, enum.MoveNormal),
		movegen.NewMove(enum.SD3, enum.SD2, 0, enum.MoveNormal),
		movegen.NewMove(enum.SD7, enum.SD8, 0, enum.MoveNormal),
		movegen.NewMove(enum.SD2, enum.SD1, 0, enum.MoveNormal),
		movegen.NewMove(enum.SD8, enum.SD7, 0, enum.MoveNormal),
		movegen.NewMove(enum.SD1, enum.SD2, 0, enum.MoveNormal),
	}

	for i, move := range moveStack {
		game.PushMove(move)
		game.LegalMoves = movegen.GenLegalMoves(game.Bitboards, game.ActiveColor, game.CastlingRights,
			game.EnPassantTarget)
		if i < len(moveStack)-1 {
			positionKey := position{game.LegalMoves, game.Bitboards, game.ActiveColor, game.CastlingRights}.repetitionKey()
			game.Repetitions[positionKey]++
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
