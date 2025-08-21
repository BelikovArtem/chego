// types.go contains declarations of custom types and predefined constants.

package chego

/*
Move represents a chess move, encoded as a 16 bit unsigned integer:
  - 0-5:   To (destination) square index.
  - 6-11:  From (origin/source) square index.
  - 12-13: Promotion piece (see [PromotionFlag]).
  - 14-15: Move type (see [MoveType]).
*/
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

/*
MoveList is used to store moves.  The main idea behind it is to preallocate
an array with enough capacity to store all possible moves and avoid dynamic
memory allocations.
*/
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
	// PieceSymbols maps each piece type to its symbol.
	PieceSymbols = [12]byte{
		'P', 'p', 'N', 'n', 'B', 'b',
		'R', 'r', 'Q', 'q', 'K', 'k',
	}
	// Square2String maps each board square to its string representation.
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
	PieceBPawn
	PieceWKnight
	PieceBKnight
	PieceWBishop
	PieceBBishop
	PieceWRook
	PieceBRook
	PieceWQueen
	PieceBQueen
	PieceWKing
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

/*
CastlingRights defines the player's rights to perform castlings.
  - 0 bit: white king can O-O.
  - 1 bit: white king can O-O-O.
  - 2 bit: black king can O-O.
  - 3 bit: black king can O-O-O.
*/
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
