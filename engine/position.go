package engine

import (
	"fmt"
)

type Position struct {
	Board Board

	Status PositionStatus
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

func (p *Position) GetHalfmoveClock() uint8 {
	return p.Status.HalfmoveClock
}
func (p *Position) GetFullmoveCounter() uint {
	return p.Status.FullmoveCounter
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

			EnPassant:       nil, //TODO
			HalfmoveClock:   parsedFen.HalfmoveClock,
			FullmoveCounter: parsedFen.FulmoveCounter,
		},
	}
}

// TODO: Complete
func (p Position) Fen() string {
	boardFen := p.Board.Fen()

	dataFen := fmt.Sprintf("%s %c ", boardFen, p.Status.PlayerToMove.ToRune())

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
