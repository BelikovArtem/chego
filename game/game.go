// Package game impements chess game state management.
//
// Make sure to call [movegen.InitAttackTables] ONCE before using this package.
package game

import (
	"github.com/BelikovArtem/chego/bitutil"
	"github.com/BelikovArtem/chego/fen"
	"github.com/BelikovArtem/chego/movegen"
	"github.com/BelikovArtem/chego/types"
)

// Game represents a single chess game state.
type Game struct {
	Position    types.Position
	LegalMoves  types.MoveList
	MoveStack   []CompletedMove
	Repetitions map[string]int
	// Keep track of all captured pieces.
	Captured []types.Piece
}

// CompletedMove represents a completed move.
type CompletedMove struct {
	// Game state after completing the move to enable move undo and state restoration.
	FenString string
	// Move itself.
	Move types.Move
}

// NewGame creates a new game initialized with the default chess position.
// Generates legal moves.
func NewGame() *Game {
	g := &Game{
		Position:    fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
		MoveStack:   make([]CompletedMove, 0, 15),
		Repetitions: make(map[string]int),
		Captured:    make([]types.Piece, 0, 15),
	}

	movegen.GenLegalMoves(g.Position, &g.LegalMoves)
	// Add beginning repitition key.
	key := repetitionKey(g.Position, g.LegalMoves)
	g.Repetitions[key]++
	return g
}

// PushMove updates the game state by performing the specified move.
// It is a caller responsibility to check if the specified move is legal.
// Generates legal moves for the next turn.
func (g *Game) PushMove(m types.Move) {
	captured := g.Position.GetPieceFromSquare(1 << m.To())

	g.Position.MakeMove(m)

	// Memorize captured piece.
	if captured != types.PieceNone {
		g.Captured = append(g.Captured, captured)
	}

	// Store the completed move.
	g.MoveStack = append(g.MoveStack, CompletedMove{
		Move:      m,
		FenString: fen.Serialize(g.Position),
	})

	// Generate legal moves for the next turn.
	movegen.GenLegalMoves(g.Position, &g.LegalMoves)

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
		pos := fen.Parse(g.MoveStack[len(g.MoveStack)-1].FenString)
		g.Position = pos
		movegen.GenLegalMoves(g.Position, &g.LegalMoves)
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
	var dark uint64 = 0xAA55AA55AA55AA55
	material := g.calculateMaterial()

	if material == 0 {
		return true
	}

	if material == 3 && g.Position.Bitboards[types.PieceWPawn] == 0 &&
		g.Position.Bitboards[types.PieceBPawn] == 0 {
		return true
	}

	if material == 6 {
		whiteBishop := g.Position.Bitboards[types.PieceWBishop]
		blackBishop := g.Position.Bitboards[types.PieceBBishop]

		// TODO: clean up this mess
		return (whiteBishop != 0 && blackBishop != 0 && ((whiteBishop&dark > 0 &&
			blackBishop&dark > 0) || (whiteBishop&dark == 0 && blackBishop&dark == 0))) ||
			(g.Position.Bitboards[types.PieceWKnight] != 0 &&
				g.Position.Bitboards[types.PieceBKnight] != 0)
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
	// Check whether the king in check.
	var occupancy uint64
	for i := types.PieceWPawn; i <= types.PieceBKing; i++ {
		occupancy |= g.Position.Bitboards[i]
	}
	kingBB := g.Position.Bitboards[types.PieceWKing]
	if g.Position.ActiveColor == types.ColorBlack {
		kingBB = g.Position.Bitboards[types.PieceBKing]
	}

	isKingInCheck := movegen.IsSquareUnderAttack(g.Position.Bitboards, bitutil.BitScan(kingBB),
		1^g.Position.ActiveColor)
	return isKingInCheck && g.LegalMoves.LastMoveIndex == 0
}

// GetLegalMoveIndex checks if the specified move is legal.
// If it is, it returns the index of the legal move in the game.LegalMoves list.
// If it isn't, it returns -1.
//
// NOTE: It also updates the promotion piece flag in the legal move,
// so the player can promote to the desired piece.
func (g *Game) GetLegalMoveIndex(m types.Move) int {
	for i, legalMove := range g.LegalMoves.Moves {
		if legalMove.From() == m.From() && legalMove.To() == m.To() {
			if legalMove.Type() == types.MovePromotion {
				promo := m.PromotionPiece()
				// Update promotion piece in case it is invalid.
				if promo < types.PromotionKnight || promo > types.PromotionQueen {
					promo = types.PromotionQueen
				}
				g.LegalMoves.Moves[i] = types.NewPromotionMove(m.To(), m.From(), promo)
			}
			return i
		}
	}
	return -1
}

// calculateMaterial calculates the piece valies of each side.
// Is used to determine a draw by insufficient material.
func (g *Game) calculateMaterial() int {
	var material int

	for pieceType := types.PieceWPawn; pieceType < types.PieceBKing; pieceType++ {
		if pieceType == types.PieceWKing {
			continue
		}

		coefficient := 1
		switch pieceType {
		case types.PieceWKnight, types.PieceBKnight,
			types.PieceWBishop, types.PieceBBishop:
			coefficient = 3
		case types.PieceWRook, types.PieceBRook:
			coefficient = 5
		case types.PieceWQueen, types.PieceBQueen:
			coefficient = 9
		}

		material += bitutil.CountBits(g.Position.Bitboards[pieceType]) * coefficient
	}

	return material
}
