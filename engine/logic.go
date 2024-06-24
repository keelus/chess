package engine

func (g *Game) GetLegalMovements() []Movement {
	return g.ComputedLegalMovements
}

func (p Position) GetPseudoMovements(color Color, doCastlingCheck bool) ([]Movement, [8][8]bool) {
	movements := make([]Movement, 0, 256)
	var attackMatrix [8][8]bool
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			attackMatrix[i][j] = false
		}
	}

	for _, row := range p.Board {
		for _, piece := range row {
			if piece.Color == color {
				p.GetPiecePseudoMovements(piece, &movements, doCastlingCheck, &attackMatrix)
			}
		}
	}

	return movements, attackMatrix
}

func (p Position) GetPiecePseudoMovements(piece Piece, movements *[]Movement, doCastlingCheck bool, attackMatrix *[8][8]bool) {
	switch piece.Kind {
	case Kind_Bishop:
		p.getDiagonalPseudoMovements(piece, movements, attackMatrix)
		return
	case Kind_Rook:
		p.getOrthogonalPseudoMovements(piece, movements, attackMatrix)
		return
	case Kind_Queen:
		p.getDiagonalPseudoMovements(piece, movements, attackMatrix)
		p.getOrthogonalPseudoMovements(piece, movements, attackMatrix)
		return
	case Kind_King:
		for i := int8(-1); i < 2; i++ {
			for j := int8(-1); j < 2; j++ {
				if i == 0 && j == 0 {
					continue
				}

				finalI, finalJ := int8(piece.Point.I)+i, int8(piece.Point.J)+j
				if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
					row := uint8(finalI)
					col := uint8(finalJ)
					pieceAt := p.Board.GetPieceAt(uint8(finalI), uint8(finalJ))

					if pieceAt.Kind == Kind_None {
						attackMatrix[row][col] = true
						*movements = append(*movements,
							*NewMovement(piece,
								piece.Point,
								NewPoint(row, col),
							))
					} else if pieceAt.Color != piece.Color {
						attackMatrix[row][col] = true
						*movements = append(*movements,
							*NewMovement(piece,
								piece.Point,
								NewPoint(row, col),
							).WithTakingPiece(pieceAt))
					}
				}
			}
		}

		opponentColor := Color_White
		if piece.Color == Color_White {
			opponentColor = Color_Black
		}

		if doCastlingCheck {
			castlingRow := 7
			if piece.Color == Color_Black {
				castlingRow = 0
			}

			_, enemyAttackBoard := p.GetPseudoMovements(opponentColor, false)
			// If king is not in check, continue
			if enemyAttackBoard[piece.Point.I][piece.Point.J] == false {
				if p.Status.CastlingRights.QueenSide[piece.Color] {
					// Check if space to rook is empty
					canCastle := true
					for j := piece.Point.J - 1; j >= piece.Point.J-3; j-- {
						if p.Board.GetPieceAt(piece.Point.I, j).Kind != Kind_None {
							canCastle = false
							break
						}
					}

					if canCastle {
						// Extra check: Position is not being attacked by enemy
						// On queen side, positions that cannot be attacked to castle:
						// On the left of the King, in d1 and c1 [castlingRow, 2], [castlingRow, 3]
						if enemyAttackBoard[castlingRow][2] == false && enemyAttackBoard[castlingRow][3] == false {
							*movements = append(*movements,
								*NewMovement(piece,
									piece.Point,
									NewPoint(piece.Point.I, piece.Point.J-2),
								).WithCastling(true, false))
						}
					}
				}

				if p.Status.CastlingRights.KingSide[piece.Color] {
					// Check if space to rook is empty
					canCastle := true
					for j := piece.Point.J + 1; j < 7; j++ {
						if p.Board.GetPieceAt(piece.Point.I, j).Kind != Kind_None {
							canCastle = false
							break
						}
					}

					if canCastle {
						// Extra check: Position is not being attacked by enemy
						// On King side, positions that cannot be attacked to castle:
						// On the right of the King, in f1 and g1 [castlingRow, 5], [castlingRow, 6]
						if enemyAttackBoard[castlingRow][5] == false && enemyAttackBoard[castlingRow][6] == false {
							*movements = append(*movements,
								*NewMovement(piece,
									piece.Point,
									NewPoint(piece.Point.I, piece.Point.J+2),
								).WithCastling(false, true))
						}
					}
				}
			}
		}

		return
	case Kind_Knight:
		dirs := [8][2]int{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}} // {i, j} -> From TopRight, clockwise untill TopRight (bottom left)

		for _, dir := range dirs {
			finalI, finalJ := int(piece.Point.I)+dir[0], int(piece.Point.J)+dir[1]

			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				row := uint8(finalI)
				col := uint8(finalJ)

				pieceAt := p.Board.GetPieceAt(row, col)

				if pieceAt.Kind == Kind_None {
					attackMatrix[row][col] = true
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(row, col),
						))
				} else if pieceAt.Color != piece.Color {
					attackMatrix[row][col] = true
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(row, col),
						).WithTakingPiece(pieceAt))
				}
			}
		}

		return

	case Kind_Pawn:
		invertMult := 1
		if piece.Color == Color_Black {
			invertMult = -1
		}

		startPawnRow := uint8(6)
		if piece.Color == Color_Black {
			startPawnRow = 1
		}

		maxDistance := -2
		if piece.Point.I != startPawnRow {
			maxDistance = -1
		}

		promotionRow := 0
		promotingKinds := [4]Kind{Kind_Queen, Kind_Rook, Kind_Bishop, Kind_Knight}
		if piece.Color == Color_Black {
			promotionRow = 7
		}

		// Straight line
		for i := -1; i >= maxDistance; i-- {
			finalI := int(piece.Point.I) + i*invertMult
			if finalI >= 0 && finalI < 8 {
				row := uint8(finalI)
				pieceAt := p.Board.GetPieceAt(row, piece.Point.J)
				if pieceAt.Kind != Kind_None {
					break
				}

				if finalI == promotionRow {
					for _, kind := range promotingKinds {
						*movements = append(*movements,
							*NewMovement(piece,
								piece.Point,
								NewPoint(row, piece.Point.J),
							).WithPawn(i == -2, false).WithPawnPromotion(kind))
					}
				} else {
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(row, piece.Point.J),
						).WithPawn(i == -2, false))
				}
			}
		}

		dirs := [2][2]int{{-1, -1}, {-1, 1}}

		for _, dir := range dirs {
			finalI, finalJ := int(piece.Point.I)+dir[0]*invertMult, int(piece.Point.J)+dir[1]*invertMult
			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				row := uint8(finalI)
				col := uint8(finalJ)

				pieceAt := p.Board.GetPieceAt(row, col)
				// TODO: Pawn attacking castling could be done in a different way:
				// PseudoMovements returns if attackingQueen, and if attackingKing (if any piece is) so is faster

				if pieceAt.Kind != Kind_None && pieceAt.Color != piece.Color {
					if finalI == promotionRow {
						for _, kind := range promotingKinds {
							attackMatrix[row][col] = true
							*movements = append(*movements,
								*NewMovement(piece,
									piece.Point,
									NewPoint(row, col),
								).WithTakingPiece(pieceAt).WithPawn(false, false).WithPawnPromotion(kind))
						}
					} else {
						attackMatrix[row][col] = true
						*movements = append(*movements,
							*NewMovement(piece,
								piece.Point,
								NewPoint(row, col),
							).WithTakingPiece(pieceAt).WithPawn(false, false))
					}
				} else {
					// If there is no piece in the diagonal, we cannot move to it, but mark the position as being attacked
					attackMatrix[row][col] = true
				}

				if pieceAt.Kind == Kind_None && p.Status.EnPassant != nil && p.Status.EnPassant.I == row && p.Status.EnPassant.J == col {
					enPassantPiecePoint := NewPoint(p.Status.EnPassant.I+1, p.Status.EnPassant.J)
					if piece.Color == Color_Black {
						enPassantPiecePoint.I = p.Status.EnPassant.I - 1
					}

					pieceAt := p.Board.GetPieceAt(enPassantPiecePoint.I, enPassantPiecePoint.J)

					attackMatrix[row][col] = true
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(row, col),
						).WithTakingPiece(pieceAt).WithPawn(false, false))
				}
			}
		}

		return
	}

	return
}

