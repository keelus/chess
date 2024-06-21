package board

import (
	"chess/piece"
	"strconv"
	"unicode"
)

type Board struct {
	Data [8][8]*piece.Piece

	PlayerToMove piece.Color
}

func NewEmptyBoard() Board {
	return Board{
		PlayerToMove: piece.Color_White,
	}
}

func NewBoardFromFen(fen string) Board {
	parsedFen, err := parseFen(fen)
	if err != nil {
		panic(err)
	}

	col := 0
	row := 0

	newBoard := Board{
		PlayerToMove: parsedFen.ActiveColor,
	}

	for _, rowData := range parsedFen.PlacementData {
		for _, colData := range rowData {
			if unicode.IsNumber(colData) {
				colsToJump, _ := strconv.Atoi(string(colData))
				col += colsToJump - 1 // Subtract current
			} else {
				kind, color := piece.KindAndColorFromRune(colData)
				newPiece := piece.NewPiece(color, kind)
				newBoard.Data[row][col] = &newPiece
			}

			col++
			if col >= 8 {
				col = 0
			}
		}
		row++
	}

	return newBoard
}

// For future implementation
// func (b Board) GetPseudoMovements() {}
// func (b Board) GetLegalMovements()  {}
// func (b *Board) MakeMovement()      {}
// func (b *Board) UnMakeMovement()    {}
