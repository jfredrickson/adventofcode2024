package main

import (
	"common"
	"fmt"
	"slices"
	"strings"
)

func main() {
	puzzle := NewPuzzle("day06/input.txt")

	fmt.Println("Visited squares:", puzzle.Process())
	fmt.Println("Possible obstructions:", puzzle.CountPossibleObstructions())
}

type Point struct {
	X, Y int
}

type Direction struct {
	Symbol string
	Offset Point
}

var (
	North = Direction{"^", Point{0, -1}}
	East  = Direction{">", Point{1, 0}}
	South = Direction{"v", Point{0, 1}}
	West  = Direction{"<", Point{-1, 0}}
)

type Cell struct {
	Content                 string
	Visited                 bool
	PreviousGuardDirections []Direction
}

func (c *Cell) Visit(guardDirection Direction) {
	c.Visited = true
	c.Content = guardDirection.Symbol
	c.PreviousGuardDirections = append(c.PreviousGuardDirections, guardDirection)
}

type Puzzle struct {
	Grid           [][]Cell
	GuardLocation  Point
	GuardDirection Direction
	directionMap   map[string]Direction
}

func NewPuzzle(filename string) *Puzzle {
	p := &Puzzle{
		directionMap: map[string]Direction{
			"^": North,
			">": East,
			"v": South,
			"<": West,
		},
	}

	// build the grid from the input file
	y := 0
	common.ProcessFile(filename, func(line string) {
		cells := make([]Cell, len(line))
		for x, char := range strings.Split(line, "") {
			cells[x] = Cell{Content: char}
			// if the cell is the guard, record its location and direction
			if direction, found := p.directionMap[char]; found {
				p.GuardLocation = Point{x, y}
				p.GuardDirection = direction
			}
		}
		p.Grid = append(p.Grid, cells)
		y++
	})

	return p
}

func (p *Puzzle) CountPossibleObstructions() int {
	count := 0
	// iterate over the grid, finding each open cell
	for y, row := range p.Grid {
		for x, cell := range row {
			if cell.Content == "." {
				testPuzzle := p.Clone()
				testPuzzle.Grid[y][x].Content = "#"
				if testPuzzle.Process() == -1 {
					count++
				}
			}
		}
	}
	return count
}

func (p *Puzzle) Clone() *Puzzle {
	pCopy := *p
	pCopy.Grid = make([][]Cell, len(p.Grid))
	for i, row := range p.Grid {
		pCopy.Grid[i] = make([]Cell, len(row))
		copy(pCopy.Grid[i], row)
	}

	return &pCopy
}

func (p *Puzzle) Process() int {
	puzzle := p.Clone()

	for {
		// mark the current space as visited
		currentCell := &puzzle.Grid[puzzle.GuardLocation.Y][puzzle.GuardLocation.X]
		currentCell.Visit(puzzle.GuardDirection)

		// get the guard's next move
		nextX := puzzle.GuardLocation.X + puzzle.GuardDirection.Offset.X
		nextY := puzzle.GuardLocation.Y + puzzle.GuardDirection.Offset.Y

		// check if that move leads off the edge of the grid
		if nextX < 0 || nextX >= len(puzzle.Grid[0]) || nextY < 0 || nextY >= len(puzzle.Grid) {
			// if so, all done
			return puzzle.countVisitedSquares()
		}

		// while there is an obstruction in that direction, rotate the guard and recalculate the next move
		for puzzle.Grid[nextY][nextX].Content == "#" {
			puzzle.rotateGuard()
			nextX = puzzle.GuardLocation.X + puzzle.GuardDirection.Offset.X
			nextY = puzzle.GuardLocation.Y + puzzle.GuardDirection.Offset.Y
		}

		// if the next cell has the guard's direction previously recorded, this is an infinite loop
		if slices.Contains(puzzle.Grid[nextY][nextX].PreviousGuardDirections, puzzle.GuardDirection) {
			return -1
		}

		// move guard forward one space
		puzzle.GuardLocation.X = nextX
		puzzle.GuardLocation.Y = nextY
	}
}

func (p *Puzzle) rotateGuard() {
	directionSequence := []Direction{North, East, South, West}
	currentDirectionIndex := slices.Index(directionSequence, p.GuardDirection)
	nextDirection := directionSequence[(currentDirectionIndex+1)%len(directionSequence)]
	p.GuardDirection = nextDirection
}

func (p *Puzzle) countVisitedSquares() int {
	visited := 0
	for _, row := range p.Grid {
		for _, cell := range row {
			if cell.Visited {
				visited++
			}
		}
	}
	return visited
}

func (p *Puzzle) Print() {
	fmt.Println()
	for _, row := range p.Grid {
		for _, cell := range row {
			fmt.Print(cell.Content)
		}
		fmt.Println()
	}
}
