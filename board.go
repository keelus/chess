package chess

import (
	"strconv"
	"strings"
	"unicode"
)

type Board [8][8]Piece

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
				kindRune := b[i][j].Kind.ToRune()
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
