package engine

func (b *Board) Perft(initialDepth, depth int, currentMove string, posMap map[string]int) int {
	var nMoves, i int
	nodes := 0

	if depth == 0 {
		return 1
	}

	moveList := b.GetLegalMovements(b.PlayerToMove)
	nMoves = len(moveList)

	for i = 0; i < nMoves; i++ {
		b.MakeMovement(moveList[i])
		nodes += b.Perft(initialDepth, depth-1, moveList[i].ToAlgebraic(), posMap)
		b.UndoMovement(moveList[i])
	}

	if depth == initialDepth-1 && currentMove != "" {
		//fmt.Printf("%s: %d\n", currentAlgebraicMovement, nodes)
		posMap[currentMove] = nodes
	}

	return nodes
}
