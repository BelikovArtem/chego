package chego

import (
	"testing"
)

func BenchmarkGenPawnAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genPawnAttacks(B4, ColorWhite)
	}
}

func BenchmarkGenKnightAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genKnightAttacks(B4)
	}
}

func BenchmarkGenKingAttakcs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genKingAttacks(B4)
	}
}

func BenchmarkGenBishopAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genBishopAttacks(D5, B3)
	}
}

func BenchmarkGenRookAttacks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genRookAttacks(D5, B3)
	}
}

func BenchmarkInitBishopOccupancy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initBishopOccupancy()
	}
}

func BenchmarkInitRookOccupancy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initRookOccupancy()
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

// func BenchmarkGenPawnMoves(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		genPawnMoves(SE4, 0x0, 0x0, 0, ColorWhite, &MoveList{})
// 	}
// }

func BenchmarkGenKingMoves(b *testing.B) {
	pos := ParseFEN("8/8/8/8/8/8/8/R3K2R w - - 0 1")

	for b.Loop() {
		genKingMoves(&pos, &MoveList{})
	}
}

func BenchmarkGenLegalMoves(b *testing.B) {
	pos := ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")

	for b.Loop() {
		lm := MoveList{}
		GenLegalMoves(pos, &lm)
	}
}

func BenchmarkInitAttackTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitAttackTables()
	}
}
