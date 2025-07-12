// Package game impements chess game state management.
// Make sure to call [movegen.InitAttackTables] ONCE
// before using other functions from this package.
package game

import (
	"github.com/BelikovArtem/chego/bitutil"
	"github.com/BelikovArtem/chego/enum"
	"github.com/BelikovArtem/chego/fen"
	"github.com/BelikovArtem/chego/movegen"
)

// Game represents a single chess game state.
type Game struct {
	LegalMoves      movegen.MoveList
	Bitboards       [12]uint64
	MoveStack       []CompletedMove
	Repetitions     map[string]int
	EnPassantTarget int
	CastlingRights  enum.CastlingFlag
	ActiveColor     enum.Color
	HalfmoveCnt     int
	FullmoveCnt     int
}

// CompletedMove represents a completed move.
type CompletedMove struct {
	// Game state after completing the move to enable move undo and state restoration.
	FenString string
	// Move itself.
	Move movegen.Move
}

// NewGame creates a new game initialized with the default chess position.
// Generates legal moves.
func NewGame() *Game {
	g := &Game{
		Bitboards:       fen.ToBitboardArray("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"),
		MoveStack:       make([]CompletedMove, 0),
		Repetitions:     make(map[string]int),
		EnPassantTarget: 0,
		CastlingRights:  0xF,
		ActiveColor:     enum.ColorWhite,
		HalfmoveCnt:     0,
		FullmoveCnt:     1,
	}

	g.LegalMoves = movegen.GenLegalMoves(g.Bitboards, g.ActiveColor, g.CastlingRights, g.EnPassantTarget)
	return g
}

// SetState sets the game state and generates legal moves.
func (g *Game) SetState(bitboards [12]uint64, enPassantTarget int, castlingRights enum.CastlingFlag,
	activeColor enum.Color, halfmoveCnt, fullmoveCnt int) {
	g.Bitboards = bitboards
	g.EnPassantTarget = enPassantTarget
	g.CastlingRights = castlingRights
	g.ActiveColor = activeColor
	g.HalfmoveCnt = halfmoveCnt
	g.FullmoveCnt = fullmoveCnt

	g.LegalMoves = movegen.GenLegalMoves(g.Bitboards, g.ActiveColor, g.CastlingRights, g.EnPassantTarget)
}

// PushMove updates the game state by performing the specified move.
// It is a caller responsibility to check if the specified move is legal.
// Generates legal moves for the next turn.
func (g *Game) PushMove(move movegen.Move) {
	movedPiece := movegen.GetPieceTypeFromSquare(g.Bitboards, 1<<move.From())
	isCapture := movegen.GetPieceTypeFromSquare(g.Bitboards, 1<<move.To()) != -1
	movegen.MakeMove(&g.Bitboards, move)

	// Reset the en passant target since the en passant capture is possible only for 1 move.
	g.EnPassantTarget = 0

	switch movedPiece {
	case enum.PieceWPawn, enum.PieceBPawn:
		// Set en passant target.
		if g.ActiveColor == enum.ColorWhite && move.From()-move.To() == -16 {
			g.EnPassantTarget = move.To() - 8
		} else if g.ActiveColor == enum.ColorBlack && move.From()-move.To() == 16 {
			g.EnPassantTarget = move.To() + 8
		}
	case enum.PieceWKing:
		// Disable castling rights for white.
		g.CastlingRights &= ^(enum.CastlingWhiteShort | enum.CastlingWhiteLong)
	case enum.PieceBKing:
		// Disable castling rights for black.
		g.CastlingRights &= ^(enum.CastlingBlackShort | enum.CastlingBlackLong)
	}

	// Disable castling rights if the rooks aren't standing on their initial positions.
	if g.Bitboards[enum.PieceWRook]&enum.A1 == 0 {
		g.CastlingRights &= ^enum.CastlingWhiteLong
	}
	if g.Bitboards[enum.PieceWRook]&enum.H1 == 0 {
		g.CastlingRights &= ^enum.CastlingWhiteShort
	}
	if g.Bitboards[enum.PieceBRook]&enum.A8 == 0 {
		g.CastlingRights &= ^enum.CastlingBlackLong
	}
	if g.Bitboards[enum.PieceBRook]&enum.H8 == 0 {
		g.CastlingRights &= ^enum.CastlingBlackShort
	}

	if isCapture ||
		movedPiece == enum.PieceWPawn || movedPiece == enum.PieceBPawn {
		g.HalfmoveCnt = 0
	} else {
		g.HalfmoveCnt++
	}

	// Increment the full move counter after black moves.
	if g.ActiveColor == enum.ColorBlack {
		g.FullmoveCnt++
	}

	// Switch the active color.
	g.ActiveColor ^= 1

	// Store the completed move.
	g.MoveStack = append(g.MoveStack, CompletedMove{
		Move: move,
		FenString: fen.Serialize(g.Bitboards, g.ActiveColor, g.CastlingRights,
			g.EnPassantTarget, g.HalfmoveCnt, g.FullmoveCnt),
	})

	// Generate legal moves for the next turn.
	g.LegalMoves = movegen.GenLegalMoves(g.Bitboards, g.ActiveColor, g.CastlingRights,
		g.EnPassantTarget)
}

