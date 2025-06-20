package game

import (
	"chego/enum"
	"chego/fen"
	"chego/movegen"
	"strings"
)

// position is used to implement the threefold-repetition rule.
type position struct {
	legalMoves     movegen.MoveList
	bitboards      [12]uint64
	activeColor    enum.Color
	castlingRights enum.CastlingFlag
}

// repetitionKey is used to compare positions and use them as map keys.
func (p position) repetitionKey() string {
	var keyBuilder strings.Builder
	keyBuilder.Grow(50)

	keyBuilder.WriteString(fen.FromBitboardArray(p.bitboards))
	keyBuilder.WriteByte(byte(p.activeColor))
	keyBuilder.WriteByte(byte(p.castlingRights))

	for _, move := range p.legalMoves.Moves {
		keyBuilder.WriteRune(rune(move))
	}

	return keyBuilder.String()
}
