package engine

func (g *Game) FilterPseudoMovements(movements *[]Movement) []Movement {
	//beginningColor := b.PlayerToMove
	filteredMovements := []Movement{}

	allyColor := g.currentPosition.Status.PlayerToMove
	opponentColor := Color_White
	if g.currentPosition.Status.PlayerToMove == Color_White {
		opponentColor = Color_Black
	}

	for _, myMovement := range *movements {
		// TODO: DO a simulateMovement
		g.simulateMovement(myMovement)
		_, opponentAttackMatrix := g.currentPosition.GetPseudoMovements(opponentColor, false)

		weGetChecked := g.currentPosition.checkForCheck(allyColor, &opponentAttackMatrix)

		if !weGetChecked {
			filteredMovements = append(filteredMovements, myMovement)
		}

		g.undoSimulatedMovement()
	}

	return filteredMovements
}

// TODO: Save king's positions (in Position{}, or get them at GetPseudoMovements() to reuse the loop, if possible)
func (p Position) checkForCheck(allyColor Color, opponentAttackMatrix *[8][8]bool) bool {
	for i := uint8(0); i < 8; i++ {
		for j := uint8(0); j < 8; j++ {
			if p.Board[i][j].Kind == Kind_King && p.Board[i][j].Color == allyColor {
				return opponentAttackMatrix[i][j]
			}
		}
	}

	// Won't get here
	return false
}
