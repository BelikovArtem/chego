// Package movegen implements move generation using Magic Bitboards approach.
package movegen

import (
	"github.com/BelikovArtem/chego/bitutil"
	"github.com/BelikovArtem/chego/types"
)

// The following constants are used to prevent bitboard overflow during move generation.
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
	// White O-O castling path. Includes the king square.
	WHITE_KING_CASTLING_PATH uint64 = 0x70
	// White O-O-O castling path. Includes the king square.
	WHITE_QUEEN_CASTLING_PATH uint64 = 0x1C
	// Black O-O castling path. Includes the king square.
	BLACK_KING_CASTLING_PATH uint64 = 0x7000000000000000
	// Black O-O-O castling path. Includes the king square.
	BLACK_QUEEN_CASTLING_PATH uint64 = 0x1C00000000000000
)

var (
	// bishopMagicNumbers is a precalculated lookup table of magic numbers for a bishop.
	// Generated by the GenMagicNumber function.
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
	// rookMagicNumbers is a precalculated lookup table of magic numbers for a rook.
	// Generated by the GenMagicNumber function.
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
	// Precalculated lookup tables used to speed up the move generation process.
	// Pawn's attack pattern depends on the color, so it is necessary to store two tables.
	pawnAttacks             [2][64]uint64
	knightAttacks           [64]uint64
	kingAttacks             [64]uint64
	bishopRelevantOccupancy [64]uint64
	rookRelevantOccupancy   [64]uint64
	// Lookup bishop attack table for every possible combination of square/occupancy.
	bishopAttacks [64][512]uint64
	// Lookup rook attack table for every possible combination of square/occupancy.
	rookAttacks [64][4096]uint64
	// Precalculated lookup table of bishop relevant occupancy bit count for every square.
	bishopRelevantOccupancyBitCount = [64]int{
		6, 5, 5, 5, 5, 5, 5, 6,
		5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5,
		6, 5, 5, 5, 5, 5, 5, 6,
	}
	// Precalculated lookup table of rook relevant occupancy bit count for every square.
	rookRelevantOccupancyBitCount = [64]int{
		12, 11, 11, 11, 11, 11, 11, 12,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		12, 11, 11, 11, 11, 11, 11, 12,
	}
)

// InitAttackTables initializes the predefined attack tables.
// Call this function ONCE as close as possible to the start of your program.
//
// NOTE: Move generation will not work if the attack tables are not initialized.
func InitAttackTables() {
	initBishopRelevantOccupancy()
	initRookRelevantOccupancy()

	for square := 0; square < 64; square++ {
		bb := uint64(1 << square)

		pawnAttacks[types.ColorWhite][square] = genPawnAttacks(bb, types.ColorWhite)
		pawnAttacks[types.ColorBlack][square] = genPawnAttacks(bb, types.ColorBlack)

		knightAttacks[square] = genKnightAttacks(bb)

		kingAttacks[square] = genKingAttacks(bb)

		bitCount := bishopRelevantOccupancyBitCount[square]
		for i := 0; i < 1<<bitCount; i++ {
			occupancy := genOccupancy(i, bitCount, bishopRelevantOccupancy[square])

			magicKey := occupancy * bishopMagicNumbers[square] >> (64 - bitCount)

			bishopAttacks[square][magicKey] = genBishopAttacks(bb, occupancy)
		}

		bitCount = rookRelevantOccupancyBitCount[square]
		for i := 0; i < 1<<bitCount; i++ {
			occupancy := genOccupancy(i, bitCount, rookRelevantOccupancy[square])

			magicKey := occupancy * rookMagicNumbers[square] >> (64 - bitCount)

			rookAttacks[square][magicKey] = genRookAttacks(bb, occupancy)
		}
	}
}

// IsSquareUnderAttack checks whether the specified square is attacked
// by pieces of the specified color in the given position using attack tables.
func IsSquareUnderAttack(bitboards [15]uint64, square int, c types.Color) bool {
	occupancy := bitboards[14]

	offs := 6 * c
	// If attacked by pawns.
	return (pawnAttacks[c^1][square]&bitboards[offs+types.PieceWPawn] != 0) ||
		// If attacked by knights.
		(knightAttacks[square]&bitboards[offs+types.PieceWKnight] != 0) ||
		// If attacked by bishops.
		(lookupBishopAttacks(square, occupancy)&bitboards[offs+types.PieceWBishop] != 0) ||
		// If attacked by rooks.
		(lookupRookAttacks(square, occupancy)&bitboards[offs+types.PieceWRook] != 0) ||
		// If attacked by queens.
		(lookupQueenAttacks(square, occupancy)&bitboards[offs+types.PieceWQueen] != 0) ||
		// If attacked by king.
		(kingAttacks[square]&bitboards[offs+types.PieceWKing] != 0)
}

