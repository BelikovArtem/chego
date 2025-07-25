package bitutil

import "testing"

func TestBitScan(t *testing.T) {
	for i := 0; i < 64; i++ {
		var bitboard uint64 = 1 << i

		got := BitScan(bitboard)
		if got != i {
			t.Fatalf("Expected: %d got %d", i, got)
		}
	}
}

func TestPopLSB(t *testing.T) {
	for i := 0; i < 64; i++ {
		var bitboard uint64 = 1 << i

		got := PopLSB(&bitboard)
		if got != i {
			t.Fatalf("Expected %d got %d", i, got)
		}
	}

	var bitboard uint64 = 0
	got := PopLSB(&bitboard)
	if got != -1 {
		t.Fatalf("Expected 0 got %d", got)
	}
}

func TestCountBits(t *testing.T) {
	var got int

	got = CountBits(0x8000000000000000)
	if got != 1 {
		t.Fatalf("Expected 1 got %d", got)
	}

	got = CountBits(0x0)
	if got != 0 {
		t.Fatalf("Expected 1 got %d", got)
	}

	got = CountBits(0xFFFFFFFFFFFFFFFF)
	if got != 64 {
		t.Fatalf("Expected 64 got %d", got)
	}
}

func BenchmarkBitScan(b *testing.B) {
	for b.Loop() {
		BitScan(0x8000000000000000)
	}
}

func BenchmarkPopLSB(b *testing.B) {
	var bitboard uint64 = 0xFFFFFFFFFFFFFFFF

	for b.Loop() {
		PopLSB(&bitboard)
	}
}

func BenchmarkCountBits(b *testing.B) {
	for b.Loop() {
		CountBits(0xFFFFFFFFFFFFFFFF)
	}
}
