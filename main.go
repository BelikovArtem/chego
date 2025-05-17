// Chego implements chess logic.
package main

import (
	"fmt"
)

type Piece int

const (
	WPawn Piece = iota
	WKnight
	WBishop
	WRook
	WQueen
	WKing
	BPawn
	BKnight
	BBishop
	BRook
	BQueen
	BKing
)

var (
	pieceSymbols = [12]rune{
		'♙', '♘', '♗', '♖', '♕', '♔',
		'♟', '♞', '♝', '♜', '♛', '♚',
	}
)

func main() {
	printBitboard(0xFF, WPawn)
}

// printBitboard prints the specified bitboard of the specified piece type.
func printBitboard(bitboard uint64, piece Piece) {
	fmt.Printf("   a  b  c  d  e  f  g  h\n")

	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("%d  ", rank+1)

		for file := 0; file < 8; file++ {
			squareIndex := uint64(1 << (8*rank + file))

			symbol := pieceSymbols[piece]
			if bitboard&squareIndex == 0 {
				symbol = '.'
			}

			fmt.Printf("%c  ", symbol)
		}

		fmt.Printf("%d\n", rank+1)
	}

	fmt.Printf("   a  b  c  d  e  f  g  h\n")
}
