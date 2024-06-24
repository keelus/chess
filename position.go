package engine

import (
	"strconv"
	"strings"
)

type Position struct {
	Board    Board
	Status   PositionStatus
	Captures []Piece // Only used via API. Perft ignores this
}

type CastlingRights struct {
	QueenSide map[Color]bool
	KingSide  map[Color]bool
}

func (cr *CastlingRights) clone() CastlingRights {
	return CastlingRights{
		QueenSide: map[Color]bool{
			Color_White: cr.QueenSide[Color_White],
			Color_Black: cr.QueenSide[Color_Black],
		},
		KingSide: map[Color]bool{
			Color_White: cr.KingSide[Color_White],
			Color_Black: cr.KingSide[Color_Black],
		},
	}
}

type PositionStatus struct {
	PlayerToMove    Color
	CastlingRights  CastlingRights
	EnPassant       *Square
	HalfmoveClock   uint8
	FullmoveCounter uint
}

func (p Position) GetHalfmoveClock() uint8 {
	return p.Status.HalfmoveClock
}
func (p Position) GetFullmoveCounter() uint {
	return p.Status.FullmoveCounter
}

func (p Position) GetCaptures() []Piece {
	return p.Captures
}

func (ps *PositionStatus) clone() PositionStatus {
	return PositionStatus{
		PlayerToMove:    ps.PlayerToMove,
		CastlingRights:  ps.CastlingRights.clone(),
		EnPassant:       ps.EnPassant,
		HalfmoveClock:   ps.HalfmoveClock,
		FullmoveCounter: ps.FullmoveCounter,
	}
}

func newPositionFromFen(fen string) Position {
	parsedFen, err := parseFen(fen)
	if err != nil {
		panic(err)
	}

	return Position{
		Board: newBoardFromFen(parsedFen.PlacementData),
		Status: PositionStatus{
			PlayerToMove: parsedFen.ActiveColor,

			CastlingRights: CastlingRights{
				QueenSide: map[Color]bool{
					Color_White: parsedFen.WhiteCanQueenSideCastling,
					Color_Black: parsedFen.BlackCanQueenSideCastling,
				},
				KingSide: map[Color]bool{
					Color_White: parsedFen.WhiteCanKingSideCastling,
					Color_Black: parsedFen.BlackCanKingSideCastling,
				},
			},

			EnPassant:       parsedFen.EnPassant,
			HalfmoveClock:   parsedFen.HalfmoveClock,
			FullmoveCounter: parsedFen.FulmoveCounter,
		},
		Captures: make([]Piece, 0),
	}
}

