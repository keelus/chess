package textures

import (
	"chess/piece"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var loadedTextures map[piece.Color]map[piece.Kind]rl.Texture2D

const CELL_SIZE int32 = 100

func loadPieceTexture(filename string) rl.Texture2D {
	image := rl.LoadImage(filename)
	defer rl.UnloadImage(image)

	rl.ImageResize(image, CELL_SIZE, CELL_SIZE)

	return rl.LoadTextureFromImage(image)
}

func LoadTextures() {
	loadedTextures = make(map[piece.Color]map[piece.Kind]rl.Texture2D)

	for c := 0; c < piece.COLOR_AMOUNT; c++ {
		castedColor := piece.Color(c)
		loadedTextures[castedColor] = make(map[piece.Kind]rl.Texture2D)

		for k := 0; k < piece.KIND_AMOUNT; k++ {
			castedKind := piece.Kind(k)

			kindRune := castedKind.ToRune()
			colorRune := castedColor.ToRune()

			filename := fmt.Sprintf("./media/%c_%c.png", colorRune, kindRune)
			loadedTexture := loadPieceTexture(filename)

			loadedTextures[castedColor][castedKind] = loadedTexture
		}
	}
}

func GetPieceTexture(color piece.Color, kind piece.Kind) rl.Texture2D {
	return loadedTextures[color][kind]
}
