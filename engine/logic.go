package engine

func (b *Board) GetLegalMovements(color Color) []Movement {
	pseudoMovements := b.GetPseudoMovements(color)
	legalMovements := b.FilterPseudoMovements(pseudoMovements)
	//fmt.Printf("Pseudo vs legal: %d vs %d\n", len(pseudoMovements), len(legalMovements))
	return legalMovements
}

func (b Board) GetPseudoMovements(color Color) []Movement {
	movements := []Movement{}

	for _, row := range b.Data {
		for _, p := range row {
			if p.Color == color {
				movements = append(movements, b.GetPiecePseudoMovements(p)...)
			}
		}
	}

	return movements
}

func (b Board) GetPiecePseudoMovements(p Piece) []Movement {
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

					if pieceAt.Kind == Kind_None {
						movements = append(movements,
							*NewMovement(p,
								p.Position,
								NewPosition(finalI, finalJ),
								b.EnPassant,
								b.CanQueenCastling[Color_White],
								b.CanKingCastling[Color_White],
								b.CanQueenCastling[Color_Black],
								b.CanKingCastling[Color_Black],
							))
					} else if pieceAt.Color != p.Color {
						movements = append(movements,
							*NewMovement(p,
								p.Position,
								NewPosition(finalI, finalJ),
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
			for j := p.Position.J - 1; j >= p.Position.J-3; j-- {
				if b.GetPieceAt(p.Position.I, j).Kind != Kind_None {
					canCastle = false
					break
				}
			}

			if canCastle {
				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(p.Position.I, p.Position.J-2),
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
			for j := p.Position.J + 1; j < 7; j++ {
				if b.GetPieceAt(p.Position.I, j).Kind != Kind_None {
					canCastle = false
					break
				}
			}

			if canCastle {
				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(p.Position.I, p.Position.J+2),
						b.EnPassant,
						b.CanQueenCastling[Color_White],
						b.CanKingCastling[Color_White],
						b.CanQueenCastling[Color_Black],
						b.CanKingCastling[Color_Black],
					).WithCastling(false, true))
			}

		}

		return movements
	case Kind_Knight:
		movements := []Movement{}
		dirs := [8][2]int{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}} // {i, j} -> From TopRight, clockwise untill TopRight (bottom left)

		for _, dir := range dirs {
			finalI, finalJ := p.Position.I+dir[0], p.Position.J+dir[1]

			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := b.GetPieceAt(finalI, finalJ)

				if pieceAt.Kind == Kind_None {
					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(finalI, finalJ),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						))
				} else if pieceAt.Color != p.Color {
					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(finalI, finalJ),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt))
				}
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
				if pieceAt.Kind != Kind_None {
					break
				}

				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(finalI, p.Position.J),
						b.EnPassant,
						b.CanQueenCastling[Color_White],
						b.CanKingCastling[Color_White],
						b.CanQueenCastling[Color_Black],
						b.CanKingCastling[Color_Black],
					).WithPawn(i == -2, false))
			}
		}

		dirs := [2][2]int{{-1, -1}, {-1, 1}}

		for _, dir := range dirs {
			finalI, finalJ := p.Position.I+dir[0]*invertMult, p.Position.J+dir[1]*invertMult
			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := b.GetPieceAt(finalI, finalJ)
				// TODO: Pawn attacking castling could be done in a different way:
				// PseudoMovements returns if attackingQueen, and if attackingKing (if any piece is) so is faster

				if pieceAt.Kind != Kind_None && pieceAt.Color != p.Color {
					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(finalI, finalJ),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt).WithPawn(false, false))
				} else {
					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(finalI, finalJ),
							b.EnPassant,
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithPawn(false, true))
				}

				if pieceAt.Kind == Kind_None && b.EnPassant != nil && b.EnPassant.I == finalI && b.EnPassant.J == finalJ {
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
							b.CanQueenCastling[Color_White],
							b.CanKingCastling[Color_White],
							b.CanQueenCastling[Color_Black],
							b.CanKingCastling[Color_Black],
						).WithTakingPiece(pieceAt).WithPawn(false, false))
				}
			}
		}

		return movements
	}

	return []Movement{}
}

func (b Board) getOrthogonalPseudoMovements(p Piece) []Movement {
	movements := []Movement{}

	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // {i, j} -> Top, Right, Bottom, Left

	for _, dir := range dirs {
		for i, j := p.Position.I+dir[0], p.Position.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			pieceAt := b.GetPieceAt(i, j)

			if pieceAt.Kind == Kind_None {
				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(i, j),
						b.EnPassant,
						b.CanQueenCastling[Color_White],
						b.CanKingCastling[Color_White],
						b.CanQueenCastling[Color_Black],
						b.CanKingCastling[Color_Black],
					))
			} else {
				if pieceAt.Color != p.Color {
					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(i, j),
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

	return movements
}

func (b Board) getDiagonalPseudoMovements(p Piece) []Movement {
	movements := []Movement{}

	dirs := [4][2]int{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}} // {i, j} -> TopLeft, TopRight, BottomRight, BottomLeft

	for _, dir := range dirs {
		for i, j := p.Position.I+dir[0], p.Position.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			pieceAt := b.GetPieceAt(i, j)

			if pieceAt.Kind == Kind_None {
				movements = append(movements,
					*NewMovement(p,
						p.Position,
						NewPosition(i, j),
						b.EnPassant,
						b.CanQueenCastling[Color_White],
						b.CanKingCastling[Color_White],
						b.CanQueenCastling[Color_Black],
						b.CanKingCastling[Color_Black],
					))
			} else {
				if pieceAt.Color != p.Color {
					movements = append(movements,
						*NewMovement(p,
							p.Position,
							NewPosition(i, j),
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

	return movements
}
