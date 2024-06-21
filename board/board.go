package board

import (
	"chess/piece"
	"chess/position"
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

func NewStartingBoard() Board {
	return NewBoardFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
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
				newPiece := piece.NewPiece(color, kind, position.NewPosition(row, col))
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

func (b Board) GetPieceAt(i, j int) *piece.Piece {
	return b.Data[i][j]
}

// For future implementation
// func (b Board) GetPseudoMovements() {}
// func (b Board) GetLegalMovements()  {}
// func (b *Board) MakeMovement()      {}
// func (b *Board) UnMakeMovement()    {}
