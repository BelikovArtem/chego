// Package cli provides functions to print a game state.
package cli

import (
	"chego/enum"
	"fmt"
)

// pieceSymbols is an array of chess piece runes.
var pieceSymbols = [12]rune{
	'♙', '♘', '♗', '♖', '♕', '♔',
	'♟', '♞', '♝', '♜', '♛', '♚',
}

// PrintBitboard prints the bitboard of the specified piece type as a chessboard.
func PrintBitboard(bitboard uint64, pieceType enum.Piece) {
	fmt.Printf("   a  b  c  d  e  f  g  h\n")

	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("%d  ", rank+1)

		for file := 0; file < 8; file++ {
			squareIndex := uint64(1 << (8*rank + file))

			symbol := pieceSymbols[pieceType]
			if bitboard&squareIndex == 0 {
				symbol = '.'
			}

			fmt.Printf("%c  ", symbol)
		}

		fmt.Printf("%d\n", rank+1)
	}

	fmt.Printf("   a  b  c  d  e  f  g  h\n")
}

// PrintBitboardsArray prints the specified array of bitboards as a chessboard.
func PrintBitboardsArray(bitboards [12]uint64) {
	fmt.Printf("   a  b  c  d  e  f  g  h\n")

	board := [8][8]rune{
		{'.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.'},
	}

	for pieceType, bitboard := range bitboards {
		for rank := 0; rank < 8; rank++ {
			for file := 0; file < 8; file++ {
				squareIndex := uint64(1 << (8*rank + file))

				if bitboard&squareIndex != 0 {
					board[rank][file] = pieceSymbols[pieceType]
				}
			}
		}
	}

	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("%d  ", rank+1)

		for file := 0; file < 8; file++ {
			fmt.Printf("%c  ", board[rank][file])
		}

		fmt.Printf("%d\n", rank+1)
	}

	fmt.Printf("   a  b  c  d  e  f  g  h\n")
}
