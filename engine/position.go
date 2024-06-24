package engine

import (
	"fmt"
	"unicode"
)

type Position struct {
	Board Board

	Status PositionStatus
}

type CastlingRights struct {
	QueenSide map[Color]bool
	KingSide  map[Color]bool
}

func (cr *CastlingRights) DeepCopy() CastlingRights {
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
	PlayerToMove Color

	CastlingRights CastlingRights

	EnPassant *Point
}

func (ps *PositionStatus) DeepCopy() PositionStatus {
	return PositionStatus{
		PlayerToMove:   ps.PlayerToMove,
		CastlingRights: ps.CastlingRights.DeepCopy(),
		EnPassant:      ps.EnPassant,
	}
}

func NewPositionFromFen(fen string) Position {
	parsedFen, err := parseFen(fen)
	if err != nil {
		panic(err)
	}

	return Position{
		Board: NewBoardFromFen(parsedFen.PlacementData),
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

			EnPassant: nil, //TODO
		},
	}
}

// TODO: Complete
func (p Position) ToFen() string {
	dataFen := ""
	spaceAccum := 0

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if p.Board[i][j].Kind == Kind_None {
				spaceAccum++
			} else {
				if spaceAccum > 0 {
					dataFen = fmt.Sprintf("%s%d", dataFen, spaceAccum)
					spaceAccum = 0
				}
				kindRune := p.Board[i][j].Kind.ToRune()
				if p.Board[i][j].Color == Color_White {
					kindRune = unicode.ToUpper(kindRune)
				}
				dataFen = fmt.Sprintf("%s%c", dataFen, kindRune)
			}
		}

		if spaceAccum > 0 {
			dataFen = fmt.Sprintf("%s%d", dataFen, spaceAccum)
			spaceAccum = 0
		}
		dataFen = fmt.Sprintf("%s/", dataFen)
	}

	dataFen = fmt.Sprintf("%s", dataFen[:len(dataFen)-1])

	dataFen = fmt.Sprintf("%s %c ", dataFen, p.Status.PlayerToMove.ToRune())

	if p.Status.CastlingRights.KingSide[Color_White] {
		dataFen = fmt.Sprintf("%sK", dataFen)
	}
	if p.Status.CastlingRights.QueenSide[Color_White] {
		dataFen = fmt.Sprintf("%sQ", dataFen)
	}
	if p.Status.CastlingRights.KingSide[Color_Black] {
		dataFen = fmt.Sprintf("%sk", dataFen)
	}
	if p.Status.CastlingRights.QueenSide[Color_Black] {
		dataFen = fmt.Sprintf("%sq", dataFen)
	}

	if p.Status.EnPassant != nil {
		dataFen = fmt.Sprintf("%s %s", dataFen, p.Status.EnPassant.ToAlgebraic())
	} else {
		dataFen = fmt.Sprintf("%s -", dataFen)
	}

	dataFen = fmt.Sprintf("%s 0 1", dataFen)

	return dataFen
}
