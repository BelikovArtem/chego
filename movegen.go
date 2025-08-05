// movegen.go implements move generation using Magic Bitboards approach.

package chego

const (
	// Bitmask of all files except the A.
	NOT_A_FILE uint64 = 0xFEFEFEFEFEFEFEFE
	// Bitmask of all files except the H.
	NOT_H_FILE uint64 = 0x7F7F7F7F7F7F7F7F
	// Bitmask of all files except the A and B.
	NOT_AB_FILE uint64 = 0xFCFCFCFCFCFCFCFC
	// Bitmask of all files except the G and H.
	NOT_GH_FILE uint64 = 0x3F3F3F3F3F3F3F3F
	// Bitmask of all ranks except first.
	NOT_1ST_RANK uint64 = 0xFFFFFFFFFFFFFF00
	// Bitmask of all ranks except eighth.
	NOT_8TH_RANK uint64 = 0x00FFFFFFFFFFFFFF
	// Bitmask of the first rank.
	RANK_1 uint64 = 0xFF
	// Bitmask of the second rank.
	RANK_2 uint64 = 0xFF00
	// Bitmask of the seventh rank.
	RANK_7 uint64 = 0xFF000000000000
	// Bitmask of the eighth rank.
	RANK_8 uint64 = 0xFF00000000000000
)

// InitAttackTables initializes the predefined attack tables.
// Call this function ONCE as close as possible to the start of your program.
//
// NOTE: Move generation will not work if the attack tables are not initialized.
func InitAttackTables() {
	initBishopOccupancy()
	initRookOccupancy()

	for square := range 64 {
		bb := uint64(1 << square)

		pawnAttacks[ColorWhite][square] = genPawnAttacks(bb, ColorWhite)
		pawnAttacks[ColorBlack][square] = genPawnAttacks(bb, ColorBlack)

		knightAttacks[square] = genKnightAttacks(bb)

		kingAttacks[square] = genKingAttacks(bb)

		bitCount := bishopBitCount[square]
		for i := 0; i < 1<<bitCount; i++ {
			occupancy := genOccupancy(i, bitCount, bishopOccupancy[square])

			key := occupancy * bishopMagicNumbers[square] >> (64 - bitCount)

			bishopAttacks[square][key] = genBishopAttacks(bb, occupancy)
		}

		bitCount = rookBitCount[square]
		for i := 0; i < 1<<bitCount; i++ {
			occupancy := genOccupancy(i, bitCount, rookOccupancy[square])

			key := occupancy * rookMagicNumbers[square] >> (64 - bitCount)

			rookAttacks[square][key] = genRookAttacks(bb, occupancy)
		}
	}
}

// GenLegalMoves generates legal moves for the currently active color
// using copy-make approach.
func GenLegalMoves(p Position, l *MoveList) {
	l.LastMoveIndex = 0

	genKingMoves(p, l)

	if GenChecksCounter(p.Bitboards, 1^p.ActiveColor) > 2 {
		return
	}

	pseudoLegal := MoveList{}

	genPawnMoves(p, &pseudoLegal)

	genNormalMoves(p, &pseudoLegal)

	prev := p

	for i := range pseudoLegal.LastMoveIndex {

		p.MakeMove(pseudoLegal.Moves[i])

		if GenChecksCounter(p.Bitboards, 1^prev.ActiveColor) == 0 {
			l.Push(pseudoLegal.Moves[i])
		}

		p = prev
	}
}

// GenChecksCounter returns the number of the pieces of the
// specified color that are delivering a check to the enemy king.
func GenChecksCounter(bitboards [15]uint64, c Color) (cnt int) {
	king := bitScan(bitboards[PieceWKing+(1^c)])

	if pawnAttacks[1^c][king]&bitboards[PieceWPawn+c] != 0 {
		cnt++
	}

	if knightAttacks[king]&bitboards[PieceWKnight+c] != 0 {
		cnt++
	}

	if lookupBishopAttacks(king, bitboards[14])&bitboards[PieceWBishop+c] != 0 {
		cnt++
	}

	if lookupRookAttacks(king, bitboards[14])&bitboards[PieceWRook+c] != 0 {
		cnt++
	}

	if lookupQueenAttacks(king, bitboards[14])&bitboards[PieceWQueen+c] != 0 {
		cnt++
	}

	return cnt
}

