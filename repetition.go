package chego

import (
	"strings"
)

// repetitionKey generates a compact string representation of a position with legal moves.
// This allows positions to be used as map keys and compared efficiently.
func repetitionKey(p Position, legalMoves MoveList) string {
	var keyBuilder strings.Builder
	keyBuilder.Grow(50)

	keyBuilder.WriteString(FromBitboardArray(p.Bitboards))
	keyBuilder.WriteByte(byte(p.ActiveColor))
	keyBuilder.WriteByte(byte(p.CastlingRights))

	for i := range legalMoves.LastMoveIndex {
		move := legalMoves.Moves[i]
		// Skip empty moves.
		if move != 0 {
			keyBuilder.WriteRune(rune(move))
		}
	}

	return keyBuilder.String()
}
