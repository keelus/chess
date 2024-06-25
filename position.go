package chess

import (
	"errors"
	"strconv"
	"strings"
)

// Position represents a specific position of a chess game, including piece
// placement in board, active en passant, player turn, move counter, etc.
type Position struct {
	board Board

	playerToMove    Color
	castlingRights  CastlingRights
	enPassantSq     *Square
	halfmoveClock   uint8
	fullmoveCounter uint

	captures []Piece // Only used via API. Perft ignores this.
}

// Note: enPassantSq is not cloned, as it's not needed.
func (p *Position) clone() Position {
	return Position{
		board: p.board,

		playerToMove:    p.playerToMove,
		castlingRights:  p.castlingRights.clone(),
		enPassantSq:     nil,
		halfmoveClock:   p.halfmoveClock,
		fullmoveCounter: p.fullmoveCounter,

		captures: p.captures,
	}
}

// CastlingRights represents the position's current castling rights,
// of both players.
type CastlingRights struct {
	queenSide map[Color]bool
	kingSide  map[Color]bool
}

func (cr *CastlingRights) clone() CastlingRights {
	return CastlingRights{
		queenSide: map[Color]bool{
			Color_White: cr.queenSide[Color_White],
			Color_Black: cr.queenSide[Color_Black],
		},
		kingSide: map[Color]bool{
			Color_White: cr.kingSide[Color_White],
			Color_Black: cr.kingSide[Color_Black],
		},
	}
}

// CastlingRights returns the position's current castling rights,
// of both players.
func (p Position) CastlingRights() CastlingRights {
	return p.castlingRights
}

// QueenSide returns whether the passed color (player/side) has
// castling rights on queenside or not.
func (cr CastlingRights) QueenSide(color Color) bool {
	return cr.queenSide[color]
}

// KingSide returns whether the passed color (player/side) has
// castling rights on kingside or not.
func (cr CastlingRights) KingSide(color Color) bool {
	return cr.kingSide[color]
}

// HalfmoveClock returns the position's halfmove clock.
func (p Position) HalfmoveClock() uint8 {
	return p.halfmoveClock
}

// FullmoveCounter returns the position's fullmove counter.
func (p Position) FullmoveCounter() uint {
	return p.fullmoveCounter
}

// Captures returns a slice containing all the Pieces captured
// until the position, in order.
func (p Position) Captures() []Piece {
	return p.captures
}

// HasActiveEnPassant reports whether this possition has an active
// en passsant oportunity.
func (p Position) HasActiveEnPassant() bool {
	return p.enPassantSq != nil
}

// EnPassantSquare returns the current en passant square.
//
// If there is no en passant in the current position, it will
// return an empty Square and the error.
func (p Position) EnPassantSquare() (Square, error) {
	if p.enPassantSq == nil {
		return Square{}, errors.New("There is no en passant in the current position.")
	}
	return *p.enPassantSq, nil
}

// Turn returns the position's player/side to move.
func (p Position) Turn() Color {
	return p.playerToMove
}

// Fen returns the position's Forsyth–Edwards Notation, as an string, containing
// the board piece placement, player to move, castling rights, en passant, halfmove clock
// and fullmove counter.
//
// For example, for a starting chess position:
//
//	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
func (p Position) Fen() string {
	var sb strings.Builder

	sb.WriteRune(' ')
	sb.WriteString(p.board.Fen())

	sb.WriteRune(p.playerToMove.ToRune())

	if p.castlingRights.queenSide[Color_White] && p.castlingRights.kingSide[Color_White] && p.castlingRights.queenSide[Color_Black] && p.castlingRights.kingSide[Color_Black] {
		sb.WriteRune(' ')

		if p.castlingRights.kingSide[Color_White] {
			sb.WriteRune('K')
		}
		if p.castlingRights.queenSide[Color_White] {
			sb.WriteRune('Q')
		}
		if p.castlingRights.kingSide[Color_Black] {
			sb.WriteRune('k')
		}
		if p.castlingRights.queenSide[Color_Black] {
			sb.WriteRune('q')
		}

		sb.WriteRune(' ')
	} else {
		sb.WriteString(" - ")
	}

	if p.enPassantSq != nil {
		sb.WriteString(p.enPassantSq.Algebraic())
	} else {
		sb.WriteRune('-')
	}

	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(int(p.halfmoveClock)))
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(int(p.fullmoveCounter)))

	return sb.String()
}

