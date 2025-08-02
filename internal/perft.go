// Package main provides debugging and testing functions.
// It is excluded from the chego package, as it is only used
// for testing purposes. The chego users won't be able to import this package.
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

// result information will be printed is the perft is executed with the
// verbose flag.
type result struct {
	nodes        int
	captures     int
	epCaptures   int
	castles      int
	promotions   int
	checks       int
	doubleChecks int
	checkmates   int
}

// perft is a debugging function that walks through the move generation
// tree of strictly legal moves to a given depth and counts the number of
// visited leaf nodes. The resulting count is then compared to
// predetermined values.
//
// See https://www.chessprogramming.org/Perft_Results
func perft(p *chego.Position, depth int) int {
	if depth == 0 {
		return 1
	}

	nodes := 0
	l := chego.MoveList{}

	chego.GenLegalMoves(*p, &l)

	for i := range l.LastMoveIndex {
		m := l.Moves[i]

		p.MakeMove(m)

		nodes += perft(p, depth-1)

		p.UndoMove()
	}

	return nodes
}

// perftVerbose follows the same principle as the perft function, except it
// writes detailed move debugging information to r. Use this function to debug
// and find invalid branches in the move generation tree,
// not to measure performance.
func perftVerbose(p *chego.Position, depth int, r *result, isRoot bool) int {
	if depth == 0 {
		return 1
	}

	l := chego.MoveList{}

	chego.GenLegalMoves(*p, &l)
	if l.LastMoveIndex == 0 {
		r.checkmates++
	}

	nodes := 0

	c := p.ActiveColor
	for i := range l.LastMoveIndex {
		m := l.Moves[i]

		if p.GetPieceFromSquare(1<<m.To()) != chego.PieceNone {
			r.captures++
		}

		p.MakeMove(m)

		checkers := chego.GenCheckingPieces(p.Bitboards, 1^c)
		if checkers > 0 {
			r.checks++
		}
		if chego.CountBits(checkers) > 1 {
			r.doubleChecks++
		}

		cnt := perftVerbose(p, depth-1, r, false)
		if isRoot {
			log.Printf("%s %d", chego.Move2UCI(m), cnt)
		}
		nodes += cnt

		switch m.Type() {
		case chego.MoveCastling:
			r.castles++
		case chego.MoveEnPassant:
			r.epCaptures++
		case chego.MovePromotion:
			r.promotions++
		}

		p.UndoMove()
	}

	return nodes
}

// main runs the perft and measures it's execution time.
func main() {
	// It is important to initialize the attack tables.
	// Otherwise, perft will not work.
	chego.InitAttackTables()

	depth := flag.Int("depth", 2, "Performance test depth")
	verbose := flag.Bool("verbose", false, "Wether to print the debug info")

	flag.Parse()

	r := &result{}

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)

		if *verbose {
			log.Printf("\nRoot position:\n%s\n\n\t%s\n\n",
				position(chego.ParseFEN(initFEN)), initFEN)
			log.Printf("\t%d\t%d\t\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t",
				*depth,
				r.nodes,
				r.captures,
				r.epCaptures,
				r.castles,
				r.promotions,
				r.checks,
				r.doubleChecks,
				r.checkmates,
			)
			log.Printf("Elapsed time: %d ns", elapsed.Nanoseconds())
		} else {

			log.Printf("Nodes reached: %d", r.nodes)
			log.Printf("Elapsed time: %d ns", elapsed.Nanoseconds())
		}
	}()

	p := chego.ParseFEN(initFEN)

	if *verbose {
		r.nodes = perftVerbose(&p, *depth, r, true)
	} else {
		r.nodes = perft(&p, *depth)
	}
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
