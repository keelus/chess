package engine

func (b Board) GetLegalMovements() []Movement {
	//pseudoMovements := b.GetPseudoMovements()
	panic("TODO")
	return []Movement{}
}

func (b Board) GetPseudoMovements() []Movement {
	movements := []Movement{}

	for _, row := range b.Data {
		for _, p := range row {
			if p != nil {
				movements = append(movements, b.GetPiecePseudoMovements(p)...)
			}
		}
	}

	return movements
}

func (b Board) GetPiecePseudoMovements(p *Piece) []Movement {
	switch p.Kind {
	case Kind_Bishop:
		return b.getDiagonalPseudoMovements(p)
	case Kind_Rook:
		return b.getOrthogonalPseudoMovements(p)
	case Kind_Queen:
		movements := b.getDiagonalPseudoMovements(p)
		return append(movements, b.getOrthogonalPseudoMovements(p)...)
	case Kind_King:
		movements := []Movement{}
		for i := -1; i < 2; i++ {
			for j := -1; j < 2; j++ {
				if i == 0 && j == 0 {
					continue
				}

				finalI, finalJ := p.Position.I+i, p.Position.J+j

				if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
					pieceAt := b.GetPieceAt(finalI, finalJ)

					if pieceAt != nil && pieceAt.Color == b.PlayerToMove {
						continue
					}

					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(finalI, finalJ),
							b.EnPassant,
							b.CanKingCastling[p.Color],
							b.CanQueenCastling[p.Color],
						).WithTakingPiece(pieceAt))
				}
			}
		}

		// TODO: Black

		handleCastling := func(jDelta, jStop int) {
			// Check if space to rook is empty
			canCastle := true
			for j := p.Position.J + jDelta; j != jStop; j += jDelta {
				if b.GetPieceAt(p.Position.I, j) != nil {
					canCastle = false
					break
				}
			}

			if canCastle {
				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(p.Position.I, p.Position.J+jDelta*2),
						b.EnPassant,
						b.CanKingCastling[p.Color],
						b.CanQueenCastling[p.Color],
					).WithCastling(jDelta == -1, jDelta == 1))
			}
		}

		if b.CanQueenCastling[p.Color] {
			handleCastling(-1, 1)
		}

		if b.CanKingCastling[p.Color] {
			handleCastling(1, 7)
		}

		return movements
	case Kind_Knight:
		movements := []Movement{}
		dirs := [8][2]int{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}} // {i, j} -> From TopRight, clockwise untill TopRight (bottom left)

		for _, dir := range dirs {
			finalI, finalJ := p.Position.I+dir[0], p.Position.J+dir[1]

			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := b.GetPieceAt(finalI, finalJ)

				if pieceAt != nil && pieceAt.Color == b.PlayerToMove {
					continue
				}

				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(finalI, finalJ),
						b.EnPassant,
						b.CanKingCastling[p.Color],
						b.CanQueenCastling[p.Color],
					).WithTakingPiece(pieceAt))
			}
		}

		return movements

	case Kind_Pawn:
		movements := []Movement{}

		invertMult := 1
		if p.Color == Color_Black {
			invertMult = -1
		}

		maxDistance := -2

		if !p.IsPawnFirstMovement {
			maxDistance = -1
		}

		// Straight line
		for i := -1; i >= maxDistance; i-- {
			finalI := p.Position.I + i*invertMult
			if finalI >= 0 && finalI < 8 {
				pieceAt := b.GetPieceAt(finalI, p.Position.J)
				if pieceAt != nil {
					break
				}

				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(finalI, p.Position.J),
						b.EnPassant,
						b.CanKingCastling[p.Color],
						b.CanQueenCastling[p.Color],
					).WithTakingPiece(pieceAt).WithPawn(i == -2))
			}
		}

		dirs := [2][2]int{{-1, -1}, {-1, 1}}

		for _, dir := range dirs {
			finalI, finalJ := p.Position.I+dir[0]*invertMult, p.Position.J+dir[1]*invertMult
			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := b.GetPieceAt(finalI, finalJ)

				if pieceAt != nil && pieceAt.Color != p.Color {
					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(finalI, finalJ),
							b.EnPassant,
							b.CanKingCastling[p.Color],
							b.CanQueenCastling[p.Color],
						).WithTakingPiece(pieceAt).WithPawn(false))
				} else if pieceAt == nil && b.EnPassant != nil && b.EnPassant.I == finalI && b.EnPassant.J == finalJ {
					enPassantPiecePosition := NewPosition(b.EnPassant.I+1, b.EnPassant.J)
					if p.Color == Color_Black {
						enPassantPiecePosition.I = b.EnPassant.I - 1
					}

					pieceAt := b.GetPieceAt(enPassantPiecePosition.I, enPassantPiecePosition.J)

					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(finalI, finalJ),
							b.EnPassant,
							b.CanKingCastling[p.Color],
							b.CanQueenCastling[p.Color],
						).WithTakingPiece(pieceAt).WithPawn(false))
				}
			}
		}

		return movements
		// Two diagonals
	}

	return []Movement{}
}

func (b Board) getOrthogonalPseudoMovements(p *Piece) []Movement {
	movements := []Movement{}

	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // {i, j} -> Top, Right, Bottom, Left

	for _, dir := range dirs {
		for i, j := p.Position.I+dir[0], p.Position.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			pieceAt := b.GetPieceAt(i, j)

			if pieceAt == nil {
				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(i, j),
						b.EnPassant,
						b.CanKingCastling[p.Color],
						b.CanQueenCastling[p.Color],
					).WithTakingPiece(pieceAt))
			} else {
				if pieceAt.Color != b.PlayerToMove {
					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(i, j),
							b.EnPassant,
							b.CanKingCastling[p.Color],
							b.CanQueenCastling[p.Color],
						).WithTakingPiece(pieceAt))
				}

				break
			}
		}
	}

	return movements
}

func (b Board) getDiagonalPseudoMovements(p *Piece) []Movement {
	movements := []Movement{}

	dirs := [4][2]int{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}} // {i, j} -> TopLeft, TopRight, BottomRight, BottomLeft

	for _, dir := range dirs {
		for i, j := p.Position.I+dir[0], p.Position.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			pieceAt := b.GetPieceAt(i, j)

			if pieceAt == nil {
				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(i, j),
						b.EnPassant,
						b.CanKingCastling[p.Color],
						b.CanQueenCastling[p.Color],
					).WithTakingPiece(pieceAt))
			} else {
				if pieceAt.Color != b.PlayerToMove {
					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(i, j),
							b.EnPassant,
							b.CanKingCastling[p.Color],
							b.CanQueenCastling[p.Color],
						).WithTakingPiece(pieceAt))
				}

				break
			}
		}
	}

	return movements
}
