package chess

import (
	"strconv"
	"strings"
	"unicode"
)

// Board represents a 8x8 matrix, that consists of rows
// of Pieces.
//
// Note: If you intend in creating a new Board, use NewBoard()
// function, or else the created board's piece's position will
// be (0, 0) by default.
type Board [8][8]Piece

// PieceAtSquare returns a copy of the Piece placed in the
// board.
//
// If there is no piece at that position, it's Kind and Color
// will be none (Kind_None and Color_None).
//
// Note: It assumes that the square is correct. If you create a Square
// instance without the NewSquare, make sure it's valid.
func (b Board) PieceAtSquare(square Square) Piece {
	return b[square.I][square.J]
}

// PieceAtSquareAlgebraic returns a copy of the Piece placed in the
// board.
//
// If the algebraic position is invalid, it will return an empty Piece and it's error.
//
// Examples:
//
//	PieceAtSquareAlgebraic("d2") // returns Piece{Kind:Kind_Pawn, Color:Color_White, ...}, nil
//	PieceAtSquareAlgebraic("aa") // returns Piece{}, error
func (b Board) PieceAtSquareAlgebraic(algebraic string) (Piece, error) {
	square, err := NewSquareFromAlgebraic(algebraic)
	if err != nil {
		return Piece{}, err
	}

	return b.PieceAtSquare(square), nil
}

// Fen returns a string of the board's piece placement
// in Forsyth–Edwards Notation.
//
// For example, for a starting chess game would return:
// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR
//
// Note: This Fen() function returns solely the FEN information
// related to the board's piece placement. You should use
// the Position type's Fen() function if you intend to
// get the turn, halfmove clock, etc.
func (b *Board) Fen() string {
	var sb strings.Builder
	spaceAccum := 0

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b[i][j].Kind == Kind_None {
				spaceAccum++
			} else {
				if spaceAccum > 0 {
					sb.WriteString(strconv.Itoa(spaceAccum))
					spaceAccum = 0
				}
				kindRune := b[i][j].Kind.Rune()
				if b[i][j].Color == Color_White {
					kindRune = unicode.ToUpper(kindRune)
				}
				sb.WriteRune(kindRune)
			}
		}

		if spaceAccum > 0 {
			sb.WriteString(strconv.Itoa(spaceAccum))
			spaceAccum = 0
		}

		if i != 7 {
			sb.WriteRune('/')
		}
	}

	return sb.String()
}

// Unicode returns the unicode representation of the board, with
// ranks, files and pieces.
func (b Board) Unicode() string {
	s := ""
	for i := 0; i < 8; i++ {
		s += strconv.Itoa(i) + " "
		for j := 0; j < 8; j++ {
			s += string([]rune{b[i][j].Unicode(), ' '})
		}
		s += "\n"
	}
	return s + "  a b c d e f g h"
}

func newBoardEmpty() Board {
	var board Board
	for i := uint8(0); i < 8; i++ {
		for j := uint8(0); j < 8; j++ {
			board[i][j].Square = newSquare(i, j)
			board[i][j].Color = Color_None
			board[i][j].Kind = Kind_None
		}
	}

	return board
}

func newBoardFromFen(placementFenData [8]string) Board {
	board := newBoardEmpty()

	col := uint8(0)
	row := uint8(0)

	for _, rowData := range placementFenData {
		for _, colData := range rowData {
			if unicode.IsNumber(colData) {
				colsToJump, _ := strconv.Atoi(string(colData))
				col += uint8(colsToJump) - 1 // Subtract current
			} else {
				kind, color := KindAndColorFromRune(colData)
				board.createPieceAt(color, kind, row, col)
			}

			col++
			if col >= 8 {
				col = 0
			}
		}
		row++
		col = 0
	}

	return board
}

func (b *Board) createPieceAt(color Color, kind Kind, i, j uint8) {
	b[i][j] = newPiece(color, kind, newSquare(i, j))
}
