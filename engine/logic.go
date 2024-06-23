package engine

func (g *Game) GetLegalMovements() []Movement {
	pseudoMovements := g.CurrentPosition.GetPseudoMovements(g.CurrentPosition.Status.PlayerToMove)
	legalMovements := g.FilterPseudoMovements(&pseudoMovements)
	//fmt.Printf("Pseudo vs legal: %d vs %d\n", len(pseudoMovements), len(legalMovements))
	return legalMovements
}

func (p Position) GetPseudoMovements(color Color) []Movement {
	movements := make([]Movement, 0, 256)

	for _, row := range p.Board {
		for _, piece := range row {
			if piece.Color == color {
				p.GetPiecePseudoMovements(piece, &movements)
				//movements = append(movements, p.Status.GetPiecePseudoMovements(p)...)
			}
		}
	}

	return movements
}

func (p Position) GetPiecePseudoMovements(piece Piece, movements *[]Movement) {
	switch piece.Kind {
	case Kind_Bishop:
		p.getDiagonalPseudoMovements(piece, movements)
		return
	case Kind_Rook:
		p.getOrthogonalPseudoMovements(piece, movements)
		return
	case Kind_Queen:
		p.getDiagonalPseudoMovements(piece, movements)
		p.getOrthogonalPseudoMovements(piece, movements)
		return
	case Kind_King:
		for i := -1; i < 2; i++ {
			for j := -1; j < 2; j++ {
				if i == 0 && j == 0 {
					continue
				}

				finalI, finalJ := piece.Point.I+i, piece.Point.J+j

				if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
					pieceAt := p.Board.GetPieceAt(finalI, finalJ)

					if pieceAt.Kind == Kind_None {
						*movements = append(*movements,
							*NewMovement(piece,
								piece.Point,
								NewPoint(finalI, finalJ),
								// p.Status.EnPassant,
								// p.Status.CanQueenCastling[Color_White],
								// p.Status.CanKingCastling[Color_White],
								// p.Status.CanQueenCastling[Color_Black],
								// p.Status.CanKingCastling[Color_Black],
							))
					} else if pieceAt.Color != piece.Color {
						*movements = append(*movements,
							*NewMovement(piece,
								piece.Point,
								NewPoint(finalI, finalJ),
								// p.Status.EnPassant,
								// p.Status.CanQueenCastling[Color_White],
								// p.Status.CanKingCastling[Color_White],
								// p.Status.CanQueenCastling[Color_Black],
								// p.Status.CanKingCastling[Color_Black],
							).WithTakingPiece(pieceAt))
					}
				}
			}
		}

		if p.Status.CastlingRights.CanQueenCastling(piece.Color) {
			// Check if space to rook is empty
			canCastle := true
			for j := piece.Point.J - 1; j >= piece.Point.J-3; j-- {
				if p.Board.GetPieceAt(piece.Point.I, j).Kind != Kind_None {
					canCastle = false
					break
				}
			}

			if canCastle {
				*movements = append(*movements,
					*NewMovement(piece,
						piece.Point,
						NewPoint(piece.Point.I, piece.Point.J-2),
						// p.Status.EnPassant,
						// p.Status.CanQueenCastling[Color_White],
						// p.Status.CanKingCastling[Color_White],
						// p.Status.CanQueenCastling[Color_Black],
						// p.Status.CanKingCastling[Color_Black],
					).WithCastling(true, false))
			}

		}

		if p.Status.CastlingRights.CanKingCastling(piece.Color) {
			// Check if space to rook is empty
			canCastle := true
			for j := piece.Point.J + 1; j < 7; j++ {
				if p.Board.GetPieceAt(piece.Point.I, j).Kind != Kind_None {
					canCastle = false
					break
				}
			}

			if canCastle {
				*movements = append(*movements,
					*NewMovement(piece,
						piece.Point,
						NewPoint(piece.Point.I, piece.Point.J+2),
						// p.Status.EnPassant,
						// p.Status.CanQueenCastling[Color_White],
						// p.Status.CanKingCastling[Color_White],
						// p.Status.CanQueenCastling[Color_Black],
						// p.Status.CanKingCastling[Color_Black],
					).WithCastling(false, true))
			}

		}

		return
	case Kind_Knight:
		dirs := [8][2]int{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}} // {i, j} -> From TopRight, clockwise untill TopRight (bottom left)

		for _, dir := range dirs {
			finalI, finalJ := piece.Point.I+dir[0], piece.Point.J+dir[1]

			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := p.Board.GetPieceAt(finalI, finalJ)

				if pieceAt.Kind == Kind_None {
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(finalI, finalJ),
							// p.Status.EnPassant,
							// p.Status.CanQueenCastling[Color_White],
							// p.Status.CanKingCastling[Color_White],
							// p.Status.CanQueenCastling[Color_Black],
							// p.Status.CanKingCastling[Color_Black],
						))
				} else if pieceAt.Color != piece.Color {
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(finalI, finalJ),
							// p.Status.EnPassant,
							// p.Status.CanQueenCastling[Color_White],
							// p.Status.CanKingCastling[Color_White],
							// p.Status.CanQueenCastling[Color_Black],
							// p.Status.CanKingCastling[Color_Black],
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

		maxDistance := -2

		if !piece.IsPawnFirstMovement {
			maxDistance = -1
		}

		promotionRow := 0
		promotingKinds := [4]Kind{Kind_Queen, Kind_Rook, Kind_Bishop, Kind_Knight}
		if piece.Color == Color_Black {
			promotionRow = 7
		}

		// Straight line
		for i := -1; i >= maxDistance; i-- {
			finalI := piece.Point.I + i*invertMult
			if finalI >= 0 && finalI < 8 {
				pieceAt := p.Board.GetPieceAt(finalI, piece.Point.J)
				if pieceAt.Kind != Kind_None {
					break
				}

				if finalI == promotionRow {
					for _, kind := range promotingKinds {
						*movements = append(*movements,
							*NewMovement(piece,
								piece.Point,
								NewPoint(finalI, piece.Point.J),
								// p.Status.EnPassant,
								// p.Status.CanQueenCastling[Color_White],
								// p.Status.CanKingCastling[Color_White],
								// p.Status.CanQueenCastling[Color_Black],
								// p.Status.CanKingCastling[Color_Black],
							).WithPawn(i == -2, false).WithPawnPromotion(kind))
					}
				} else {
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(finalI, piece.Point.J),
							// p.Status.EnPassant,
							// p.Status.CanQueenCastling[Color_White],
							// p.Status.CanKingCastling[Color_White],
							// p.Status.CanQueenCastling[Color_Black],
							// p.Status.CanKingCastling[Color_Black],
						).WithPawn(i == -2, false))
				}
			}
		}

		dirs := [2][2]int{{-1, -1}, {-1, 1}}

		for _, dir := range dirs {
			finalI, finalJ := piece.Point.I+dir[0]*invertMult, piece.Point.J+dir[1]*invertMult
			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := p.Board.GetPieceAt(finalI, finalJ)
				// TODO: Pawn attacking castling could be done in a different way:
				// PseudoMovements returns if attackingQueen, and if attackingKing (if any piece is) so is faster

				if pieceAt.Kind != Kind_None && pieceAt.Color != piece.Color {
					if finalI == promotionRow {
						for _, kind := range promotingKinds {
							*movements = append(*movements,
								*NewMovement(piece,
									piece.Point,
									NewPoint(finalI, finalJ),
									// p.Status.EnPassant,
									// p.Status.CanQueenCastling[Color_White],
									// p.Status.CanKingCastling[Color_White],
									// p.Status.CanQueenCastling[Color_Black],
									// p.Status.CanKingCastling[Color_Black],
								).WithTakingPiece(pieceAt).WithPawn(false, false).WithPawnPromotion(kind))
						}
					} else {
						*movements = append(*movements,
							*NewMovement(piece,
								piece.Point,
								NewPoint(finalI, finalJ),
								// p.Status.EnPassant,
								// p.Status.CanQueenCastling[Color_White],
								// p.Status.CanKingCastling[Color_White],
								// p.Status.CanQueenCastling[Color_Black],
								// p.Status.CanKingCastling[Color_Black],
							).WithTakingPiece(pieceAt).WithPawn(false, false))
					}
				} else {
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(finalI, finalJ),
							// p.Status.EnPassant,
							// p.Status.CanQueenCastling[Color_White],
							// p.Status.CanKingCastling[Color_White],
							// p.Status.CanQueenCastling[Color_Black],
							// p.Status.CanKingCastling[Color_Black],
						).WithPawn(false, true))
				}

				if pieceAt.Kind == Kind_None && p.Status.EnPassant != nil && p.Status.EnPassant.I == finalI && p.Status.EnPassant.J == finalJ {
					enPassantPiecePoint := NewPoint(p.Status.EnPassant.I+1, p.Status.EnPassant.J)
					if piece.Color == Color_Black {
						enPassantPiecePoint.I = p.Status.EnPassant.I - 1
					}

					pieceAt := p.Board.GetPieceAt(enPassantPiecePoint.I, enPassantPiecePoint.J)

					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(finalI, finalJ),
							// p.Status.EnPassant,
							// p.Status.CanQueenCastling[Color_White],
							// p.Status.CanKingCastling[Color_White],
							// p.Status.CanQueenCastling[Color_Black],
							// p.Status.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt).WithPawn(false, false))
				}
			}
		}

		return
	}

	return
}

