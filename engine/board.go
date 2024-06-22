package engine

import (
	"fmt"
	"strconv"
	"unicode"
)

type Board struct {
	Data [8][8]*Piece

	PlayerToMove Color

	CanKingCastling  map[Color]bool
	CanQueenCastling map[Color]bool

	EnPassant *Position
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
	b.EnPassant = nil

	if movement.IsQueenSideCastling != nil || movement.IsKingSideCastling != nil { // Handle castling
		b.CanQueenCastling[movement.MovingPiece.Color] = false
		b.CanKingCastling[movement.MovingPiece.Color] = false

		castlingRow := 7
		if movement.MovingPiece.Color == Color_Black {
			castlingRow = 0
		}

		if *movement.IsQueenSideCastling {
			rookPiece := b.Data[castlingRow][0]
			kingPiece := b.Data[castlingRow][4]

			rookPiece.Position = NewPosition(castlingRow, 3)
			b.Data[castlingRow][3] = rookPiece
			b.Data[castlingRow][0] = nil

			kingPiece.Position = NewPosition(castlingRow, 2)
			b.Data[castlingRow][2] = kingPiece
			b.Data[castlingRow][4] = nil
		} else if *movement.IsKingSideCastling {
			rookPiece := b.Data[castlingRow][7]
			kingPiece := b.Data[castlingRow][4]

			rookPiece.Position = NewPosition(castlingRow, 5)
			b.Data[castlingRow][5] = rookPiece
			b.Data[castlingRow][7] = nil

			kingPiece.Position = NewPosition(castlingRow, 6)
			b.Data[castlingRow][6] = kingPiece
			b.Data[castlingRow][4] = nil
		}
	} else {
		if movement.MovingPiece.Kind == Kind_Pawn {
			movement.MovingPiece.IsPawnFirstMovement = false

			if *movement.PawnIsDoublePositionMovement {
				invertSum := -1
				if movement.MovingPiece.Color == Color_Black {
					invertSum = +1
				}

				newEnPassantPosition := NewPosition(movement.From.I+invertSum, movement.From.J)
				b.EnPassant = &newEnPassantPosition
			}
		} else if movement.MovingPiece.Kind == Kind_King {
			b.CanQueenCastling[movement.MovingPiece.Color] = false
			b.CanKingCastling[movement.MovingPiece.Color] = false
		} else if movement.MovingPiece.Kind == Kind_Rook {
			rookRow := 7
			if movement.MovingPiece.Color == Color_Black {
				rookRow = 0
			}

			if pieceAt := b.GetPieceAt(rookRow, 0); pieceAt != nil && pieceAt == movement.MovingPiece { // Queen side
				b.CanQueenCastling[movement.MovingPiece.Color] = false
			} else if pieceAt := b.GetPieceAt(rookRow, 7); pieceAt != nil && pieceAt == movement.MovingPiece { // King side
				b.CanKingCastling[movement.MovingPiece.Color] = false
			}
		}

		if movement.TakingPiece != nil {
			b.Data[movement.TakingPiece.Position.I][movement.TakingPiece.Position.J] = nil // Do it this way, so it's en passant compatible
		}

		movement.MovingPiece.Position = NewPosition(movement.To.I, movement.To.J)
		b.Data[movement.To.I][movement.To.J] = movement.MovingPiece
		b.Data[movement.From.I][movement.From.J] = nil
	}
}

func (b *Board) UndoMovement(movement Movement) {
	fmt.Println("Undo movement")
	// Remove the moved piece
	b.Data[movement.To.I][movement.To.J] = nil

	// Create the moved piece into the old position
	movedPieceCopy := movement.MovingPieceCopy.DeepCopy()
	b.Data[movement.MovingPieceCopy.Position.I][movement.MovingPieceCopy.Position.J] = &movedPieceCopy

	// Create the taken piece (if aplicable) into the old position
	if movement.TakingPieceCopy != nil {
		takenPieceCopy := movement.TakingPieceCopy.DeepCopy()
		b.Data[movement.TakingPieceCopy.Position.I][movement.TakingPieceCopy.Position.J] = &takenPieceCopy
	}

	if movement.IsQueenSideCastling != nil && *movement.IsQueenSideCastling {
		// Move castle
		castlingRow := 7
		if b.PlayerToMove == Color_Black {
			castlingRow = 0
		}

		// Delete castle
		b.Data[castlingRow][3] = nil
		b.CreatePieceAt(b.PlayerToMove, Kind_Rook, castlingRow, 0)
	} else if movement.IsKingSideCastling != nil && *movement.IsKingSideCastling {
		// Move castle
		castlingRow := 7
		if b.PlayerToMove == Color_Black {
			castlingRow = 0
		}

		// Delete castle
		b.Data[castlingRow][5] = nil
		b.CreatePieceAt(b.PlayerToMove, Kind_Rook, castlingRow, 7)
	}

	b.CanQueenCastling[movement.MovingPieceCopy.Color] = movement.CanQueenSideCastling
	b.CanKingCastling[movement.MovingPieceCopy.Color] = movement.CanKingSideCastling

	// If is Pawn, set it's variables. Althought this might not be necessary
	// if movement.MovingPieceCopy.Kind == Kind_Pawn {
	// 	b.Data[movement.TakingPieceCopy.Position.I][movement.TakingPieceCopy.Position.J].IsPawnFirstMovement = *movement.PawnIsFirstMove
	// 	//movement.MovingPieceCopy.IsPawnFirstMovement = *movement.PawnIsFirstMove
	// }

	b.EnPassant = movement.EnPassant
}

// For future implementation
// func (b Board) GetPseudoMovements() {}
// func (b Board) GetLegalMovements()  {}
// func (b *Board) MakeMovement()      {}
// func (b *Board) UnMakeMovement()    {}
