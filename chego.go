// Package chego (chess in go) implements chess logic.
package chego

import (
	"chego/enum"
	"math/bits"
)

// Move represents a chess move, encoded as a 16 bit unsigned integer:
//
//	0-5: To (destination) square index;
//	6-11: From (origin/source) square index;
//	12-13: Promotion piece (00 - knight, 01 - bishop, 10 - rook, 11 - queen);
//	14-15: Move type (see [enum.MoveType]).
type Move uint16

func NewMove(to, from, promotionPiece, moveType int) Move {
	return Move(to | (from << 6) | (promotionPiece << 12) | (moveType << 14))
}
func (m Move) To() int                    { return int(m & 0x3F) }
func (m Move) From() int                  { return int(m>>6) & 0x3F }
func (m Move) PromotionPiece() enum.Piece { return enum.Piece(m>>12) & 0x3 }
func (m Move) Type() enum.MoveType        { return enum.MoveType(m>>14) & 0x3 }

// The following block of constants defines the bitmasks needed to
// calculate possible moves by performing bitwise operations on a bitboard.
const (
	not_A_File   uint64 = 0xFEFEFEFEFEFEFEFE // All files except the A.
	not_H_File   uint64 = 0x7F7F7F7F7F7F7F7F // All files except the H.
	not_AB_File  uint64 = 0xFCFCFCFCFCFCFCFC // All files except the A and B.
	not_GH_File  uint64 = 0x3F3F3F3F3F3F3F3F // All files except the G and H.
	not_1st_Rank uint64 = 0xFFFFFFFFFFFFFF00 // All ranks except first.
	not_8th_Rank uint64 = 0x00FFFFFFFFFFFFFF // All ranks except eighth.
)

// PopLSB pops the Least Significant Bit from a bitboard.
func PopLSB(bitboard *uint64) int {
	lsb := bits.TrailingZeros64(*bitboard)
	*bitboard &= *bitboard - 1
	return lsb
}

// Precalculated attack tables used to speed up the move generation process.
var (
	// Pawn's attack pattern depends on the color, so it is necessary to store two tables.
	PawnAttacks   [2][64]uint64
	KnightAttacks [64]uint64
	KingAttacks   [64]uint64
)

// GenPawnAttacks returns a bitboard of squares attacked by a pawn.
func GenPawnAttacks(pawn uint64, color enum.Color) uint64 {
	if color == enum.ColorWhite {
		return (pawn & not_A_File << 7) | (pawn & not_H_File << 9)
	}
	// Handle black pawns.
	return (pawn & not_A_File >> 9) | (pawn & not_H_File >> 7)
}

// GenKnightAttacks returns a bitboard of squares attacked by a knight.
func GenKnightAttacks(knight uint64) uint64 {
	return (knight & not_A_File >> 17) |
		(knight & not_H_File >> 15) |
		(knight & not_AB_File >> 10) |
		(knight & not_GH_File >> 6) |
		(knight & not_AB_File << 6) |
		(knight & not_GH_File << 10) |
		(knight & not_A_File << 15) |
		(knight & not_H_File << 17)
}

// GenKingAttacks returns a bitboard of squares attacked by a king.
func GenKingAttacks(king uint64) uint64 {
	return (king & not_A_File >> 9) |
		(king >> 8) |
		(king & not_H_File >> 7) |
		(king & not_A_File >> 1) |
		(king & not_H_File << 1) |
		(king & not_A_File << 7) |
		(king << 8) |
		(king & not_H_File << 9)
}

// GenBishopAttacks returns a bitboard of squares attacked by a bishop.
// Occupied squares that block movement in each direction are taken into account.
// The resulting bitboard also includes occupied squares.
func GenBishopAttacks(bishop uint64, occupancy uint64) uint64 {
	var attacks uint64

	for i := bishop & not_A_File >> 9; i&not_H_File != 0; i >>= 9 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := bishop & not_H_File >> 7; i&not_A_File != 0; i >>= 7 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := bishop & not_A_File << 7; i&not_H_File != 0; i <<= 7 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := bishop & not_H_File << 9; i&not_A_File != 0; i <<= 9 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	return attacks
}

// GenRookAttacks returns a bitboard of squares attacked by a rook.
// Occupied squares that block movement in each direction are taken into account.
// The resulting bitboard also includes occupied squares.
func GenRookAttacks(rook uint64, occupancy uint64) uint64 {
	var attacks uint64

	for i := rook & not_A_File >> 1; i&not_H_File != 0; i >>= 1 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := rook & not_H_File << 1; i&not_A_File != 0; i <<= 1 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := rook & not_1st_Rank >> 8; i&not_8th_Rank != 0; i >>= 8 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := rook & not_8th_Rank << 8; i&not_1st_Rank != 0; i <<= 8 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	return attacks
}

// GenBishopRelevantOccupancy returns a bitboard of "relevant occupancy squares".
// They are the only squares whose occupancy matters when generating the legal moves of
// a bishop. This function is used to generate magic bitboards.
func GenBishopRelevantOccupancy(bishop uint64) uint64 {
	var occupancy uint64

	not_A_not_1st := not_A_File & not_1st_Rank
	not_H_not_1st := not_H_File & not_1st_Rank
	not_A_not_8th := not_A_File & not_8th_Rank
	not_H_not_8th := not_H_File & not_8th_Rank

	for i := bishop & not_A_File >> 9; i&not_A_not_1st != 0; i >>= 9 {
		occupancy |= i
	}

	for i := bishop & not_H_File >> 7; i&not_H_not_1st != 0; i >>= 7 {
		occupancy |= i
	}

	for i := bishop & not_A_File << 7; i&not_A_not_8th != 0; i <<= 7 {
		occupancy |= i
	}

	for i := bishop & not_H_File << 9; i&not_H_not_8th != 0; i <<= 9 {
		occupancy |= i
	}

	return occupancy
}

// GenRookRelevantOccupancy returns a bitboard of "relevant occupancy squares".
// They are the only squares whose occupancy matters when generating the legal moves of
// a rook. This function is used to generate magic bitboards.
func GenRookRelevantOccupancy(rook uint64) uint64 {
	var occupancy uint64

	for i := rook & not_1st_Rank >> 8; i&not_1st_Rank != 0; i >>= 8 {
		occupancy |= i
	}

	for i := rook & not_A_File >> 1; i&not_A_File != 0; i >>= 1 {
		occupancy |= i
	}

	for i := rook & not_H_File << 1; i&not_H_File != 0; i <<= 1 {
		occupancy |= i
	}

	for i := rook & not_8th_Rank << 8; i&not_8th_Rank != 0; i <<= 8 {
		occupancy |= i
	}

	return occupancy
}

// InitAttackTables initializes the predefined attack tables.
func InitAttackTables() {
	for square := 0; square < 64; square++ {
		PawnAttacks[enum.ColorWhite][square] = GenPawnAttacks(1<<square, enum.ColorWhite)
		PawnAttacks[enum.ColorBlack][square] = GenPawnAttacks(1<<square, enum.ColorBlack)

		KnightAttacks[square] = GenKnightAttacks(1 << square)

		KingAttacks[square] = GenKingAttacks(1 << square)
	}
}