// GenLegalMoves generates legal moves for the currently active color.
//
// This involves two steps:
//  1. Generate pseudo-legal moves for all pieces of the active color.
//  2. Filter out illegal moves using the copy-make approach.
func GenLegalMoves(p types.Position, l *types.MoveList) {
	l.LastMoveIndex = 0
	pseudoLegal := types.MoveList{}

	genPseudoLegalMoves(p, &pseudoLegal)

	c := p.ActiveColor
	cOff := 6 * p.ActiveColor

	// Memorize board state.
	var bitboards [15]uint64
	ep := p.EPTarget
	cr := p.CastlingRights
	copy(bitboards[:], p.Bitboards[:])

	for i := 0; i < int(pseudoLegal.LastMoveIndex); i++ {
		m := pseudoLegal.Moves[i]

		p.MakeMove(m)

		// If the allied king is not in check, the move is legal.
		king := bitutil.BitScan(p.Bitboards[types.PieceWKing+cOff])
		if !IsSquareUnderAttack(p.Bitboards, king, 1^c) {
			l.Push(pseudoLegal.Moves[i])
		}

		// Restore board state.
		p.EPTarget = ep
		p.CastlingRights = cr
		copy(p.Bitboards[:], bitboards[:])
	}
}

func genPseudoLegalMoves(p types.Position, l *types.MoveList) {
	allies := p.Bitboards[12+p.ActiveColor]
	enemies := p.Bitboards[12+(1^p.ActiveColor)]
	occupancy := p.Bitboards[14]
	cOff := 6 * p.ActiveColor

	pawns := p.Bitboards[types.PieceWPawn+cOff]
	for pawns > 0 {
		pawn := bitutil.PopLSB(&pawns)
		genPawnMoves(pawn, enemies, occupancy, 1<<p.EPTarget, p.ActiveColor, l)
	}

	for i := types.PieceWKnight + cOff; i <= types.PieceWQueen+cOff; i++ {
		pieces := p.Bitboards[i]
		for pieces > 0 {
			piece := bitutil.PopLSB(&pieces)
			genNormalPieceMoves(piece, i, allies, occupancy, l)
		}
	}
	king := bitutil.BitScan(p.Bitboards[types.PieceWKing+cOff])
	genKingMoves(p, king, l)
}

// genPawnMoves generates pseudo-legal moves for a pawn and appends them
// to the given move list.
func genPawnMoves(pawn int, enemies, occupancy, epTarget uint64, c types.Color,
	l *types.MoveList) {
	bb := uint64(1 << pawn)
	delta, initRank, promoRank := 8, RANK_2, RANK_8
	if c == types.ColorBlack {
		delta = -8
		initRank = RANK_7
		promoRank = RANK_1
	}

	fwd, dblFwd := pawn+delta, pawn+2*delta

	// If the pawn can move forward.
	var fwdBB uint64 = 1 << fwd
	if fwdBB&occupancy == 0 {
		// Check if the move is promotion.
		if fwdBB&promoRank != 0 {
			l.Push(types.NewMove(fwd, pawn, types.MovePromotion))
		} else {
			l.Push(types.NewMove(fwd, pawn, types.MoveNormal))
		}

		// If the pawn is standing on its initial rank and can move double forward.
		if bb&initRank != 0 && 1<<dblFwd&occupancy == 0 {
			l.Push(types.NewMove(dblFwd, pawn, types.MoveNormal))
		}
	}

	// Handle pawn attacks.
	attacks := pawnAttacks[c][pawn] & (enemies | epTarget)
	for attacks > 0 {
		atackSq := bitutil.PopLSB(&attacks)

		if 1<<atackSq&promoRank != 0 {
			l.Push(types.NewMove(atackSq, pawn, types.MovePromotion))
		} else if 1<<atackSq&epTarget != 0 {
			l.Push(types.NewMove(atackSq, pawn, types.MoveEnPassant))
		} else {
			l.Push(types.NewMove(atackSq, pawn, types.MoveNormal))
		}
	}
}