// TODO: Complete
func (p Position) Fen() string {
	var sb strings.Builder

	sb.WriteRune(' ')
	sb.WriteString(p.Board.Fen())

	sb.WriteRune(p.Status.PlayerToMove.ToRune())

	if p.Status.CastlingRights.QueenSide[Color_White] && p.Status.CastlingRights.KingSide[Color_White] && p.Status.CastlingRights.QueenSide[Color_Black] && p.Status.CastlingRights.KingSide[Color_Black] {
		sb.WriteRune(' ')

		if p.Status.CastlingRights.KingSide[Color_White] {
			sb.WriteRune('K')
		}
		if p.Status.CastlingRights.QueenSide[Color_White] {
			sb.WriteRune('Q')
		}
		if p.Status.CastlingRights.KingSide[Color_Black] {
			sb.WriteRune('k')
		}
		if p.Status.CastlingRights.QueenSide[Color_Black] {
			sb.WriteRune('q')
		}

		sb.WriteRune(' ')
	} else {
		sb.WriteString(" - ")
	}

	if p.Status.EnPassant != nil {
		sb.WriteString(p.Status.EnPassant.ToAlgebraic())
	} else {
		sb.WriteRune('-')
	}

	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(int(p.Status.HalfmoveClock)))
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(int(p.Status.FullmoveCounter)))

	return sb.String()
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
		p.getDirectionPseudoMovements(piece, movements, attackMatrix, bishopDirections)
		return
	case Kind_Rook:
		p.getDirectionPseudoMovements(piece, movements, attackMatrix, rookDirections)
		return
	case Kind_Queen:
		p.getDirectionPseudoMovements(piece, movements, attackMatrix, bishopDirections)
		p.getDirectionPseudoMovements(piece, movements, attackMatrix, rookDirections)
		return
	case Kind_King:
		for _, offset := range kingOffsets {
			targetRow, targetCol := int8(piece.Square.I)+offset[0], int8(piece.Square.J)+offset[1]
			if targetRow >= 0 && targetCol >= 0 && targetRow < 8 && targetCol < 8 {
				row := uint8(targetRow)
				col := uint8(targetCol)
				pieceAt := p.Board[row][col]

				if pieceAt.Kind == Kind_None {
					attackMatrix[row][col] = true
					*movements = append(*movements,
						*newMovement(piece,
							piece.Square,
							newSquare(row, col),
						))
				} else if pieceAt.Color != piece.Color {
					attackMatrix[row][col] = true
					*movements = append(*movements,
						*newMovement(piece,
							piece.Square,
							newSquare(row, col),
						).withTakingPiece(pieceAt))
				}
			}

		}

		if doCastlingCheck {
			castlingRow := 7
			if piece.Color == Color_Black {
				castlingRow = 0
			}

			_, enemyAttackBoard := p.GetPseudoMovements(piece.Color.Opposite(), false)
			// If king is not in check, continue
			if enemyAttackBoard[piece.Square.I][piece.Square.J] == false {
				if p.Status.CastlingRights.QueenSide[piece.Color] {
					// Check if space to rook is empty
					canCastle := true
					for j := piece.Square.J - 1; j >= piece.Square.J-3; j-- {
						if p.Board[piece.Square.I][j].Kind != Kind_None {
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
								*newMovement(piece,
									piece.Square,
									newSquare(piece.Square.I, piece.Square.J-2),
								).withCastling(true, false))
						}
					}
				}

				if p.Status.CastlingRights.KingSide[piece.Color] {
					// Check if space to rook is empty
					canCastle := true
					for j := piece.Square.J + 1; j < 7; j++ {
						if p.Board[piece.Square.I][j].Kind != Kind_None {
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
								*newMovement(piece,
									piece.Square,
									newSquare(piece.Square.I, piece.Square.J+2),
								).withCastling(false, true))
						}
					}
				}
			}
		}

		return
	case Kind_Knight:
		for _, offset := range knightOffsets {
			targetRow := int8(piece.Square.I) + offset[0]
			targetCol := int8(piece.Square.J) + offset[1]

			if targetRow >= 0 && targetCol >= 0 && targetRow < 8 && targetCol < 8 {
				row := uint8(targetRow)
				col := uint8(targetCol)

				pieceAt := p.Board[row][col]

				if pieceAt.Kind == Kind_None {
					attackMatrix[row][col] = true
					*movements = append(*movements,
						*newMovement(piece,
							piece.Square,
							newSquare(row, col),
						))
				} else if pieceAt.Color != piece.Color {
					attackMatrix[row][col] = true
					*movements = append(*movements,
						*newMovement(piece,
							piece.Square,
							newSquare(row, col),
						).withTakingPiece(pieceAt))
				}
			}
		}

		return

	case Kind_Pawn:
		// Straight moves
		var maxDistance int8 = 2
		if piece.Square.I != pawnStartingRows[piece.Color] {
			maxDistance = 1
		}

		promotionRow := uint8(0)
		if piece.Color == Color_Black {
			promotionRow = 7
		}

		for i := int8(1); i <= maxDistance; i++ {
			targetRow := int8(piece.Square.I) + pawnMoveRowDirections[piece.Color]*i
			if targetRow >= 0 && targetRow < 8 {
				row := uint8(targetRow)
				pieceAt := p.Board[row][piece.Square.J]
				if pieceAt.Kind != Kind_None {
					break
				}

				if row == promotionRow {
					for _, kind := range promotableKinds {
						*movements = append(*movements,
							*newMovement(piece,
								piece.Square,
								newSquare(row, piece.Square.J),
							).withPawn(i == 2).withPawnPromotion(kind))
					}
				} else {
					*movements = append(*movements,
						*newMovement(piece,
							piece.Square,
							newSquare(row, piece.Square.J),
						).withPawn(i == 2))
				}
			}
		}

		// Diagonal move
		for _, offset := range pawnAttackOffsets[piece.Color] {
			targetRow := int8(piece.Square.I) + offset[0]
			targetCol := int8(piece.Square.J) + offset[1]

			if targetRow >= 0 && targetCol >= 0 && targetRow < 8 && targetCol < 8 {
				row := uint8(targetRow)
				col := uint8(targetCol)

				pieceAt := p.Board[row][col]

				if pieceAt.Kind != Kind_None && pieceAt.Color != piece.Color {
					if row == pawnPromotionRows[piece.Color] {
						for _, kind := range promotableKinds {
							attackMatrix[row][col] = true
							*movements = append(*movements,
								*newMovement(piece,
									piece.Square,
									newSquare(row, col),
								).withTakingPiece(pieceAt).withPawn(false).withPawnPromotion(kind))
						}
					} else {
						attackMatrix[row][col] = true
						*movements = append(*movements,
							*newMovement(piece,
								piece.Square,
								newSquare(row, col),
							).withTakingPiece(pieceAt).withPawn(false))
					}
				} else {
					// If there is no piece in the diagonal, we cannot move to it, but mark the position as being attacked
					attackMatrix[row][col] = true
				}

				// En passant available check in current diagonal
				if pieceAt.Kind == Kind_None && p.Status.EnPassant != nil && p.Status.EnPassant.I == row && p.Status.EnPassant.J == col {
					enPassantPieceSquare := newSquare(p.Status.EnPassant.I+1, p.Status.EnPassant.J)
					if piece.Color == Color_Black {
						enPassantPieceSquare.I = p.Status.EnPassant.I - 1
					}

					pieceAt := p.Board[enPassantPieceSquare.I][enPassantPieceSquare.J]

					attackMatrix[row][col] = true
					*movements = append(*movements,
						*newMovement(piece,
							piece.Square,
							newSquare(row, col),
						).withTakingPiece(pieceAt).withPawn(false))
				}
			}
		}

		return
	}
}

func (p Position) getDirectionPseudoMovements(piece Piece, movements *[]Movement, attackMatrix *[8][8]bool, directions [4][2]int8) {
	for _, dir := range directions {
		i := int8(piece.Square.I) + dir[0]
		j := int8(piece.Square.J) + dir[1]

		for i >= 0 && j >= 0 && i < 8 && j < 8 {
			row := uint8(i)
			col := uint8(j)

			pieceAt := p.Board[row][col]

			if pieceAt.Kind == Kind_None {
				attackMatrix[row][col] = true
				*movements = append(*movements,
					*newMovement(piece,
						piece.Square,
						newSquare(row, col),
					))
			} else {
				if pieceAt.Color != piece.Color {
					attackMatrix[row][col] = true
					*movements = append(*movements,
						*newMovement(piece,
							piece.Square,
							newSquare(row, col),
						).withTakingPiece(pieceAt))
				}

				break
			}

			i += dir[0]
			j += dir[1]
		}
	}

	return
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
