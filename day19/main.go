package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	available, desired := loadData("day19/input.txt")

	// Create a map to track segments that can be built from available patterns
	validSegments := make(map[string]bool, 0)

	// Available patterns themselves are always valid segments
	for _, a := range available {
		validSegments[a] = true
	}

	count := 0
	for _, d := range desired {
		found := design(d, available, &validSegments)
		if found {
			count++
		}
	}

	fmt.Println("Number of possible designs:", count)
}

func design(desired string, available []string, validSegments *map[string]bool) (possible bool) {
	// Done, no part of the desired pattern remains to check for
	if len(desired) == 0 {
		return true
	}

	// Check if the desired segment is already in the map of valid segments
	if desiredIsValid, found := (*validSegments)[desired]; found {
		return desiredIsValid
	}

	// Check if the desired segment begins with any of the available patterns
	for _, a := range available {
		if strings.HasPrefix(desired, a) {
			// If so, continue recursing on the remaining segment
			if design(desired[len(a):], available, validSegments) {
				// If the remaining segment is valid, cache that it's valid
				(*validSegments)[desired] = true
				return true
			}
		}
	}

	// At this point, we know the desired segment can't be built from available patterns
	(*validSegments)[desired] = false

	return false
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
