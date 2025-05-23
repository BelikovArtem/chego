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
	not_A_File  uint64 = 0xFEFEFEFEFEFEFEFE // All files except the A.
	not_H_File  uint64 = 0x7F7F7F7F7F7F7F7F // All files except the H.
	not_AB_File uint64 = 0xFCFCFCFCFCFCFCFC // All files except the A and B.
	not_GH_File uint64 = 0x3F3F3F3F3F3F3F3F // All files except the G and H.
)

// PopLSB pops the Least Significant Bit from a bitboard.
func PopLSB(bitboard *uint64) int {
	lsb := bits.TrailingZeros64(*bitboard)
	*bitboard &= *bitboard - 1
	return lsb
}

// Precalculated attack tables used to speed up the move generation process.
var (
	// Pawn's attack pattern depends on the color, so it is necessary to have two tables.
	pawnAttacks   [2][64]uint64
	knightAttacks [64]uint64
	kingAttacks   [64]uint64
)

// GenPawnAttacks returns a bitboard of attacked by a pawn squares.
func GenPawnAttacks(pawn uint64, color enum.Color) uint64 {
	if color == enum.ColorWhite {
		return (pawn & not_A_File << 7) | (pawn & not_H_File << 9)
	}
	// Handle black pawns.
	return (pawn & not_A_File >> 9) | (pawn & not_H_File >> 7)
}

// GenKnightAttacks returns a bitboard of attacked by a knight squares.
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

// GenKingAttacks returns a bitboard of attacked by a king squares.
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

func InitAttackTables() {
	for square := 0; square < 64; square++ {
		pawnAttacks[enum.ColorWhite][square] = GenPawnAttacks(1<<square, enum.ColorWhite)
		pawnAttacks[enum.ColorBlack][square] = GenPawnAttacks(1<<square, enum.ColorBlack)

		knightAttacks[square] = GenKnightAttacks(1 << square)

		kingAttacks[square] = GenKingAttacks(1 << square)
	}
}
