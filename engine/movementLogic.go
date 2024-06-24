package engine

import "fmt"

func (g *Game) ComputeLegalMovements() {
	pseudoMovements, _ := g.CurrentPosition.GetPseudoMovements(g.CurrentPosition.Status.PlayerToMove, true)
	legalMovements := g.FilterPseudoMovements(&pseudoMovements)
	g.ComputedLegalMovements = legalMovements
}

func (g *Game) MakeMovement(movement Movement, recomputeLegalMovements bool) {
	// Suppose is legal. TODO: Check if illegal with the computed list
	newPosition := Position{
		Board:  g.CurrentPosition.Board,
		Status: g.CurrentPosition.Status.DeepCopy(),
	}

	newPosition.Status.EnPassant = nil

	if movement.IsQueenSideCastling || movement.IsKingSideCastling {
		newPosition.Status.CastlingRights.QueenSide[movement.MovingPiece.Color] = false
		newPosition.Status.CastlingRights.KingSide[movement.MovingPiece.Color] = false

		castlingRow := 7
		if movement.MovingPiece.Color == Color_Black {
			castlingRow = 0
		}

		if movement.IsQueenSideCastling {
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
		} else if movement.IsKingSideCastling {
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
			if movement.PawnIsDoublePointMovement {
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

func (g *Game) UndoMovement(recomputeLegalMovements bool) {
	if len(g.Positions) == 0 {
		fmt.Println("Can't undo more.")
		return
	} else {
		g.CurrentPosition = g.Positions[len(g.Positions)-1]
		g.Positions = g.Positions[:len(g.Positions)-1]

		if recomputeLegalMovements {
			g.ComputeLegalMovements()
		}
	}
}
