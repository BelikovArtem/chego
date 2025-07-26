// Package format provides functions to format a chess boards and positions.
// It is used mainly to visualize test cases.
package format

import (
	"strings"

	"github.com/BelikovArtem/chego/types"
)

var pieceSymbols = [12]rune{
	'♙', '♘', '♗', '♖', '♕', '♔',
	'♟', '♞', '♝', '♜', '♛', '♚',
}

var squareString = [64]string{
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
}

// Bitboard formats a single bitboard into a string.
func Bitboard(bitboard uint64, pieceType types.Piece) string {
	var bitboardStr strings.Builder

	for rank := 7; rank >= 0; rank-- {
		bitboardStr.WriteByte(byte(rank) + 1 + '0')
		bitboardStr.WriteString("  ")

		for file := 0; file < 8; file++ {
			square := uint64(1 << (8*rank + file))

			symbol := pieceSymbols[pieceType]
			if bitboard&square == 0 {
				symbol = '.'
			}

			bitboardStr.WriteRune(symbol)
			bitboardStr.WriteString("  ")
		}
		bitboardStr.WriteByte('\n')
	}
	bitboardStr.WriteString("   a  b  c  d  e  f  g  h\n")

	return bitboardStr.String()
}

// Position formats a full chess position into a string.
func Position(p types.Position) string {
	var positionStr strings.Builder

	for rank := 7; rank >= 0; rank-- {
		positionStr.WriteByte(byte(rank) + 1 + '0')
		positionStr.WriteString("  ")

		for file := 0; file < 8; file++ {
			square := uint64(1 << (8*rank + file))

			var symbol rune = '.'

			for i := types.PieceWPawn; i <= types.PieceBKing; i++ {
				if square&p.Bitboards[i] != 0 {
					symbol = pieceSymbols[i]
					break
				}
			}

			positionStr.WriteRune(symbol)
			positionStr.WriteString("  ")
		}
		positionStr.WriteByte('\n')
	}

	positionStr.WriteString("   a  b  c  d  e  f  g  h\nActive color: ")

	if p.ActiveColor == types.ColorWhite {
		positionStr.WriteString("white\nEn passant: ")
	} else {
		positionStr.WriteString("black\nEn passant: ")
	}

	if p.EPTarget == 0 {
		positionStr.WriteString("none\nCastling rights: ")
	} else {
		positionStr.WriteString(squareString[p.EPTarget])
		positionStr.WriteString("\nCastling rights: ")
	}

	if p.CastlingRights&types.CastlingWhiteShort != 0 {
		positionStr.WriteByte('K')
	}
	if p.CastlingRights&types.CastlingWhiteLong != 0 {
		positionStr.WriteByte('Q')
	}
	if p.CastlingRights&types.CastlingBlackShort != 0 {
		positionStr.WriteByte('k')
	}
	if p.CastlingRights&types.CastlingBlackLong != 0 {
		positionStr.WriteByte('q')
	}

	return positionStr.String()
}
