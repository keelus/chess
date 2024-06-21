package engine

import (
	"strconv"
	"unicode"
)

type Board struct {
	Data [8][8]*Piece

	PlayerToMove Color

	CanKingCastling  map[Color]bool
	CanQueenCastling map[Color]bool
}

func NewEmptyBoard() Board {
	return Board{
		PlayerToMove: Color_White,

		CanKingCastling:  map[Color]bool{Color_White: true, Color_Black: true},
		CanQueenCastling: map[Color]bool{Color_White: true, Color_Black: true},
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

		CanKingCastling: map[Color]bool{
			Color_White: parsedFen.WhiteCanKingSideCastling,
			Color_Black: parsedFen.BlackCanKingSideCastling,
		},
		CanQueenCastling: map[Color]bool{
			Color_White: parsedFen.WhiteCanQueenSideCastling,
			Color_Black: parsedFen.BlackCanQueenSideCastling,
		},
	}

	for _, rowData := range parsedFen.PlacementData {
		for _, colData := range rowData {
			if unicode.IsNumber(colData) {
				colsToJump, _ := strconv.Atoi(string(colData))
				col += colsToJump - 1 // Subtract current
			} else {
				kind, color := KindAndColorFromRune(colData)
				newPiece := NewPiece(color, kind, NewPosition(row, col))
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

func (b Board) GetPieceAt(i, j int) *Piece {
	return b.Data[i][j]
}

// For future implementation
// func (b Board) GetPseudoMovements() {}
// func (b Board) GetLegalMovements()  {}
// func (b *Board) MakeMovement()      {}
// func (b *Board) UnMakeMovement()    {}