func newPositionFromFen(fen string) (Position, error) {
	parsedFen, err := parseFen(fen)
	if err != nil {
		return Position{}, err
	}

	return Position{
		board: newBoardFromFen(parsedFen.placementData),

		playerToMove: parsedFen.activeColor,

		castlingRights: CastlingRights{
			queenSide: map[Color]bool{
				Color_White: parsedFen.whiteCanQueenSideCastling,
				Color_Black: parsedFen.blackCanQueenSideCastling,
			},
			kingSide: map[Color]bool{
				Color_White: parsedFen.whiteCanKingSideCastling,
				Color_Black: parsedFen.blackCanKingSideCastling,
			},
		},

		enPassantSq:     parsedFen.enPassantSq,
		halfmoveClock:   parsedFen.halfmoveClock,
		fullmoveCounter: parsedFen.fulmoveCounter,

		captures: make([]Piece, 0),
	}, nil
}

func (p Position) computePseudoMovements(color Color, doCastlingCheck bool) ([]Movement, [8][8]bool) {
	movements := make([]Movement, 0, 256)
	var attackMatrix [8][8]bool
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			attackMatrix[i][j] = false
		}
	}

	for _, row := range p.board {
		for _, piece := range row {
			if piece.Color == color {
				p.computePiecePseudoMovements(piece, &movements, doCastlingCheck, &attackMatrix)
			}
		}
	}

	return movements, attackMatrix
}

func (p Position) computePiecePseudoMovements(piece Piece, movements *[]Movement, doCastlingCheck bool, attackMatrix *[8][8]bool) {
	switch piece.Kind {
	case Kind_Bishop:
		p.computeDirectionPseudoMovements(piece, movements, attackMatrix, bishopDirections)
		return
	case Kind_Rook:
		p.computeDirectionPseudoMovements(piece, movements, attackMatrix, rookDirections)
		return
	case Kind_Queen:
		p.computeDirectionPseudoMovements(piece, movements, attackMatrix, bishopDirections)
		p.computeDirectionPseudoMovements(piece, movements, attackMatrix, rookDirections)
		return
	case Kind_King:
		for _, offset := range kingOffsets {
			targetRow, targetCol := int8(piece.Square.I)+offset[0], int8(piece.Square.J)+offset[1]
			if targetRow >= 0 && targetCol >= 0 && targetRow < 8 && targetCol < 8 {
				row := uint8(targetRow)
				col := uint8(targetCol)
				pieceAt := p.board[row][col]

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

			_, enemyAttackBoard := p.computePseudoMovements(piece.Color.Opposite(), false)
			// If king is not in check, continue
			if enemyAttackBoard[piece.Square.I][piece.Square.J] == false {
				if p.castlingRights.queenSide[piece.Color] {
					// Check if space to rook is empty
					canCastle := true
					for j := piece.Square.J - 1; j >= piece.Square.J-3; j-- {
						if p.board[piece.Square.I][j].Kind != Kind_None {
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

				if p.castlingRights.kingSide[piece.Color] {
					// Check if space to rook is empty
					canCastle := true
					for j := piece.Square.J + 1; j < 7; j++ {
						if p.board[piece.Square.I][j].Kind != Kind_None {
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

				pieceAt := p.board[row][col]

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
				pieceAt := p.board[row][piece.Square.J]
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

				pieceAt := p.board[row][col]

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
				if pieceAt.Kind == Kind_None && p.enPassantSq != nil && p.enPassantSq.I == row && p.enPassantSq.J == col {
					enPassantPieceSquare := newSquare(p.enPassantSq.I+1, p.enPassantSq.J)
					if piece.Color == Color_Black {
						enPassantPieceSquare.I = p.enPassantSq.I - 1
					}

					pieceAt := p.board[enPassantPieceSquare.I][enPassantPieceSquare.J]

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

func (p Position) computeDirectionPseudoMovements(piece Piece, movements *[]Movement, attackMatrix *[8][8]bool, directions [4][2]int8) {
	for _, dir := range directions {
		i := int8(piece.Square.I) + dir[0]
		j := int8(piece.Square.J) + dir[1]

		for i >= 0 && j >= 0 && i < 8 && j < 8 {
			row := uint8(i)
			col := uint8(j)

			pieceAt := p.board[row][col]

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

// TODO: Save king's positions (in Position{}, or get them at computePseudoMovements() to reuse the loop, if possible)
func (p Position) checkForCheck(allyColor Color, opponentAttackMatrix *[8][8]bool) bool {
	for i := uint8(0); i < 8; i++ {
		for j := uint8(0); j < 8; j++ {
			if p.board[i][j].Kind == Kind_King && p.board[i][j].Color == allyColor {
				return opponentAttackMatrix[i][j]
			}
		}
	}

	// Won't get here, unless there is no King piece ??
	return false
}
