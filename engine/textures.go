package engine

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var loadedTextures map[Color]map[Kind]rl.Texture2D

const CELL_SIZE int32 = 100

func loadPieceTexture(filename string) rl.Texture2D {
	image := rl.LoadImage(filename)
	defer rl.UnloadImage(image)

	rl.ImageResize(image, CELL_SIZE, CELL_SIZE)

	return rl.LoadTextureFromImage(image)
}

func LoadTextures() {
	loadedTextures = make(map[Color]map[Kind]rl.Texture2D)

	for c := 0; c < COLOR_AMOUNT; c++ {
		castedColor := Color(c)
		loadedTextures[castedColor] = make(map[Kind]rl.Texture2D)

		for k := 0; k < KIND_AMOUNT; k++ {
			castedKind := Kind(k)

			kindRune := castedKind.ToRune()
			colorRune := castedColor.ToRune()

			filename := fmt.Sprintf("./media/%c_%c.png", colorRune, kindRune)
			loadedTexture := loadPieceTexture(filename)

			loadedTextures[castedColor][castedKind] = loadedTexture
		}
	}
}

func GetPieceTexture(color Color, kind Kind) rl.Texture2D {
	return loadedTextures[color][kind]
}