// genKingMoves appends legal moves for the king on
// the given position to the specified move list.
// Handles special king move - castling.
func genKingMoves(p Position, l *MoveList) {
	kingBB := p.Bitboards[PieceWKing+p.ActiveColor]
	p.removePiece(PieceWKing+p.ActiveColor, kingBB)
	attacks := genAttacks(p.Bitboards, 1^p.ActiveColor)
	p.removePiece(PieceWKing+p.ActiveColor, kingBB)
	king := bitScan(kingBB)

	dests := kingAttacks[king] & (^attacks) & (^p.Bitboards[12+p.ActiveColor])

	for dests > 0 {
		l.Push(NewMove(popLSB(&dests), king, MoveNormal))
	}

	p.Bitboards[14] ^= kingBB
	// Handle castling.
	if p.ActiveColor == ColorWhite {
		if p.canCastle(CastlingWhiteShort, attacks, p.Bitboards[14]) &&
			p.Bitboards[PieceWRook]&H1 != 0 {
			l.Push(NewMove(SG1, king, MoveCastling))
		}
		if p.canCastle(CastlingWhiteLong, attacks, p.Bitboards[14]) &&
			p.Bitboards[PieceWRook]&A1 != 0 {
			l.Push(NewMove(SC1, king, MoveCastling))
		}
	} else {
		if p.canCastle(CastlingBlackShort, attacks, p.Bitboards[14]) &&
			p.Bitboards[PieceBRook]&H8 != 0 {
			l.Push(NewMove(SG8, king, MoveCastling))
		}
		if p.canCastle(CastlingBlackLong, attacks, p.Bitboards[14]) &&
			p.Bitboards[PieceBRook]&A8 != 0 {
			l.Push(NewMove(SC8, king, MoveCastling))
		}
	}
}

// genPawnMoves appends pseudo-legal moves for a pawns to the given move list.
// Handles special pawn move - en passant.
func genPawnMoves(p Position, l *MoveList) {
	occupancy := p.Bitboards[14]
	ep := uint64(0)
	if p.EPTarget > 0 {
		ep = 1 << p.EPTarget
	}
	enemies := p.Bitboards[12+(1^p.ActiveColor)]
	pawns := p.Bitboards[PieceWPawn+p.ActiveColor]

	// Determine movement direction.
	dir, initRank, promoRank := 8, RANK_2, RANK_8
	if p.ActiveColor == ColorBlack {
		dir = -8
		initRank = RANK_7
		promoRank = RANK_1
	}

	for pawns > 0 {
		pawn := popLSB(&pawns)
		square := uint64(1 << pawn)

		fwd, dblFwd := pawn+dir, pawn+2*dir
		// If the pawn can move forward.
		fwdBB := uint64(1 << fwd)
		if fwdBB&occupancy == 0 {
			// Check if the move is promotion.
			if fwdBB&promoRank != 0 {
				l.Push(NewPromotionMove(fwd, pawn, PromotionKnight))
				l.Push(NewPromotionMove(fwd, pawn, PromotionBishop))
				l.Push(NewPromotionMove(fwd, pawn, PromotionRook))
				l.Push(NewPromotionMove(fwd, pawn, PromotionQueen))
			} else {
				l.Push(NewMove(fwd, pawn, MoveNormal))
			}
			// If the pawn is standing on its initial rank and can move
			// double forward.
			if square&initRank != 0 && 1<<dblFwd&occupancy == 0 {
				l.Push(NewMove(dblFwd, pawn, MoveNormal))
			}
		}

		// Handle pawn attacks. Pawn can only capture enemy pieces
		// or the en passant target square.
		attacks := pawnAttacks[p.ActiveColor][pawn] & (enemies | ep)
		for attacks > 0 {
			to := popLSB(&attacks)
			// Handle capture promotion.
			if 1<<to&promoRank != 0 {
				l.Push(NewPromotionMove(to, pawn, PromotionKnight))
				l.Push(NewPromotionMove(to, pawn, PromotionBishop))
				l.Push(NewPromotionMove(to, pawn, PromotionRook))
				l.Push(NewPromotionMove(to, pawn, PromotionQueen))
			} else if 1<<to&ep != 0 {
				l.Push(NewMove(to, pawn, MoveEnPassant))
			} else {
				l.Push(NewMove(to, pawn, MoveNormal))
			}
		}
	}
}

// genPawnMoves appends pseudo-legal moves for knights, bishops,
// rooks, and queens to the given move list.
func genNormalMoves(p Position, l *MoveList) {
	c := p.ActiveColor
	allies := p.Bitboards[12+c]
	occupancy := p.Bitboards[14]

	for i := PieceWKnight + c; i <= PieceWQueen+c; i += 2 {
		pieces := p.Bitboards[i]
		for pieces > 0 {
			from := popLSB(&pieces)

			dests := uint64(0)
			switch i {
			case PieceWKnight, PieceBKnight:
				dests |= knightAttacks[from]
			case PieceWBishop, PieceBBishop:
				dests |= lookupBishopAttacks(from, occupancy)
			case PieceWRook, PieceBRook:
				dests |= lookupRookAttacks(from, occupancy)
			case PieceWQueen, PieceBQueen:
				dests |= lookupQueenAttacks(from, occupancy)
			}

			dests &= ^allies
			for dests > 0 {
				l.Push(NewMove(popLSB(&dests), from, MoveNormal))
			}
		}
	}
}

