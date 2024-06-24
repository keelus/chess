<h1 align="center">Chess</h1>

<p align="center">
  <a href="./LICENSE"><img src="https://img.shields.io/badge/⚖️ license-MIT-blue" alt="MIT License"></a>
</p>

## ℹ️ Description
A chess engine written in Go.

For now, I'm using Raylib for graphics and user handling, but this will be removed once the engine is fully worked, making the "GUI" in a different repo, using the engine.

> [!NOTE]
> The project is in very early development stage.


## 📸 Screenshots
<img src="https://github.com/keelus/chess/assets/86611436/4900e816-3c28-45ca-bb5d-4f37358ececf" width=400 />


##  Tests passed
The tests are being manually made and compared with the Perft Results provided [here](https://www.chessprogramming.org/Perft_Results).
Currently, the tests are only compared using the `nodes` count.

| Position name       | Depth 1   | Depth 2   | Depth 3   | Depth 4   | Depth 5   | Depth 6   |
|---------------------|-----------|-----------|-----------|-----------|-----------|-----------|
| Initial position    | ✅       | ✅        | ✅        | ✅        | ✅        | ❓        |
| Position 2          | ✅       | ✅        | ✅        | ✅        | ✅        | ❓        |
| Position 3          | ✅       | ✅        | ✅        | ✅        | ✅        | ✅        |
| Position 4          | ✅       | ✅        | ✅        | ✅        | ✅        | ❓        |
| Position 5          | ✅       | ✅        | ✅        | ✅        | ✅        | ❓        |
| Position 6          | ✅       | ✅        | ✅        | ✅        | ✅        | ❓        |

To make a test or bench a given `.epd` file (check `engine/game_test.go` to see the file format and command arguments to the test):


#### Test
While inside `engine` folder:
```
go test -v -epd filename.epd -maxDepth 4
```
To verbose each movement's nodes and FEN,
```
go test -v -epd filename.epd -maxDepth 4 -positionVerbose
```

#### Bench
While inside `engine` folder:
```
go test -v -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out 4 -epd .\perft_tests\filename.epd -maxDepth 4
```

## ⚖️ License
This project is open source under the terms of the [MIT License](./LICENSE)

<br />
Made by <a href="https://github.com/keelus">keelus</a> ✌️

