package chess

import (
	"errors"
	"fmt"
)

type Game struct {
	positions       []Position
	currentPosition Position

	//HasEnded bool

	computedLegalMovements []Movement

	outcome Outcome

	positionMap     map[string]int // only used by api, not perft (as a map of fen positions is too slow to compute)
	movementHistory []Movement     //only used by api, not perft
}

func NewGame(fen string) Game {
	if fen == "" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}

	newGame := Game{
		positions:       make([]Position, 0),
		currentPosition: newPositionFromFen(fen),

		positionMap: make(map[string]int),

		outcome: Outcome_None,
		//HasEnded: false,
	}

	newGame.ComputeLegalMovements()
	return newGame
}

func (g *Game) PieceAtSquareAlgebraic(algebraic string) (Piece, error) {
	square, err := NewSquareFromAlgebraic(algebraic)
	if err != nil {
		return Piece{}, err
	}

	return g.PieceAtSquare(square), nil
}

func (g *Game) PieceAtSquare(square Square) Piece {
	return g.currentPosition.board[square.I][square.J]
}

func (g Game) Turn() Color {
	return g.currentPosition.playerToMove
}

func (g *Game) LegalMovements() []Movement {
	return g.computedLegalMovements
}

func (g Game) LegalMovementsAlgebraic() []string {
	movementList := make([]string, len(g.computedLegalMovements))
	for i, legalMovement := range g.computedLegalMovements {
		movementList[i] = legalMovement.Algebraic()
	}
	return movementList
}

func (g Game) LegalMovementsOfPiece(square Square) []Movement {
	legalMovementsOfPiece := make([]Movement, 0)
	for _, legalMovement := range g.computedLegalMovements {
		if legalMovement.fromSq.IsEqualTo(square) {
			legalMovementsOfPiece = append(legalMovementsOfPiece, legalMovement)
		}
	}

	return legalMovementsOfPiece
}

func (g Game) LegalMovementsOfPieceAlgebraic(square Square) []string {
	legalMovementsOfPiece := make([]string, 0)
	for _, legalMovement := range g.computedLegalMovements {
		if legalMovement.fromSq.IsEqualTo(square) {
			legalMovementsOfPiece = append(legalMovementsOfPiece, legalMovement.Algebraic())
		}
	}

	return legalMovementsOfPiece
}

func (g *Game) IsMovementLegal(movement Movement) bool {
	return g.IsMovementLegalAlgebraic(movement.Algebraic())
}

func (g *Game) IsMovementLegalAlgebraic(algebraicMovement string) bool {
	for _, legalMovement := range g.computedLegalMovements {
		if legalMovement.Algebraic() == algebraicMovement {
			return true
		}
	}

	return false
}

func (g *Game) MakeMovement(movement Movement) error {
	return g.MakeMovementAlgebraic(movement.Algebraic())
}

func (g *Game) MakeMovementAlgebraic(algebraicMovement string) error {
	for _, legalMovement := range g.computedLegalMovements {
		if legalMovement.Algebraic() == algebraicMovement {
			g.movementHistory = append(g.movementHistory, legalMovement)
			g.forceMovement(legalMovement, true)
			return nil
		}
	}

	return errors.New("That movement is not allowed.")
}

func (g *Game) StartingFen() string {
	pos, _ := g.PositionAtIndex(0)
	return pos.Fen()
}

func (g *Game) CurrentFen() string {
	return g.CurrentPosition().Fen()
}

func (g *Game) CurrentPosition() Position {
	return g.currentPosition
}

func (g *Game) PositionAtIndex(index int) (Position, error) {
	if index >= 0 && index < len(g.positions) {
		return g.positions[index], nil
	}

	return Position{}, errors.New("That index is invalid or out of bounds.")
}

func (g *Game) MovementHistory() []Movement {
	return g.movementHistory
}

type Outcome string

const (
	Outcome_None            Outcome = "None"
	Outcome_Checkmate_White         = "Checkmate_White"
	Outcome_Checkmate_Black         = "Checkmate_Black"

	Outcome_Draw_Stalemate = "Stalemate"
	//Outcome_Draw_InsufficientMaterial = "Insufficient material"
	Outcome_Draw_50Move = "Fifty move rule"
	Outcome_Draw_3Rep   = "Threefold repetition"
)

func (g *Game) Outcome() Outcome {
	return g.outcome
}

