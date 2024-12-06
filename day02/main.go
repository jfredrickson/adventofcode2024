package main

import (
	"common"
	"fmt"
	"slices"
	"strings"
)

func main() {
	inputFile := "day02/input.txt"

	fmt.Println("Safe reports:", getSafeReports(inputFile))
	fmt.Println("Dampened safe reports:", getDampenedSafeReports(inputFile))
}

func getSafeReports(inputFile string) int {
	safe := 0

	common.ProcessFile(inputFile, func(line string) {
		levels := common.ToInts(strings.Split(line, " "))
		currentSafe := true

		// If the first level is higher than the last, reverse the list
		if levels[0] > levels[len(levels)-1] {
			slices.Reverse(levels)
		}

		// Check each level difference to determine safety
		currentSafe = checkSafety(levels)

		if currentSafe {
			safe++
		}
	})

	return safe
}

func getDampenedSafeReports(inputFile string) int {
	safe := 0

	common.ProcessFile(inputFile, func(line string) {
		levels := common.ToInts(strings.Split(line, " "))
		currentSafe := true

		// If the first level is higher than the last, reverse the list
		if levels[0] > levels[len(levels)-1] {
			slices.Reverse(levels)
		}

		// Check each level difference to determine safety
		currentSafe = checkSafety(levels)

		// If unsafe, try a second pass that tries removing each level one by one
		if !currentSafe {
			for i := range len(levels) {
				testLevels := append([]int(nil), levels...)
				testLevels = slices.Delete(testLevels, i, i+1)
				testSafe := checkSafety(testLevels)
				if testSafe {
					currentSafe = true
					break
				}
			}
		}

		if currentSafe {
			safe++
		}
	})

	return safe
}

func checkSafety(levels []int) bool {
	safe := true

	for i := range levels {
		if i == 0 {
			continue
		}

		difference := levels[i] - levels[i-1]
		if difference < 1 || difference > 3 {
			safe = false
		}
	}

	return safe
}
