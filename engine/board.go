package engine

import (
	"fmt"
	"strconv"
	"unicode"
)

type Board struct {
	Data [8][8]Piece

	PlayerToMove Color

	CanKingCastling  map[Color]bool
	CanQueenCastling map[Color]bool

	EnPassant *Position
}

func NewEmptyBoard() Board {
	return Board{
		Data:         createEmptyBoardData(),
		PlayerToMove: Color_White,

		CanKingCastling:  map[Color]bool{Color_White: true, Color_Black: true},
		CanQueenCastling: map[Color]bool{Color_White: true, Color_Black: true},
	}
}

func createEmptyBoardData() [8][8]Piece {
	var boardData [8][8]Piece
	for i := range boardData {
		for j := range boardData {
			boardData[i][j].Position = NewPosition(i, j)
			boardData[i][j].Color = Color_None
			boardData[i][j].Kind = Kind_None
			boardData[i][j].IsPawnFirstMovement = false
		}
	}

	return boardData
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
		Data:         createEmptyBoardData(),
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
				if kind == Kind_Pawn {
					pawnRow := 6
					if color == Color_Black {
						pawnRow = 1
					}

					if row != pawnRow {
						newBoard.Data[row][col].IsPawnFirstMovement = false
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

	return newBoard
}

func (b *Board) CreatePieceAt(color Color, kind Kind, i, j int) {
	b.Data[i][j] = NewPiece(color, kind, NewPosition(i, j))
}

func (b Board) GetPieceAt(i, j int) Piece {
	return b.Data[i][j]
}

// Suppose is legal
func (b *Board) MakeMovement(movement Movement) {
	//fmt.Printf("Do: %s\n", movement.ToString())
	b.EnPassant = nil

	if movement.IsQueenSideCastling != nil || movement.IsKingSideCastling != nil { // Handle castling
		b.CanQueenCastling[movement.MovingPiece.Color] = false
		b.CanKingCastling[movement.MovingPiece.Color] = false

		castlingRow := 7
		if movement.MovingPiece.Color == Color_Black {
			castlingRow = 0
		}
		// ERROR CURRENTLY ON SIDE CASTLING QUEEN, ON GETTING MOVEMENTS

		// RECOMMENDATION MAYBE? DONT USE ANY POINTER, ALL VIA POSITION FROM AND TO (WITH COPIES OF TAKING AND MOVING PIECES)
		// IF DOING THIS, WE HAVE TO SAVE ALSO CASTLING POSITION

		if *movement.IsQueenSideCastling {
			rookPiece := b.Data[castlingRow][0]
			kingPiece := b.Data[castlingRow][4]

			// Set new rook
			b.Data[castlingRow][3].Kind = rookPiece.Kind
			b.Data[castlingRow][3].Color = rookPiece.Color

			// Delete old rook
			b.Data[castlingRow][0].Kind = Kind_None
			b.Data[castlingRow][0].Color = Color_None

			// Set new king
			b.Data[castlingRow][2].Kind = kingPiece.Kind
			b.Data[castlingRow][2].Color = kingPiece.Color

			// Delete old king
			b.Data[castlingRow][4].Kind = Kind_None
			b.Data[castlingRow][4].Color = Color_None
		} else if *movement.IsKingSideCastling {
			rookPiece := b.Data[castlingRow][7]
			kingPiece := b.Data[castlingRow][4]

			// Set new rook
			b.Data[castlingRow][5].Kind = rookPiece.Kind
			b.Data[castlingRow][5].Color = rookPiece.Color

			// Delete old rook
			b.Data[castlingRow][7].Kind = Kind_None
			b.Data[castlingRow][7].Color = Color_None

			// Set new king
			b.Data[castlingRow][6].Kind = kingPiece.Kind
			b.Data[castlingRow][6].Color = kingPiece.Color

			// Delete old king
			b.Data[castlingRow][4].Kind = Kind_None
			b.Data[castlingRow][4].Color = Color_None
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
			// Check if currently moving rook is from queen or king side
			if b.CanQueenCastling[movement.MovingPiece.Color] {
				if movement.MovingPiece.Position.J == 0 {
					b.CanQueenCastling[movement.MovingPiece.Color] = false
				}
			}
			if b.CanKingCastling[movement.MovingPiece.Color] {
				if movement.MovingPiece.Position.J == 7 {
					b.CanKingCastling[movement.MovingPiece.Color] = false
				}
			}
		}

		if movement.IsTakingPiece {
			// Do it this wad, so it's en passant compatible
			b.Data[movement.TakingPiece.Position.I][movement.TakingPiece.Position.J].Kind = Kind_None
			b.Data[movement.TakingPiece.Position.I][movement.TakingPiece.Position.J].Color = Color_None

			if movement.TakingPiece.Kind == Kind_Rook {
				if b.CanQueenCastling[movement.TakingPiece.Color] {
					castlingRow := 7
					if movement.TakingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.TakingPiece.Position.I == castlingRow && movement.TakingPiece.Position.J == 0 {
						b.CanQueenCastling[movement.TakingPiece.Color] = false
					}
				}
				if b.CanKingCastling[movement.TakingPiece.Color] {
					castlingRow := 7
					if movement.TakingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.TakingPiece.Position.I == castlingRow && movement.TakingPiece.Position.J == 7 {
						b.CanKingCastling[movement.TakingPiece.Color] = false
					}
				}
			}
		}

		// movement.MovingPiece.Position = NewPosition(movement.To.I, movement.To.J)
		// b.Data[movement.To.I][movement.To.J] = movement.MovingPiece
		// b.Data[movement.From.I][movement.From.J] = nil

		// Update data of the new piece
		b.Data[movement.To.I][movement.To.J].Kind = movement.MovingPiece.Kind
		b.Data[movement.To.I][movement.To.J].Color = movement.MovingPiece.Color

		// Delete this piece's previous position
		b.Data[movement.From.I][movement.From.J].Kind = Kind_None
		b.Data[movement.From.I][movement.From.J].Color = Color_None
		b.Data[movement.From.I][movement.From.J].IsPawnFirstMovement = false
	}

	b.PlayerToMove = Color_White
	if movement.MovingPiece.Color == Color_White {
		b.PlayerToMove = Color_Black
	}
}

func (b *Board) UndoMovement(movement Movement) {
	//fmt.Printf("Undo: %s\n", movement.ToString())

	// Remove the moved piece
	b.Data[movement.To.I][movement.To.J].Kind = Kind_None
	b.Data[movement.To.I][movement.To.J].Color = Color_None

	// Create the moved piece into the old position
	//b.Data[movement.From.I][movement.From.J] = movement.MovingPiece
	b.Data[movement.From.I][movement.From.J].Kind = movement.MovingPiece.Kind
	b.Data[movement.From.I][movement.From.J].Color = movement.MovingPiece.Color
	b.Data[movement.From.I][movement.From.J].IsPawnFirstMovement = movement.MovingPiece.IsPawnFirstMovement
	// TODO: THIS

	// Create the taken piece (if aplicable) into the old position
	if movement.IsTakingPiece {
		b.Data[movement.TakingPiece.Position.I][movement.TakingPiece.Position.J].Kind = movement.TakingPiece.Kind
		b.Data[movement.TakingPiece.Position.I][movement.TakingPiece.Position.J].Color = movement.TakingPiece.Color
		b.Data[movement.TakingPiece.Position.I][movement.TakingPiece.Position.J].IsPawnFirstMovement = movement.TakingPiece.IsPawnFirstMovement

	}

	if movement.IsQueenSideCastling != nil && *movement.IsQueenSideCastling {
		// Move castle
		castlingRow := 7
		if movement.MovingPiece.Color == Color_Black {
			castlingRow = 0
		}

		// Delete moved castle
		b.Data[castlingRow][3].Kind = Kind_None
		b.Data[castlingRow][3].Color = Color_None

		// Create old castle
		b.Data[castlingRow][0].Kind = Kind_Rook
		b.Data[castlingRow][0].Color = movement.MovingPiece.Color
	} else if movement.IsKingSideCastling != nil && *movement.IsKingSideCastling {
		// Move castle
		castlingRow := 7
		if movement.MovingPiece.Color == Color_Black {
			castlingRow = 0
		}

		// Delete moved castle
		b.Data[castlingRow][5].Kind = Kind_None
		b.Data[castlingRow][5].Color = Color_None

		// Create old castle
		b.Data[castlingRow][7].Kind = Kind_Rook
		b.Data[castlingRow][7].Color = movement.MovingPiece.Color
	}

	// b.CanQueenCastling[movement.MovingPiece.Color] = movement.CanQueenSideCastling
	// b.CanKingCastling[movement.MovingPiece.Color] = movement.CanKingSideCastling

	b.CanQueenCastling[Color_White] = movement.CanWhiteQueenSideCastling
	b.CanKingCastling[Color_White] = movement.CanWhiteKingSideCastling
	b.CanQueenCastling[Color_Black] = movement.CanBlackQueenSideCastling
	b.CanKingCastling[Color_Black] = movement.CanBlackKingSideCastling

	b.EnPassant = movement.EnPassant

	b.PlayerToMove = Color_White
	if movement.MovingPiece.Color == Color_Black {
		b.PlayerToMove = Color_Black
	}
}

func (b *Board) FilterPseudoMovements(movements []Movement) []Movement {
	//beginningColor := b.PlayerToMove
	filteredMovements := []Movement{}

	opponentColor := Color_White
	if b.PlayerToMove == Color_White {
		opponentColor = Color_Black
	}

	for _, myMovement := range movements {
		if myMovement.PawnIsAttackingButNotTakingDiagonal != nil && *myMovement.PawnIsAttackingButNotTakingDiagonal {
			continue
		}

		isCastlingLegal := true
		if (myMovement.IsQueenSideCastling != nil && *myMovement.IsQueenSideCastling) ||
			(myMovement.IsKingSideCastling != nil && *myMovement.IsKingSideCastling) {
			currentOpponentPseudo := b.GetPseudoMovements(opponentColor)
			// Cols that must not be being attacked
			var colFrom, colTo int

			if myMovement.IsQueenSideCastling != nil && *myMovement.IsQueenSideCastling {
				colFrom = myMovement.From.J - 1
				colTo = myMovement.From.J - 2
			} else if myMovement.IsKingSideCastling != nil && *myMovement.IsKingSideCastling {
				colFrom = myMovement.From.J + 1
				colTo = myMovement.From.J + 2
			}

			isCastlingLegal = b.CheckForCastlingLegal(myMovement.From.I, colFrom, colTo, currentOpponentPseudo)
		}

		if !isCastlingLegal {
			continue
		}

		b.MakeMovement(myMovement)
		opponentPseudoMovements := b.GetPseudoMovements(opponentColor)

		weGetChecked := b.CheckForCheck(opponentPseudoMovements)

		if !weGetChecked {
			filteredMovements = append(filteredMovements, myMovement)
		}
		// Check for check
		b.UndoMovement(myMovement)
	}

	//endColor := b.PlayerToMove

	//fmt.Printf("%c vs %c\n", beginningColor.ToRune(), endColor.ToRune())

	return filteredMovements
}

func (b Board) CheckForCheck(opponentPseudoMovements []Movement) bool {
	for _, movement := range opponentPseudoMovements {
		if movement.IsTakingPiece && movement.TakingPiece.Kind == Kind_King {
			return true
		}
	}

	return false
}

func (b Board) CheckForCastlingLegal(row, colFrom, colTo int, opponentPseudoMovements []Movement) bool {
	jFrom, jTo := colFrom, colTo

	if jFrom > jTo {
		jTo, jFrom = jFrom, jTo
	}

	for j := jFrom; j <= jTo; j++ {
		for _, movement := range opponentPseudoMovements {
			if movement.To.I == row && movement.To.J == j {
				return false
			}
		}
	}
	return true
}

// TODO: Complete
func (b Board) ToFen() string {
	dataFen := ""
	spaceAccum := 0

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.Data[i][j].Kind == Kind_None {
				spaceAccum++
			} else {
				if spaceAccum > 0 {
					dataFen = fmt.Sprintf("%s%d", dataFen, spaceAccum)
					spaceAccum = 0
				}
				kindRune := b.Data[i][j].Kind.ToRune()
				if b.Data[i][j].Color == Color_White {
					kindRune = unicode.ToUpper(kindRune)
				}
				dataFen = fmt.Sprintf("%s%c", dataFen, kindRune)
			}
		}

		if spaceAccum > 0 {
			dataFen = fmt.Sprintf("%s%d", dataFen, spaceAccum)
			spaceAccum = 0
		}
		dataFen = fmt.Sprintf("%s/", dataFen)
	}

	dataFen = fmt.Sprintf("%s", dataFen[:len(dataFen)-1])

	dataFen = fmt.Sprintf("%s %c ", dataFen, b.PlayerToMove.ToRune())

	if b.CanKingCastling[Color_White] {
		dataFen = fmt.Sprintf("%sK", dataFen)
	}
	if b.CanQueenCastling[Color_White] {
		dataFen = fmt.Sprintf("%sQ", dataFen)
	}
	if b.CanKingCastling[Color_Black] {
		dataFen = fmt.Sprintf("%sk", dataFen)
	}
	if b.CanQueenCastling[Color_Black] {
		dataFen = fmt.Sprintf("%sq", dataFen)
	}

	if b.EnPassant != nil {
		dataFen = fmt.Sprintf("%s %s", dataFen, b.EnPassant.ToAlgebraic())
	} else {
		dataFen = fmt.Sprintf("%s -", dataFen)
	}

	dataFen = fmt.Sprintf("%s 0 1", dataFen)

	return dataFen
}

// We have White pseudo movements
// Do each movement, check if any enemy

// For future implementation
// func (b Board) GetPseudoMovements() {}
// func (b Board) GetLegalMovements()  {}
// func (b *Board) MakeMovement()      {}
// func (b *Board) UnMakeMovement()    {}