func (g *Game) FilterPseudoMovements(movements *[]Movement) []Movement {
	//beginningColor := b.playerToMove
	filteredMovements := []Movement{}

	// Ensure we use this colors (and not others, as CurrentPosition will change on GetPseudoMovements)
	allyColor := g.currentPosition.playerToMove
	opponentColor := g.currentPosition.playerToMove.Opposite()

	for _, myMovement := range *movements {
		// TODO: DO a simulateMovement
		g.simulateMovement(myMovement)
		_, opponentAttackMatrix := g.currentPosition.computePseudoMovements(opponentColor, false)

		weGetChecked := g.currentPosition.checkForCheck(allyColor, &opponentAttackMatrix)

		if !weGetChecked {
			filteredMovements = append(filteredMovements, myMovement)
		}

		g.undoSimulatedMovement()
	}

	return filteredMovements
}

func (g *Game) ComputeLegalMovements() {
	pseudoMovements, _ := g.currentPosition.computePseudoMovements(g.currentPosition.playerToMove, true)
	legalMovements := g.FilterPseudoMovements(&pseudoMovements)
	g.computedLegalMovements = legalMovements
}

// Used by perft
func (g *Game) simulateMovement(movement Movement) {
	g.forceMovement(movement, false)
}

func (g *Game) forceMovement(movement Movement, recomputeLegalMovements bool) {
	newPosition := g.currentPosition.clone()

	if movement.isQueenSideCastling || movement.isKingSideCastling {
		newPosition.castlingRights.queenSide[movement.movingPiece.Color] = false
		newPosition.castlingRights.kingSide[movement.movingPiece.Color] = false

		castlingRow := 7
		if movement.movingPiece.Color == Color_Black {
			castlingRow = 0
		}

		if movement.isQueenSideCastling {
			rookPiece := newPosition.board[castlingRow][0]
			kingPiece := newPosition.board[castlingRow][4]

			// Set new rook
			newPosition.board[castlingRow][3].Kind = rookPiece.Kind
			newPosition.board[castlingRow][3].Color = rookPiece.Color

			// Delete old rook
			newPosition.board[castlingRow][0].Kind = Kind_None
			newPosition.board[castlingRow][0].Color = Color_None

			// Set new king
			newPosition.board[castlingRow][2].Kind = kingPiece.Kind
			newPosition.board[castlingRow][2].Color = kingPiece.Color

			// Delete old king
			newPosition.board[castlingRow][4].Kind = Kind_None
			newPosition.board[castlingRow][4].Color = Color_None
		} else if movement.isKingSideCastling {
			rookPiece := newPosition.board[castlingRow][7]
			kingPiece := newPosition.board[castlingRow][4]

			// Set new rook
			newPosition.board[castlingRow][5].Kind = rookPiece.Kind
			newPosition.board[castlingRow][5].Color = rookPiece.Color

			// Delete old rook
			newPosition.board[castlingRow][7].Kind = Kind_None
			newPosition.board[castlingRow][7].Color = Color_None

			// Set new king
			newPosition.board[castlingRow][6].Kind = kingPiece.Kind
			newPosition.board[castlingRow][6].Color = kingPiece.Color

			// Delete old king
			newPosition.board[castlingRow][4].Kind = Kind_None
			newPosition.board[castlingRow][4].Color = Color_None
		}
	} else {
		if movement.movingPiece.Kind == Kind_Pawn {
			if movement.pawnIsDoubleSquareMovement {
				invertSum := -1
				if movement.movingPiece.Color == Color_Black {
					invertSum = +1
				}

				// Uint8 from that sum/rest, as it will never be negative in a starting double pawn
				newEnPassantSquare := newSquare(uint8(int(movement.fromSq.I)+invertSum), movement.fromSq.J)
				newPosition.enPassantSq = &newEnPassantSquare
			}
		} else if movement.movingPiece.Kind == Kind_King {
			newPosition.castlingRights.queenSide[movement.movingPiece.Color] = false
			newPosition.castlingRights.kingSide[movement.movingPiece.Color] = false
		} else if movement.movingPiece.Kind == Kind_Rook {
			// Check if currently moving rook is from queen or king side
			if newPosition.castlingRights.queenSide[movement.movingPiece.Color] {
				if movement.movingPiece.Square.J == 0 {
					newPosition.castlingRights.queenSide[movement.movingPiece.Color] = false
				}
			}
			if newPosition.castlingRights.kingSide[movement.movingPiece.Color] {
				if movement.movingPiece.Square.J == 7 {
					newPosition.castlingRights.kingSide[movement.movingPiece.Color] = false
				}
			}
		}

		if movement.isTakingPiece {
			newPosition.board[movement.takingPiece.Square.I][movement.takingPiece.Square.J].Kind = Kind_None
			newPosition.board[movement.takingPiece.Square.I][movement.takingPiece.Square.J].Color = Color_None

			if recomputeLegalMovements {
				newPosition.captures = append(newPosition.captures, movement.takingPiece)
			}

			if movement.takingPiece.Kind == Kind_Rook {
				if newPosition.castlingRights.queenSide[movement.takingPiece.Color] {
					castlingRow := uint8(7)
					if movement.takingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.takingPiece.Square.I == castlingRow && movement.takingPiece.Square.J == 0 {
						newPosition.castlingRights.queenSide[movement.takingPiece.Color] = false
					}
				}
				if newPosition.castlingRights.kingSide[movement.takingPiece.Color] {
					castlingRow := uint8(7)
					if movement.takingPiece.Color == Color_Black {
						castlingRow = 0
					}

					if movement.takingPiece.Square.I == castlingRow && movement.takingPiece.Square.J == 7 {
						newPosition.castlingRights.kingSide[movement.takingPiece.Color] = false
					}
				}
			}
		}

		newPosition.board[movement.toSq.I][movement.toSq.J].Color = movement.movingPiece.Color
		if movement.pawnPromotionTo == nil {
			// Update data of the new piece
			newPosition.board[movement.toSq.I][movement.toSq.J].Kind = movement.movingPiece.Kind
		} else {
			// Promote the piece
			newPosition.board[movement.toSq.I][movement.toSq.J].Kind = *movement.pawnPromotionTo
		}

		// Delete this piece's previous position
		newPosition.board[movement.fromSq.I][movement.fromSq.J].Kind = Kind_None
		newPosition.board[movement.fromSq.I][movement.fromSq.J].Color = Color_None
	}

	// Handle halfmove clock
	if movement.movingPiece.Kind == Kind_Pawn || movement.isTakingPiece {
		newPosition.halfmoveClock = 0
	} else {
		newPosition.halfmoveClock++
	}

	// Handle fullmove counter
	if movement.movingPiece.Color == Color_White {
		newPosition.playerToMove = Color_Black
	} else if movement.movingPiece.Color == Color_Black {
		newPosition.playerToMove = Color_White
		newPosition.fullmoveCounter++
	}

	// Switch positions
	g.positions = append(g.positions, g.currentPosition)
	g.currentPosition = newPosition

	if recomputeLegalMovements {
		if _, ok := g.positionMap[g.currentPosition.Fen()]; !ok {
			g.positionMap[g.currentPosition.Fen()] = 1
		} else {
			g.positionMap[g.currentPosition.Fen()]++
			//fmt.Println(g.positionMap[g.currentPosition.Fen()])
			// if recomputeLegalMovements {
			// 	fmt.Printf("%s has %d\n", g.currentPosition.Fen(), g.positionMap[g.currentPosition.Fen()])
			// }
			if g.positionMap[g.currentPosition.Fen()] == 3 {
				//fmt.Printf("%s ends with %d\n", g.currentPosition.Fen(), g.positionMap[g.currentPosition.Fen()])
				g.Terminate(Outcome_Draw_3Rep)
			}
		}
	}

	// Recomputing will take place:
	// 		- After making a move (via game.MakeMove())
	// 		- Manually via computeLegalMovements(), called by Perft
	if recomputeLegalMovements {
		g.ComputeLegalMovements()

		if len(g.computedLegalMovements) == 0 {
			_, opponentAttackMatrix := g.currentPosition.computePseudoMovements(g.currentPosition.playerToMove.Opposite(), false)
			isGettingChecked := g.currentPosition.checkForCheck(g.currentPosition.playerToMove, &opponentAttackMatrix)

			if isGettingChecked {
				if g.currentPosition.playerToMove == Color_White {
					g.Terminate(Outcome_Checkmate_Black)
				} else {
					g.Terminate(Outcome_Checkmate_White)
				}
			} else {
				g.Terminate(Outcome_Draw_Stalemate)
			}
		} else if g.currentPosition.halfmoveClock >= 100 {
			g.Terminate(Outcome_Draw_50Move)
		}
	}
}

func (g *Game) Terminate(outcome Outcome) {
	g.outcome = outcome
}

func (g *Game) undoSimulatedMovement() {
	if len(g.positions) == 0 {
		fmt.Println("Can't undo more.")
		return
	} else {
		g.currentPosition = g.positions[len(g.positions)-1]
		g.positions = g.positions[:len(g.positions)-1]
	}
}
