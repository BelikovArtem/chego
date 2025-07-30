// bitutil.go implements useful bit utilities which
// are used in move generation and game management logic.

package chego

// Precalculated magic used to form indices for the BitScanLookup array.
const bitscanMagic uint64 = 0x07EDD5E59A4E28C2

// Precalculated lookup table of LSB indices for 64-bit unsigned integers.
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

// bitScan returns the index of the LSB withing the bitboard.
// bitboard & -bitboard gives the LSB which is then run through
// the hashing scheme to index a lookup.
//
// NOTE: bitScan returns 63 for the empty bitboard.
func bitScan(bitboard uint64) int {
	return bitScanLookup[bitboard&-bitboard*bitscanMagic>>58]
}

// popLSB removes (pops) the LSB from the bitboard and returns its index.
//
// NOTE: popLSB returns 63 for the empty bitboard.
func popLSB(bitboard *uint64) int {
	lsb := bitScan(*bitboard)
	*bitboard &= *bitboard - 1
	return lsb
}

// CountBits returns the number of bits set within the bitboard.
func CountBits(bitboard uint64) int {
	cnt := 0
	for ; bitboard > 0; cnt++ {
		bitboard &= bitboard - 1
	}
	return cnt
}
