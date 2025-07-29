// uci.go implements Universal Chess Interface.

package chego

import "strings"

// Move2UCI converts the move into long algebraic notation string.
// Examples: e2e4, e7e5, e1g1 (white short castling), e7e8q (for promotion).
func Move2UCI(m Move) string {
	var b strings.Builder
	b.Grow(4)

	b.WriteString(Square2String[m.From()])
	b.WriteString(Square2String[m.To()])

	if m.Type() == MovePromotion {
		switch m.PromoPiece() {
		case PromotionKnight:
			b.WriteByte('n')
		case PromotionBishop:
			b.WriteByte('b')
		case PromotionRook:
			b.WriteByte('r')
		case PromotionQueen:
			b.WriteByte('q')
		}
	}

	return b.String()
}
