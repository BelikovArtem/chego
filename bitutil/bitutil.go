// Package bitutil implements helpful bit utilities used in move generation and game management logic.
package bitutil

// Precalculated magic used to form indices for the BitScanLookup array.
const BITSCAN_MAGIC uint64 = 0x07EDD5E59A4E28C2

// Precalculated lookup table of LSB indices for 64 uints.
// See http://pradu.us/old/Nov27_2008/Buzz/research/magic/Bitboards.pdf section 3.2.
var bitScanLookup = [64]int{
	63, 0, 58, 1, 59, 47, 53, 2,
	60, 39, 48, 27, 54, 33, 42, 3,
	61, 51, 37, 40, 49, 18, 28, 20,
	55, 30, 34, 11, 43, 14, 22, 4,
	62, 57, 46, 52, 38, 26, 32, 41,
	50, 36, 17, 19, 29, 10, 13, 21,
	56, 45, 25, 31, 35, 16, 9, 12,
	44, 24, 15, 8, 23, 7, 6, 5,
}

// BitScan returns the index of the Least Significant Bit (LSB) withing the bitboard.
// bitboard&-bitboard gives the LSB which is then run through the hashing scheme to index a lookup.
func BitScan(bitboard uint64) int { return bitScanLookup[bitboard&-bitboard*BITSCAN_MAGIC>>58] }

// PopLSB removes (pops) the least significant bit from the bitboard and returns its index.
// If the bitboard is empty, it returns -1.
func PopLSB(bitboard *uint64) int {
	if *bitboard == 0 {
		return -1
	}

	lsb := BitScan(*bitboard)
	*bitboard &= *bitboard - 1
	return lsb
}

// CountBits returns the number of bits set within the bitboard.
func CountBits(bitboard uint64) int {
	var cnt int
	for bitboard > 0 {
		cnt++
		bitboard &= bitboard - 1
	}
	return cnt
}
