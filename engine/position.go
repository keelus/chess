package engine

import (
	"fmt"
	"unicode"
)

type Position struct {
	Board Board

	Status PositionStatus
}

type CastlingRights struct {
	CanWhiteKingCastle  bool
	CanWhiteQueenCastle bool
	CanBlackKingCastle  bool
	CanBlackQueenCastle bool
}

func (c CastlingRights) CanQueenCastling(color Color) bool {
	if color == Color_White {
		return c.CanWhiteQueenCastle
	}
	return c.CanBlackQueenCastle
}

func (c CastlingRights) CanKingCastling(color Color) bool {
	if color == Color_White {
		return c.CanWhiteKingCastle
	}
	return c.CanBlackKingCastle
}

func (c *CastlingRights) SetQueenCastling(color Color, canCastling bool) {
	if color == Color_White {
		c.CanWhiteQueenCastle = canCastling
	} else {
		c.CanBlackQueenCastle = canCastling
	}
}

func (c *CastlingRights) SetKingCastling(color Color, canCastling bool) {
	if color == Color_White {
		c.CanWhiteKingCastle = canCastling
	} else {
		c.CanBlackKingCastle = canCastling
	}
}

type PositionStatus struct {
	PlayerToMove Color

	CastlingRights CastlingRights
	// CanKingCastling  map[Color]bool
	// CanQueenCastling map[Color]bool

	EnPassant *Point
}

func NewPositionFromFen(fen string) Position {
	parsedFen, err := parseFen(fen)
	if err != nil {
		panic(err)
	}

	return Position{
		Board: NewBoardFromFen(parsedFen.PlacementData),
		Status: PositionStatus{
			PlayerToMove: parsedFen.ActiveColor,

			CastlingRights: CastlingRights{
				CanWhiteQueenCastle: parsedFen.WhiteCanQueenSideCastling,
				CanWhiteKingCastle:  parsedFen.WhiteCanKingSideCastling,

				CanBlackQueenCastle: parsedFen.BlackCanQueenSideCastling,
				CanBlackKingCastle:  parsedFen.BlackCanKingSideCastling,
			},
			// CanKingCastling: map[Color]bool{
			// 	Color_White: parsedFen.WhiteCanKingSideCastling,
			// 	Color_Black: parsedFen.BlackCanKingSideCastling,
			// },
			// CanQueenCastling: map[Color]bool{
			// 	Color_White: parsedFen.WhiteCanQueenSideCastling,
			// 	Color_Black: parsedFen.BlackCanQueenSideCastling,
			// },

			EnPassant: nil, //TODO
		},
	}
}

