package engine

import (
	"errors"
	"strconv"
	"strings"
)

type FenData struct {
	PlacementData [8]string
	ActiveColor   Color
	Fullmoves     int
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

	// TODO: Castling (parts[2])
	// TODO: En passant (parts[3])
	// TODO: Halfmove clock (parts[4])

	fullmoves, err := strconv.ParseInt(parts[5], 10, 0)
	if err != nil || fullmoves < 1 {
		return FenData{}, errors.New("The provided FEN does not have a valid fullmove number.")
	}

	return FenData{
		PlacementData: [8]string(lacementParts),
		ActiveColor:   ColorFromRune(activeColor),
		Fullmoves:     int(fullmoves),
	}, nil
}
