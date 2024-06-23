package engine

func (b *Board) GetLegalMovements(color Color) []Movement {
	pseudoMovements := b.GetPseudoMovements(color)
	legalMovements := b.FilterPseudoMovements(&pseudoMovements)
	//fmt.Printf("Pseudo vs legal: %d vs %d\n", len(pseudoMovements), len(legalMovements))
	return legalMovements
}

func (b Board) GetPseudoMovements(color Color) []Movement {
	movements := make([]Movement, 0, 256)

	for _, row := range b.Data {
		for _, p := range row {
			if p.Color == color {
				b.GetPiecePseudoMovements(p, &movements)
				//movements = append(movements, b.GetPiecePseudoMovements(p)...)
			}
		}
	}

	return movements
}

func (b Board) GetPiecePseudoMovements(p Piece, movements *[]Movement) {
	switch p.Kind {
	case Kind_Bishop:
		b.getDiagonalPseudoMovements(p, movements)
		return
	case Kind_Rook:
		b.getOrthogonalPseudoMovements(p, movements)
		return
	case Kind_Queen:
		b.getDiagonalPseudoMovements(p, movements)
		b.getOrthogonalPseudoMovements(p, movements)
		return
	case Kind_King:
		for i := -1; i < 2; i++ {
			for j := -1; j < 2; j++ {
				if i == 0 && j == 0 {
					continue
				}

				finalI, finalJ := p.Point.I+i, p.Point.J+j

				if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
					pieceAt := b.GetPieceAt(finalI, finalJ)

					if pieceAt.Kind == Kind_None {
						*movements = append(*movements,
							*NewMovement(p,
								p.Point,
								NewPoint(finalI, finalJ),
								b.EnPassant,
								b.CanQueenCastling[Color_White],
								b.CanKingCastling[Color_White],
								b.CanQueenCastling[Color_Black],
								b.CanKingCastling[Color_Black],
							))
					} else if pieceAt.Color != p.Color {
						*movements = append(*movements,
							*NewMovement(p,
								p.Point,
								NewPoint(finalI, finalJ),
								b.EnPassant,
								b.CanQueenCastling[Color_White],
								b.CanKingCastling[Color_White],
								b.CanQueenCastling[Color_Black],
								b.CanKingCastling[Color_Black],
							).WithTakingPiece(pieceAt))
					}
				}
			}
		}

		if b.CanQueenCastling[p.Color] {
			// Check if space to rook is empty
			canCastle := true
			for j := p.Point.J - 1; j >= p.Point.J-3; j-- {
				if b.GetPieceAt(p.Point.I, j).Kind != Kind_None {
					canCastle = false
					break
				}
			}

			if canCastle {
				*movements = append(*movements,
					*NewMovement(p,
						p.Point,
						NewPoint(p.Point.I, p.Point.J-2),
						b.EnPassant,
						b.CanQueenCastling[Color_White],
						b.CanKingCastling[Color_White],
						b.CanQueenCastling[Color_Black],
						b.CanKingCastling[Color_Black],
					).WithCastling(true, false))
			}

		}

		if b.CanKingCastling[p.Color] {
			// Check if space to rook is empty
			canCastle := true
			for j := p.Point.J + 1; j < 7; j++ {
				if b.GetPieceAt(p.Point.I, j).Kind != Kind_None {
					canCastle = false
					break
				}
			}

			if canCastle {
				*movements = append(*movements,
					*NewMovement(p,
						p.Point,
						NewPoint(p.Point.I, p.Point.J+2),
						b.EnPassant,
						b.CanQueenCastling[Color_White],
						b.CanKingCastling[Color_White],
						b.CanQueenCastling[Color_Black],
						b.CanKingCastling[Color_Black],
					).WithCastling(false, true))
			}

		}

		return
	case Kind_Knight:
		dirs := [8][2]int{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}} // {i, j} -> From TopRight, clockwise untill TopRight (bottom left)

		for _, dir := range dirs {
			finalI, finalJ := p.Point.I+dir[0], p.Point.J+dir[1]

			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := b.GetPieceAt(finalI, finalJ)

				if pieceAt.Kind == Kind_None {
					*movements = append(*movements,
						*NewMovement(p,
							p.Point,
							NewPoint(finalI, finalJ),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						))
				} else if pieceAt.Color != p.Color {
					*movements = append(*movements,
						*NewMovement(p,
							p.Point,
							NewPoint(finalI, finalJ),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt))
				}
			}
		}

		return

	case Kind_Pawn:
		invertMult := 1
		if p.Color == Color_Black {
			invertMult = -1
		}

		maxDistance := -2

		if !p.IsPawnFirstMovement {
			maxDistance = -1
		}

		promotionRow := 0
		promotingKinds := [4]Kind{Kind_Queen, Kind_Rook, Kind_Bishop, Kind_Knight}
		if p.Color == Color_Black {
			promotionRow = 7
		}

		// Straight line
		for i := -1; i >= maxDistance; i-- {
			finalI := p.Point.I + i*invertMult
			if finalI >= 0 && finalI < 8 {
				pieceAt := b.GetPieceAt(finalI, p.Point.J)
				if pieceAt.Kind != Kind_None {
					break
				}

				if finalI == promotionRow {
					for _, kind := range promotingKinds {
						*movements = append(*movements,
							*NewMovement(p,
								p.Point,
								NewPoint(finalI, p.Point.J),
								b.EnPassant,
								b.CanQueenCastling[Color_White],
								b.CanKingCastling[Color_White],
								b.CanQueenCastling[Color_Black],
								b.CanKingCastling[Color_Black],
							).WithPawn(i == -2, false).WithPawnPromotion(kind))
					}
				} else {
					*movements = append(*movements,
						*NewMovement(p,
							p.Point,
							NewPoint(finalI, p.Point.J),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithPawn(i == -2, false))
				}
			}
		}

		dirs := [2][2]int{{-1, -1}, {-1, 1}}

		for _, dir := range dirs {
			finalI, finalJ := p.Point.I+dir[0]*invertMult, p.Point.J+dir[1]*invertMult
			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := b.GetPieceAt(finalI, finalJ)
				// TODO: Pawn attacking castling could be done in a different way:
				// PseudoMovements returns if attackingQueen, and if attackingKing (if any piece is) so is faster

				if pieceAt.Kind != Kind_None && pieceAt.Color != p.Color {
					if finalI == promotionRow {
						for _, kind := range promotingKinds {
							*movements = append(*movements,
								*NewMovement(p,
									p.Point,
									NewPoint(finalI, finalJ),
									b.EnPassant,
									b.CanQueenCastling[Color_White],
									b.CanKingCastling[Color_White],
									b.CanQueenCastling[Color_Black],
									b.CanKingCastling[Color_Black],
								).WithTakingPiece(pieceAt).WithPawn(false, false).WithPawnPromotion(kind))
						}
					} else {
						*movements = append(*movements,
							*NewMovement(p,
								p.Point,
								NewPoint(finalI, finalJ),
								b.EnPassant,
								b.CanQueenCastling[Color_White],
								b.CanKingCastling[Color_White],
								b.CanQueenCastling[Color_Black],
								b.CanKingCastling[Color_Black],
							).WithTakingPiece(pieceAt).WithPawn(false, false))
					}
				} else {
					*movements = append(*movements,
						*NewMovement(p,
							p.Point,
							NewPoint(finalI, finalJ),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithPawn(false, true))
				}

				if pieceAt.Kind == Kind_None && b.EnPassant != nil && b.EnPassant.I == finalI && b.EnPassant.J == finalJ {
					enPassantPiecePoint := NewPoint(b.EnPassant.I+1, b.EnPassant.J)
					if p.Color == Color_Black {
						enPassantPiecePoint.I = b.EnPassant.I - 1
					}

					pieceAt := b.GetPieceAt(enPassantPiecePoint.I, enPassantPiecePoint.J)

					*movements = append(*movements,
						*NewMovement(p,
							p.Point,
							NewPoint(finalI, finalJ),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt).WithPawn(false, false))
				}
			}
		}

		return
	}

	return
}

