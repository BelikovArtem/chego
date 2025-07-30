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

var (
	// bishopMagicNumbers is a precalculated lookup table of magic
	// numbers for a bishop.
	// See commit c0cfb607e20a1e469d2fc94b26146645bc1fc9a1 for details.
	bishopMagicNumbers = [64]uint64{
		0x11410121040100,
		0x2084820928010,
		0xa010208481080040,
		0x214240082000610,
		0x4d104000400480,
		0x1012010804408,
		0x42044101452000c,
		0x2844804050104880,
		0x814204290a0a00,
		0x10280688224500,
		0x1080410101010084,
		0x10020a108408004,
		0x2482020210c80080,
		0x480104a0040400,
		0x411006404200810,
		0x1024010908024292,
		0x1004401001011a,
		0x810006081220080,
		0x1040404206004100,
		0x58080000820041ce,
		0x3406000422010890,
		0x1a004100520210,
		0x202a000048040400,
		0x225004441180110,
		0x8064240102240,
		0x1424200404010402,
		0x1041100041024200,
		0x8082002012008200,
		0x1010008104000,
		0x8808004000806000,
		0x380a000080c400,
		0x31040100042d0101,
		0x110109008082220,
		0x4010880204201,
		0x4006462082100300,
		0x4002010040140041,
		0x40090200250880,
		0x2010100c40c08040,
		0x12800ac01910104,
		0x10b20051020100,
		0x210894104828c000,
		0x50440220004800,
		0x1002011044180800,
		0x4220404010410204,
		0x1002204a2020401,
		0x21021001000210,
		0x4880081009402,
		0xc208088c088e0040,
		0x4188464200080,
		0x3810440618022200,
		0xc020310401040420,
		0x2000008208800e0,
		0x4c910240020,
		0x425100a8602a0,
		0x20c4206a0c030510,
		0x4c10010801184000,
		0x200202020a026200,
		0x6000004400841080,
		0xc14004121082200,
		0x400324804208800,
		0x1802200040504100,
		0x1820000848488820,
		0x8620682a908400,
		0x8010600084204240,
	}
	// rookMagicNumbers is a precalculated lookup table of magic
	// numbers for a rook.
	// See commit c0cfb607e20a1e469d2fc94b26146645bc1fc9a1 for details.
	rookMagicNumbers = [64]uint64{
		0x2080008040002010,
		0x40200010004000,
		0x100090010200040,
		0x2080080010000480,
		0x880040080080102,
		0x8200106200042108,
		0x410041000408b200,
		0x100009a00402100,
		0x5800800020804000,
		0x848404010002000,
		0x101001820010041,
		0x10a0040100420080,
		0x8a02002006001008,
		0x926000844110200,
		0x8000800200800100,
		0x28060001008c2042,
		0x10818002204000,
		0x10004020004001,
		0x110002008002400,
		0x11a020010082040,
		0x2001010008000410,
		0x42010100080400,
		0x4004040008020110,
		0x820000840041,
		0x400080208000,
		0x2080200040005000,
		0x8000200080100080,
		0x4400080180500080,
		0x4900080080040080,
		0x4004004480020080,
		0x8006000200040108,
		0xc481000100006396,
		0x1000400080800020,
		0x201004400040,
		0x10008010802000,
		0x204012000a00,
		0x800400800802,
		0x284000200800480,
		0x3000403000200,
		0x840a6000514,
		0x4080c000228012,
		0x10002000444010,
		0x620001000808020,
		0xc210010010009,
		0x100c001008010100,
		0xc10020004008080,
		0x20100802040001,
		0x808008305420014,
		0xc010800840043080,
		0x208401020890100,
		0x10b0081020028280,
		0x6087001001220900,
		0xc080011000500,
		0x9810200040080,
		0x2000010882100400,
		0x2000050880540200,
		0x800020104200810a,
		0x6220250242008016,
		0x9180402202900a,
		0x40210500100009,
		0x6000814102026,
		0x410100080a040013,
		0x10405008022d1184,
		0x1000009400410822,
	}

	// Precalculated lookup tables used to speed up
	// the move generation process.

	// Pawn's attack pattern depends on the color,
	// so it is necessary to store two tables.
	pawnAttacks     [2][64]uint64
	knightAttacks   [64]uint64
	kingAttacks     [64]uint64
	bishopOccupancy [64]uint64
	rookOccupancy   [64]uint64
	// Lookup bishop attack table for every possible
	// combination of square/occupancy.
	bishopAttacks [64][512]uint64
	// Lookup rook attack table for every possible
	// combination of square/occupancy.
	rookAttacks [64][4096]uint64
	// Precalculated lookup table of bishop relevant occupancy
	// bit count for every square.
	bishopBitCount = [64]int{
		6, 5, 5, 5, 5, 5, 5, 6,
		5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5,
		6, 5, 5, 5, 5, 5, 5, 6,
	}
	// Precalculated lookup table of rook relevant occupancy
	// bit count for every square.
	rookBitCount = [64]int{
		12, 11, 11, 11, 11, 11, 11, 12,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		12, 11, 11, 11, 11, 11, 11, 12,
	}
	// Each path includes the king square.
	// 0 : White O-O castling path.
	// 1 : White O-O-O castling path.
	// 2 : Black O-O castling path.
	// 3 : Black O-O-O castling path.
	castlingPath = [4]uint64{
		0x70, 0x1C, 0x7000000000000000, 0x1C00000000000000,
	}
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

// GenLegalMoves generates legal moves for the currently active color.
// It generetes strictly legal moves at the first place, since filtering
// down pseudo-legal moves is an expensive task.
func GenLegalMoves(p Position, l *MoveList) {
	l.LastMoveIndex = 0

	off := 6 * p.ActiveColor
	eoff := 6 * (1 ^ p.ActiveColor)

	king := bitScan(p.Bitboards[PieceWKing+off])

	// 1. Generate king legal moves, since if we are in double
	// check, nobody except the king can move.
	genKingMoves(&p, l)

	pieces := GenCheckingPieces(p.Bitboards, 1^p.ActiveColor)
	// Only king can move if there is a double check.
	// See https://en.wikipedia.org/wiki/Double_check
	if pieces > 1 {
		return
	}

	legalDests := uint64(ALL_SQUARES)

	if pieces == 1 {
		// The are 3 types of legal moves in a single check position:
		//  1. Move the king out of check.
		//  2. Capture the checking piece.
		//  3. Block the checking piece (if being checked by a sliding piece).
		legalDests = pieces

		piece := p.GetPieceFromSquare(pieces)
		// If being checked by a sliding piece.
		eoff := 6 * (1 ^ p.ActiveColor)
		if piece >= PieceWBishop+eoff && piece <= PieceWQueen+eoff {
			checker := bitScan(pieces)
			// Add squares which help to block the checking piece.
			legalDests |= genAttackRay(king, checker)
		}
	}

	// 2. Generate pinned pieces legal moves.
	pinned := genPinnedPieces(p.Bitboards, p.ActiveColor)
	for pinned > 0 {
		from := popLSB(&pinned)
		square := uint64(1 << from)
		piece := p.GetPieceFromSquare(square)

		// Remove the pinned piece from the board.
		p.Bitboards[piece] ^= square

		// Generate the attacks from the enemy's pinning piece.
		checkers := genCheckingPieces(p.Bitboards, 1^p.ActiveColor)
		for checkers > 0 {
			checker := popLSB(&checkers)
			ray := genAttackRay(king, checker)

			// Append legal moves for the pinned piece.
			for ray > 0 {
				l.Push(NewMove(popLSB(&ray), from, MoveNormal))
			}
		}

		// Restore the board state.
		p.Bitboards[piece] |= square
	}

	// 3. Generate moves for non-pinned pieces.
	allies := p.Bitboards[12+p.ActiveColor]
	occupancy := p.Bitboards[14]

	// Genereate legal moves for pawns, excluding the pinned ones.
	pawns := p.Bitboards[PieceWPawn+off] & (^pinned)
	for pawns > 0 {
		from := popLSB(&pawns)
		genPawnMoves(p.Bitboards, from, 1<<p.EPTarget, p.ActiveColor, l)
	}

	// Generate legal moves for knights, bishops, rooks, and queens,
	// excluding the pinned ones.
	for i := PieceWKnight + off; i <= PieceWQueen+off; i++ {
		pieces := p.Bitboards[i] & (^pinned)
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
			// Take into accout only legal destinations if the
			// king is being checked.
			dests &= legalDests

			for dests > 0 {
				to := popLSB(&dests)
				if 1<<to&allies == 0 {
					l.Push(NewMove(to, from, MoveNormal))
				}
			}
			genNormalMoves(from, i, allies, occupancy, l)
		}
	}
}