// PopMove pops the last completed move and restores the game state.
// If there is no completed moves, this function is no-op.
// Also regenerates legal moves.
func (g *Game) PopMove() {
	if len(g.MoveStack) == 0 {
		return
	}

	g.MoveStack = g.MoveStack[:len(g.MoveStack)-1]

	if len(g.MoveStack) > 0 {
		// bitboards, activeColor, enPassantTarget, halfmoveCnt, fullmoveCnt
		b, a, c, e, h, f := fen.Parse(g.MoveStack[len(g.MoveStack)-1].FenString)
		g.SetState(b, e, c, a, h, f)
		movegen.GenLegalMoves(b, a, c, e)
	} else {
		// Restore initial game state.
		newGame := NewGame()
		g.SetState(newGame.Bitboards, newGame.EnPassantTarget, newGame.CastlingRights,
			newGame.ActiveColor, newGame.HalfmoveCnt, newGame.FullmoveCnt)
		g.LegalMoves = newGame.LegalMoves
	}
}

// IsThreefoldRepetition checks whether the last completed move has resulted in a threefold repetition.
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
	currentPosKey := position{
		g.LegalMoves,
		g.Bitboards,
		g.ActiveColor,
		g.CastlingRights,
	}.repetitionKey()
	// Increment the repetition count.
	g.Repetitions[currentPosKey]++
	return g.Repetitions[currentPosKey] == 3
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

	if material == 3 && g.Bitboards[enum.PieceWPawn] == 0 &&
		g.Bitboards[enum.PieceBPawn] == 0 {
		return true
	}

	if material == 6 {
		whiteBishop := g.Bitboards[enum.PieceWBishop]
		blackBishop := g.Bitboards[enum.PieceBBishop]

		if whiteBishop != 0 && blackBishop != 0 && ((whiteBishop&dark > 0 &&
			blackBishop&dark > 0) || (whiteBishop&dark == 0 && blackBishop&dark == 0)) {
			return true
		} else if g.Bitboards[enum.PieceWKnight] != 0 && g.Bitboards[enum.PieceBKnight] != 0 {
			return true
		}
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
	for i := enum.PieceWPawn; i <= enum.PieceBKing; i++ {
		occupancy |= g.Bitboards[i]
	}
	kingBB := g.Bitboards[enum.PieceWKing]
	if g.ActiveColor == enum.ColorBlack {
		kingBB = g.Bitboards[enum.PieceBKing]
	}

	isKingInCheck := movegen.IsSquareUnderAttack(g.Bitboards, occupancy, bitutil.BitScan(kingBB),
		1^g.ActiveColor)
	return isKingInCheck && g.LegalMoves.LastMoveIndex == 0
}

// GetLegalMoveIndex checks if the specified move is legal.
// If it is, it returns the index of the legal move in the game.LegalMoves list.
// If it isn't, it returns -1.
//
// NOTE: It also updates the promotion piece flag in the legal move,
// so the player can promote to the desired piece.
func (g *Game) GetLegalMoveIndex(to, from, promotionPiece int) int {
	for i, legalMove := range g.LegalMoves.Moves {
		if legalMove.From() == from && legalMove.To() == to {
			if legalMove.Type() == enum.MovePromotion {
				// Update promotion piece in case it is invalid.
				if promotionPiece > enum.PromotionKnight && promotionPiece < enum.PromotionQueen {
					promotionPiece = enum.PromotionQueen
				}
				g.LegalMoves.Moves[i] = movegen.NewMove(to, from, promotionPiece, enum.MovePromotion)
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

	for pieceType := enum.PieceWPawn; pieceType < enum.PieceBKing; pieceType++ {
		if pieceType == enum.PieceWKing {
			continue
		}

		coefficient := 1
		switch pieceType {
		case enum.PieceWKnight, enum.PieceBKnight,
			enum.PieceWBishop, enum.PieceBBishop:
			coefficient = 3
		case enum.PieceWRook, enum.PieceBRook:
			coefficient = 5
		case enum.PieceWQueen, enum.PieceBQueen:
			coefficient = 9
		}

		material += bitutil.CountBits(g.Bitboards[pieceType]) * coefficient
	}

	return material
}
