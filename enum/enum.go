// Package enum contains custom type declarations and predefined constants.
// Used to avoid the "magic numbers" antipattern.
package enum

// Piece is allias type to avoid bothersome conversion between int and Piece.
type Piece = int

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

// Color is allias type to avoid bothersome conversion between int and Color.
type Color = int

const (
	ColorWhite Color = iota
	ColorBlack
)

type MoveType int

const (
	// Quite & capture moves.
	MoveNormal MoveType = iota
	// King & queen castle.
	MoveCastling
	// Knight & Bishop & Rook & Queen promotions.
	MovePromotion
	// Special pawn move.
	MoveEnPassant
)

// CastlingFlag defines the player's rights to perform castlings:
//
// 	0 bit: white king can O-O.
//  1 bit: white king can O-O-O.
//  2 bit: black king can O-O.
//  3 bit: black king can O-O-O.
type CastlingFlag int

const (
	CastlingWhiteKing  CastlingFlag = 1
	CastlingWhiteQueen CastlingFlag = 2
	CastlingBlackKing  CastlingFlag = 4
	CastlingBlackQueen CastlingFlag = 8
)

// Bitboards of each square. Used to simplify tests.
const (
	// To distinguish the absence of the en passant target.
	NoSquare        = -1
	A1       uint64 = 1 << (iota - 1)
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)
