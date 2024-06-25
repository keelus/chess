package chess

import (
	"errors"
	"strconv"
	"strings"
)

type FenData struct {
	PlacementData [8]string
	ActiveColor   Color

	WhiteCanKingSideCastling  bool
	WhiteCanQueenSideCastling bool
	BlackCanKingSideCastling  bool
	BlackCanQueenSideCastling bool

	EnPassant *Square

	HalfmoveClock  uint8
	FulmoveCounter uint
}

func parseFen(fen string) (FenData, error) {
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return FenData{}, errors.New("The provided FEN does not have 6 parts.")
	}

	lacementParts := strings.Split(parts[0], "/")
	if len(lacementParts) != 8 {
		return FenData{}, errors.New("The provided FEN does not have 8 placement positions in part 1.")
	}

	activeColor := rune(parts[1][0])
	if activeColor != 'w' && activeColor != 'b' {
		return FenData{}, errors.New("The provided FEN does not have a valid active color.")
	}

	whiteCanKingSideCastling := false
	whiteCanQueenSideCastling := false
	blackCanKingSideCastling := false
	blackCanQueenSideCastling := false

	if parts[2] != "-" {
		if strings.Contains(parts[2], "K") {
			whiteCanKingSideCastling = true
		}
		if strings.Contains(parts[2], "Q") {
			whiteCanQueenSideCastling = true
		}

		if strings.Contains(parts[2], "k") {
			blackCanKingSideCastling = true
		}
		if strings.Contains(parts[2], "q") {
			blackCanQueenSideCastling = true
		}
	}

	var enPassant *Square = nil
	if parts[3] != "-" {
		square, err := NewSquareFromAlgebraic(parts[3])
		if err != nil {
			return FenData{}, errors.New("The provided FEN does not have a valid en passant. It must be in algebraic (e.g: d6).")
		}

		enPassant = &square
	}

	halfmoveClock, err := strconv.ParseUint(parts[4], 10, 0)
	if err != nil {
		return FenData{}, errors.New("The provided FEN does not have a valid halfmove number.")
	}

	fullmoveCounter, err := strconv.ParseUint(parts[5], 10, 0)
	if err != nil || fullmoveCounter < 1 {
		return FenData{}, errors.New("The provided FEN does not have a valid fullmove number.")
	}

	return FenData{
		PlacementData:  [8]string(lacementParts),
		ActiveColor:    ColorFromRune(activeColor),
		EnPassant:      enPassant,
		HalfmoveClock:  uint8(halfmoveClock),
		FulmoveCounter: uint(fullmoveCounter),

		WhiteCanKingSideCastling:  whiteCanKingSideCastling,
		WhiteCanQueenSideCastling: whiteCanQueenSideCastling,
		BlackCanKingSideCastling:  blackCanKingSideCastling,
		BlackCanQueenSideCastling: blackCanQueenSideCastling,
	}, nil
}

func IsFenValid(fen string) bool {
	_, err := parseFen(fen)
	return err == nil
}
