package main

import (
	"chess/engine"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BOARD_SIZE int32 = 8
	CELL_SIZE  int32 = 100

	BOTTOM_BAR int32 = 100

	SCREEN_WIDTH  int32 = BOARD_SIZE * CELL_SIZE
	SCREEN_HEIGHT int32 = SCREEN_WIDTH + BOTTOM_BAR
)

var (
	BROWN_LIGHT = rl.NewColor(241, 217, 181, 255)
	BROWN_DARK  = rl.NewColor(181, 136, 99, 255)
)

func init() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Chess game")
	rl.SetTargetFPS(60)

	engine.LoadTextures()
}

func main() {
	//board := engine.NewBoardFromFen("8/3p4/2P1P3/8/8/8/8/8 b KQkq - 0 1")
	//board := engine.NewBoardFromFen("r3k2r/8/8/6p4/3P3/8/7P/8 b KQkq - 0 1")
	board := engine.NewBoardFromFen("8/8/8/8/4p3/8/3P4/8 w KQkq - 0 1")

	var activePiece *engine.Piece = nil
	var lastMovement *engine.Movement = nil

	for !rl.WindowShouldClose() {
		currentMovements := board.GetPseudoMovements()
		//currentMovements := board.GetLegalMovements()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			i := int(math.Floor(float64(rl.GetMousePosition().Y) / float64(CELL_SIZE)))
			j := int(math.Floor(float64(rl.GetMousePosition().X) / float64(CELL_SIZE)))

			clickedAMovement := false
			if activePiece != nil {
				for _, movement := range currentMovements {
					if movement.MovingPiece == activePiece {
						if movement.To.I == i && movement.To.J == j {
							clickedAMovement = true
							board.MakeMovement(movement)
							lastMovement = &movement
							activePiece = nil
						}
					}
				}
			}

			if !clickedAMovement && i >= 0 && j >= 0 && i < 8 && j < 8 {
				activePiece = board.GetPieceAt(i, j)
			}
		}

		if rl.IsKeyPressed(rl.KeyU) {
			if lastMovement != nil {
				board.UndoMovement(*lastMovement)
				lastMovement = nil
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(32, 32, 32, 255))

		for i := int32(0); i < BOARD_SIZE; i++ {
			for j := int32(0); j < BOARD_SIZE; j++ {
				var cellColor rl.Color
				if i%2 != 0 && j%2 != 0 {
					cellColor = BROWN_LIGHT
				} else if i%2 != 0 && j%2 == 0 {
					cellColor = BROWN_DARK
				} else if i%2 == 0 && j%2 != 0 {
					cellColor = BROWN_DARK
				} else if i%2 == 0 && j%2 == 0 {
					cellColor = BROWN_LIGHT
				}

				rl.DrawRectangle(j*CELL_SIZE, i*CELL_SIZE, CELL_SIZE, CELL_SIZE, cellColor)

				if board.Data[i][j] != nil {
					currentPiece := board.Data[i][j]
					rl.DrawTexture(engine.GetPieceTexture(currentPiece.Color, currentPiece.Kind), j*CELL_SIZE, i*CELL_SIZE, rl.RayWhite)
				}
			}
		}

		if activePiece != nil {
			totalMovements := 0

			for _, m := range currentMovements {
				if m.MovingPiece == activePiece {
					//fmt.Println(m)
					totalMovements++

					cellColor := rl.NewColor(209, 121, 27, 127)
					if m.TakingPiece != nil {
						cellColor = rl.NewColor(209, 42, 27, 127)
					}

					rl.DrawRectangle(int32(m.To.J)*CELL_SIZE, int32(m.To.I)*CELL_SIZE, CELL_SIZE, CELL_SIZE, cellColor)
				}
			}

			fmt.Printf("This piece has %d movements available.\n", totalMovements)
		}

		// Draw team information
		drawPlayerInformation := func(xPos int32, color engine.Color) {
			rl.DrawText(fmt.Sprintf("Player %c", color.ToRune()), xPos+5, SCREEN_HEIGHT-BOTTOM_BAR+2, 20, rl.RayWhite)
			rl.DrawText(
				fmt.Sprintf("-Can queen castle: %t", board.CanQueenCastling[engine.Color_White]),
				xPos+5, SCREEN_HEIGHT-BOTTOM_BAR+2+20,
				16,
				func() rl.Color {
					if board.CanQueenCastling[color] {
						return rl.DarkGreen
					}
					return rl.Maroon
				}(),
			)
			rl.DrawText(
				fmt.Sprintf("-Can king castle: %t", board.CanKingCastling[engine.Color_White]),
				xPos+5, SCREEN_HEIGHT-BOTTOM_BAR+2+20+16,
				16,
				func() rl.Color {
					if board.CanKingCastling[color] {
						return rl.DarkGreen
					}
					return rl.Maroon
				}(),
			)
		}

		drawPlayerInformation(0, engine.Color_White)
		drawPlayerInformation(300, engine.Color_Black)

		rl.DrawText(
			fmt.Sprintf("En passant: %t", board.EnPassant != nil),
			600+5, SCREEN_HEIGHT-BOTTOM_BAR+2,
			16,
			func() rl.Color {
				if board.EnPassant != nil {
					return rl.DarkGreen
				}
				return rl.Maroon
			}(),
		)
		if board.EnPassant != nil {
			rl.DrawText(
				fmt.Sprintf("{i: %d, j: %d}", board.EnPassant.I, board.EnPassant.J),
				600+5, SCREEN_HEIGHT-BOTTOM_BAR+2+16,
				16,
				rl.DarkGreen,
			)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
