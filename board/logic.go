package board

import (
	"chess/movement"
	"chess/piece"
	"chess/position"
	"fmt"
)

func (b Board) GetPseudoMovements() []movement.Movement {
	movements := []movement.Movement{}

	for _, row := range b.Data {
		for _, p := range row {
			if p != nil {
				movements = append(movements, b.GetPiecePseudoMovements(p)...)
			}
		}
	}

	return movements
}

func (b Board) GetPiecePseudoMovements(p *piece.Piece) []movement.Movement {
	switch p.Kind {
	case piece.Kind_Bishop:
		return b.getDiagonalPseudoMovements(p)
	case piece.Kind_Rook:
		return b.getOrthogonalPseudoMovements(p)
	case piece.Kind_Queen:
		movements := b.getDiagonalPseudoMovements(p)
		return append(movements, b.getOrthogonalPseudoMovements(p)...)
	case piece.Kind_King:
		movements := []movement.Movement{}
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

					movements = append(movements, movement.NewMovement(
						p,
						pieceAt,
						p.Position,
						position.NewPosition(finalI, finalJ),
						false,
					))
				}
			}
		}

		return movements
	case piece.Kind_Knight:
		movements := []movement.Movement{}
		dirs := [8][2]int{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}} // {i, j} -> From TopRight, clockwise untill TopRight (bottom left)

		for _, dir := range dirs {
			finalI, finalJ := p.Position.I+dir[0], p.Position.J+dir[1]

			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := b.GetPieceAt(finalI, finalJ)

				if pieceAt != nil && pieceAt.Color == b.PlayerToMove {
					continue
				}

				movements = append(movements, movement.NewMovement(
					p,
					pieceAt,
					p.Position,
					position.NewPosition(finalI, finalJ),
					false,
				))
			}
		}

		return movements

	case piece.Kind_Pawn:
		movements := []movement.Movement{}

		invertMult := 1
		if p.Color == piece.Color_Black {
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

				movements = append(movements, movement.NewMovement(
					p,
					pieceAt,
					p.Position,
					position.NewPosition(finalI, p.Position.J),
					p.IsPawnFirstMovement,
				))
			}
		}

		dirs := [2][2]int{{-1, -1}, {-1, 1}}

		for _, dir := range dirs {
			finalI, finalJ := p.Position.I+dir[0]*invertMult, p.Position.J+dir[1]*invertMult
			if finalI >= 0 && finalJ >= 0 && finalI < 8 && finalJ < 8 {
				pieceAt := b.GetPieceAt(finalI, finalJ)
				if pieceAt != nil && pieceAt.Color != p.Color {
					movements = append(movements, movement.NewMovement(
						p,
						pieceAt,
						p.Position,
						position.NewPosition(finalI, finalJ),
						false,
					))
				}
			}
		}

		return movements
		// Two diagonals
	}

	return []movement.Movement{}
}

func (b Board) getOrthogonalPseudoMovements(p *piece.Piece) []movement.Movement {
	movements := []movement.Movement{}

	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // {i, j} -> Top, Right, Bottom, Left

	for _, dir := range dirs {
		for i, j := p.Position.I+dir[0], p.Position.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			fmt.Println(i, j)
			pieceAt := b.GetPieceAt(i, j)

			if pieceAt == nil {
				movements = append(movements, movement.NewMovement(
					p,
					nil,
					p.Position,
					position.NewPosition(i, j),
					false,
				))
			} else {
				if pieceAt.Color != b.PlayerToMove {
					movements = append(movements, movement.NewMovement(
						p,
						pieceAt,
						p.Position,
						position.NewPosition(i, j),
						false,
					))
				}

				break
			}
		}
	}

	return movements
}

func (b Board) getDiagonalPseudoMovements(p *piece.Piece) []movement.Movement {
	movements := []movement.Movement{}

	dirs := [4][2]int{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}} // {i, j} -> TopLeft, TopRight, BottomRight, BottomLeft

	for _, dir := range dirs {
		for i, j := p.Position.I+dir[0], p.Position.J+dir[1]; i >= 0 && j >= 0 && i < 8 && j < 8; i, j = i+dir[0], j+dir[1] {
			fmt.Println(i, j)
			pieceAt := b.GetPieceAt(i, j)

			if pieceAt == nil {
				movements = append(movements, movement.NewMovement(
					p,
					nil,
					p.Position,
					position.NewPosition(i, j),
					false,
				))
			} else {
				if pieceAt.Color != b.PlayerToMove {
					movements = append(movements, movement.NewMovement(
						p,
						pieceAt,
						p.Position,
						position.NewPosition(i, j),
						false,
					))
				}

				break
			}
		}
	}

	return movements
}
