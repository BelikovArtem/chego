// game.go impements chess game state management.

package chego

import "time"

// Game represents a single chess game state.
//
// Make sure to call [InitAttackTables] and [InitZobristKeys]
// ONCE before creating a [Game].
type Game struct {
	LegalMoves MoveList
	Position   Position
	MoveStack  []CompletedMove
	// Keep track of all captured pieces.
	Captured []Piece
	// Keep track of all repeated Zobrist keys to detect
	// a threefold repetition.
	Repetitions map[uint64]int
	// Remaining time on a white player's clock in seconds.
	// Set to 0 by default.
	WhiteTime int
	// Remaining time on a black player's clock in seconds.
	// Set to 0 by default.
	BlackTime int
	// Clock will send a tick signal every second. By default the Clock is
	// stopped. The caller should call SetClock to apply the time limit
	// for a game.
	Clock  *time.Ticker
	Result Result
}

// CompletedMove represents a completed move.
type CompletedMove struct {
	// Board state after completing the move
	// to enable move undo and state restoration.
	FenString string
	// Move itself.
	Move Move
	// Remaining time on a player's clock in seconds.
	TimeLeft int
}

// NewGame creates a new game initialized with the default chess position.
// Generates legal moves.
func NewGame() *Game {
	g := &Game{
		MoveStack:   make([]CompletedMove, 0, 15),
		Repetitions: make(map[uint64]int),
		Captured:    make([]Piece, 0, 15),
		Clock:       time.NewTicker(time.Second),
	}

	g.Clock.Stop()

	g.Position = ParseFEN(InitialPos)

	GenLegalMoves(g.Position, &g.LegalMoves)

	// Add initial repetition key.
	g.Repetitions[zobristKey(g.Position)]++
	return g
}

// PushMove updates the game state by performing the specified move.
// It is a caller responsibility to check if the specified move is legal.
// Generates legal moves for the next turn.
func (g *Game) PushMove(m Move) {
	moved := g.Position.GetPieceFromSquare(1 << m.From())
	captured := g.Position.GetPieceFromSquare(1 << m.To())

	g.Position.MakeMove(m)

	// Memorize the captured piece and clear the repetitions
	// map after applying irreversible moves.
	// See https://www.chessprogramming.org/Irreversible_Moves
	if captured != PieceNone {
		g.Captured = append(g.Captured, captured)
		clear(g.Repetitions)
	} else if m.Type() == MoveCastling || m.Type() == MovePromotion ||
		moved <= PieceBPawn {
		clear(g.Repetitions)
	}

	// Store the clock value.
	tl := g.WhiteTime
	if moved%2 != 0 {
		tl = g.BlackTime
	}

	// Store the completed move.
	g.MoveStack = append(g.MoveStack, CompletedMove{
		Move:      m,
		FenString: SerializeFEN(g.Position),
		TimeLeft:  tl,
	})

	// Generate legal moves for the next turn.
	GenLegalMoves(g.Position, &g.LegalMoves)

	// Is the en passant capture is not possible, clear the en passant
	// target, since it can break the threefold-repetition detection
	// by corrupting the Zobrist hash.
	// See [IsThreefoldRepetition] commentary
	ep := 0
	for i := range g.LegalMoves.LastMoveIndex {
		if g.LegalMoves.Moves[i].Type() == MoveEnPassant {
			ep = g.Position.EPTarget
		}
	}
	g.Position.EPTarget = ep

	// Add repetition key to detect repetitions.
	g.Repetitions[zobristKey(g.Position)]++
}

// PopMove pops the last completed move and restores the game state.
// If there are no completed moves, this function is no-op.
func (g *Game) PopMove() {
	if len(g.MoveStack) == 0 {
		return
	}

	// Decrement repetition key.
	g.Repetitions[zobristKey(g.Position)]--

	// Pop move from the stack.
	g.MoveStack = g.MoveStack[:len(g.MoveStack)-1]

	if len(g.MoveStack) == 0 { // No moves left.
		// Restore position.
		p := ParseFEN(g.MoveStack[len(g.MoveStack)-1].FenString)
		g.Position = p
		// Restore time on the player timers.
		// Since there are no more completed moves, to restore the initial
		// clock values just assign a WhiteTime to a BlackTime.
		g.WhiteTime = g.BlackTime
	} else if len(g.MoveStack)%2 == 0 {
		// White player has moved.
		last := g.MoveStack[len(g.MoveStack)-1]
		g.Position = ParseFEN(last.FenString)
		g.WhiteTime = last.TimeLeft
	} else {
		// Black player has moved.
		last := g.MoveStack[len(g.MoveStack)-1]
		g.Position = ParseFEN(last.FenString)
		g.BlackTime = last.TimeLeft
	}

	// Restore legal moves.
	GenLegalMoves(g.Position, &g.LegalMoves)
}

