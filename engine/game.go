package engine

import (
	"errors"
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

func (g *Game) GetPieceAt(square Square) Piece {
	return g.currentPosition.Board[square.I][square.J]
}

func (g *Game) GetPieceAtAlgebraic(algebraic string) (Piece, error) {
	square, err := NewSquareFromAlgebraic(algebraic)
	if err != nil {
		return Piece{}, err
	}

	return g.GetPieceAt(square), nil
}

func (g Game) Turn() Color {
	return g.currentPosition.Status.PlayerToMove
}

// public func
func (g *Game) GetLegalMovements() []string {
	movementList := make([]string, len(g.computedLegalMovements))
	for i, legalMovement := range g.computedLegalMovements {
		movementList[i] = legalMovement.Algebraic()
	}
	return movementList
}

func (g *Game) GetLegalMovementsOfPiece(square Square) []string {
	movementList := make([]string, 0)
	for _, legalMovement := range g.computedLegalMovements {
		if legalMovement.from.I == square.I && legalMovement.from.J == square.J {
			movementList = append(movementList, legalMovement.Algebraic())
		}
	}
	return movementList
}

func (g *Game) getLegalMovementsMap() map[string]Movement {
	movementMap := make(map[string]Movement)
	for _, legalMovement := range g.computedLegalMovements {
		movementMap[legalMovement.Algebraic()] = legalMovement
	}
	return movementMap
}

func (g *Game) IsMovementLegal(movement string) bool {
	legalMovementsMap := g.getLegalMovementsMap()
	_, ok := legalMovementsMap[movement]
	return ok
}

// public func
func (g *Game) MakeMovement(movement string) error {
	legalMovementsMap := g.getLegalMovementsMap()
	if movementValue, ok := legalMovementsMap[movement]; ok {
		g.movementHistory = append(g.movementHistory, movementValue)
		g.forceMovement(movementValue, true)
		return nil
	}

	return errors.New("That movement is not allowed.")
}

// TODO: Public field / getters
func (g *Game) MovementInformation(movement string) (Movement, error) {
	legalMovementsMap := g.getLegalMovementsMap()
	if movementValue, ok := legalMovementsMap[movement]; ok {
		return movementValue, nil
	}

	return Movement{}, errors.New("That movement is not allowed or is invalid.")
}

func (g *Game) PrintStartingFen() string {
	pos, _ := g.GetPositionAtIndex(0)
	return pos.Fen()
}

func (g *Game) PrintCurrentFen() string {
	return g.GetCurrentPosition().Fen()
}

func (g *Game) GetCurrentPosition() Position {
	return g.currentPosition
}

func (g *Game) GetPositionAtIndex(index int) (Position, error) {
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