// genNormalPieceMoves generates pseudo-legal moves for a piece that can't
// perform special moves (knights, bishops, rooks, and queens) and appends
// them to the given move list.
func genNormalPieceMoves(piece int, pieceType types.Piece, allies, occupancy uint64,
	l *types.MoveList) {
	var dests uint64

	switch pieceType {
	case types.PieceWKnight, types.PieceBKnight:
		dests = knightAttacks[piece]

	case types.PieceWBishop, types.PieceBBishop:
		dests = lookupBishopAttacks(piece, occupancy)

	case types.PieceWRook, types.PieceBRook:
		dests = lookupRookAttacks(piece, occupancy)

	case types.PieceWQueen, types.PieceBQueen:
		dests = lookupQueenAttacks(piece, occupancy)
	}

	for dests > 0 {
		square := bitutil.PopLSB(&dests)
		if (1<<square&occupancy != 0 &&
			1<<square&allies == 0) || 1<<square&occupancy == 0 {
			l.Push(types.NewMove(square, piece, types.MoveNormal))
		}
	}
}

// genKingMoves appends pseudo-legal moves (quiet moves and captures) for the king on
// the given position to the specified move list.
func genKingMoves(p types.Position, king int, l *types.MoveList) {
	enemies := p.Bitboards[12+(1^p.ActiveColor)]
	occupancy := p.Bitboards[14]
	// Lookup normal moves.
	attacks := kingAttacks[king]
	for attacks > 0 {
		square := bitutil.PopLSB(&attacks)
		if 1<<square&enemies != 0 || 1<<square&occupancy == 0 {
			l.Push(types.NewMove(square, king, types.MoveNormal))
		}
	}

	attacked := genAttacks(p.Bitboards, 1^p.ActiveColor)
	occupancy ^= 1 << king

	// Handle castling.
	if p.ActiveColor == types.ColorWhite {
		if p.CastlingRights&types.CastlingWhiteShort != 0 &&
			occupancy&WHITE_KING_CASTLING_PATH == 0 &&
			attacked&WHITE_KING_CASTLING_PATH == 0 {
			// White O-O.
			l.Push(types.NewMove(types.SG1, king, types.MoveCastling))
		}
		if p.CastlingRights&types.CastlingWhiteLong != 0 &&
			occupancy&WHITE_QUEEN_CASTLING_PATH == 0 &&
			attacked&WHITE_QUEEN_CASTLING_PATH == 0 {
			// White O-O-O.
			l.Push(types.NewMove(types.SC1, king, types.MoveCastling))
		}
	} else {
		if p.CastlingRights&types.CastlingBlackShort != 0 &&
			occupancy&BLACK_KING_CASTLING_PATH == 0 &&
			attacked&BLACK_KING_CASTLING_PATH == 0 {
			// Black O-O.
			l.Push(types.NewMove(types.SG8, king, types.MoveCastling))
		}
		if p.CastlingRights&types.CastlingBlackLong != 0 &&
			occupancy&BLACK_QUEEN_CASTLING_PATH == 0 &&
			attacked&BLACK_QUEEN_CASTLING_PATH == 0 {
			// Black O-O-O.
			l.Push(types.NewMove(types.SC8, king, types.MoveCastling))
		}
	}
}

// genAttacks generates a bitboard of squares attacked by a pieces of the specified color.
func genAttacks(bitboards [15]uint64, c types.Color) uint64 {
	occupancy := bitboards[14]

	attacked := genPawnAttacks(bitboards[types.PieceWPawn+6*c], c)
	attacked |= genKnightAttacks(bitboards[types.PieceWKnight+6*c])
	attacked |= genKingAttacks(bitboards[types.PieceWKing+6*c])

	for i := types.PieceWBishop + 6*c; i <= types.PieceWKing+6*c; i++ {
		for bitboards[i] > 0 {
			square := bitutil.PopLSB(&bitboards[i])

			switch i {
			case types.PieceWBishop, types.PieceBBishop:
				attacked |= lookupBishopAttacks(square, occupancy)
			case types.PieceWRook, types.PieceBRook:
				attacked |= lookupRookAttacks(square, occupancy)
			case types.PieceWQueen, types.PieceBQueen:
				attacked |= lookupQueenAttacks(square, occupancy)
			}
		}
	}

	return attacked
}

