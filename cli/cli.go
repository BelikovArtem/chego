// Package cli provides functions to print a game state.
package cli

import (
	"chego/enum"
	"strings"
)

// pieceSymbols is an array of chess piece runes.
var pieceSymbols = [12]rune{
	'♙', '♘', '♗', '♖', '♕', '♔',
	'♟', '♞', '♝', '♜', '♛', '♚',
}

func FormatBitboard(bitboard uint64, pieceType enum.Piece) string {
	var bitboardStr strings.Builder

	for rank := 7; rank >= 0; rank-- {
		bitboardStr.WriteByte(byte(rank) + 1 + '0')
		bitboardStr.WriteString("  ")

		for file := 0; file < 8; file++ {
			squareIndex := uint64(1 << (8*rank + file))

			symbol := pieceSymbols[pieceType]
			if bitboard&squareIndex == 0 {
				symbol = '.'
			}

			bitboardStr.WriteRune(symbol)
			bitboardStr.WriteString("  ")
		}
		bitboardStr.WriteString("\n")
	}
	bitboardStr.WriteString("   a  b  c  d  e  f  g  h\n")

	return bitboardStr.String()
}
