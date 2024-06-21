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
				newBoard.CreatePieceAt(color, kind, row, col)
			}

			col++
			if col >= 8 {
				col = 0
			}
		}
		row++
		col = 0
	}

	return newBoard
}

func (b *Board) CreatePieceAt(color Color, kind Kind, i, j int) {
	newPiece := NewPiece(color, kind, NewPosition(i, j))
	b.Data[i][j] = &newPiece
}

func (b Board) GetPieceAt(i, j int) *Piece {
	return b.Data[i][j]
}

// Suppose is legal
func (b *Board) MakeMovement(movement Movement) {
	if movement.IsQueenSideCastling != nil || movement.IsKingSideCastling != nil { // Handle castling
		b.CanQueenCastling[movement.MovingPiece.Color] = false
		b.CanKingCastling[movement.MovingPiece.Color] = false
		if *movement.IsQueenSideCastling {
			// TODO: Do not hardcode this
			rookPiece := b.Data[7][0]
			kingPiece := b.Data[7][4]

			rookPiece.Position = NewPosition(7, 3)
			b.Data[7][3] = rookPiece
			b.Data[7][0] = nil

			kingPiece.Position = NewPosition(7, 2)
			b.Data[7][2] = kingPiece
			b.Data[7][4] = nil
		} else if *movement.IsKingSideCastling {
			// TODO: Do not hardcode this
			rookPiece := b.Data[7][7]
			kingPiece := b.Data[7][4]

			rookPiece.Position = NewPosition(7, 5)
			b.Data[7][5] = rookPiece
			b.Data[7][7] = nil

			kingPiece.Position = NewPosition(7, 6)
			b.Data[7][6] = kingPiece
			b.Data[7][4] = nil

		}
	} else {
		if movement.MovingPiece.Kind == Kind_Pawn {
			movement.MovingPiece.IsPawnFirstMovement = false
		}

		movement.MovingPiece.Position = NewPosition(movement.To.I, movement.To.J)
		b.Data[movement.To.I][movement.To.J] = movement.MovingPiece
		b.Data[movement.From.I][movement.From.J] = nil
	}
}

func (b *Board) UndoMovement(movement Movement) {
}

// For future implementation
// func (b Board) GetPseudoMovements() {}
// func (b Board) GetLegalMovements()  {}
// func (b *Board) MakeMovement()      {}
// func (b *Board) UnMakeMovement()    {}