// Suppose is legal
func (g *Game) MakeMovement(movement Movement) {
	newPosition := Position{
		Board:  g.CurrentPosition.Board,
		Status: g.CurrentPosition.Status,
	}

	//fmt.Printf("Do: %s\n", movement.ToString())
	newPosition.Status.EnPassant = nil

	if movement.IsQueenSideCastling != nil || movement.IsKingSideCastling != nil { // Handle castling
		newPosition.Status.CastlingRights.SetQueenCastling(movement.MovingPiece.Color, false)
		newPosition.Status.CastlingRights.SetKingCastling(movement.MovingPiece.Color, false)

		castlingRow := 7
		if movement.MovingPiece.Color == Color_Black {
			castlingRow = 0
		}
		// ERROR CURRENTLY ON SIDE CASTLING QUEEN, ON GETTING MOVEMENTS

		// RECOMMENDATION MAYBE? DONT USE ANY POINTER, ALL VIA POSITION FROM AND TO (WITH COPIES OF TAKING AND MOVING PIECES)
		// IF DOING THIS, WE HAVE TO SAVE ALSO CASTLING POSITION

		if *movement.IsQueenSideCastling {
			rookPiece := newPosition.Board[castlingRow][0]
			kingPiece := newPosition.Board[castlingRow][4]

			// Set new rook
			newPosition.Board[castlingRow][3].Kind = rookPiece.Kind
			newPosition.Board[castlingRow][3].Color = rookPiece.Color

			// Delete old rook
			newPosition.Board[castlingRow][0].Kind = Kind_None
			newPosition.Board[castlingRow][0].Color = Color_None

			// Set new king
			newPosition.Board[castlingRow][2].Kind = kingPiece.Kind
			newPosition.Board[castlingRow][2].Color = kingPiece.Color

			// Delete old king
			newPosition.Board[castlingRow][4].Kind = Kind_None
			newPosition.Board[castlingRow][4].Color = Color_None
		} else if *movement.IsKingSideCastling {
			rookPiece := newPosition.Board[castlingRow][7]
			kingPiece := newPosition.Board[castlingRow][4]

			// Set new rook
			newPosition.Board[castlingRow][5].Kind = rookPiece.Kind
			newPosition.Board[castlingRow][5].Color = rookPiece.Color

			// Delete old rook
			newPosition.Board[castlingRow][7].Kind = Kind_None
			newPosition.Board[castlingRow][7].Color = Color_None

			// Set new king
			newPosition.Board[castlingRow][6].Kind = kingPiece.Kind
			newPosition.Board[castlingRow][6].Color = kingPiece.Color

			// Delete old king
			newPosition.Board[castlingRow][4].Kind = Kind_None
			newPosition.Board[castlingRow][4].Color = Color_None
		}
	} else {
		if movement.MovingPiece.Kind == Kind_Pawn {
			movement.MovingPiece.IsPawnFirstMovement = false

			if *movement.PawnIsDoublePointMovement {
				invertSum := -1
				if movement.MovingPiece.Color == Color_Black {
					invertSum = +1
				}

				newEnPassantPoint := NewPoint(movement.From.I+invertSum, movement.From.J)
				newPosition.Status.EnPassant = &newEnPassantPoint
			}
		} else if movement.MovingPiece.Kind == Kind_King {
			newPosition.Status.CastlingRights.SetQueenCastling(movement.MovingPiece.Color, false)
			newPosition.Status.CastlingRights.SetKingCastling(movement.MovingPiece.Color, false)
		} else if movement.MovingPiece.Kind == Kind_Rook {
			// Check if currently moving rook is from queen or king side
			if newPosition.Status.CastlingRights.CanQueenCastling(movement.MovingPiece.Color) {
				if movement.MovingPiece.Point.J == 0 {
					newPosition.Status.CastlingRights.SetQueenCastling(movement.MovingPiece.Color, false)
				}
			}
			if newPosition.Status.CastlingRights.CanKingCastling(movement.MovingPiece.Color) {
				if movement.MovingPiece.Point.J == 7 {
					newPosition.Status.CastlingRights.SetKingCastling(movement.MovingPiece.Color, false)
				}
			}
		}

		if movement.IsTakingPiece {
			// Do it this wad, so it's en passant compatible
			newPosition.Board[movement.TakingPiece.Point.I][movement.TakingPiece.Point.J].Kind = Kind_None
			newPosition.Board[movement.TakingPiece.Point.I][movement.TakingPiece.Point.J].Color = Color_None

			if movement.TakingPiece.Kind == Kind_Rook {
				if newPosition.Status.CastlingRights.CanQueenCastling(movement.TakingPiece.Color) {
					castlingRow := 7
					if movement.TakingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.TakingPiece.Point.I == castlingRow && movement.TakingPiece.Point.J == 0 {
						newPosition.Status.CastlingRights.SetQueenCastling(movement.TakingPiece.Color, false)
					}
				}
				if newPosition.Status.CastlingRights.CanKingCastling(movement.TakingPiece.Color) {
					castlingRow := 7
					if movement.TakingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.TakingPiece.Point.I == castlingRow && movement.TakingPiece.Point.J == 7 {
						newPosition.Status.CastlingRights.SetKingCastling(movement.TakingPiece.Color, false)
					}
				}
			}
		}

		// movement.MovingPiece.Point = NewPoint(movement.To.I, movement.To.J)
		// p.Board[movement.To.I][movement.To.J] = movement.MovingPiece
		// p.Board[movement.From.I][movement.From.J] = nil

		newPosition.Board[movement.To.I][movement.To.J].Color = movement.MovingPiece.Color
		if movement.PawnPromotionTo == nil {
			// Update data of the new piece
			newPosition.Board[movement.To.I][movement.To.J].Kind = movement.MovingPiece.Kind
		} else {
			// Promote the piece
			newPosition.Board[movement.To.I][movement.To.J].Kind = *movement.PawnPromotionTo
		}

		// Delete this piece's previous position
		newPosition.Board[movement.From.I][movement.From.J].Kind = Kind_None
		newPosition.Board[movement.From.I][movement.From.J].Color = Color_None
		newPosition.Board[movement.From.I][movement.From.J].IsPawnFirstMovement = false
	}

	newPosition.Status.PlayerToMove = Color_White
	if movement.MovingPiece.Color == Color_White {
		newPosition.Status.PlayerToMove = Color_Black
	}

	g.Positions = append(g.Positions, g.CurrentPosition)
	g.CurrentPosition = newPosition
}

