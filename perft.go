package chess

import (
	"fmt"
)

// Used for testing via game_test.go
func (g *Game) perft(initialDepth, depth int, currentMove string, positionVerbose bool) int {
	var nMoves, i int
	nodes := 0

	if depth == 0 {
		if positionVerbose && initialDepth == 1 {
			fmt.Printf("\t%s: %d\n", currentMove, 1)
		}
		return 1
	}

	g.computeLegalMovements()
	moveList := g.computedLegalMovements
	nMoves = len(moveList)

	for i = 0; i < nMoves; i++ {
		g.simulateMovement(moveList[i])
		nodes += g.perft(initialDepth, depth-1, moveList[i].Algebraic(), positionVerbose)
		g.undoSimulatedMovement()
	}

	if positionVerbose && depth == initialDepth-1 && currentMove != "" {
		fmt.Printf("\t%s: %d\n", currentMove, nodes)
		fmt.Printf("\t^ %s\n", g.currentPosition.Fen())
	}

	return nodes
}
