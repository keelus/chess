package position

type Position struct {
	I int
	J int
}

func NewPosition(i, j int) Position {
	return Position{
		I: i,
		J: j,
	}
}
