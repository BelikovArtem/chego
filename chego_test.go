package chego

import (
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
			t.Fatalf("test %s failed: expected %x, got %x", tc.name, tc.expected, got)
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
			t.Fatalf("test %s failed: expected %x, got %x", tc.name, tc.expected, got)
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
			t.Fatalf("test %s failed: expected %x, got %x", tc.name, tc.expected, got)
		}
	}
}

func BenchmarkGenPawnAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenPawnAttacks(enum.B4, enum.ColorWhite)
	}
}

func BenchmarkGenKnightsAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenKnightAttacks(enum.B4)
	}
}

func BenchmarkInitAttackTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitAttackTables()
	}
}