// genKingMoves appends legal moves for the king on
// the given position to the specified move list.
// Handles special king move - castling.
func genKingMoves(p *Position, l *MoveList) {
	c := p.ActiveColor

	attacks := genAttacks(p.Bitboards, c^1)

	king := bitScan(p.Bitboards[PieceWKing+6*c])
	occupancy := p.Bitboards[14] & ^(1 << king)
	allies := p.Bitboards[12+c]

	// Exclude the attacked squares to forbit
	// the king to move on them.
	dests := kingAttacks[king] & (^attacks) & (^allies)
	for dests > 0 {
		to := popLSB(&dests)
		l.Push(NewMove(to, king, MoveNormal))
	}

	// Handle castling.
	if c == ColorWhite {
		if p.CanCastle(CastlingWhiteShort, attacks, occupancy) {
			l.Push(NewMove(SG1, king, MoveCastling))
		}
		if p.CanCastle(CastlingWhiteLong, attacks, occupancy) {
			l.Push(NewMove(SC1, king, MoveCastling))
		}
	} else {
		if p.CanCastle(CastlingBlackShort, attacks, occupancy) {
			l.Push(NewMove(SG8, king, MoveCastling))
		}
		if p.CanCastle(CastlingBlackLong, attacks, occupancy) {
			l.Push(NewMove(SC1, king, MoveCastling))
		}
	}
}

