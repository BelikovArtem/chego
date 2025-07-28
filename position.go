// position.go defines the Position structure and it's methods
// for representing and modifying the chessboard state.

package chego

// Position represents a chessboard state that can be
// converted to or parsed from a FEN string.
type Position struct {
	Bitboards      [15]uint64
	ActiveColor    Color
	CastlingRights CastlingRights
	EPTarget       int
	HalfmoveCnt    int
	FullmoveCnt    int
	lastPos        *Position
}

// MakeMove modifies the position by applying the specified move.
// It's the caller's responsibility to ensure the move is legal.
// The current position is saved in the lastPos field before making the move.
//
// Not only is the piece placement updated, but also the entire position,
// including castling rights, en passant target, move counters, and the active color.
func (p *Position) MakeMove(m Move) {
	var from, to uint64 = 1 << m.From(), 1 << m.To()
	fromTo := from ^ to
	movedPiece := p.GetPieceFromSquare(from)

	p.lastPos = &Position{
		Bitboards:      p.Bitboards,
		ActiveColor:    p.ActiveColor,
		CastlingRights: p.CastlingRights,
		EPTarget:       p.EPTarget,
		HalfmoveCnt:    p.HalfmoveCnt,
		FullmoveCnt:    p.FullmoveCnt,
	}

	switch m.Type() {
	case MoveNormal:
		// If the move is capture.
		capturedPiece := p.GetPieceFromSquare(to)
		if capturedPiece != PieceNone {
			// Remove the captured piece from the board.
			p.Bitboards[capturedPiece] ^= to
			p.Bitboards[12+(1^p.ActiveColor)] ^= to
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
		p.Bitboards[12+(1^p.ActiveColor)] ^= to

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
		capturedPiece := p.GetPieceFromSquare(to)
		if capturedPiece != PieceNone {
			// Remove the captured piece from the board.
			p.Bitboards[capturedPiece] ^= to
			p.Bitboards[12+(1^p.ActiveColor)] ^= to
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
	// Update allies bitboard.
	p.Bitboards[12+p.ActiveColor] ^= fromTo
	// Update occupancy bitboard.
	p.Bitboards[14] ^= fromTo

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

// UndoMove restores the position to the last known state, it it's not nil.
func (p *Position) UndoMove() {
	if p.lastPos != nil {
		*p = *p.lastPos
	}
}

// GetPieceFromSquare returns the type of the piece that stands
// on the specified square, or [PieceNone] if the square is empty.
func (p *Position) GetPieceFromSquare(square uint64) Piece {
	for piece, bitboard := range p.Bitboards {
		if square&bitboard != 0 {
			return piece
		}
	}
	return PieceNone
}
