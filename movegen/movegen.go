// Package movegen implements move generation using Magic Bitboards approach.
package movegen

import (
	"math/rand/v2"

	"github.com/BelikovArtem/chego/bitutil"
	"github.com/BelikovArtem/chego/enum"
)

// Move represents a chess move, encoded as a 16 bit unsigned integer:
//
//	0-5:   To (destination) square index;
//	6-11:  From (origin/source) square index;
//	12-13: Promotion piece (see [enum.PromotionFlag]);
//	14-15: Move type (see [enum.MoveType]).
type Move uint16

func NewMove(to, from, promotionPiece, moveType int) Move {
	return Move(to | (from << 6) | (promotionPiece << 12) | (moveType << 14))
}
func (m Move) To() int                            { return int(m & 0x3F) }
func (m Move) From() int                          { return int(m>>6) & 0x3F }
func (m Move) PromotionPiece() enum.PromotionFlag { return enum.PromotionFlag(m>>12) & 0x3 }
func (m Move) Type() enum.MoveType                { return enum.MoveType(m>>14) & 0x3 }

// MoveList is used to store moves.
type MoveList struct {
	// Maximum number of moves per chess position is equal to 218, hence 218 elements.
	Moves [218]Move
	// To keep track of the next move index.
	LastMoveIndex byte
}

// Push adds the move to the end of the move list.
func (ml *MoveList) Push(move Move) {
	ml.Moves[ml.LastMoveIndex] = move
	ml.LastMoveIndex++
}

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
	// bishopMagicNumberLookup is a precalculated lookup table of magic numbers for a bishop.
	// Generated by the GenMagicNumber function.
	bishopMagicNumberLookup = [64]uint64{
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
	// rookMagicNumberLookup is a precalculated lookup table of magic numbers for a rook.
	// Generated by the GenMagicNumber function.
	rookMagicNumberLookup = [64]uint64{
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
		var squareBB uint64 = 1 << square

		pawnAttacks[enum.ColorWhite][square] = genPawnAttacks(squareBB, enum.ColorWhite)
		pawnAttacks[enum.ColorBlack][square] = genPawnAttacks(squareBB, enum.ColorBlack)

		knightAttacks[square] = genKnightAttacks(squareBB)

		kingAttacks[square] = genKingAttacks(squareBB)

		bitCount := bishopRelevantOccupancyBitCount[square]
		for i := 0; i < 1<<bitCount; i++ {
			occupancy := genOccupancy(i, bitCount, bishopRelevantOccupancy[square])

			magicKey := occupancy * bishopMagicNumberLookup[square] >> (64 - bitCount)

			bishopAttacks[square][magicKey] = genBishopAttacks(squareBB, occupancy)
		}

		bitCount = rookRelevantOccupancyBitCount[square]
		for i := 0; i < 1<<bitCount; i++ {
			occupancy := genOccupancy(i, bitCount, rookRelevantOccupancy[square])

			magicKey := occupancy * rookMagicNumberLookup[square] >> (64 - bitCount)

			rookAttacks[square][magicKey] = genRookAttacks(squareBB, occupancy)
		}
	}
}

// GenLegalMoves generates legal moves for the currently active color.
// This involves two steps:
//
//  1. Generate pseudo-legal moves for all pieces of the active color.
//  2. Filter out illegal moves using the copy-make approach.
func GenLegalMoves(bitboards [12]uint64, color enum.Color, castlingRights enum.CastlingFlag,
	enPassantTarget int) MoveList {
	var pseudoLegal, legal MoveList

	genPseudoLegalMoves(bitboards, color, &pseudoLegal, castlingRights, 1<<enPassantTarget)

	var occupancy uint64
	for i := enum.PieceWPawn; i <= enum.PieceBKing; i++ {
		occupancy |= bitboards[i]
	}

	var i byte
	for i = 0; i < pseudoLegal.LastMoveIndex; i++ {
		copy := bitboards

		move := pseudoLegal.Moves[i]

		MakeMove(&copy, move)
		// Update occupancy.
		var from, to uint64 = 1 << move.From(), 1 << move.To()
		fromTo := from ^ to
		_occupancy := occupancy ^ fromTo

		kingSquare := bitutil.BitScan(copy[enum.PieceWKing+6*color])

		// If the allied king is not in check, the move is legal.
		if !IsSquareUnderAttack(copy, _occupancy, kingSquare, 1^color) {
			legal.Push(move)
		}
	}

	return legal
}

