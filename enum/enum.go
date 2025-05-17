// Package enum contains predefined constants of custom types.
// Used to avoid the "magic numbers" antipattern.
package enum

type Piece int

const (
	PieceWPawn Piece = iota
	PieceWKnight
	PieceWBishop
	PieceWRook
	PieceWQueen
	PieceWKing
	PieceBPawn
	PieceBKnight
	PieceBBishop
	PieceBRook
	PieceBQueen
	PieceBKing
)
