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

// PromotionFlag is allias type to avoid bothersome conversion between int and Color.
type PromotionFlag = int

// 00 - knight, 01 - bishop, 10 - rook, 11 - queen
const (
	PromotionKnight PromotionFlag = iota
	PromotionBishop
	PromotionRook
	PromotionQueen
)

// Color is allias type to avoid bothersome conversion between int and Color.
type Color = int

const (
	ColorWhite Color = iota
	ColorBlack
)

// MoveType is allias type to avoid bothersome conversion between int and MoveType.
type MoveType = int

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

// Each square.
const (
	SA1 int = iota
	SB1
	SC1
	SD1
	SE1
	SF1
	SG1
	SH1
	SA2
	SB2
	SC2
	SD2
	SE2
	SF2
	SG2
	SH2
	SA3
	SB3
	SC3
	SD3
	SE3
	SF3
	SG3
	SH3
	SA4
	SB4
	SC4
	SD4
	SE4
	SF4
	SG4
	SH4
	SA5
	SB5
	SC5
	SD5
	SE5
	SF5
	SG5
	SH5
	SA6
	SB6
	SC6
	SD6
	SE6
	SF6
	SG6
	SH6
	SA7
	SB7
	SC7
	SD7
	SE7
	SF7
	SG7
	SH7
	SA8
	SB8
	SC8
	SD8
	SE8
	SF8
	SG8
	SH8
)
