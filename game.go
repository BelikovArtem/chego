// game.go impements chess game state management.
// Make sure to call [InitAttackTables] ONCE before using
// functions from this file.

package chego

import "strings"

// Game represents a single chess game state.
type Game struct {
	Position    Position
	LegalMoves  MoveList
	MoveStack   []CompletedMove
	Repetitions map[string]int
	// Keep track of all captured pieces.
	Captured []Piece
}

// CompletedMove represents a completed move.
type CompletedMove struct {
	// Board state after completing the move
	// to enable move undo and state restoration.
	FenString string
	// Move itself.
	Move Move
}

// NewGame creates a new game initialized with the default chess position.
// Generates legal moves.
func NewGame() *Game {
	g := &Game{
		MoveStack:   make([]CompletedMove, 0, 15),
		Repetitions: make(map[string]int),
		Captured:    make([]Piece, 0, 15),
	}
	g.Position = ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	GenLegalMoves(g.Position, &g.LegalMoves)
	// Add beginning repitition key.
	key := repetitionKey(g.Position, g.LegalMoves)
	g.Repetitions[key]++
	return g
}

// PushMove updates the game state by performing the specified move.
// It is a caller responsibility to check if the specified move is legal.
// Generates legal moves for the next turn.
func (g *Game) PushMove(m Move) {
	captured := g.Position.GetPieceFromSquare(1 << m.To())

	g.Position.MakeMove(m)

	// Memorize captured piece.
	if captured != PieceNone {
		g.Captured = append(g.Captured, captured)
	}

	// Store the completed move.
	g.MoveStack = append(g.MoveStack, CompletedMove{
		Move:      m,
		FenString: SerializeFEN(g.Position),
	})

	// Generate legal moves for the next turn.
	GenLegalMoves(g.Position, &g.LegalMoves)

	// Add repetition key to detect repetitions.
	key := repetitionKey(g.Position, g.LegalMoves)
	g.Repetitions[key]++
}

// PopMove pops the last completed move and restores the game state.
// If there is no completed moves, this function is no-op.
// Also regenerates legal moves.
func (g *Game) PopMove() {
	if len(g.MoveStack) == 0 {
		return
	}

	// Decrement repetition.
	key := repetitionKey(g.Position, g.LegalMoves)
	g.Repetitions[key]--

	g.MoveStack = g.MoveStack[:len(g.MoveStack)-1]

	if len(g.MoveStack) > 0 {
		pos := ParseFEN(g.MoveStack[len(g.MoveStack)-1].FenString)
		g.Position = pos
		GenLegalMoves(g.Position, &g.LegalMoves)
	} else {
		// Restore initial game state.
		newGame := NewGame()
		g.Position = newGame.Position
		g.LegalMoves = newGame.LegalMoves
	}
}

// IsThreefoldRepetition checks whether the game has reached a threefold repetition.
//
// A position is considered identical if all of the following conditions are met:
//  1. Active colors are the same.
//  2. Pieces occupy the same squares.
//  3. Legal moves are the same.
//  4. Castling rights are identical.
//
// NOTE: Positions are identical even if the en passant target square differs,
// provided that no en passant capture is possible.
func (g *Game) IsThreefoldRepetition() bool {
	for _, reps := range g.Repetitions {
		if reps >= 3 {
			return true
		}
	}
	return false
}

// IsInsufficientMaterial returns true if one of the following statements is true:
//
//  1. Both sides have a bare king.
//  2. One side has a king and a minor piece against a bare king.
//  3. Both sides have a king and a bishop, the bishops standing on the same color.
//  4. Both sides have a king and a knight.
func (g *Game) IsInsufficientMaterial() bool {
	// Bitmask for all dark squares.
	dark := uint64(0xAA55AA55AA55AA55)
	material := g.calculateMaterial()

	if material == 0 {
		return true
	}

	if material == 3 && g.Position.Bitboards[PieceWPawn] == 0 &&
		g.Position.Bitboards[PieceBPawn] == 0 {
		return true
	}

	if material == 6 {
		whiteBishop := g.Position.Bitboards[PieceWBishop]
		blackBishop := g.Position.Bitboards[PieceBBishop]

		// TODO: clean up this mess
		return (whiteBishop != 0 && blackBishop != 0 && ((whiteBishop&dark > 0 &&
			blackBishop&dark > 0) || (whiteBishop&dark == 0 && blackBishop&dark == 0))) ||
			(g.Position.Bitboards[PieceWKnight] != 0 &&
				g.Position.Bitboards[PieceBKnight] != 0)
	}
	return false
}

// IsCheckmate returns true if one of the following statements is true:
//
//  1. There are no legal moves available for the current turn.
//  2. The king of the side to move is in check.
//
// NOTE: If there are no legal moves, but the king is not in check, the position
// is a stalemate.
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

// repetitionKey generates a compact string representation of a
// position with legal moves. This allows positions to be used as
// map keys and compared efficiently.
func repetitionKey(p Position, legalMoves MoveList) string {
	var keyBuilder strings.Builder
	keyBuilder.Grow(50)

	keyBuilder.WriteString(SerializeBitboards(p.Bitboards))
	keyBuilder.WriteByte(byte(p.ActiveColor))
	keyBuilder.WriteByte(byte(p.CastlingRights))

	for i := range legalMoves.LastMoveIndex {
		if legalMoves.Moves[i] == 0 {
			break
		}

		keyBuilder.WriteRune(rune(legalMoves.Moves[i]))
	}

	return keyBuilder.String()
}
