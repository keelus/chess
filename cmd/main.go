package main

import (
	"chess/board"
	"chess/piece"
	"chess/textures"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BOARD_SIZE int32 = 8
	CELL_SIZE  int32 = 100

	SCREEN_WIDTH  int32 = BOARD_SIZE * CELL_SIZE
	SCREEN_HEIGHT int32 = SCREEN_WIDTH
)

var (
	BROWN_LIGHT = rl.NewColor(241, 217, 181, 255)
	BROWN_DARK  = rl.NewColor(181, 136, 99, 255)
)

func init() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Chess game")
	rl.SetTargetFPS(60)

	textures.LoadTextures()
}

func main() {
	currentBoard := board.NewBoardFromFen("8/8/2p1p3/3P4/8/8/8/8 w KQkq - 0 1")

	var activePiece *piece.Piece = nil

	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			i := int(math.Floor(float64(rl.GetMousePosition().Y) / float64(CELL_SIZE)))
			j := int(math.Floor(float64(rl.GetMousePosition().X) / float64(CELL_SIZE)))

			if i >= 0 && j >= 0 && i < 8 && j < 8 {
				activePiece = currentBoard.GetPieceAt(i, j)
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

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

				if currentBoard.Data[i][j] != nil {
					currentPiece := currentBoard.Data[i][j]
					rl.DrawTexture(textures.GetPieceTexture(currentPiece.Color, currentPiece.Kind), j*CELL_SIZE, i*CELL_SIZE, rl.RayWhite)
				}
			}
		}

		if activePiece != nil {
			totalMovements := 0

			for _, m := range currentBoard.GetPseudoMovements() {
				if m.MovingPiece == activePiece {
					totalMovements++
					rl.DrawRectangle(int32(m.To.J)*CELL_SIZE, int32(m.To.I)*CELL_SIZE, CELL_SIZE, CELL_SIZE, rl.NewColor(209, 121, 27, 255))
				}
			}

			fmt.Printf("This piece has %d movements available.\n", totalMovements)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