func (p Position) getOrthogonalPseudoMovements(piece Piece, movements *[]Movement, attackMatrix *[8][8]bool) {
	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // {i, j} -> Top, Right, Bottom, Left

	for _, dir := range dirs {
		for i, j := int(piece.Point.I)+dir[0], int(piece.Point.J)+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			row := uint8(i)
			col := uint8(j)

			pieceAt := p.Board.GetPieceAt(row, col)

			if pieceAt.Kind == Kind_None {
				attackMatrix[row][col] = true
				*movements = append(*movements,
					*NewMovement(piece,
						piece.Point,
						NewPoint(row, col),
					))
			} else {
				if pieceAt.Color != piece.Color {
					attackMatrix[row][col] = true
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(row, col),
						).WithTakingPiece(pieceAt))
				}

				break
			}
		}
	}

	return
}

func (p Position) getDiagonalPseudoMovements(piece Piece, movements *[]Movement, attackMatrix *[8][8]bool) {
	dirs := [4][2]int{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}} // {i, j} -> TopLeft, TopRight, BottomRight, BottomLeft

	for _, dir := range dirs {
		for i, j := int(piece.Point.I)+dir[0], int(piece.Point.J)+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			row := uint8(i)
			col := uint8(j)

			pieceAt := p.Board.GetPieceAt(row, col)

			if pieceAt.Kind == Kind_None {
				attackMatrix[row][col] = true
				*movements = append(*movements,
					*NewMovement(piece,
						piece.Point,
						NewPoint(row, col),
					))
			} else {
				if pieceAt.Color != piece.Color {
					attackMatrix[row][col] = true
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(row, col),
						).WithTakingPiece(pieceAt))
				}

				break
			}
		}
	}

	return
}
