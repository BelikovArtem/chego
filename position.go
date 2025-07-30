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
	// prevPos is nil by default.
	prevPos *Position
}

// MakeMove modifies the position by applying the specified move.
// It's a caller responsibility to ensure the legality of the move.
// The position is saved in the prevPos field before making the move to
// allow undo functionality.
//
// Not only is the piece placement updated, but also the entire position,
// including castling rights, en passant target, halfmove counter, fullmove
// counter, and the active color.
func (p *Position) MakeMove(m Move) {
	// Save the position to be able to undo the move.
	p.prevPos = &Position{
		ActiveColor:    p.ActiveColor,
		CastlingRights: p.CastlingRights,
		EPTarget:       p.EPTarget,
		HalfmoveCnt:    p.HalfmoveCnt,
		FullmoveCnt:    p.FullmoveCnt,
		prevPos:        p.prevPos,
	}
	copy(p.prevPos.Bitboards[:], p.Bitboards[:])

	to := uint64(1 << m.To())
	from := uint64(1 << m.From())
	piece := p.GetPieceFromSquare(from)
	captured := p.GetPieceFromSquare(to)

	// Clear the origin square.
	p.removePiece(piece, from)

	// Increment halfmove counter to detect 50-move rule draw.
	// This will be reseted if the move is capture, or a pawn push.
	p.HalfmoveCnt++

	// Remove the captured piece from the board.
	// This skips en passant captures, since the captured
	// pawn does not occupy the square the capturing piece moves to.
	if captured != PieceNone {
		p.removePiece(captured, to)
		// Reset the halfmove counter after capture.
		p.HalfmoveCnt = 0
	}

	switch m.Type() {
	case MoveNormal:
		p.placePiece(piece, to)

	case MoveEnPassant:
		p.placePiece(piece, to)
		// Remove the captured piece from the board.
		if piece == PieceWPawn {
			p.removePiece(PieceBPawn, to>>8)
		} else {
			p.removePiece(PieceWPawn, to<<8)
		}

	case MoveCastling:
		p.placePiece(piece, to)
		// Update the rook position.
		switch to {
		case G1: // White O-O.
			p.removePiece(PieceWRook, H1)
			p.placePiece(PieceWRook, F1)
		case G8: // Black O-O.
			p.removePiece(PieceBRook, H8)
			p.placePiece(PieceBRook, F8)
		case C1: // White O-O-O.
			p.removePiece(PieceWRook, A1)
			p.placePiece(PieceWRook, D1)
		case C8: // Black O-O-O.
			p.removePiece(PieceBRook, A8)
			p.placePiece(PieceBRook, D8)
		}

	case MovePromotion:
		// Color offset to correctly promote a piece.
		off := 1
		if piece == PieceBPawn {
			off = 7
		}
		p.placePiece(m.PromoPiece()+off, to)
	}

	// Reset the en passant target since the en passant capture
	// is only legal for 1 move.
	p.EPTarget = 0

	switch piece {
	// Set en passant target square is case of double pawn push.
	case PieceWPawn, PieceBPawn:
		if m.From()-m.To() == -16 {
			p.EPTarget = m.To() - 8
		} else if m.From()-m.To() == 16 {
			p.EPTarget = m.To() + 8
		}
		// Reset the halfmove counter after pawn moves.
		p.HalfmoveCnt = 0
	// The king cannot castle with a rook that has already moved.
	case PieceWRook:
		switch m.From() {
		case SA1:
			p.CastlingRights &= ^CastlingWhiteLong
		case SH1:
			p.CastlingRights &= ^CastlingWhiteShort
		}
	// The king cannot castle with a rook that has already moved.
	case PieceBRook:
		switch m.From() {
		case SA8:
			p.CastlingRights &= ^CastlingBlackLong
		case SH8:
			p.CastlingRights &= ^CastlingBlackShort
		}
	// Disable white castling rights.
	case PieceWKing:
		p.CastlingRights &= ^(CastlingWhiteShort | CastlingWhiteLong)
	// Disable black castling rights.
	case PieceBKing:
		p.CastlingRights &= ^(CastlingBlackShort | CastlingBlackLong)
	}

	// Increment the full move counter after black moves.
	if p.ActiveColor == ColorBlack {
		p.FullmoveCnt++
	}

	// Switch the active color.
	p.ActiveColor ^= 1
}

// UndoMove restores the position to the previous state, if it's not nil.
func (p *Position) UndoMove() {
	if p.prevPos != nil {
		*p = *p.prevPos
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

// CanCastle checks whether the king can peform
// castling in the specified direction.
// side == 1 -> White O-O.
// side == 2 -> White O-O-O.
// side == 4 -> Black O-O.
// side == 8 -> Black O-O-O.
func (p *Position) CanCastle(side int, attacks, occupancy uint64) bool {
	path := castlingPath[bitScan(uint64(side))]
	return p.CastlingRights&side != 0 &&
		attacks&path == 0 &&
		occupancy&path == 0
}

// placePiece places the piece on the specified square and
// updates the occupancy and allies bitboards.
func (p *Position) placePiece(piece Piece, square uint64) {
	// Place the piece.
	p.Bitboards[piece] |= square
	// Update allies bitboard.
	if piece <= PieceWKing { // White piece.
		p.Bitboards[12] |= square
	} else { // Black piece.
		p.Bitboards[13] |= square
	}
	// Update occupancy bitboard.
	p.Bitboards[14] |= square
}

// removePiece removes the piece from the specified square and
// updates the occupancy and allies bitboards.
//
// NOTE: if there is no piece of the specified type on the
// specified square, this function will place the piece
// instead of removing it.
func (p *Position) removePiece(piece Piece, square uint64) {
	// Remove the piece.
	p.Bitboards[piece] ^= square
	// Update allies bitboard.
	if piece < PieceWKing { // White piece.
		p.Bitboards[12] ^= square
	} else {
		p.Bitboards[13] ^= square
	}
	// Update occupancy bitboard.
	p.Bitboards[14] ^= square
}
