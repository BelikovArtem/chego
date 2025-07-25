// Package types contains declarations of custom types and predefined constants.
package types

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

// NewPromotionMove creates a new move with the promotion type and specified promotion piece.
func NewPromotionMove(to, from, promotionPiece int) Move {
	return Move(to | (from << 6) | (promotionPiece << 12) | (MovePromotion << 14))
}

func (m Move) To() int                       { return int(m & 0x3F) }
func (m Move) From() int                     { return int(m>>6) & 0x3F }
func (m Move) PromotionPiece() PromotionFlag { return PromotionFlag(m>>12) & 0x3 }
func (m Move) Type() MoveType                { return MoveType(m>>14) & 0x3 }

// Position represents a chessboard state that can be converted to or parsed from a FEN string.
type Position struct {
	// TODO: include all white pieces, all black pieces and occupancy.
	Bitboards      [12]uint64
	ActiveColor    Color
	CastlingRights CastlingRights
	EPTarget       int
	HalfmoveCnt    int
	FullmoveCnt    int
}

// MakeMove modifies the position by applying the specified move.
// It’s the caller’s responsibility to ensure the move is legal.
//
// Not only is the piece placement updated, but also the entire position,
// including castling rights, en passant target, move counters, and active color.
func (p *Position) MakeMove(m Move) {
	var from, to uint64 = 1 << m.From(), 1 << m.To()
	fromTo := from ^ to
	movedPiece := p.GetPieceFromSquare(from)

	switch m.Type() {
	case MoveNormal:
		// If the move is capture.
		capturedPiece := p.GetPieceFromSquare(to)
		if capturedPiece != PieceNone {
			// Remove the captured piece from the board.
			p.Bitboards[capturedPiece] ^= to
			// Reset the halfmove counter after captures.
			p.HalfmoveCnt = 0
		} else {
			p.HalfmoveCnt++
		}
		p.Bitboards[movedPiece] ^= fromTo

	case MoveEnPassant:
		// Remove the captured pawn from the board.
		if movedPiece == PieceWPawn {
			p.Bitboards[PieceBPawn] ^= to >> 8
		} else {
			p.Bitboards[PieceWPawn] ^= to << 8
		}
		p.Bitboards[movedPiece] ^= fromTo

	case MoveCastling:
		switch to {
		case G1, G8: // O-O
			p.Bitboards[movedPiece-2] ^= (to << 1) ^ (to >> 1)
		case C1, C8: // O-O-O
			p.Bitboards[movedPiece-2] ^= (to >> 2) ^ (to << 1)
		}
		p.Bitboards[movedPiece] ^= fromTo

	case MovePromotion:
		// If the move is capture-promotion.
		capturedPieceType := p.GetPieceFromSquare(to)
		if capturedPieceType != PieceNone {
			// Remove the captured piece from the board.
			p.Bitboards[capturedPieceType] ^= to
		}

		// Remove a promoted pawn from the board.
		p.Bitboards[movedPiece] ^= from
		// Place a new piece.
		if movedPiece == PieceWPawn {
			p.Bitboards[m.PromotionPiece()+1] ^= to
		} else {
			p.Bitboards[m.PromotionPiece()+7] ^= to
		}
	}

	// Reset the en passant target since the en passant capture is possible only for 1 move.
	p.EPTarget = 0

	switch movedPiece {
	// Set en passant targets in case of pawn double pushes.
	case PieceWPawn, PieceBPawn:
		if m.From()-m.To() == -16 {
			p.EPTarget = m.To() - 8
		} else if m.From()-m.To() == 16 {
			p.EPTarget = m.To() + 8
		}
		// Reset the halfmove counter after pawn moves.
		p.HalfmoveCnt = 0

	// Disable white castling rigts.
	case PieceWKing:
		p.CastlingRights &= ^(CastlingWhiteShort | CastlingWhiteLong)

	// Disable black castling rigts.
	case PieceBKing:
		p.CastlingRights &= ^(CastlingBlackShort | CastlingBlackLong)

	// Disable castling rights if the white rooks aren't standing on their initial positions.
	case PieceWRook:
		if p.Bitboards[PieceWRook]&A1 == 0 {
			p.CastlingRights &= ^CastlingWhiteLong
		}
		if p.Bitboards[PieceWRook]&H1 == 0 {
			p.CastlingRights &= ^CastlingWhiteShort
		}

	// Disable castling rights if the black rooks aren't standing on their initial positions.
	case PieceBRook:
		if p.Bitboards[PieceBRook]&A8 == 0 {
			p.CastlingRights &= ^CastlingBlackLong
		}
		if p.Bitboards[PieceBRook]&H8 == 0 {
			p.CastlingRights &= ^CastlingBlackShort
		}
	}

	// Increment the full move counter after black moves.
	if p.ActiveColor == ColorBlack {
		p.FullmoveCnt++
	}

	// Switch the active color.
	p.ActiveColor ^= 1
}

// GetPieceFromSquare returns the type of the piece that stands on the specified square.
// Returns [PieceNone] if there is no piece on the square.
func (p *Position) GetPieceFromSquare(square uint64) Piece {
	for piece, bitboard := range p.Bitboards {
		if square&bitboard != 0 {
			return piece
		}
	}
	return PieceNone
}

// MoveList is used to store moves. The main idea behind it is to preallocate
// an array with enough capacity to store all possible moves and avoid dynamic
// memory allocations.
type MoveList struct {
	// Maximum number of moves per chess position is equal to 218, hence 218 elements.
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

// Piece is an allias type to avoid bothersome conversion between int and Piece.
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

// PromotionFlag is an allias type to avoid bothersome conversion between int and Color.
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

// MoveType is an allias type to avoid bothersome conversion between int and MoveType.
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
type CastlingRights int

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
