// types.go contains declarations of custom types and predefined constants.

package chego

// Move represents a chess move, encoded as a 16 bit unsigned integer:
//
//	0-5:   To (destination) square index;
//	6-11:  From (origin/source) square index;
//	12-13: Promotion piece (see [PromotionFlag]);
//	14-15: Move type (see [MoveType]).
type Move uint16

// NewMove creates a new move with the promotion piece set to [PromotionQueen].
func NewMove(to, from, moveType int) Move {
	return Move(to | (from << 6) | (PromotionQueen << 12) | (moveType << 14))
}

// NewPromotionMove creates a new move with the promotion type and specified
// promotion piece.
func NewPromotionMove(to, from, promoPiece int) Move {
	return Move(to | (from << 6) | (promoPiece << 12) | (MovePromotion << 14))
}

func (m Move) To() int                   { return int(m & 0x3F) }
func (m Move) From() int                 { return int(m>>6) & 0x3F }
func (m Move) PromoPiece() PromotionFlag { return PromotionFlag(m>>12) & 0x3 }
func (m Move) Type() MoveType            { return MoveType(m>>14) & 0x3 }

// MoveList is used to store moves. The main idea behind it is to preallocate
// an array with enough capacity to store all possible moves and avoid dynamic
// memory allocations.
type MoveList struct {
	// Maximum number of moves per chess position is equal to 218,
	// hence 218 elements.
	// See https://www.talkchess.com/forum/viewtopic.php?t=61792
	Moves [218]Move
	// To keep track of the next move index.
	LastMoveIndex byte
}

// Push adds the move to the end of the move list.
func (l *MoveList) Push(m Move) {
	l.Moves[l.LastMoveIndex] = m
	l.LastMoveIndex++
}

var (
	// PieceSymbols is used in fen, format, and uci packages.
	PieceSymbols = [12]byte{
		'P', 'N', 'B', 'R', 'Q', 'K',
		'p', 'n', 'b', 'r', 'q', 'k',
	}
	// Square2String is used in format and uci packages.
	Square2String = [64]string{
		"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
		"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
		"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
		"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
		"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
		"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
		"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
		"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	}
)

// Piece is an allias type to avoid bothersome conversion between
// int and Piece.
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
	// To avoid magic numbers.
	PieceNone = -1
)

// PromotionFlag is an allias type to avoid bothersome conversion between
// int and Color.
type PromotionFlag = int

// 00 - knight, 01 - bishop, 10 - rook, 11 - queen.
const (
	PromotionKnight PromotionFlag = iota
	PromotionBishop
	PromotionRook
	PromotionQueen
)

// Color is an allias type to avoid bothersome conversion between int and Color.
type Color = int

const (
	ColorWhite Color = iota
	ColorBlack
	ColorBoth
)

// MoveType is an allias type to avoid bothersome conversion between
// int and MoveType.
type MoveType = int

const (
	// Quite & capture moves.
	MoveNormal MoveType = iota
	// King & queen castling.
	MoveCastling
	// Knight & Bishop & Rook & Queen promotions.
	MovePromotion
	// Special pawn move.
	MoveEnPassant
)

// CastlingRights defines the player's rights to perform castlings.
//
// 	0 bit: white king can O-O.
//  1 bit: white king can O-O-O.
//  2 bit: black king can O-O.
//  3 bit: black king can O-O-O.
type CastlingRights = int

const (
	CastlingWhiteShort CastlingRights = 1
	CastlingWhiteLong  CastlingRights = 2
	CastlingBlackShort CastlingRights = 4
	CastlingBlackLong  CastlingRights = 8
)

// Result represents the possible outcomes of a chess game.
type Result int

const (
	ResultUnscored Result = iota // Default value: the game isn't finished yet.
	ResultCheckmate
	ResultTimeout
	ResultStalemate
	ResultInsufficientMaterial
	ResultFiftyMove
	ResultThreefoldRepetition
	ResultResignation
	ResultDrawByAgreement
)

// Bitboards of each square. Used to simplify tests.
const (
	A1 uint64 = 1 << iota
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
	ALL_SQUARES = 0xFFFFFFFFFFFFFFFF
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