func (g *Game) FilterPseudoMovements(movements *[]Movement) []Movement {
	//beginningColor := b.PlayerToMove
	filteredMovements := []Movement{}

	opponentColor := Color_White
	if g.CurrentPosition.Status.PlayerToMove == Color_White {
		opponentColor = Color_Black
	}

	for _, myMovement := range *movements {
		// Check this ???
		// This attack does not have to be evaluated, as its not taking, and pawns cant move diagonally while not attacking
		// The purpose of this is to, when checking a castling, having the movements (these "attacking but not taking") diagonals
		if myMovement.PawnIsAttackingButNotTakingDiagonal != nil && *myMovement.PawnIsAttackingButNotTakingDiagonal {
			continue
		}

		isCastlingLegal := true
		if (myMovement.IsQueenSideCastling != nil && *myMovement.IsQueenSideCastling) ||
			(myMovement.IsKingSideCastling != nil && *myMovement.IsKingSideCastling) {
			currentOpponentPseudo := g.CurrentPosition.GetPseudoMovements(opponentColor)
			for _, opponentMovement := range currentOpponentPseudo {
				if opponentMovement.IsTakingPiece && opponentMovement.TakingPiece.Kind == Kind_King {
					isCastlingLegal = false
					break
				}
			}
			if !isCastlingLegal {
				continue
			}

			// Cols that must not be being attacked
			var colFrom, colTo int

			if myMovement.IsQueenSideCastling != nil && *myMovement.IsQueenSideCastling {
				colFrom = myMovement.From.J - 1
				colTo = myMovement.From.J - 2
			} else if myMovement.IsKingSideCastling != nil && *myMovement.IsKingSideCastling {
				colFrom = myMovement.From.J + 1
				colTo = myMovement.From.J + 2
			}

			isCastlingLegal = g.CurrentPosition.CheckForCastlingLegal(myMovement.From.I, colFrom, colTo, currentOpponentPseudo)
		}

		if !isCastlingLegal {
			continue
		}

		g.MakeMovement(myMovement)
		opponentPseudoMovements := g.CurrentPosition.GetPseudoMovements(opponentColor)

		weGetChecked := g.CurrentPosition.CheckForCheck(opponentPseudoMovements)

		if !weGetChecked {
			filteredMovements = append(filteredMovements, myMovement)
		}
		// Check for check
		//g.UndoMovement(myMovement)
		g.UndoMovement()
	}

	//endColor := b.PlayerToMove

	//fmt.Printf("%c vs %c\n", beginningColor.ToRune(), endColor.ToRune())

	return filteredMovements
}

func (p Position) CheckForCheck(opponentPseudoMovements []Movement) bool {
	for _, movement := range opponentPseudoMovements {
		if movement.IsTakingPiece && movement.TakingPiece.Kind == Kind_King {
			return true
		}
	}

	return false
}

func (p Position) CheckForCastlingLegal(row, colFrom, colTo int, opponentPseudoMovements []Movement) bool {
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
func (p Position) ToFen() string {
	dataFen := ""
	spaceAccum := 0

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if p.Board[i][j].Kind == Kind_None {
				spaceAccum++
			} else {
				if spaceAccum > 0 {
					dataFen = fmt.Sprintf("%s%d", dataFen, spaceAccum)
					spaceAccum = 0
				}
				kindRune := p.Board[i][j].Kind.ToRune()
				if p.Board[i][j].Color == Color_White {
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

	dataFen = fmt.Sprintf("%s %c ", dataFen, p.Status.PlayerToMove.ToRune())

	if p.Status.CastlingRights.CanKingCastling(Color_White) {
		dataFen = fmt.Sprintf("%sK", dataFen)
	}
	if p.Status.CastlingRights.CanQueenCastling(Color_White) {
		dataFen = fmt.Sprintf("%sQ", dataFen)
	}
	if p.Status.CastlingRights.CanKingCastling(Color_Black) {
		dataFen = fmt.Sprintf("%sk", dataFen)
	}
	if p.Status.CastlingRights.CanQueenCastling(Color_Black) {
		dataFen = fmt.Sprintf("%sq", dataFen)
	}

	if p.Status.EnPassant != nil {
		dataFen = fmt.Sprintf("%s %s", dataFen, p.Status.EnPassant.ToAlgebraic())
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
