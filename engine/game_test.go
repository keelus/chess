package engine

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	_ "net/http"
	_ "net/http/pprof"
)

type PerftTest struct {
	fen      string
	depthMap map[int]int
	maxDepth int
}

var loadedPerftTests []PerftTest

var maxDepth int // Default: 1

// Each line of the file should contain only one test position, along with it's depth result.
// If you want to ignore a test, write # in the first char of the line.
// Example line-test, with FEN position, and depth with results from 1 to 6:
// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1;D1 20;D2 400;D3 8902;D4 197281;D5 4865609;D6 119060324
var perftEdpFile string // Default: "perft_test.epd"

// If you wish to print each movement's nodes for each depth done, set this flag.
var positionVerbose bool // Default: false

func init() {
	flag.IntVar(&maxDepth, "maxDepth", 1, "The max amount of depth to explore in a test")
	flag.StringVar(&perftEdpFile, "epd", "perft_tests.epd", "The name of the file in which the tests are written. Check code for more information about the file format.")
	flag.BoolVar(&positionVerbose, "positionVerbose", false, "If set, prints each movement's nodes, for each depth done. Note: The tests will be runned one by one instead of in parallel.")
}

func parsePerftFile() {
	loadedPerftTests = make([]PerftTest, 0)

	file, err := os.Open(perftEdpFile)
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Errorf("The perft test file with the name \"%s\" was not found. Use -epd=<filename> to specify a perft test file.", perftEdpFile))
		} else {
			panic(err)
		}
	}
	defer file.Close()

	//perftBegin := time.Now()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		if strings.HasPrefix(lineText, "#") {
			continue
		}

		parts := strings.SplitN(lineText, ";", 2)

		fen := strings.TrimSpace(parts[0])
		depths := strings.Split(strings.TrimSpace(parts[1]), ";")

		depthMap := make(map[int]int)

		maxDepthFound := 0

		for _, depth := range depths {
			depthParts := strings.Split(strings.TrimSpace(depth), " ")
			if len(depthParts) != 2 {
				panic(fmt.Errorf("There was an error parsing the depth \"%s\".", depth))
			}

			depthAmount, _ := strconv.Atoi(string(depthParts[0][1]))
			nodeAmount, _ := strconv.Atoi(depthParts[1])

			if depthAmount > maxDepthFound && depthAmount <= maxDepth {
				maxDepthFound = depthAmount
			}

			depthMap[depthAmount] = nodeAmount
		}

		loadedPerftTests = append(loadedPerftTests, PerftTest{
			fen:      fen,
			depthMap: depthMap,
			maxDepth: maxDepthFound,
		})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Sort Perft tests in increasing amount of nodes at max depth
	slices.SortFunc(loadedPerftTests, func(a, b PerftTest) int {
		if a.depthMap[a.maxDepth] < b.depthMap[a.maxDepth] {
			return -1
		}
		return 1
	})

	fmt.Printf("Running %d Perfts at max-depth of %d", len(loadedPerftTests), maxDepth)
}

func TestPerft(t *testing.T) {
	parsePerftFile()
	for i, perftTest := range loadedPerftTests {
		i := i
		perftTest := perftTest
		t.Run(fmt.Sprintf("Test %d: '%s'", i, perftTest.fen), func(t *testing.T) {
			game := NewGame(perftTest.fen)
			if !positionVerbose {
				t.Parallel()
			}

			for depth := 1; depth <= perftTest.maxDepth; depth++ {
				if positionVerbose {
					fmt.Printf("Evaluation depth: %d\n", depth)
				}

				result := game.Perft(depth, depth, "", positionVerbose)

				if val, ok := perftTest.depthMap[depth]; ok && val != result {
					t.Fatalf("\tTest_%d: failed at depth: %d. Expected nodes: %d, got: %d\n\n", i, depth, perftTest.depthMap[depth], result)
				}
			}
		})
	}
}

func BenchPerft(b *testing.B) {
	parsePerftFile()
	b.ResetTimer()
	for i, perftTest := range loadedPerftTests {
		b.Run(fmt.Sprintf("Test %d: '%s'", i, perftTest.fen), func(bb *testing.B) {
			game := NewGame(perftTest.fen)

			for depth := 1; depth <= perftTest.maxDepth; depth++ {
				result := game.Perft(depth, depth, "", positionVerbose)

				if val, ok := perftTest.depthMap[depth]; ok && val != result {
					b.Fatalf("failed at depth: %d. Expected nodes: %d, got: %d\n", depth, val, result)
				}
			}
		})
	}
}
