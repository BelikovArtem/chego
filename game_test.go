package chego

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup.
	InitAttackTables()
	InitZobristKeys()
	// Tests and benchmarks execution.
	os.Exit(m.Run())
}

func TestIsThreefoldRepetition(t *testing.T) {
	testcases := []struct {
		moveStack []Move
		expected  bool
	}{
		{[]Move{
			NewMove(SH3, SH2, MoveNormal),
			NewMove(SH6, SH7, MoveNormal),
			NewMove(SH2, SH1, MoveNormal),
			NewMove(SH7, SH8, MoveNormal),
			NewMove(SH1, SH2, MoveNormal),
			NewMove(SH8, SH7, MoveNormal),
			NewMove(SH2, SH1, MoveNormal),
			NewMove(SH7, SH8, MoveNormal),
			NewMove(SH1, SH2, MoveNormal),
			NewMove(SH8, SH7, MoveNormal),
			NewMove(SH2, SH1, MoveNormal),
			NewMove(SH7, SH8, MoveNormal),
		}, true},
		{[]Move{
			NewMove(SH3, SH2, MoveNormal),
			NewMove(SH6, SH7, MoveNormal),
			NewMove(SH2, SH1, MoveNormal),
			NewMove(SH7, SH8, MoveNormal),
			NewMove(SH1, SH2, MoveNormal),
			NewMove(SH8, SH7, MoveNormal),
			NewMove(SH2, SH1, MoveNormal),
			NewMove(SH7, SH8, MoveNormal),
			NewMove(SH1, SH2, MoveNormal),
			NewMove(SH8, SH7, MoveNormal),
			NewMove(SH2, SH1, MoveNormal),
			NewMove(SD6, SD7, MoveNormal),
			NewMove(SD3, SD2, MoveNormal),
			NewMove(SD7, SD8, MoveNormal),
			NewMove(SD2, SD1, MoveNormal),
			NewMove(SD8, SD7, MoveNormal),
			NewMove(SD1, SD2, MoveNormal),
		}, false},
		{[]Move{
			NewMove(SE4, SE2, MoveNormal),
			NewMove(SE5, SE7, MoveNormal),
			NewMove(SF3, SG1, MoveNormal),
			NewMove(SC6, SB8, MoveNormal),
			NewMove(SG1, SF3, MoveNormal),
			NewMove(SB8, SC6, MoveNormal),
			NewMove(SF3, SG1, MoveNormal),
			NewMove(SC6, SB8, MoveNormal),
			NewMove(SG1, SF3, MoveNormal),
			NewMove(SB8, SC6, MoveNormal),
		}, true},
		{[]Move{
			NewMove(SD4, SD2, MoveNormal),
			NewMove(SD6, SD7, MoveNormal),
			NewMove(SD5, SD4, MoveNormal),
			NewMove(SA6, SA7, MoveNormal),
			NewMove(SA3, SA2, MoveNormal),
			NewMove(SC5, SC7, MoveNormal),
			NewMove(SA2, SA1, MoveNormal),
			NewMove(SA7, SA8, MoveNormal),
			NewMove(SA1, SA2, MoveNormal),
			NewMove(SA8, SA7, MoveNormal),
			NewMove(SA2, SA1, MoveNormal),
			NewMove(SA7, SA8, MoveNormal),
			NewMove(SA1, SA2, MoveNormal),
			NewMove(SA8, SA7, MoveNormal),
		}, false},
		{[]Move{
			NewMove(SD4, SD2, MoveNormal),
			NewMove(SD6, SD7, MoveNormal),
			NewMove(SD5, SD4, MoveNormal),
			NewMove(SA6, SA7, MoveNormal),
			NewMove(SA3, SA2, MoveNormal),
			NewMove(SC5, SC7, MoveNormal),
			NewMove(SA2, SA1, MoveNormal),
			NewMove(SA7, SA8, MoveNormal),
			NewMove(SA1, SA2, MoveNormal),
			NewMove(SA8, SA7, MoveNormal),
			NewMove(SA2, SA1, MoveNormal),
			NewMove(SA7, SA8, MoveNormal),
			NewMove(SA1, SA2, MoveNormal),
			NewMove(SA8, SA7, MoveNormal),
			NewMove(SA2, SA1, MoveNormal),
			NewMove(SA7, SA8, MoveNormal),
		}, true},
	}

	for i, tc := range testcases {
		g := NewGame()

		for _, move := range tc.moveStack {
			g.PushMove(move)
		}

		got := g.IsThreefoldRepetition()
		if tc.expected != got {
			t.Fatalf("case %d failed: expected %t, got %t", i, tc.expected, got)
		}
	}
}

func TestIsInsufficientMaterial(t *testing.T) {
	testcases := []struct {
		fen      string
		expected bool
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
		game.Position.Bitboards = ParseBitboards(tc.fen)

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
		game.Position = ParseFEN(tc.fenString)
		GenLegalMoves(game.Position, &game.LegalMoves)

		got := game.IsCheckmate()
		if got != tc.expected {
			t.Fatalf("expected: %t, got: %t", tc.expected, got)
		}
	}
}

func BenchmarkPushMove(b *testing.B) {
	game := NewGame()
	pos := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	for b.Loop() {
		game.Position = pos
		game.PushMove(NewMove(SE4, SE2, MoveNormal))
	}
}

func BenchmarkIsThreefoldRepetition(b *testing.B) {
	game := NewGame()
	moveStack := []Move{
		NewMove(SH3, SH2, MoveNormal),
		NewMove(SH6, SH7, MoveNormal),
		NewMove(SH2, SH1, MoveNormal),
		NewMove(SH7, SH8, MoveNormal),
		NewMove(SH1, SH2, MoveNormal),
		NewMove(SH8, SH7, MoveNormal),
		NewMove(SH2, SH1, MoveNormal),
		NewMove(SH7, SH8, MoveNormal),
		NewMove(SH1, SH2, MoveNormal),
		NewMove(SH8, SH7, MoveNormal),
		NewMove(SH2, SH1, MoveNormal),
		NewMove(SD6, SD7, MoveNormal),
		NewMove(SD3, SD2, MoveNormal),
		NewMove(SD7, SD8, MoveNormal),
		NewMove(SD2, SD1, MoveNormal),
		NewMove(SD8, SD7, MoveNormal),
		NewMove(SD1, SD2, MoveNormal),
	}

	for i, move := range moveStack {
		game.PushMove(move)
		GenLegalMoves(game.Position, &game.LegalMoves)
		if i < len(moveStack)-1 {
			game.Repetitions[zobristKey(game.Position)]++
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
