package engine

type Game struct {
	Positions       []Position
	CurrentPosition Position

	HasEnded bool
}

func NewGame(fen string) {

}