func (p Position) getOrthogonalPseudoMovements(piece Piece, movements *[]Movement) {
	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // {i, j} -> Top, Right, Bottom, Left

	for _, dir := range dirs {
		for i, j := piece.Point.I+dir[0], piece.Point.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			pieceAt := p.Board.GetPieceAt(i, j)

			if pieceAt.Kind == Kind_None {
				*movements = append(*movements,
					*NewMovement(piece,
						piece.Point,
						NewPoint(i, j),
						// p.Status.EnPassant,
						// p.Status.CanQueenCastling[Color_White],
						// p.Status.CanKingCastling[Color_White],
						// p.Status.CanQueenCastling[Color_Black],
						// p.Status.CanKingCastling[Color_Black],
					))
			} else {
				if pieceAt.Color != piece.Color {
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(i, j),
							// p.Status.EnPassant,
							// p.Status.CanQueenCastling[Color_White],
							// p.Status.CanKingCastling[Color_White],
							// p.Status.CanQueenCastling[Color_Black],
							// p.Status.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt))
				}

				break
			}
		}
	}

	return
}

func (p Position) getDiagonalPseudoMovements(piece Piece, movements *[]Movement) {
	dirs := [4][2]int{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}} // {i, j} -> TopLeft, TopRight, BottomRight, BottomLeft

	for _, dir := range dirs {
		for i, j := piece.Point.I+dir[0], piece.Point.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			pieceAt := p.Board.GetPieceAt(i, j)

			if pieceAt.Kind == Kind_None {
				*movements = append(*movements,
					*NewMovement(piece,
						piece.Point,
						NewPoint(i, j),
						// p.Status.EnPassant,
						// p.Status.CanQueenCastling[Color_White],
						// p.Status.CanKingCastling[Color_White],
						// p.Status.CanQueenCastling[Color_Black],
						// p.Status.CanKingCastling[Color_Black],
					))
			} else {
				if pieceAt.Color != piece.Color {
					*movements = append(*movements,
						*NewMovement(piece,
							piece.Point,
							NewPoint(i, j),
							// p.Status.EnPassant,
							// p.Status.CanQueenCastling[Color_White],
							// p.Status.CanKingCastling[Color_White],
							// p.Status.CanQueenCastling[Color_Black],
							// p.Status.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt))
				}

				break
			}
		}
	}

	return
}
