<p align="center">
  <img src="https://github.com/keelus/chess/assets/86611436/4207ba13-3004-44ba-8325-7533f08a4e47" alt="Gopher carrying a Chess Queen piece." width="200" />
</p>
<h1 align="center">Chess</h1>

<p align="center">
<a href="./LICENSE"><img src="https://img.shields.io/badge/⚖️ license-MIT-blue" alt="MIT License"></a>
    <a href="https://godoc.org/github.com/keelus/chess"><img src="https://godoc.org/github.com/keelus/chess?status.svg" alt="GoDoc"></a>
</p>

## ℹ️ Description
This repository contains my own **Chess library** for **Golang**, made from scratch, with 0 dependencies.

You can use this library to handle a Chess game's logic. It supports **all the common chess rules** (including Checkmate, En Passant, Castling, etc).

This library handles the movement generation and legality checks, movement making, game turns, movement history, game's outcome (checkmate, stalemate, etc), unicode and FEN representation, and more.

It has been tuned and tested on multiple board positions and edge cases via **Perft** tests.

This is a hobby project I made for fun in the weekend, and I have decided to keep it simple (without using bitboards or magic numbers, at least for now), so its move generation it's not the fastest. Although it's robust and safely returns a position's legal movements within microseconds.

Feel free to use this library along with your own UI/app implementation.

## ⬇️ Installation
To install and use the library, use:
```sh
go get github.com/keelus/chess
```

## ♟️ Basic example
This is a basic hardcoded example, playing a game with random moves, and showcasing some functions.

You can read the documentation <a href="https://godoc.org/github.com/keelus/chess">here</a> to learn more
```Go
package main

import (
	"fmt"
	"math/rand"

	"github.com/keelus/chess"
)

func main() {
	game, _ := chess.NewGame("")
	fmt.Println(game.CurrentPosition().Board().Unicode())

	// Main game loop (by random movements)
	for game.Outcome() == chess.Outcome_None {
		legalMoves := game.LegalMovements()

		move := legalMoves[rand.Intn(len(legalMoves))]
		if err := game.MakeMovement(move); err != nil {
			// Handle error if needed
		}

		if game.IsMovementLegalAlgebraic("f7f5") {
			game.MakeMovementAlgebraic("f7f5")
		}

		pos := game.CurrentPosition()
		if pos.Turn() == chess.Color_White {
			// ...
		}

		// ...
	}

	fmt.Println("Game has ended.")
	fmt.Printf("\t- Total captures: %d\n", len(game.CurrentPosition().Captures()))
	fmt.Printf("\t- Outcome: %s\n", game.Outcome())
	fmt.Printf("\t- Final FEN: \"%s\"", game.CurrentFen())
}
```

## ⚒ Testing and benchmarking
To perform a test or a benchmark in movement generation, you can use the provided game_test.go, using the golang's builtin test/benchmarking tool.
This test supports the following arguments:
- `-epd <filename>`: The EPD filename where the tests and depth results resides.
- `-positionVerbose`: If set, each depth's movement's node count will be outputted. This is useful for debugging.
- `-maxDepth <amount>`: The maximum depth to perform the tests.

For example:
```
go test/bench -v -epd filename.epd -maxDepth 4
```

The EPD file should contain one test per line, specifying the FEN and each depth's node values. For example:
```
rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1;D1 20;D2 400;D3 8902;D4 197281;D5 4865609;D6 119060324
r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1;D1 48;D2 2039;D3 97862;D4 4085603;D5 193690690;D6 8031647685
```

*epd examples based on the data from [chessprogramming.com's Perft results](https://www.chessprogramming.org/Perft_Results).*

## ⚖️ License
This project is open source under the terms of the [MIT License](./LICENSE)

<br />
Made by <a href="https://github.com/keelus">keelus</a> ✌️
