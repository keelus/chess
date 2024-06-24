package engine

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

func (g *Game) GetPieceAt(i, j uint8) Piece {
	return g.CurrentPosition.Board.GetPieceAt(i, j)
}

func (g Game) GetPlayerToMove() Color {
	return g.CurrentPosition.Status.PlayerToMove
}

func (g *Game) ForceSetPlayerToMove(color Color) {
	g.CurrentPosition.Status.PlayerToMove = color
}

func (g *Game) GetLegalMovements() []Movement {
	return g.ComputedLegalMovements
}
