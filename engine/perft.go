package engine

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func (b *Board) Perft(initialDepth, depth int, currentMove string, positionVerbose bool) int {
	var nMoves, i int
	nodes := 0

	if depth == 0 {
		if positionVerbose && initialDepth == 1 {
			fmt.Printf("%s: %d\n", currentMove, 1)
		}
		return 1
	}

	moveList := b.GetLegalMovements(b.PlayerToMove)
	nMoves = len(moveList)

	for i = 0; i < nMoves; i++ {
		b.MakeMovement(moveList[i])
		nodes += b.Perft(initialDepth, depth-1, moveList[i].ToAlgebraic(), positionVerbose)
		b.UndoMovement(moveList[i])
	}

	if positionVerbose && depth == initialDepth-1 && currentMove != "" {
		fmt.Printf("\t%s: %d\n", currentMove, nodes)
		fmt.Println(b.ToFen())
	}

	return nodes
}

func RunPerftTest(perftTest PerftTest, outputMode string, positionVerbose bool) {
	if outputMode == "short" {
		fmt.Printf("- Perft test: '%s'. Passed: ", perftTest.fen)
	} else {
		fmt.Printf("## Running a Perft [max-depth: %d] ##\n", perftTest.maxDepth)
		fmt.Printf("\tFen: '%s'\n", perftTest.fen)
		fmt.Println("\tResults:")
	}
	board := NewBoardFromFen(perftTest.fen)

	testBegin := time.Now()

	for depth := 1; depth <= perftTest.maxDepth; depth++ {
		currentDepthBegin := time.Now()
		result := board.Perft(depth, depth, "", positionVerbose)
		currentDepthSpentMs := time.Now().Sub(currentDepthBegin).Milliseconds()

		if outputMode != "short" {
			fmt.Printf("\t\tDepth %d -> %d (%d milliseconds)\n", depth, result, currentDepthSpentMs)
		}

		if val, ok := perftTest.depthMap[depth]; ok && val != result {
			fmt.Printf("\t\t\t ❌ failed, stopping Perft. Expected nodes: %d, got: %d at depth %d\n", val, result, depth)
			os.Exit(1)
		}
	}

	testSpentS := time.Now().Sub(testBegin).Seconds()
	if outputMode == "short" {
		fmt.Println("true")
	} else {
		fmt.Printf("\t\t\t ✔ Test succeeded in %02f seconds.\n", testSpentS)
	}
}

type PerftTest struct {
	fen      string
	depthMap map[int]int
	maxDepth int
}

func RunPerftsFromEpdFile(filename, outputMode string, maxDepth int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	perftBegin := time.Now()

	perftTests := make([]PerftTest, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		if strings.HasPrefix(lineText, "#") {
			continue
		}

		parts := strings.SplitN(lineText, ";", 2)

		fen := strings.TrimSpace(parts[0])
		depths := strings.Split(strings.TrimSpace(parts[1]), " ;")

		depthMap := make(map[int]int)

		maxDepthFound := 0

		for _, depth := range depths {
			depthParts := strings.Split(depth, " ")

			depthAmount, _ := strconv.Atoi(string(depthParts[0][1]))
			nodeAmount, _ := strconv.Atoi(depthParts[1])

			if depthAmount > maxDepthFound && depthAmount <= maxDepth {
				maxDepthFound = depthAmount
			}

			depthMap[depthAmount] = nodeAmount
		}

		perftTests = append(perftTests, PerftTest{
			fen:      fen,
			depthMap: depthMap,
			maxDepth: maxDepthFound,
		})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	slices.SortFunc(perftTests, func(a, b PerftTest) int {
		if a.depthMap[a.maxDepth] < b.depthMap[a.maxDepth] {
			return -1
		}
		return 1
	})

	fmt.Printf("### %d TESTS FOUND ###\n", len(perftTests))
	for _, perftTest := range perftTests {
		RunPerftTest(perftTest, outputMode, false)
	}

	perftSpentNs := time.Now().Sub(perftBegin).Nanoseconds()
	perftSpentUs := time.Now().Sub(perftBegin).Microseconds()
	perftSpentMs := time.Now().Sub(perftBegin).Milliseconds()
	perftSpentS := time.Now().Sub(perftBegin).Seconds()

	fmt.Printf("Time spent: %fs | %dms | %dus | %dns\n", perftSpentS, perftSpentMs, perftSpentUs, perftSpentNs)
	os.Exit(1)
}