func (b Board) getOrthogonalPseudoMovements(p Piece, movements *[]Movement) {
	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // {i, j} -> Top, Right, Bottom, Left

	for _, dir := range dirs {
		for i, j := p.Point.I+dir[0], p.Point.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			pieceAt := b.GetPieceAt(i, j)

			if pieceAt.Kind == Kind_None {
				*movements = append(*movements,
					*NewMovement(p,
						p.Point,
						NewPoint(i, j),
						b.EnPassant,
						b.CanQueenCastling[Color_White],
						b.CanKingCastling[Color_White],
						b.CanQueenCastling[Color_Black],
						b.CanKingCastling[Color_Black],
					))
			} else {
				if pieceAt.Color != p.Color {
					*movements = append(*movements,
						*NewMovement(p,
							p.Point,
							NewPoint(i, j),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt))
				}

				break
			}
		}
	}

	return
}

func (b Board) getDiagonalPseudoMovements(p Piece, movements *[]Movement) {
	dirs := [4][2]int{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}} // {i, j} -> TopLeft, TopRight, BottomRight, BottomLeft

	for _, dir := range dirs {
		for i, j := p.Point.I+dir[0], p.Point.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			pieceAt := b.GetPieceAt(i, j)

			if pieceAt.Kind == Kind_None {
				*movements = append(*movements,
					*NewMovement(p,
						p.Point,
						NewPoint(i, j),
						b.EnPassant,
						b.CanQueenCastling[Color_White],
						b.CanKingCastling[Color_White],
						b.CanQueenCastling[Color_Black],
						b.CanKingCastling[Color_Black],
					))
			} else {
				if pieceAt.Color != p.Color {
					*movements = append(*movements,
						*NewMovement(p,
							p.Point,
							NewPoint(i, j),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt))
				}

				break
			}
		}
	}

	return
}
