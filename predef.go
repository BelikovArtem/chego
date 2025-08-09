// predef.go contains declarations of predefined attack tables, magic numbers,
// and other useful constants.

package chego

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
		0x70, 0x1E, 0x7000000000000000, 0x1E00000000000000,
	}
	castlingAttackPath = [4]uint64{
		0x70, 0x1C, 0x7000000000000000, 0x1C00000000000000,
	}
)

const InitialPos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Bitboards of each square.
const (
	A1 uint64 = 1 << iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
	ALL_SQUARES = 0xFFFFFFFFFFFFFFFF
)

// Indicies of each square.
const (
	SA1 int = iota
	SB1
	SC1
	SD1
	SE1
	SF1
	SG1
	SH1
	SA2
	SB2
	SC2
	SD2
	SE2
	SF2
	SG2
	SH2
	SA3
	SB3
	SC3
	SD3
	SE3
	SF3
	SG3
	SH3
	SA4
	SB4
	SC4
	SD4
	SE4
	SF4
	SG4
	SH4
	SA5
	SB5
	SC5
	SD5
	SE5
	SF5
	SG5
	SH5
	SA6
	SB6
	SC6
	SD6
	SE6
	SF6
	SG6
	SH6
	SA7
	SB7
	SC7
	SD7
	SE7
	SF7
	SG7
	SH7
	SA8
	SB8
	SC8
	SD8
	SE8
	SF8
	SG8
	SH8
)