// genPawnMoves appends legal moves for a pawn to the given move list.
// Handles special pawn move - en passant. Pinned pawns must be excluded
// before calling this function.
func genPawnMoves(bitboards [15]uint64, pawn int, epTarget uint64,
	c Color, l *MoveList) {
	enemies := bitboards[12+(1^c)]
	occupancy := bitboards[14]
	square := uint64(1 << pawn)
	// Determine movement direction.
	dir, initRank, promoRank := 8, RANK_2, RANK_8
	if c == ColorBlack {
		dir = -8
		initRank = RANK_7
		promoRank = RANK_8
	}

	fwd, dblFwd := pawn+dir, pawn+2*dir

	// If the pawn can move forward.
	fwdBB := uint64(1 << fwd)
	if fwdBB&occupancy == 0 {
		// Check if the move is promotion.
		if fwdBB&promoRank != 0 {
			l.Push(NewMove(fwd, pawn, MovePromotion))
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
	attacks := pawnAttacks[c][pawn] & (enemies | epTarget)
	for attacks > 0 {
		to := popLSB(&attacks)
		// Handle capture promotion.
		if 1<<to&promoRank != 0 {
			l.Push(NewMove(to, pawn, MovePromotion))
		} else if 1<<to&epTarget != 0 {
			// Making en passant move may expose the discovered check.
			// There may be two pawns, so the genPinnedPieces won't determine
			// that there is a pin.

			// Remove the moving pawn and captured pawn from the board.
			// If there are no checkers, add the move.
			if genCheckingPieces(bitboards, c) == 0 {
				l.Push(NewMove(to, pawn, MoveEnPassant))
			}
		} else {
			l.Push(NewMove(to, pawn, MoveNormal))
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
func genAttacks(bitboards [15]uint64, c Color) uint64 {
	return genSliderAttacks(bitboards, c) | genLeaperAttacks(bitboards, c)
}

// GenCheckingPieces generates the bitboard of all pieces of the
// specified color that are delivering a check to the enemy king.
func GenCheckingPieces(bitboards [15]uint64, c Color) (checkers uint64) {
	occupancy := bitboards[14]
	// Color offset.
	off := 6 * c

	king := bitScan(bitboards[PieceWKing+6*(1^c)])

	checkers |= pawnAttacks[c][king] & bitboards[PieceWPawn+off]

	checkers |= knightAttacks[king] & bitboards[PieceWKnight+off]

	checkers |= lookupBishopAttacks(king, occupancy) &
		bitboards[PieceWBishop+off]

	checkers |= lookupRookAttacks(king, occupancy) &
		bitboards[PieceWRook+off]

	checkers |= lookupQueenAttacks(king, occupancy) &
		bitboards[PieceWQueen+off]

	return checkers
}

// genAttackRay generates the bitboard of squares between
// the checker and the king, excluding the king and the checking piece.
// The main purpose of this function is to generate a bitboard
// of squares that pieces can move to in order to evade check.
func genAttackRay(king int, checker int) (ray uint64) {
	dstFile, dstRank := king%8, king/8
	srcFile, srcRank := checker%8, checker/8
	// Determine ray direction.
	dirRank := 0
	if dstRank > srcRank {
		dirRank = 1
	} else if dstRank < srcRank {
		dirRank = -1
	}
	dirFile := 0
	if dstFile > srcFile {
		dirFile = 1
	} else {
		dirFile = -1
	}

	dst := 8*(dstRank-dirRank) + (dstFile - dirFile)
	square := 8*(srcRank+dirRank) + (srcFile + dirFile)

	for square != dst {
		ray |= 1 << square
		srcRank += dirRank
		srcFile += dirFile
		square = 8*(srcRank+dirRank) + (srcFile + dirFile)
	}

	return ray
}

// genPinnedPiece generates the bitboard of all pinned pieces of the
// specified color. The main purpose of this function is to generate
// a bitboard of pieces which can expose the king to check while moving.
// TODO: fix
func genPinnedPieces(bitboards [15]uint64, c Color) uint64 {
	// Generate attacks for the enemy sliding pieces.
	attacks := genSliderAttacks(bitboards, 1^c)

	// Generate sliding pseudo attacks from the allied king position.
	king := bitboards[PieceWKing+6*c]
	pseudoAttacks := lookupQueenAttacks(bitScan(king), bitboards[14])

	// Calculate the intersection between the attack rays which
	// land on the allied pieces.
	intersect := attacks & pseudoAttacks & bitboards[12+c]
	if intersect != 0 {
		// Determine the direction of the attack ray.
		pin := bitScan(intersect)
		kingSq := bitScan(king)
		pinnerSq := bitScan()

		dir := dstRank - d
		if dstRank {

		}

	}

	return intersect
}

func genSliderAttacks(bitboards [15]uint64, c Color) (attacks uint64) {
	off := 6 * c

	for i := PieceWBishop + off; i <= PieceWQueen+off; i++ {
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

	return attacks
}

func genLeaperAttacks(bitboards [15]uint64, c Color) (attacks uint64) {
	off := 6 * c
	allies := bitboards[12+c]
	enemies := bitboards[12+(1^c)]

	//  Exclude empty squares and squares occupied by allied pieces.
	attacks |= genPawnAttacks(bitboards[PieceWPawn+off], c) & enemies
	// Exclude squares occupied by allied pieces.
	attacks |= genKnightAttacks(bitboards[PieceWKnight+off]) & (^allies)
	//  Exclude squares occupied by allied pieces.
	attacks |= genKingAttacks(bitboards[PieceWKing+off]) & (^allies)

	return attacks
}

// genPawnAttacks returns a bitboard of squares attacked by pawns.
//
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
