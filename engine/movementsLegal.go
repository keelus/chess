package engine

func (g *Game) FilterPseudoMovements(movements *[]Movement) []Movement {
	//beginningColor := b.PlayerToMove
	filteredMovements := []Movement{}

	allyColor := g.CurrentPosition.Status.PlayerToMove
	opponentColor := Color_White
	if g.CurrentPosition.Status.PlayerToMove == Color_White {
		opponentColor = Color_Black
	}

	for _, myMovement := range *movements {
		g.MakeMovement(myMovement, false)
		_, opponentAttackMatrix := g.CurrentPosition.GetPseudoMovements(opponentColor, false)

		weGetChecked := g.CurrentPosition.CheckForCheck(allyColor, &opponentAttackMatrix)

		if !weGetChecked {
			filteredMovements = append(filteredMovements, myMovement)
		}

		g.UndoMovement(false)
	}

	return filteredMovements
}

// TODO: Save king's positions (in Position{}, or get them at GetPseudoMovements() to reuse the loop, if possible)
func (p Position) CheckForCheck(allyColor Color, opponentAttackMatrix *[8][8]bool) bool {
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