// IsThreefoldRepetition checks whether the game has reached
// a threefold repetition.
//
// A position is considered identical if all of the following
// conditions are met:
//  1. Active colors are the same.
//  2. Pieces occupy the same squares.
//  3. Legal moves are the same.
//  4. Castling rights are identical.
//
// NOTE: Positions are identical even if the en passant target
// square differs, provided that no en passant capture is possible.
func (g *Game) IsThreefoldRepetition() bool {
	for _, reps := range g.Repetitions {
		if reps >= 3 {
			return true
		}
	}
	return false
}

// IsInsufficientMaterial returns true if one of the following
// statements is true:
//  1. Both sides have a bare king.
//  2. One side has a king and a minor piece against a bare king.
//  3. Both sides have a king and a bishop, the bishops standing
//     on the same color.
//  4. Both sides have a king and a knight.
func (g *Game) IsInsufficientMaterial() bool {
	// Bitmask for all dark squares.
	dark := uint64(0xAA55AA55AA55AA55)
	mat := g.calculateMaterial()

	if mat == 0 {
		return true
	}

	if mat == 3 && g.Position.Bitboards[PieceWPawn] == 0 &&
		g.Position.Bitboards[PieceBPawn] == 0 {
		return true
	}

	if mat == 6 {
		wb := g.Position.Bitboards[PieceWBishop]
		bb := g.Position.Bitboards[PieceBBishop]

		// If there are two bishops both standing on the same colored squares.
		return (wb != 0 && bb != 0 && ((wb&dark > 0 && bb&dark > 0) ||
			(wb&dark == 0 && bb&dark == 0))) ||
			// Or if there are two knights.
			(g.Position.Bitboards[PieceWKnight] != 0 &&
				g.Position.Bitboards[PieceBKnight] != 0)
	}
	return false
}

// IsCheckmate returns true if one of the following statements is true:
//  1. There are no legal moves available for the current turn.
//  2. The king of the side to move is in check.
//
// NOTE: If there are no legal moves, but the king is not in check,
// the position is a stalemate.
func (g *Game) IsCheckmate() bool {
	isKingInCheck := GenChecksCounter(g.Position.Bitboards,
		1^g.Position.ActiveColor) > 0
	return isKingInCheck && g.LegalMoves.LastMoveIndex == 0
}

// IsMoveLegal checks if the specified move is legal.
func (g *Game) IsMoveLegal(m Move) bool {
	for _, move := range g.LegalMoves.Moves {
		if move == 0x0 {
			return false
		}
		if move.From() == m.From() && move.To() == m.To() &&
			move.Type() == m.Type() && move.PromoPiece() == m.PromoPiece() {
			return true
		}
	}
	return false
}

// SetClock sets each player's clock value to dur and starts the [Clock],
// applying a time limit to both players.
func (g *Game) SetClock(dur int) {
	g.WhiteTime = dur
	g.BlackTime = dur
	g.Clock.Reset(time.Second)
}

// TimeTick decrements the clock value of the currently active player
// in response to ticks from the [Clock]. It is the caller's responsibility
// to end the game if one of the players runs out of time.
func (g *Game) TimeTick() {
	if g.Position.ActiveColor == ColorWhite {
		g.WhiteTime--
	} else {
		g.BlackTime--
	}
}

// calculateMaterial calculates the piece valies of each side.
// Is used to determine a draw by insufficient material.
func (g *Game) calculateMaterial() (mat int) {
	coeff := 1
	for pieceType := range PieceWKing {
		switch pieceType {
		case PieceWKnight, PieceBKnight,
			PieceWBishop, PieceBBishop:
			coeff = 3
		case PieceWRook, PieceBRook:
			coeff = 5
		case PieceWQueen, PieceBQueen:
			coeff = 9
		}

		mat += CountBits(g.Position.Bitboards[pieceType]) * coeff
	}

	return mat
}
