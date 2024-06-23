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
	QueenSide map[Color]bool
	KingSide  map[Color]bool
}

func (cr *CastlingRights) DeepCopy() CastlingRights {
	return CastlingRights{
		QueenSide: map[Color]bool{
			Color_White: cr.QueenSide[Color_White],
			Color_Black: cr.QueenSide[Color_Black],
		},
		KingSide: map[Color]bool{
			Color_White: cr.KingSide[Color_White],
			Color_Black: cr.KingSide[Color_Black],
		},
	}
}

type PositionStatus struct {
	PlayerToMove Color

	CastlingRights CastlingRights

	EnPassant *Point
}

func (ps *PositionStatus) DeepCopy() PositionStatus {
	return PositionStatus{
		PlayerToMove:   ps.PlayerToMove,
		CastlingRights: ps.CastlingRights.DeepCopy(),
		EnPassant:      ps.EnPassant,
	}
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
				QueenSide: map[Color]bool{
					Color_White: parsedFen.WhiteCanQueenSideCastling,
					Color_Black: parsedFen.BlackCanQueenSideCastling,
				},
				KingSide: map[Color]bool{
					Color_White: parsedFen.WhiteCanKingSideCastling,
					Color_Black: parsedFen.BlackCanKingSideCastling,
				},
			},

			EnPassant: nil, //TODO
		},
	}
}

// Suppose is legal
func (g *Game) MakeMovement(movement Movement, recomputeLegalMovements bool) {
	newPosition := Position{
		Board:  g.CurrentPosition.Board,
		Status: g.CurrentPosition.Status.DeepCopy(),
	}

	//fmt.Printf("Do: %s\n", movement.ToString())
	newPosition.Status.EnPassant = nil

	if movement.IsQueenSideCastling != nil || movement.IsKingSideCastling != nil { // Handle castling
		newPosition.Status.CastlingRights.QueenSide[movement.MovingPiece.Color] = false
		newPosition.Status.CastlingRights.KingSide[movement.MovingPiece.Color] = false

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

			if *movement.PawnIsDoublePointMovement {
				invertSum := -1
				if movement.MovingPiece.Color == Color_Black {
					invertSum = +1
				}

				// Uint8 from that sum/rest, as it will never be negative in a starting double pawn
				newEnPassantPoint := NewPoint(uint8(int(movement.From.I)+invertSum), movement.From.J)
				newPosition.Status.EnPassant = &newEnPassantPoint
			}
		} else if movement.MovingPiece.Kind == Kind_King {
			newPosition.Status.CastlingRights.QueenSide[movement.MovingPiece.Color] = false
			newPosition.Status.CastlingRights.KingSide[movement.MovingPiece.Color] = false
		} else if movement.MovingPiece.Kind == Kind_Rook {
			// Check if currently moving rook is from queen or king side
			if newPosition.Status.CastlingRights.QueenSide[movement.MovingPiece.Color] {
				if movement.MovingPiece.Point.J == 0 {
					newPosition.Status.CastlingRights.QueenSide[movement.MovingPiece.Color] = false
				}
			}
			if newPosition.Status.CastlingRights.KingSide[movement.MovingPiece.Color] {
				if movement.MovingPiece.Point.J == 7 {
					newPosition.Status.CastlingRights.KingSide[movement.MovingPiece.Color] = false
				}
			}
		}

		if movement.IsTakingPiece {
			// Do it this wad, so it's en passant compatible
			newPosition.Board[movement.TakingPiece.Point.I][movement.TakingPiece.Point.J].Kind = Kind_None
			newPosition.Board[movement.TakingPiece.Point.I][movement.TakingPiece.Point.J].Color = Color_None

			if movement.TakingPiece.Kind == Kind_Rook {
				if newPosition.Status.CastlingRights.QueenSide[movement.TakingPiece.Color] {
					castlingRow := uint8(7)
					if movement.TakingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.TakingPiece.Point.I == castlingRow && movement.TakingPiece.Point.J == 0 {
						newPosition.Status.CastlingRights.QueenSide[movement.TakingPiece.Color] = false
					}
				}
				if newPosition.Status.CastlingRights.KingSide[movement.TakingPiece.Color] {
					castlingRow := uint8(7)
					if movement.TakingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.TakingPiece.Point.I == castlingRow && movement.TakingPiece.Point.J == 7 {
						newPosition.Status.CastlingRights.KingSide[movement.TakingPiece.Color] = false
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
	}

	newPosition.Status.PlayerToMove = Color_White
	if movement.MovingPiece.Color == Color_White {
		newPosition.Status.PlayerToMove = Color_Black
	}

	g.Positions = append(g.Positions, g.CurrentPosition)
	g.CurrentPosition = newPosition
	if recomputeLegalMovements {
		g.ComputeLegalMovements()
	}
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
			var colFrom, colTo uint8

			// TODO: Hard code positions to less operations
			if myMovement.IsQueenSideCastling != nil && *myMovement.IsQueenSideCastling {
				colFrom = uint8(int(myMovement.From.J) - 1)
				colTo = uint8(int(myMovement.From.J) - 2)
			} else if myMovement.IsKingSideCastling != nil && *myMovement.IsKingSideCastling {
				colFrom = uint8(int(myMovement.From.J) + 1)
				colTo = uint8(int(myMovement.From.J) + 2)
			}

			isCastlingLegal = g.CurrentPosition.CheckForCastlingLegal(myMovement.From.I, colFrom, colTo, currentOpponentPseudo)
		}

		if !isCastlingLegal {
			continue
		}

		g.MakeMovement(myMovement, false)
		opponentPseudoMovements := g.CurrentPosition.GetPseudoMovements(opponentColor)

		weGetChecked := g.CurrentPosition.CheckForCheck(opponentPseudoMovements)

		if !weGetChecked {
			filteredMovements = append(filteredMovements, myMovement)
		}
		// Check for check
		//g.UndoMovement(myMovement)
		g.UndoMovement(false)
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

func (p Position) CheckForCastlingLegal(row, colFrom, colTo uint8, opponentPseudoMovements []Movement) bool {
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

	if p.Status.CastlingRights.KingSide[Color_White] {
		dataFen = fmt.Sprintf("%sK", dataFen)
	}
	if p.Status.CastlingRights.QueenSide[Color_White] {
		dataFen = fmt.Sprintf("%sQ", dataFen)
	}
	if p.Status.CastlingRights.KingSide[Color_Black] {
		dataFen = fmt.Sprintf("%sk", dataFen)
	}
	if p.Status.CastlingRights.QueenSide[Color_Black] {
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
