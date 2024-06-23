package main

import (
	"chess/engine"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var loadedTextures map[engine.Color]map[engine.Kind]rl.Texture2D

func loadPieceTexture(filename string) rl.Texture2D {
	image := rl.LoadImage(filename)
	defer rl.UnloadImage(image)

	rl.ImageResize(image, CELL_SIZE, CELL_SIZE)

	return rl.LoadTextureFromImage(image)
}

func LoadTextures() {
	loadedTextures = make(map[engine.Color]map[engine.Kind]rl.Texture2D)

	for c := 0; c < engine.COLOR_AMOUNT; c++ {
		castedColor := engine.Color(c)
		loadedTextures[castedColor] = make(map[engine.Kind]rl.Texture2D)

		for k := 0; k < engine.KIND_AMOUNT; k++ {
			castedKind := engine.Kind(k)

			kindRune := castedKind.ToRune()
			colorRune := castedColor.ToRune()

			filename := fmt.Sprintf("./media/%c_%c.png", colorRune, kindRune)
			loadedTexture := loadPieceTexture(filename)

			loadedTextures[castedColor][castedKind] = loadedTexture
		}
	}
}

func GetPieceTexture(color engine.Color, kind engine.Kind) rl.Texture2D {
	return loadedTextures[color][kind]
}
