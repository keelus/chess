package engine

import "fmt"

func (g *Game) ComputeLegalMovements() {
	pseudoMovements, _ := g.currentPosition.GetPseudoMovements(g.currentPosition.Status.PlayerToMove, true)
	legalMovements := g.FilterPseudoMovements(&pseudoMovements)
	g.computedLegalMovements = legalMovements
}

// Used by perft
func (g *Game) simulateMovement(movement Movement) {
	g.forceMovement(movement, false)
}

func (g *Game) undoSimulatedMovement() {
	g.undoMovement(false)
}

func (g *Game) forceMovement(movement Movement, recomputeLegalMovements bool) {
	// Suppose is legal. TODO: Check if illegal with the computed list
	newPosition := Position{
		Board:  g.currentPosition.Board,
		Status: g.currentPosition.Status.clone(),
	}

	newPosition.Status.EnPassant = nil

	if movement.isQueenSideCastling || movement.isKingSideCastling {
		newPosition.Status.CastlingRights.QueenSide[movement.movingPiece.Color] = false
		newPosition.Status.CastlingRights.KingSide[movement.movingPiece.Color] = false

		castlingRow := 7
		if movement.movingPiece.Color == Color_Black {
			castlingRow = 0
		}

		if movement.isQueenSideCastling {
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
		} else if movement.isKingSideCastling {
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
		if movement.movingPiece.Kind == Kind_Pawn {
			if movement.pawnIsDoubleSquareMovement {
				invertSum := -1
				if movement.movingPiece.Color == Color_Black {
					invertSum = +1
				}

				// Uint8 from that sum/rest, as it will never be negative in a starting double pawn
				newEnPassantSquare := newSquare(uint8(int(movement.from.I)+invertSum), movement.from.J)
				newPosition.Status.EnPassant = &newEnPassantSquare
			}
		} else if movement.movingPiece.Kind == Kind_King {
			newPosition.Status.CastlingRights.QueenSide[movement.movingPiece.Color] = false
			newPosition.Status.CastlingRights.KingSide[movement.movingPiece.Color] = false
		} else if movement.movingPiece.Kind == Kind_Rook {
			// Check if currently moving rook is from queen or king side
			if newPosition.Status.CastlingRights.QueenSide[movement.movingPiece.Color] {
				if movement.movingPiece.Square.J == 0 {
					newPosition.Status.CastlingRights.QueenSide[movement.movingPiece.Color] = false
				}
			}
			if newPosition.Status.CastlingRights.KingSide[movement.movingPiece.Color] {
				if movement.movingPiece.Square.J == 7 {
					newPosition.Status.CastlingRights.KingSide[movement.movingPiece.Color] = false
				}
			}
		}

		if movement.isTakingPiece {
			newPosition.Board[movement.takingPiece.Square.I][movement.takingPiece.Square.J].Kind = Kind_None
			newPosition.Board[movement.takingPiece.Square.I][movement.takingPiece.Square.J].Color = Color_None

			if movement.takingPiece.Kind == Kind_Rook {
				if newPosition.Status.CastlingRights.QueenSide[movement.takingPiece.Color] {
					castlingRow := uint8(7)
					if movement.takingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.takingPiece.Square.I == castlingRow && movement.takingPiece.Square.J == 0 {
						newPosition.Status.CastlingRights.QueenSide[movement.takingPiece.Color] = false
					}
				}
				if newPosition.Status.CastlingRights.KingSide[movement.takingPiece.Color] {
					castlingRow := uint8(7)
					if movement.takingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.takingPiece.Square.I == castlingRow && movement.takingPiece.Square.J == 7 {
						newPosition.Status.CastlingRights.KingSide[movement.takingPiece.Color] = false
					}
				}
			}
		}

		newPosition.Board[movement.to.I][movement.to.J].Color = movement.movingPiece.Color
		if movement.pawnPromotionTo == nil {
			// Update data of the new piece
			newPosition.Board[movement.to.I][movement.to.J].Kind = movement.movingPiece.Kind
		} else {
			// Promote the piece
			newPosition.Board[movement.to.I][movement.to.J].Kind = *movement.pawnPromotionTo
		}

		// Delete this piece's previous position
		newPosition.Board[movement.from.I][movement.from.J].Kind = Kind_None
		newPosition.Board[movement.from.I][movement.from.J].Color = Color_None
	}

	newPosition.Status.PlayerToMove = Color_White
	if movement.movingPiece.Color == Color_White {
		newPosition.Status.PlayerToMove = Color_Black
	}

	g.positions = append(g.positions, g.currentPosition)
	g.currentPosition = newPosition
	if recomputeLegalMovements {
		g.ComputeLegalMovements()
	}
}

func (g *Game) undoMovement(recomputeLegalMovements bool) {
	if len(g.positions) == 0 {
		fmt.Println("Can't undo more.")
		return
	} else {
		g.currentPosition = g.positions[len(g.positions)-1]
		g.positions = g.positions[:len(g.positions)-1]

		if recomputeLegalMovements {
			g.ComputeLegalMovements()
		}
	}
}
