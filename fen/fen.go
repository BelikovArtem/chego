// Package fen implements conversions between the Forsyth-Edwards Notation strings and bitboard arrays.
// fen expects that the passed FEN and bitboard arrays are always valid and may panic if they are not.
package fen

import (
	"chego/enum"
	"strconv"

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
			var pieceType enum.Piece // enum.PieceWPawn by default.
			// Manual switch construction is ~3x faster than map approach.
			switch char {
			case 'N':
				pieceType = enum.PieceWKnight
			case 'B':
				pieceType = enum.PieceWBishop
			case 'R':
				pieceType = enum.PieceWRook
			case 'Q':
				pieceType = enum.PieceWQueen
			case 'K':
				pieceType = enum.PieceWKing
			case 'p':
				pieceType = enum.PieceBPawn
			case 'n':
				pieceType = enum.PieceBKnight
			case 'b':
				pieceType = enum.PieceBBishop
			case 'r':
				pieceType = enum.PieceBRook
			case 'q':
				pieceType = enum.PieceBQueen
			case 'k':
				pieceType = enum.PieceBKing
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
func Parse(fenStr string) ([12]uint64, enum.Color, enum.CastlingFlag, int, int, int) {
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
	bitboards := ToBitboardArray(fields[0])
	// Parse active color.
	var activeColor enum.Color
	if fields[1] == "b" {
		activeColor = enum.ColorBlack
	}
	// Parse castling rights.
	var castlingRights enum.CastlingFlag
	for i := 0; i < len(fields[2]); i++ {
		switch fields[2][i] {
		case 'K':
			castlingRights |= enum.CastlingWhiteShort
		case 'Q':
			castlingRights |= enum.CastlingWhiteLong
		case 'k':
			castlingRights |= enum.CastlingBlackShort
		case 'q':
			castlingRights |= enum.CastlingBlackLong
		}
	}
	// Parse en passant target square.
	enPassantTarget := squareFromString(fields[3])
	// Parse halfmove counter.
	halfmoveCnt, err := strconv.Atoi(fields[4])
	if err != nil {
		panic("cannot parse halfmove counter from FEN string")
	}
	// Parse fullmove counter.
	fullmoveCnt, err := strconv.Atoi(fields[5])
	if err != nil {
		panic("cannot parse fullmove counter from FEN string")
	}

	return bitboards, activeColor, castlingRights, enPassantTarget, halfmoveCnt, fullmoveCnt
}

// Serialize serializes the game state into a FEN string.
// FEN string contains six fields, each separated by a space.
func Serialize(bitboards [12]uint64, activeColor enum.Color, castlingRights enum.CastlingFlag,
	enPassantTarget, halfmoveCnt, fullmoveCnt int) string {
	var fenStr strings.Builder
	fenStr.Grow(64)

	// 1 field: piece placement.
	fenStr.WriteString(FromBitboardArray(bitboards))
	// 2 field: active color.
	if activeColor == enum.ColorWhite {
		fenStr.WriteString(" w ")
	} else {
		fenStr.WriteString(" b ")
	}
	// 3 field: castling rights.
	cnt := 4
	if castlingRights&enum.CastlingWhiteShort != 0 {
		fenStr.WriteByte('K')
		cnt--
	}
	if castlingRights&enum.CastlingWhiteLong != 0 {
		fenStr.WriteByte('Q')
		cnt--
	}
	if castlingRights&enum.CastlingBlackShort != 0 {
		fenStr.WriteByte('k')
		cnt--
	}
	if castlingRights&enum.CastlingBlackLong != 0 {
		fenStr.WriteByte('q')
		cnt--
	}
	if cnt == 4 {
		fenStr.WriteByte('-')
	}
	fenStr.WriteByte(' ')
	// 4 field: en passant target square.
	if enPassantTarget == 0 {
		fenStr.WriteString("- ")
	} else {
		files := "abcdefgh"
		fenStr.WriteByte(files[enPassantTarget%8])
		fenStr.WriteByte('0' + byte(enPassantTarget/8+1))
		fenStr.WriteByte(' ')
	}
	// 5 field: the number of halfmoves.
	fenStr.WriteString(strconv.Itoa(halfmoveCnt))
	fenStr.WriteByte(' ')
	// 6 field: the number of fullmoves.
	fenStr.WriteString(strconv.Itoa(fullmoveCnt))
	return fenStr.String()
}