// genAttacks generates the bitboard of squares attacked
// by pieces of the specified color.
// The main purpose of this function is to generate a bitboard
// of squares to which the king is forbidden to move.
//
// NOTE: The king must be excluded from the occupancy (bitboards[14])
// bitboard to avoid blocking the attacks of slider pieces.
// Otherwise, the king may appear to be able to move into check.
func genAttacks(bitboards [15]uint64, c Color) (attacks uint64) {
	for i := PieceWBishop + c; i <= PieceWQueen+c; i += 2 {
		bitboard := bitboards[i]
		for bitboard > 0 {
			slider := popLSB(&bitboard)

			switch i {
			case PieceWBishop, PieceBBishop:
				attacks |= lookupBishopAttacks(slider, bitboards[14])
			case PieceWRook, PieceBRook:
				attacks |= lookupRookAttacks(slider, bitboards[14])
			case PieceWQueen, PieceBQueen:
				attacks |= lookupQueenAttacks(slider, bitboards[14])
			}
		}
	}

	//  Exclude empty squares and squares occupied by allied pieces.
	attacks |= genPawnAttacks(bitboards[PieceWPawn+c], c)
	// Exclude squares occupied by allied pieces.
	attacks |= genKnightAttacks(bitboards[PieceWKnight+c])
	//  Exclude squares occupied by allied pieces.
	attacks |= genKingAttacks(bitboards[PieceWKing+c])

	return attacks
}

// Use this function only to generate attacks for multiple pawns
// simultaneously. To get attacks for a single pawn, use the
// pawnAttacks lookup table.
func genPawnAttacks(pawn uint64, color Color) uint64 {
	if color == ColorWhite {
		return (pawn & NOT_A_FILE << 7) | (pawn & NOT_H_FILE << 9)
	}
	// Handle black pawns.
	return (pawn & NOT_A_FILE >> 9) | (pawn & NOT_H_FILE >> 7)
}

// genKnightAttacks returns a bitboard of squares attacked by knights.
//
// Use this function only to generate attacks for multiple knights
// simultaneously. To get attacks for a single knight, use the
// knightAttacks lookup table.
func genKnightAttacks(knight uint64) uint64 {
	return (knight & NOT_A_FILE >> 17) |
		(knight & NOT_H_FILE >> 15) |
		(knight & NOT_AB_FILE >> 10) |
		(knight & NOT_GH_FILE >> 6) |
		(knight & NOT_AB_FILE << 6) |
		(knight & NOT_GH_FILE << 10) |
		(knight & NOT_A_FILE << 15) |
		(knight & NOT_H_FILE << 17)
}

// genKingAttacks returns a bitboard of squares attacked by a king.
func genKingAttacks(king uint64) uint64 {
	return (king & NOT_A_FILE >> 9) |
		(king >> 8) |
		(king & NOT_H_FILE >> 7) |
		(king & NOT_A_FILE >> 1) |
		(king & NOT_H_FILE << 1) |
		(king & NOT_A_FILE << 7) |
		(king << 8) |
		(king & NOT_H_FILE << 9)
}

// genBishopAttacks returns a bitboard of squares
// attacked by a bishop. Occupied squares that block
// movement in each direction are taken into account.
// The resulting bitboard includes the occupied squares.
//
// This function cannot generate attacks for multiple bishops simultaneously.
func genBishopAttacks(bishop, occupancy uint64) (attacks uint64) {
	for i := bishop & NOT_A_FILE >> 9; i&NOT_H_FILE != 0; i >>= 9 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := bishop & NOT_H_FILE >> 7; i&NOT_A_FILE != 0; i >>= 7 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := bishop & NOT_A_FILE << 7; i&NOT_H_FILE != 0; i <<= 7 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := bishop & NOT_H_FILE << 9; i&NOT_A_FILE != 0; i <<= 9 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	return attacks
}

