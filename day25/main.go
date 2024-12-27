package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	locks, keys := loadSchematics("day25/input.txt")

	fitCount := 0
	for _, lock := range locks {
		for _, key := range keys {
			if fits(lock, key) {
				fitCount++
			}
		}
	}

	fmt.Println("Fitting locks and keys:", fitCount)
}

func fits(lock, key []int) bool {
	for i := range lock {
		if lock[i]+key[i] > 5 {
			return false
		}
	}
	return true
}

func loadSchematics(filename string) (locks [][]int, keys [][]int) {
	data, _ := os.ReadFile(filename)
	locks = make([][]int, 0)
	keys = make([][]int, 0)
	for _, schematicData := range strings.Split(string(data), "\n\n") {
		rows := strings.Split(schematicData, "\n")

		heights := make([]int, 0)
		for col := range len(rows[0]) {
			height := -1
			for _, row := range rows {
				if row[col] == '#' {
					height++
				}
			}
			heights = append(heights, height)
		}

		if rows[0][0] == '#' {
			locks = append(locks, heights)
		} else {
			keys = append(keys, heights)
		}
	}
	return
}
