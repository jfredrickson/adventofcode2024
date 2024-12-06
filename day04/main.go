package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	puzzle := Puzzle{}

	common.ProcessFile("day04/input.txt", func(line string) {
		puzzle.Grid = append(puzzle.Grid, strings.Split(line, ""))
	})

	fmt.Println("Occurrences:", puzzle.Count("XMAS"))
	fmt.Println("Occurrences X-MAS:", puzzle.CountX())
}

type Puzzle struct {
	Grid [][]string
}

func (p *Puzzle) Count(word string) int {
	occurrences := 0

	for ri, row := range p.Grid {
		for ci := range row {
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if p.Check(word, ci, ri, i, j) {
						occurrences++
					}
				}
			}
		}
	}

	return occurrences
}

func (p *Puzzle) Check(word string, startX, startY, dirX, dirY int) bool {
	if dirX == 0 && dirY == 0 {
		return false
	}

	found := true

	for i := 0; i < len(word); i++ {
		nextX := startX + (i * dirX)
		nextY := startY + (i * dirY)
		if nextX < 0 || nextY < 0 || nextX >= len(p.Grid[0]) || nextY >= len(p.Grid) {
			return false
		}
		if p.Grid[startY+(i*dirY)][startX+(i*dirX)] != string(word[i]) {
			found = false
		}
	}

	return found
}

func (p *Puzzle) CountX() int {
	occurrences := 0

	for ri, row := range p.Grid {
		for ci := range row {
			if p.Grid[ri][ci] == "A" {
				if p.CheckX(ci, ri) {
					occurrences++
				}
			}
		}
	}

	return occurrences
}

type OffsetPair struct {
	X1, Y1, X2, Y2 int
}

func (o *OffsetPair) Inverted() *OffsetPair {
	return &OffsetPair{
		X1: -o.X1,
		Y1: -o.Y1,
		X2: -o.X2,
		Y2: -o.Y2,
	}
}

func (p *Puzzle) CheckX(startX, startY int) bool {
	// All four possible locations for the Ms
	mPairs := []OffsetPair{
		{-1, -1, -1, 1},
		{-1, -1, 1, -1},
		{1, 1, -1, 1},
		{1, 1, 1, -1},
	}

	for _, mPair := range mPairs {
		sPair := mPair.Inverted()

		// Bounds checks
		if startX+mPair.X1 < 0 || startX+mPair.X2 < 0 || startX+mPair.X1 >= len(p.Grid[0]) || startX+mPair.X2 >= len(p.Grid[0]) {
			continue
		}
		if startY+mPair.Y1 < 0 || startY+mPair.Y2 < 0 || startY+mPair.Y1 >= len(p.Grid) || startY+mPair.Y2 >= len(p.Grid) {
			continue
		}
		if startX+sPair.X1 < 0 || startX+sPair.X2 < 0 || startX+sPair.X1 >= len(p.Grid[0]) || startX+sPair.X2 >= len(p.Grid[0]) {
			continue
		}
		if startY+sPair.Y1 < 0 || startY+sPair.Y2 < 0 || startY+sPair.Y1 >= len(p.Grid) || startY+sPair.Y2 >= len(p.Grid) {
			continue
		}

		// Check if the Ms are there
		if p.Grid[startY+mPair.Y1][startX+mPair.X1] != "M" || p.Grid[startY+mPair.Y2][startX+mPair.X2] != "M" {
			continue
		}

		// Check if the Ss are there
		if p.Grid[startY+sPair.Y1][startX+sPair.X1] != "S" || p.Grid[startY+sPair.Y2][startX+sPair.X2] != "S" {
			continue
		}

		// If loop hasn't exited at this point, the Ms and Ss are in the right places
		return true
	}

	return false
}
