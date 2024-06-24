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
