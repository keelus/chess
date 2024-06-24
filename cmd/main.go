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

	LoadTextures()
}

func main() {
	game := engine.NewGame("")

	var activeSquare *engine.Square = nil
	_ = activeSquare

	for !rl.WindowShouldClose() {
		outcome := game.Outcome()
		if outcome != engine.Outcome_None {
			break
		}
		currentMovements := []string{}
		if activeSquare != nil {
			currentMovements = game.GetLegalMovements()
		}
		_ = currentMovements
		//currentMovements := board.GetLegalMovements()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			i := uint8(math.Floor(float64(rl.GetMousePosition().Y) / float64(CELL_SIZE)))
			j := uint8(math.Floor(float64(rl.GetMousePosition().X) / float64(CELL_SIZE)))

			clickedAMovement := false
			// TODO: From.
			if activeSquare != nil {
				for _, m := range currentMovements {
					movementInformation, _ := game.MovementInformation(m)
					if movementInformation.From().I == activeSquare.I && movementInformation.From().J == activeSquare.J {
						if movementInformation.To().I == i && movementInformation.To().J == j {
							clickedAMovement = true
							game.MakeMovement(m)
							//lastMovement = &movement
							activeSquare = nil
							break
						}
					}
				}
			}

			if !clickedAMovement && i >= 0 && j >= 0 && i < 8 && j < 8 {
				square, _ := engine.NewSquare(i, j)
				if game.GetPieceAt(square).Kind == engine.Kind_None {
					activeSquare = nil
				} else {
					activeSquare = &square
				}
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

				square, _ := engine.NewSquare(uint8(i), uint8(j))
				if pieceAt := game.GetPieceAt(square); pieceAt.Kind != engine.Kind_None {
					rl.DrawTexture(GetPieceTexture(pieceAt.Color, pieceAt.Kind), j*CELL_SIZE, i*CELL_SIZE, rl.RayWhite)
				}
			}
		}

		// TODO: from. if.
		if activeSquare != nil {
			totalThisPieceMovements := 0

			for _, m := range game.GetLegalMovementsOfPiece(*activeSquare) {
				movementInformation, _ := game.MovementInformation(m)
				if movementInformation.From().I == activeSquare.I && movementInformation.From().J == activeSquare.J {
					totalThisPieceMovements++

					cellColor := rl.NewColor(209, 121, 27, 127)
					if movementInformation.IsCapturing() {
						cellColor = rl.NewColor(209, 42, 27, 127)
					}

					rl.DrawRectangle(int32(movementInformation.To().J)*CELL_SIZE, int32(movementInformation.To().I)*CELL_SIZE, CELL_SIZE, CELL_SIZE, cellColor)
				}
			}

			fmt.Printf("This piece has %d movements available.\n", totalThisPieceMovements)
		}

		// totalMovements := len(board.GetLegalMovements(board.PlayerToMove))
		// fmt.Printf("\n##### TOTAL MOVEMENTS NOW: %d #####\n\n", totalMovements)

		// Draw team information [[ Debug purposes. TODO: Remove ]]
		drawPlayerInformation := func(xPos int32, color engine.Color) {
			rl.DrawText(fmt.Sprintf("Player %c", color.ToRune()), xPos+5, SCREEN_HEIGHT-BOTTOM_BAR+2, 20, rl.RayWhite)
			rl.DrawText(
				fmt.Sprintf("-Can queen castle: %t", game.GetCurrentPosition().Status.CastlingRights.QueenSide[color]),
				xPos+5, SCREEN_HEIGHT-BOTTOM_BAR+2+20,
				16,
				func() rl.Color {
					if game.GetCurrentPosition().Status.CastlingRights.QueenSide[color] {
						return rl.DarkGreen
					}
					return rl.Maroon
				}(),
			)
			rl.DrawText(
				fmt.Sprintf("-Can king castle: %t", game.GetCurrentPosition().Status.CastlingRights.KingSide[color]),
				xPos+5, SCREEN_HEIGHT-BOTTOM_BAR+2+20+16,
				16,
				func() rl.Color {
					if game.GetCurrentPosition().Status.CastlingRights.KingSide[color] {
						return rl.DarkGreen
					}
					return rl.Maroon
				}(),
			)
		}

		drawPlayerInformation(0, engine.Color_White)
		drawPlayerInformation(300, engine.Color_Black)

		rl.DrawText(
			fmt.Sprintf("En passant: %t", game.GetCurrentPosition().Status.EnPassant != nil),
			600+5, SCREEN_HEIGHT-BOTTOM_BAR+2,
			16,
			func() rl.Color {
				if game.GetCurrentPosition().Status.EnPassant != nil {
					return rl.DarkGreen
				}
				return rl.Maroon
			}(),
		)
		if game.GetCurrentPosition().Status.EnPassant != nil {
			rl.DrawText(
				fmt.Sprintf("{i: %d, j: %d}", game.GetCurrentPosition().Status.EnPassant.I, game.GetCurrentPosition().Status.EnPassant.J),
				600+5, SCREEN_HEIGHT-BOTTOM_BAR+2+16,
				16,
				rl.DarkGreen,
			)
		}

		rl.DrawText(
			fmt.Sprintf("Turn color: %c", game.GetCurrentPosition().Status.PlayerToMove.ToRune()),
			600+5, SCREEN_HEIGHT-BOTTOM_BAR+2+20,
			16,
			rl.RayWhite,
		)

		rl.EndDrawing()
	}

	rl.CloseWindow()

	fmt.Println(game.Outcome())
}