// MakeMove modifies the piece placement in the given bitboard array by performing the specified move.
func MakeMove(bitboards *[12]uint64, move Move) {
	var from, to uint64 = 1 << move.From(), 1 << move.To()
	fromTo := from ^ to
	movedPiece := GetPieceTypeFromSquare(*bitboards, from)

	switch move.Type() {
	case enum.MoveNormal:
		// If the move is capture.
		capturedPieceType := GetPieceTypeFromSquare(*bitboards, to)
		if capturedPieceType != -1 {
			// Remove the captured piece from the board.
			bitboards[capturedPieceType] ^= to
		}

	case enum.MoveEnPassant:
		// Remove the captured pawn from the board.
		if movedPiece == enum.PieceWPawn {
			bitboards[enum.PieceBPawn] ^= to >> 8
		} else {
			bitboards[enum.PieceWPawn] ^= to << 8
		}

	case enum.MoveCastling:
		switch to {
		case enum.G1, enum.G8: // O-O
			bitboards[movedPiece-2] ^= (to << 1) ^ (to >> 1)
		case enum.C1, enum.C8: // O-O-O
			bitboards[movedPiece-2] ^= (to >> 2) ^ (to << 1)
		}

	case enum.MovePromotion:
		// If the move is capture-promotion.
		capturedPieceType := GetPieceTypeFromSquare(*bitboards, to)
		if capturedPieceType != -1 {
			// Remove the captured piece from the board.
			bitboards[capturedPieceType] ^= to
		}

		// Remove a promoted pawn from the board.
		bitboards[movedPiece] ^= from
		// Place a new piece.
		if movedPiece == enum.PieceWPawn {
			bitboards[move.PromotionPiece()+1] ^= to
		} else {
			bitboards[move.PromotionPiece()+7] ^= to
		}
		return
	}

	// Move piece from the source square to the destination square.
	bitboards[movedPiece] ^= fromTo
}

// IsSquareUnderAttack returns true if the specified square is attacked
// by pieces of the specified color in the given position.
func IsSquareUnderAttack(bitboards [12]uint64, occupancy uint64,
	square int, color enum.Color) bool {
	offset := 6 * color
	// If attacked by pawns.
	if (pawnAttacks[color^1][square]&bitboards[offset+enum.PieceWPawn] != 0) ||
		// If attacked by knights.
		(knightAttacks[square]&bitboards[offset+enum.PieceWKing] != 0) ||
		// If attacked by bishops.
		(lookupBishopAttacks(square, occupancy)&bitboards[offset+enum.PieceWBishop] != 0) ||
		// If attacked by rooks.
		(lookupRookAttacks(square, occupancy)&bitboards[offset+enum.PieceWRook] != 0) ||
		// If attacked by queens.
		(lookupQueenAttacks(square, occupancy)&bitboards[offset+enum.PieceWQueen] != 0) ||
		// If attacked by king.
		(kingAttacks[square]&bitboards[offset+enum.PieceWKing] != 0) {
		// Square is under attack.
		return true
	}

	// Square is not under attack.
	return false
}

// pseudoRandUint64FewBits returns a pseudo-random uint64 with a few set bits.
// It is used only for a magic number generation.
func pseudoRandUint64FewBits() uint64 { return rand.Uint64() & rand.Uint64() & rand.Uint64() }

