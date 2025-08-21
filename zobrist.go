/*
zobrist.go implements Zobrist hashing algorithm to detect position repetitions
(threefold repetition rule).
*/

package chego

import "math/rand/v2"

/*
Keys are used to hash each possible position into the unique number.  Each key
is generated randomly and large enough, so the probability of hash collisions is
negligible.
*/
var (
	pieceKeys [12][64]uint64
	// Used only when black is the active color.
	epKeys       [64]uint64
	castlingKeys [16]uint64
	// Used only when black is the active color.
	colorKey uint64
)

/*
InitZobristKeys initializes the pseudo-random keys used in the Zobrist hashing
scheme.  Call this function ONCE as close as possible to the start of your
program.

NOTE: Threefold repetitions will not be detected if this funtcion wasn't called.
*/
func InitZobristKeys() {
	for i := PieceWPawn; i <= PieceBKing; i++ {
		for square := range 64 {
			pieceKeys[i][square] = rand.Uint64()
		}
	}

	for square := range 64 {
		epKeys[square] = rand.Uint64()
	}

	for i := range 16 {
		castlingKeys[i] = rand.Uint64()
	}

	colorKey = rand.Uint64()
}

/*
zobristKey hashes the given position into a 64-bit unsigned integer.  This
allows positions to be used as lookup keys and stored or compared efficiently.
*/
func zobristKey(p Position) (key uint64) {
	for i := PieceWPawn; i <= PieceBKing; i++ {
		for p.Bitboards[i] > 0 {
			key ^= pieceKeys[i][popLSB(&p.Bitboards[i])]
		}
	}

	key ^= epKeys[p.EPTarget]

	key ^= castlingKeys[p.CastlingRights]

	key ^= colorKey & uint64(p.ActiveColor)

	return key
}
