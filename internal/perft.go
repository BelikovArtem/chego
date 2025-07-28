// Package main provides debugging and testing functions.
// It is excluded from the chego package, as it is only used for testing purposes.
// The chego users won't be able to import this package.
package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/BelikovArtem/chego"
)

// Test positions. See https://www.chessprogramming.org/Perft
const initFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Perft is a debugging function that walks through the move generation
// tree of strictly legal moves to a given depth and counts the number of
// visited leaf nodes. The resulting count is then compared to
// predetermined values.
//
// See https://www.chessprogramming.org/Perft_Results
func Perft(p chego.Position, depth int, isRoot bool) uint64 {
	if depth == 0 {
		return 1
	}

	nodes := uint64(0)
	l := chego.MoveList{}

	chego.GenLegalMoves(p, &l)

	for i := 0; i < int(l.LastMoveIndex); i++ {
		m := l.Moves[i]

		p.MakeMove(m)

		cnt := Perft(p, depth-1, false)
		// if isRoot {
		log.Printf("%s %d", chego.Move2UCI(m), cnt)
		// }
		nodes += cnt

		p.UndoMove()
	}

	return nodes
}

// main calls the Perft function and measures it's execution time.
func main() {
	depth := flag.Int("depth", 5, "Performance test depth")
	flag.Parse()

	nodes := uint64(0)

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Printf("Nodes reached: %d", nodes)
		log.Printf("Elapsed time: %d ns", elapsed.Nanoseconds())
	}()

	p := chego.ParseFEN(initFEN)

	nodes = Perft(p, *depth, true)
}

// position formats a full chess position into a string.
func position(p chego.Position) string {
	var posStr strings.Builder

	for rank := 7; rank >= 0; rank-- {
		posStr.WriteByte(byte(rank) + 1 + '0')
		posStr.WriteString("  ")

		for file := 0; file < 8; file++ {
			square := uint64(1 << (8*rank + file))

			symbol := byte('.')

			for i := chego.PieceWPawn; i <= chego.PieceBKing; i++ {
				if square&p.Bitboards[i] != 0 {
					symbol = chego.PieceSymbols[i]
					break
				}
			}

			posStr.WriteByte(symbol)
			posStr.WriteString("  ")
		}
		posStr.WriteByte('\n')
	}

	posStr.WriteString("   a  b  c  d  e  f  g  h\nActive color: ")

	if p.ActiveColor == chego.ColorWhite {
		posStr.WriteString("white\nEn passant: ")
	} else {
		posStr.WriteString("black\nEn passant: ")
	}

	if p.EPTarget == 0 {
		posStr.WriteString("none\nCastling rights: ")
	} else {
		posStr.WriteString(chego.Square2String[p.EPTarget])
		posStr.WriteString("\nCastling rights: ")
	}

	if p.CastlingRights&chego.CastlingWhiteShort != 0 {
		posStr.WriteByte('K')
	}
	if p.CastlingRights&chego.CastlingWhiteLong != 0 {
		posStr.WriteByte('Q')
	}
	if p.CastlingRights&chego.CastlingBlackShort != 0 {
		posStr.WriteByte('k')
	}
	if p.CastlingRights&chego.CastlingBlackLong != 0 {
		posStr.WriteByte('q')
	}

	return posStr.String()
}
