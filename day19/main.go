package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	available, desired := loadData("day19/input.txt")

	validSegments := make(map[string]int)
	for _, a := range available {
		validSegments[a] = design(a, available, validSegments)
	}

	possible := 0
	combinations := 0
	for _, d := range desired {
		found := design(d, available, validSegments)
		if found > 0 {
			possible++
		}
		combinations += found
	}

	fmt.Println("Number of possible designs:", possible)
	fmt.Println("Number of possible combinations:", combinations)
}

func design(desired string, available []string, validSegments map[string]int) int {
	// Done, no part of the desired pattern remains to check for, we've found a valid point
	if len(desired) == 0 {
		return 1
	}

	// Check if the desired segment is already in the map of valid segments
	if count, found := validSegments[desired]; found {
		return count
	}

	// Count how many ways the desired segment contains the available patterns
	for _, a := range available {
		if strings.HasPrefix(desired, a) {
			nextSegment := desired[len(a):]
			validSegments[desired] += design(nextSegment, available, validSegments)
			// fmt.Println(validSegments)
		}
	}

	return validSegments[desired]
}

func loadData(filename string) (available, desired []string) {
	processingAvailable := true
	common.ProcessFile(filename, func(line string) {
		if line == "" {
			processingAvailable = false
			return
		}

		if processingAvailable {
			available = strings.Split(line, ", ")
		} else {
			desired = append(desired, line)
		}
	})
	return available, desired
}
