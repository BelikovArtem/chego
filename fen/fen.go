// Package fen implements conversions between the Forsyth-Edwards Notation strings and bitboard arrays.
// fen expects that the passed FEN and bitboard arrays are always valid, and may panic if they are not.
package fen

import (
	"strconv"

	"github.com/BelikovArtem/chego/types"

	// bits is used to speed up the iteration over bitboards.
	"math/bits"
	// strings is used to reduce the number of memory allocations during strings concatenation.
	"strings"
)

// ToBitboardArray converts the first part of a Forsyth-Edwards Notation string into
// an array of bitboards.
func ToBitboardArray(piecePlacementData string) [12]uint64 {
	var bitboards [12]uint64
	squareIndex := 56

	// Piece placement data describes each rank beginning from the eigth.
	for i := 0; i < len(piecePlacementData); i++ {
		char := piecePlacementData[i]

		if char == '/' { // Rank separator.
			squareIndex -= 16
		} else if char >= '1' && char <= '8' { // Number of consecutive empty squares.
			// Convert byte to the integer it represents.
			squareIndex += int(char - '0')
		} else { // There is piece on a square.
			var pieceType types.Piece // types.PieceWPawn by default.
			// Manual switch construction is ~3x faster than map approach.
			switch char {
			case 'N':
				pieceType = types.PieceWKnight
			case 'B':
				pieceType = types.PieceWBishop
			case 'R':
				pieceType = types.PieceWRook
			case 'Q':
				pieceType = types.PieceWQueen
			case 'K':
				pieceType = types.PieceWKing
			case 'p':
				pieceType = types.PieceBPawn
			case 'n':
				pieceType = types.PieceBKnight
			case 'b':
				pieceType = types.PieceBBishop
			case 'r':
				pieceType = types.PieceBRook
			case 'q':
				pieceType = types.PieceBQueen
			case 'k':
				pieceType = types.PieceBKing
			}
			// Set the bit on the bitboard to place a piece.
			bitboards[pieceType] |= 1 << squareIndex
			squareIndex++
		}
	}

	return bitboards
}

// pieceSymbols is used in FromBitboardArray function.
var pieceSymbols = [12]byte{
	'P', 'N', 'B', 'R', 'Q', 'K',
	'p', 'n', 'b', 'r', 'q', 'k',
}

// FromBitboardArray converts the array of bitboards into the first part
// of Forsyth-Edwards Notation.
func FromBitboardArray(bitboards [12]uint64) string {
	// Used to add characters to a string without extra mem allocs.
	var piecePlacementData strings.Builder
	piecePlacementData.Grow(20)

	var board [8][8]byte

	for pieceType, bitboard := range bitboards {
		// Go through all pieces on a bitboard.
		for ; bitboard > 0; bitboard &= bitboard - 1 {
			squareIndex := bits.TrailingZeros64(bitboard)
			// Add piece on board.
			board[squareIndex/8][squareIndex%8] = pieceSymbols[pieceType]
		}
	}

	var numOfEmptySquares byte

	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			char := board[rank][file]

			if char == 0 { // Empty square.
				numOfEmptySquares++
			} else { // Piece on square.
				if numOfEmptySquares > 0 {
					piecePlacementData.WriteByte('0' + numOfEmptySquares)
					numOfEmptySquares = 0
				}
				piecePlacementData.WriteByte(char)
			}

			// To add rank separators.
			squareIndex := 8*rank + file
			if (squareIndex+1)%8 == 0 {
				if numOfEmptySquares > 0 {
					piecePlacementData.WriteByte('0' + numOfEmptySquares)
					numOfEmptySquares = 0
				}
				// Do not add separator in the end of the string.
				if squareIndex != 7 {
					piecePlacementData.WriteByte('/')
				}
			}
		}
	}

	return piecePlacementData.String()
}

// squareFromString is used to parse en passant target square.
// Handles '-' as A1 square.
func squareFromString(str string) int {
	var square int
	switch str[0] {
	case 'b':
		square = 1
	case 'c':
		square = 2
	case 'd':
		square = 3
	case 'e':
		square = 4
	case 'f':
		square = 5
	case 'g':
		square = 6
	case 'h':
		square = 7
	case '-':
		return 0
	}
	return square + (int(str[1]-'0')-1)*8
}

// Parse parses the given FEN string and returns the values of its fields.
func Parse(fenStr string) (p types.Position) {
	var fields [6]string
	// Separate FEN fields.
	var j, prev int
	for i := 0; i < len(fenStr); i++ {
		// Field separator.
		if fenStr[i] == ' ' {
			fields[j] = fenStr[prev:i]
			j++
			prev = i + 1
		}
	}
	fields[5] = fenStr[prev:]
	// Parce piece placement.
	p.Bitboards = ToBitboardArray(fields[0])
	// Parse active color.
	if fields[1] == "b" {
		p.ActiveColor = types.ColorBlack
	}
	// Parse castling rights.
	for i := 0; i < len(fields[2]); i++ {
		switch fields[2][i] {
		case 'K':
			p.CastlingRights |= types.CastlingWhiteShort
		case 'Q':
			p.CastlingRights |= types.CastlingWhiteLong
		case 'k':
			p.CastlingRights |= types.CastlingBlackShort
		case 'q':
			p.CastlingRights |= types.CastlingBlackLong
		}
	}
	// Parse en passant target square.
	p.EPTarget = squareFromString(fields[3])
	// Parse halfmove counter.
	var err error
	p.HalfmoveCnt, err = strconv.Atoi(fields[4])
	if err != nil {
		panic("cannot parse halfmove counter from FEN string")
	}
	// Parse fullmove counter.
	p.FullmoveCnt, err = strconv.Atoi(fields[5])
	if err != nil {
		panic("cannot parse fullmove counter from FEN string")
	}

	return p
}

// Serialize serializes the specified position into a FEN string.
// FEN string contains six fields, each separated by a space.
func Serialize(p types.Position) string {
	var fenStr strings.Builder
	fenStr.Grow(64)

	// 1 field: piece placement.
	fenStr.WriteString(FromBitboardArray(p.Bitboards))
	// 2 field: active color.
	if p.ActiveColor == types.ColorWhite {
		fenStr.WriteString(" w ")
	} else {
		fenStr.WriteString(" b ")
	}
	// 3 field: castling rights.
	cnt := 4
	if p.CastlingRights&types.CastlingWhiteShort != 0 {
		fenStr.WriteByte('K')
		cnt--
	}
	if p.CastlingRights&types.CastlingWhiteLong != 0 {
		fenStr.WriteByte('Q')
		cnt--
	}
	if p.CastlingRights&types.CastlingBlackShort != 0 {
		fenStr.WriteByte('k')
		cnt--
	}
	if p.CastlingRights&types.CastlingBlackLong != 0 {
		fenStr.WriteByte('q')
		cnt--
	}
	if cnt == 4 {
		fenStr.WriteByte('-')
	}
	fenStr.WriteByte(' ')
	// 4 field: en passant target square.
	if p.EPTarget == 0 {
		fenStr.WriteString("- ")
	} else {
		files := "abcdefgh"
		fenStr.WriteByte(files[p.EPTarget%8])
		fenStr.WriteByte('0' + byte(p.EPTarget/8+1))
		fenStr.WriteByte(' ')
	}
	// 5 field: the number of halfmoves.
	fenStr.WriteString(strconv.Itoa(p.HalfmoveCnt))
	fenStr.WriteByte(' ')
	// 6 field: the number of fullmoves.
	fenStr.WriteString(strconv.Itoa(p.FullmoveCnt))

	return fenStr.String()
}
