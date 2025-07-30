package chego

import "testing"

func TestBitScan(t *testing.T) {
	for i := 0; i < 64; i++ {
		bb := uint64(1 << i)

		got := bitScan(bb)
		if got != i {
			t.Fatalf("Expected: %d got %d", i, got)
		}
	}
}

func TestPopLSB(t *testing.T) {
	for i := 0; i < 64; i++ {
		bb := uint64(1 << i)

		got := popLSB(&bb)
		if got != i {
			t.Fatalf("Expected %d got %d", i, got)
		}
	}
}

func TestCountBits(t *testing.T) {
	bb := uint64(0)

	for i := 0; i < 64; i++ {
		bb |= uint64(1 << i)

		got := CountBits(bb)
		if got != i+1 {
			t.Fatalf("Expected: %d got %d", i+1, got)
		}
	}
}

func BenchmarkBitScan(b *testing.B) {
	for b.Loop() {
		bitScan(0x8000000000000000)
	}
}

func BenchmarkPopLSB(b *testing.B) {
	var bitboard uint64 = 0xFFFFFFFFFFFFFFFF

	for b.Loop() {
		popLSB(&bitboard)
	}
}

func BenchmarkCountBits(b *testing.B) {
	for b.Loop() {
		CountBits(0xFFFFFFFFFFFFFFFF)
	}
}
