package engine

import (
	"strconv"
	"unicode"
)

type Board [8][8]Piece

func NewBoardEmpty() Board {
	var board Board
	for i := uint8(0); i < 8; i++ {
		for j := uint8(0); j < 8; j++ {
			board[i][j].Point = NewPoint(i, j)
			board[i][j].Color = Color_None
			board[i][j].Kind = Kind_None
			board[i][j].IsPawnFirstMovement = false
		}
	}

	return board
}

func NewBoardFromFen(placementFenData [8]string) Board {
	board := NewBoardEmpty()

	col := uint8(0)
	row := uint8(0)

	for _, rowData := range placementFenData {
		for _, colData := range rowData {
			if unicode.IsNumber(colData) {
				colsToJump, _ := strconv.Atoi(string(colData))
				col += uint8(colsToJump) - 1 // Subtract current
			} else {
				kind, color := KindAndColorFromRune(colData)

				board.CreatePieceAt(color, kind, row, col)
				if kind == Kind_Pawn {
					pawnRow := uint8(6)
					if color == Color_Black {
						pawnRow = 1
					}

					if row != pawnRow {
						board[row][col].IsPawnFirstMovement = false
					}
				}
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

func (b *Board) CreatePieceAt(color Color, kind Kind, i, j uint8) {
	b[i][j] = NewPiece(color, kind, NewPoint(i, j))
}

func (b Board) GetPieceAt(i, j uint8) Piece {
	return b[i][j]
}
