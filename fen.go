package chess

import (
	"errors"
	"strconv"
	"strings"
)

type fenData struct {
	placementData [8]string
	activeColor   Color

	whiteCanKingSideCastling  bool
	whiteCanQueenSideCastling bool
	blackCanKingSideCastling  bool
	blackCanQueenSideCastling bool

	enPassantSq *Square

	halfmoveClock  uint8
	fulmoveCounter uint
}

// IsFenValid returns whether the passed FEN string
// is valid or not.
func IsFenValid(fen string) bool {
	_, err := parseFen(fen)
	return err == nil
}

func parseFen(fen string) (fenData, error) {
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return fenData{}, errors.New("The provided FEN does not have 6 parts.")
	}

	lacementParts := strings.Split(parts[0], "/")
	if len(lacementParts) != 8 {
		return fenData{}, errors.New("The provided FEN does not have 8 placement positions in part 1.")
	}

	activeColor := rune(parts[1][0])
	if activeColor != 'w' && activeColor != 'b' {
		return fenData{}, errors.New("The provided FEN does not have a valid active color.")
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
			return fenData{}, errors.New("The provided FEN does not have a valid en passant. It must be in algebraic (e.g: d6).")
		}

		enPassant = &square
	}

	halfmoveClock, err := strconv.ParseUint(parts[4], 10, 0)
	if err != nil {
		return fenData{}, errors.New("The provided FEN does not have a valid halfmove number.")
	}

	fullmoveCounter, err := strconv.ParseUint(parts[5], 10, 0)
	if err != nil || fullmoveCounter < 1 {
		return fenData{}, errors.New("The provided FEN does not have a valid fullmove number.")
	}

	return fenData{
		placementData:  [8]string(lacementParts),
		activeColor:    ColorFromRune(activeColor),
		enPassantSq:    enPassant,
		halfmoveClock:  uint8(halfmoveClock),
		fulmoveCounter: uint(fullmoveCounter),

		whiteCanKingSideCastling:  whiteCanKingSideCastling,
		whiteCanQueenSideCastling: whiteCanQueenSideCastling,
		blackCanKingSideCastling:  blackCanKingSideCastling,
		blackCanQueenSideCastling: blackCanQueenSideCastling,
	}, nil
}