// genPawnAttacks returns a bitboard of squares attacked by a pawn.
func genPawnAttacks(pawn uint64, color enum.Color) uint64 {
	if color == enum.ColorWhite {
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

// genMagicNumber returns a magic number used to hash sliding piece's possible moves.
// THERE IS NO NEED TO CALL THIS FUNCTION! The magic number tables are predifined.
// May panic if the random number was not generated correctly.
// See https://www.chessprogramming.org/Looking_for_Magics
func genMagicNumber(square int, isBishop bool) uint64 {
	var occupancies, attacks [4096]uint64

	var attack uint64
	var relevantBitCount int
	if isBishop {
		attack = bishopRelevantOccupancy[square]
		relevantBitCount = bishopRelevantOccupancyBitCount[square]
	} else {
		attack = rookRelevantOccupancy[square]
		relevantBitCount = rookRelevantOccupancyBitCount[square]
	}

	for i := 0; i < 1<<relevantBitCount; i++ {
		occupancies[i] = genOccupancy(i, relevantBitCount, attack)

		if isBishop {
			attacks[i] = genBishopAttacks(1<<square, occupancies[i])
		} else {
			attacks[i] = genRookAttacks(1<<square, occupancies[i])
		}
	}

	for i := 0; i < 100000000; i++ {
		magicNumber := pseudoRandUint64FewBits()

		if bitutil.CountBits(attack*magicNumber&0xFF00000000000000) < 6 {
			continue
		}

		var taken [4096]uint64
		var j, fail int

		for ; j < 1<<relevantBitCount; j++ {
			// Form a magic hash key.
			magicKey := int(occupancies[j] * magicNumber >> (64 - relevantBitCount))

			if taken[magicKey] == 0 {
				taken[magicKey] = attacks[j]
			} else if taken[magicKey] != attacks[j] {
				fail = 1
				break
			}
		}

		if fail == 0 {
			return magicNumber
		}
	}

	panic("failed to generate magic number")
}

// lookupBishopAttacks returns a bitboard of squares attacked by a bishop.
// The bitboard is taken from the BishopAttacks using magic hashing scheme.
func lookupBishopAttacks(square int, occupancy uint64) uint64 {
	occupancy &= bishopRelevantOccupancy[square]
	occupancy *= bishopMagicNumberLookup[square]
	occupancy >>= 64 - bishopRelevantOccupancyBitCount[square]
	return bishopAttacks[square][occupancy]
}

// lookupRookAttacks returns a bitboard of squares attacked by a rook.
// The bitboard is taken from the RookAttacks using magic hashing scheme.
func lookupRookAttacks(square int, occupancy uint64) uint64 {
	occupancy &= rookRelevantOccupancy[square]
	occupancy *= rookMagicNumberLookup[square]
	occupancy >>= 64 - rookRelevantOccupancyBitCount[square]
	return rookAttacks[square][occupancy]
}

// lookupQueenAttacks returns a bitboard of squares attacked by a queen.
// The bitboard is calculated as the logical disjunction of the bishop and rook attack bitboards.
func lookupQueenAttacks(square int, occupancy uint64) uint64 {
	return lookupBishopAttacks(square, occupancy) | lookupRookAttacks(square, occupancy)
}

// genPawnsPseudoLegalMoves appends pseudo-legal moves (quiet moves and captures) for the pawns on
// the given bitboard to the specified move list.
func genPawnsPseudoLegalMoves(bitboard, allies, enemies, enPassantTarget uint64,
	color enum.Color, moveList *MoveList) {
	var square, forward, doubleForward, delta int
	var pawn, initialRank, promotionRank uint64
	occupancy := allies | enemies

	if color == enum.ColorWhite {
		delta = 8
		initialRank = RANK_2
		promotionRank = RANK_8
	} else {
		delta = -8
		initialRank = RANK_7
		promotionRank = RANK_1
	}

	// Loop over each pawn within a bitboard.
	for bitboard > 0 {
		square = bitutil.PopLSB(&bitboard)
		pawn = 1 << square

		forward = square + delta
		doubleForward = square + delta*2

		// If a pawn can move 1 square forward.
		var forwardBB uint64 = 1 << forward
		if forwardBB&occupancy == 0 {
			// Check if the move is promotion.
			if forwardBB&promotionRank != 0 {
				moveList.Push(NewMove(forward, square, 0, enum.MovePromotion))
			} else {
				moveList.Push(NewMove(forward, square, 0, enum.MoveNormal))
			}

			// If the pawn is standing on its initial rank and can move 2 squares forward.
			if pawn&initialRank != 0 && 1<<doubleForward&occupancy == 0 {
				moveList.Push(NewMove(doubleForward, square, 0, enum.MoveNormal))
			}
		}

		// Handle pawn attacks.
		attacks := pawnAttacks[color][square] & (enemies | enPassantTarget)
		for attacks > 0 {
			attackedSquare := bitutil.PopLSB(&attacks)

			if 1<<attackedSquare&promotionRank != 0 {
				moveList.Push(NewMove(attackedSquare, square, 0, enum.MovePromotion))
			} else if 1<<attackedSquare&enPassantTarget != 0 {
				moveList.Push(NewMove(attackedSquare, square, 0, enum.MoveEnPassant))
			} else {
				moveList.Push(NewMove(attackedSquare, square, 0, enum.MoveNormal))
			}
		}
	}
}

// genNormalPseudoLegalMoves appends generated pseudo-legal moves for all
// pieces that do not have special move rules (knights, bishops, rooks, and queens).
func genNormalPseudoLegalMoves(pieceType enum.Piece, bitboard, allies,
	enemies uint64, moveList *MoveList) {
	occupancy := allies | enemies

	for bitboard > 0 {
		square := bitutil.PopLSB(&bitboard)

		var attacks uint64

		switch pieceType {
		case enum.PieceWKnight, enum.PieceBKnight:
			attacks = knightAttacks[square]

		case enum.PieceWBishop, enum.PieceBBishop:
			attacks = lookupBishopAttacks(square, occupancy)

		case enum.PieceWRook, enum.PieceBRook:
			attacks = lookupRookAttacks(square, occupancy)

		case enum.PieceWQueen, enum.PieceBQueen:
			attacks = lookupQueenAttacks(square, occupancy)
		}

		for attacks > 0 {
			attackedSquare := bitutil.PopLSB(&attacks)

			if 1<<attackedSquare&enemies != 0 ||
				1<<attackedSquare&allies == 0 {
				moveList.Push(NewMove(attackedSquare, square, 0, enum.MoveNormal))
			}
		}
	}
}

// genKingPseudoLegalMoves appends pseudo-legal moves (quiet moves and captures) for the king on
// the given bitboard to the specified move list.
//
// NOTE: the allies bitboard must exclude the allied king!
func genKingPseudoLegalMoves(square int, allies, enemies, attacked uint64,
	castlingRights enum.CastlingFlag, moveList *MoveList, color enum.Color) {
	occupancy := allies | enemies
	// Lookup normal moves.
	attacks := kingAttacks[square]
	for attacks > 0 {
		attackedSquare := bitutil.PopLSB(&attacks)
		if 1<<attackedSquare&enemies != 0 ||
			1<<attackedSquare&allies == 0 {
			moveList.Push(NewMove(attackedSquare, square, 0, enum.MoveNormal))
		}
	}

	// Handle castling.
	if color == enum.ColorWhite {
		// White O-O.
		if castlingRights&enum.CastlingWhiteShort != 0 &&
			occupancy&WHITE_KING_CASTLING_PATH == 0 &&
			attacked&WHITE_KING_CASTLING_PATH == 0 {
			moveList.Push(NewMove(enum.SG1, square, 0, enum.MoveCastling))
		}
		// White O-O-O.
		if castlingRights&enum.CastlingWhiteLong != 0 &&
			occupancy&WHITE_QUEEN_CASTLING_PATH == 0 &&
			attacked&WHITE_QUEEN_CASTLING_PATH == 0 {
			moveList.Push(NewMove(enum.SC1, square, 0, enum.MoveCastling))
		}
	} else {
		// Black O-O.
		if castlingRights&enum.CastlingBlackShort != 0 &&
			occupancy&BLACK_KING_CASTLING_PATH == 0 &&
			attacked&BLACK_KING_CASTLING_PATH == 0 {
			moveList.Push(NewMove(enum.SG8, square, 0, enum.MoveCastling))
		}
		// Black O-O-O.
		if castlingRights&enum.CastlingBlackLong != 0 &&
			occupancy&BLACK_QUEEN_CASTLING_PATH == 0 &&
			attacked&BLACK_QUEEN_CASTLING_PATH == 0 {
			moveList.Push(NewMove(enum.SC8, square, 0, enum.MoveCastling))
		}
	}
}

// genPseudoLegalMoves generates pseudo-legal moves for the pieces of the specified active color.
func genPseudoLegalMoves(bitboards [12]uint64, color enum.Color, moveList *MoveList,
	castlingRights enum.CastlingFlag, enPassantTarget uint64) {
	var allies, enemies uint64
	colorOffset, enemyColorOffset := 6*color, (1^color)*6

	for i := enum.PieceWPawn; i <= enum.PieceWKing; i++ {
		allies |= bitboards[i+colorOffset]
		enemies |= bitboards[i+enemyColorOffset]
	}

	genPawnsPseudoLegalMoves(bitboards[enum.PieceWPawn+colorOffset], allies, enemies,
		enPassantTarget, color, moveList)

	for i := enum.PieceWKnight + colorOffset; i <= enum.PieceWQueen+colorOffset; i++ {
		genNormalPseudoLegalMoves(i, bitboards[i], allies, enemies, moveList)
	}

	attacked := genAttackedSquares(bitboards, allies|enemies, 1^color)
	kingBB := bitboards[enum.PieceWKing+colorOffset]

	genKingPseudoLegalMoves(bitutil.BitScan(kingBB), allies^kingBB, enemies, attacked,
		castlingRights, moveList, color)
}

// genAttackedSquares returns a bitboard of squares attacked by a pieces of the specified color.
func genAttackedSquares(bitboards [12]uint64, occupancy uint64, color enum.Color) uint64 {
	var attacked uint64
	colorOffset := (color * 6)

	attacked |= genPawnAttacks(bitboards[enum.PieceWPawn+colorOffset], color)
	attacked |= genKnightAttacks(bitboards[enum.PieceWKnight+colorOffset])
	attacked |= genKingAttacks(bitboards[enum.PieceWKing+colorOffset])

	for i := enum.PieceWBishop + colorOffset; i <= enum.PieceWKing+colorOffset; i++ {
		for bitboards[i] > 0 {
			square := bitutil.PopLSB(&bitboards[i])

			switch i {
			case enum.PieceWBishop, enum.PieceBBishop:
				attacked |= lookupBishopAttacks(square, occupancy)
			case enum.PieceWRook, enum.PieceBRook:
				attacked |= lookupRookAttacks(square, occupancy)
			case enum.PieceWQueen, enum.PieceBQueen:
				attacked |= lookupQueenAttacks(square, occupancy)
			}
		}
	}

	return attacked
}

// GetPieceTypeFromSquare returns the type of the piece that stands on the specified square.
// Returns -1 if there is no piece on the square.
func GetPieceTypeFromSquare(bitboards [12]uint64, square uint64) enum.Piece {
	for pieceType, bitboard := range bitboards {
		if square&bitboard != 0 {
			return pieceType
		}
	}
	return -1
}
