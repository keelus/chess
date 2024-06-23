package engine

type Game struct {
	Positions       []Position
	CurrentPosition Position

	HasEnded bool
}

func NewGame(fen string) Game {
	if fen == "" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}

	return Game{
		Positions:       make([]Position, 0),
		CurrentPosition: NewPositionFromFen(fen),

		HasEnded: false,
	}
}

func (g *Game) GetLegalMovements() []Movement {
	return g.CurrentPosition.GetLegalMovements(g.CurrentPosition.Status.PlayerToMove)
}

func (g *Game) MakeMovement(movement Movement) {
	g.CurrentPosition.MakeMovement(movement)
}

func (g *Game) UndoMovement(movement Movement) {
	g.CurrentPosition.UndoMovement(movement)
}

func (g *Game) GetPieceAt(i, j int) Piece {
	return g.CurrentPosition.Board.GetPieceAt(i, j)
}

func (g Game) GetPlayerToMove() Color {
	return g.CurrentPosition.Status.PlayerToMove
}

func (g *Game) ForceSetPlayerToMove(color Color) {
	g.CurrentPosition.Status.PlayerToMove = color
}
