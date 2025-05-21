// Package fen implements conversions between the Forsyth-Edwards Notation strings and bitboard arrays.
// fen expects that the passed FEN and bitboard arrays are always valid and panics if they are not.
package fen

import (
	"chego/enum"
	// bits is used to speed up the iteration over bitboards.
	"math/bits"
	// strings is used to reduce the number of memory allocations during strings concatenation.
	"strings"
)

// ToBitboardArray converts the first part of a Forsyth-Edwards Notation string into
// an array of bitboards.
func ToBitboardArray(piecePlacementData string) [12]uint64 {
	var bitboards [12]uint64
	squareIndex := 56

	// Piece placement data describes each rank beginning from the eigth.
	for i := 0; i < len(piecePlacementData); i++ {
		char := piecePlacementData[i]

		if char == '/' { // Rank separator.
			squareIndex -= 16
		} else if char >= '1' && char <= '8' { // Number of consecutive empty squares.
			// Convert byte to the integer it represents.
			squareIndex += int(char - '0')
		} else { // There is piece on a square.
			var pieceType enum.Piece // enum.PieceWPawn by default.
			// Manual switch construction is ~3x faster than map approach.
			switch char {
			case 'N':
				pieceType = enum.PieceWKnight
			case 'B':
				pieceType = enum.PieceWBishop
			case 'R':
				pieceType = enum.PieceWRook
			case 'Q':
				pieceType = enum.PieceWQueen
			case 'K':
				pieceType = enum.PieceWKing
			case 'p':
				pieceType = enum.PieceBPawn
			case 'n':
				pieceType = enum.PieceBKnight
			case 'b':
				pieceType = enum.PieceBBishop
			case 'r':
				pieceType = enum.PieceBRook
			case 'q':
				pieceType = enum.PieceBQueen
			case 'k':
				pieceType = enum.PieceBKing
			}
			// Set the bit on the bitboard to place a piece.
			bitboards[pieceType] |= 1 << squareIndex
			squareIndex++
		}
	}

	return bitboards
}

// pieceSymbols is used in FromBitboardArray function.
var pieceSymbols = [12]byte{
	'P', 'N', 'B', 'R', 'Q', 'K',
	'p', 'n', 'b', 'r', 'q', 'k',
}

// FromBitboardArray converts the array of bitboards into the first part
// of Forsyth-Edwards Notation.
func FromBitboardArray(bitboards [12]uint64) string {
	// Used to add characters to a string without extra mem allocs.
	var piecePlacementData strings.Builder

	var board [8][8]byte

	for pieceType, bitboard := range bitboards {
		// Go through all pieces on a bitboard.
		for ; bitboard > 0; bitboard &= bitboard - 1 {
			squareIndex := bits.TrailingZeros64(bitboard)
			// Add piece on board.
			board[squareIndex/8][squareIndex%8] = pieceSymbols[pieceType]
		}
	}

	var numOfEmptySquares byte

	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			char := board[rank][file]

			if char == 0 { // Empty square.
				numOfEmptySquares++
			} else { // Piece on square.
				if numOfEmptySquares > 0 {
					piecePlacementData.WriteByte('0' + numOfEmptySquares)
					numOfEmptySquares = 0
				}
				piecePlacementData.WriteByte(char)
			}

			// To add rank separators.
			squareIndex := 8*rank + file
			if (squareIndex+1)%8 == 0 {
				if numOfEmptySquares > 0 {
					piecePlacementData.WriteByte('0' + numOfEmptySquares)
					numOfEmptySquares = 0
				}
				// Do not add separator in the end of the string.
				if squareIndex != 7 {
					piecePlacementData.WriteByte('/')
				}
			}
		}
	}

	return piecePlacementData.String()
}