// genPawnAttacks returns a bitboard of squares attacked by a pawn.
func genPawnAttacks(pawn uint64, color types.Color) uint64 {
	if color == types.ColorWhite {
		return (pawn & NOT_A_FILE << 7) | (pawn & NOT_H_FILE << 9)
	}
	// Handle black pawns.
	return (pawn & NOT_A_FILE >> 9) | (pawn & NOT_H_FILE >> 7)
}

// genKnightAttacks returns a bitboard of squares attacked by a knight.
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

// genBishopAttacks returns a bitboard of squares attacked by a bishop.
// Occupied squares that block movement in each direction are taken into account.
// The resulting bitboard includes the occupied squares.
func genBishopAttacks(bishop, occupancy uint64) uint64 {
	var attacks uint64

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
// Occupied squares that block movement in each direction are taken into account.
// The resulting bitboard includes the occupied squares.
func genRookAttacks(rook, occupancy uint64) uint64 {
	var attacks uint64

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

// initBishopRelevantOccupancy initializes the lookup table of the "relevant occupancy squares" for a bishop.
// They are the only squares whose occupancy matters when generating legal moves of a bishop.
// This function is used to initialize a predefined array of bishop attacks using magic bitboards.
func initBishopRelevantOccupancy() {
	// Helper constants.
	const not_A_not_1st = NOT_A_FILE & NOT_1ST_RANK
	const not_H_not_1st = NOT_H_FILE & NOT_1ST_RANK
	const not_A_not_8th = NOT_A_FILE & NOT_8TH_RANK
	const not_H_not_8th = NOT_H_FILE & NOT_8TH_RANK

	for square := 0; square < 64; square++ {
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

		bishopRelevantOccupancy[square] = occupancy
	}
}

// initRookRelevantOccupancy initializes the lookup table of the "relevant occupancy squares" for a rook.
// They are the only squares whose occupancy matters when generating legal moves of a rook.
// This function is used to initialize a predefined array of rook attacks using magic bitboards.
func initRookRelevantOccupancy() {
	for square := 0; square < 64; square++ {
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

		rookRelevantOccupancy[square] = occupancy
	}
}

// genOccupancy returns a bitboard of blocker pieces for the specified attack bitboard.
//
//	key <= (2^relevantBitsCount) - 1
//	relevantBitsCount = (Rook | Bishop)RelevantOccupancyBitCount[square]
//	relevantOccupancy = (Rook | Bishop)RelevantOccupancy[square]
func genOccupancy(key, relevantBitCount int, relevantOccupancy uint64) uint64 {
	var occupancy uint64

	for i := 0; i < relevantBitCount; i++ {
		square := bitutil.PopLSB(&relevantOccupancy)

		if key&(1<<i) != 0 {
			occupancy |= 1 << square
		}
	}

	return occupancy
}

// lookupBishopAttacks returns a bitboard of squares attacked by a bishop.
// The bitboard is taken from the BishopAttacks using magic hashing scheme.
func lookupBishopAttacks(square int, occupancy uint64) uint64 {
	occupancy &= bishopRelevantOccupancy[square]
	occupancy *= bishopMagicNumbers[square]
	occupancy >>= 64 - bishopRelevantOccupancyBitCount[square]
	return bishopAttacks[square][occupancy]
}

// lookupRookAttacks returns a bitboard of squares attacked by a rook.
// The bitboard is taken from the RookAttacks using magic hashing scheme.
func lookupRookAttacks(square int, occupancy uint64) uint64 {
	occupancy &= rookRelevantOccupancy[square]
	occupancy *= rookMagicNumbers[square]
	occupancy >>= 64 - rookRelevantOccupancyBitCount[square]
	return rookAttacks[square][occupancy]
}

// lookupQueenAttacks returns a bitboard of squares attacked by a queen.
// The bitboard is calculated as the logical disjunction of the bishop and rook attack bitboards.
func lookupQueenAttacks(square int, occupancy uint64) uint64 {
	return lookupBishopAttacks(square, occupancy) | lookupRookAttacks(square, occupancy)
}