// genRookAttacks returns a bitboard of squares attacked by a rook.
// Occupied squares that block movement in each direction are
// taken into account.
// The resulting bitboard includes the occupied squares.
//
// This function cannot generate attacks for multiple rooks simultaneously.
func genRookAttacks(rook, occupancy uint64) (attacks uint64) {
	for i := rook & NOT_A_FILE >> 1; i&NOT_H_FILE != 0; i >>= 1 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := rook & NOT_H_FILE << 1; i&NOT_A_FILE != 0; i <<= 1 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := rook & NOT_1ST_RANK >> 8; i&NOT_8TH_RANK != 0; i >>= 8 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	for i := rook & NOT_8TH_RANK << 8; i&NOT_1ST_RANK != 0; i <<= 8 {
		attacks |= i
		if i&occupancy != 0 {
			break
		}
	}

	return attacks
}

// initBishopOccupancy initializes the lookup table
// of the "relevant occupancy squares" for a bishop.
// They are the only squares whose occupancy matters when
// generating legal moves of a bishop. This function is used
// to initialize a predefined array of bishop attacks using magic bitboards.
func initBishopOccupancy() {
	// Helper constants.
	const not_A_not_1st = NOT_A_FILE & NOT_1ST_RANK
	const not_H_not_1st = NOT_H_FILE & NOT_1ST_RANK
	const not_A_not_8th = NOT_A_FILE & NOT_8TH_RANK
	const not_H_not_8th = NOT_H_FILE & NOT_8TH_RANK

	for square := range 64 {
		var occupancy, bishop uint64 = 0, 1 << square

		for i := bishop & NOT_A_FILE >> 9; i&not_A_not_1st != 0; i >>= 9 {
			occupancy |= i
		}

		for i := bishop & NOT_H_FILE >> 7; i&not_H_not_1st != 0; i >>= 7 {
			occupancy |= i
		}

		for i := bishop & NOT_A_FILE << 7; i&not_A_not_8th != 0; i <<= 7 {
			occupancy |= i
		}

		for i := bishop & NOT_H_FILE << 9; i&not_H_not_8th != 0; i <<= 9 {
			occupancy |= i
		}

		bishopOccupancy[square] = occupancy
	}
}

// initRookOccupancy initializes the lookup table
// of the "relevant occupancy squares" for a rook.
// They are the only squares whose occupancy matters when
// generating legal moves of a rook. This function is used
// to initialize a predefined array of rook attacks using magic bitboards.
func initRookOccupancy() {
	for square := range 64 {
		var occupancy, rook uint64 = 0, 1 << square

		for i := rook & NOT_1ST_RANK >> 8; i&NOT_1ST_RANK != 0; i >>= 8 {
			occupancy |= i
		}

		for i := rook & NOT_A_FILE >> 1; i&NOT_A_FILE != 0; i >>= 1 {
			occupancy |= i
		}

		for i := rook & NOT_H_FILE << 1; i&NOT_H_FILE != 0; i <<= 1 {
			occupancy |= i
		}

		for i := rook & NOT_8TH_RANK << 8; i&NOT_8TH_RANK != 0; i <<= 8 {
			occupancy |= i
		}

		rookOccupancy[square] = occupancy
	}
}

// genOccupancy returns a bitboard of blocker pieces
// for the specified attack bitboard.
func genOccupancy(key, relevantBitCount int,
	relevantOccupancy uint64) (occupancy uint64) {

	for i := 0; i < relevantBitCount; i++ {
		square := popLSB(&relevantOccupancy)

		if key&(1<<i) != 0 {
			occupancy |= 1 << square
		}
	}

	return occupancy
}

// lookupBishopAttacks returns a bitboard of squares attacked by a bishop.
// The bitboard is taken from the bishopAttacks using magic hashing scheme.
func lookupBishopAttacks(square int, occupancy uint64) uint64 {
	occupancy &= bishopOccupancy[square]
	occupancy *= bishopMagicNumbers[square]
	occupancy >>= 64 - bishopBitCount[square]
	return bishopAttacks[square][occupancy]
}

// lookupRookAttacks returns a bitboard of squares attacked by a rook.
// The bitboard is taken from the rookAttacks using magic hashing scheme.
func lookupRookAttacks(square int, occupancy uint64) uint64 {
	occupancy &= rookOccupancy[square]
	occupancy *= rookMagicNumbers[square]
	occupancy >>= 64 - rookBitCount[square]
	return rookAttacks[square][occupancy]
}

// lookupQueenAttacks returns a bitboard of squares attacked by a queen.
// The bitboard is calculated as the logical disjunction
// of the bishop and rook attack bitboards.
func lookupQueenAttacks(square int, occupancy uint64) uint64 {
	return lookupBishopAttacks(square, occupancy) |
		lookupRookAttacks(square, occupancy)
}
