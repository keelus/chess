package engine

import "fmt"

type Game struct {
	Positions       []Position
	CurrentPosition Position

	HasEnded bool

	ComputedLegalMovements []Movement
}

func NewGame(fen string) Game {
	if fen == "" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}

	newGame := Game{
		Positions:       make([]Position, 0),
		CurrentPosition: NewPositionFromFen(fen),

		HasEnded: false,
	}

	newGame.ComputeLegalMovements()
	return newGame
}

func (g *Game) ComputeLegalMovements() {
	pseudoMovements, _ := g.CurrentPosition.GetPseudoMovements(g.CurrentPosition.Status.PlayerToMove, true)
	legalMovements := g.FilterPseudoMovements(&pseudoMovements)
	g.ComputedLegalMovements = legalMovements
}

func (g *Game) UndoMovement(recomputeLegalMovements bool) {
	if len(g.Positions) == 0 {
		fmt.Println("Can undo more.")
		return
	} else {
		g.CurrentPosition = g.Positions[len(g.Positions)-1]
		g.Positions = g.Positions[:len(g.Positions)-1]

		if recomputeLegalMovements {
			g.ComputeLegalMovements()
		}
	}
}

func (g *Game) GetPieceAt(i, j uint8) Piece {
	return g.CurrentPosition.Board.GetPieceAt(i, j)
}

func (g Game) GetPlayerToMove() Color {
	return g.CurrentPosition.Status.PlayerToMove
}

func (g *Game) ForceSetPlayerToMove(color Color) {
	g.CurrentPosition.Status.PlayerToMove = color
}
