package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	mapData := make([]string, 0)
	moves := make([]Direction, 0)
	processingMap := true
	common.ProcessFile("day15/input.txt", func(line string) {
		if line == "" {
			processingMap = false
			return
		}

		if processingMap {
			mapData = append(mapData, line)
		} else {
			for _, m := range strings.Split(line, "") {
				moves = append(moves, Direction(m))
			}
		}
	})

	w := Warehouse{}
	w.Load(mapData)

	for _, direction := range moves {
		w.Move(direction)
	}

	fmt.Println("Total GPS coordinates:", w.TotalCoordinates())
}

type Direction string

const (
	Up    Direction = "^"
	Down  Direction = "v"
	Left  Direction = "<"
	Right Direction = ">"
)

type Cell struct {
	X, Y      int
	Neighbors map[Direction]*Cell
	Contents  string
}

type Warehouse struct {
	Cells         []*Cell
	Width, Height int
}

func (w *Warehouse) Load(lines []string) {
	// Start by creating individual cells
	for y, line := range lines {
		w.Width = len(line)
		w.Height++
		for x, contents := range strings.Split(line, "") {
			c := Cell{
				X:         x,
				Y:         y,
				Contents:  contents,
				Neighbors: make(map[Direction]*Cell, 0),
			}
			w.Cells = append(w.Cells, &c)
		}
	}

	// Build the graph
	for _, cell := range w.Cells {
		if cell.Y > 0 {
			cell.Neighbors[Up] = w.GetCellAt(cell.X, cell.Y-1)
		}
		if cell.Y < w.Height-1 {
			cell.Neighbors[Down] = w.GetCellAt(cell.X, cell.Y+1)
		}
		if cell.X > 0 {
			cell.Neighbors[Left] = w.GetCellAt(cell.X-1, cell.Y)
		}
		if cell.X < w.Width-1 {
			cell.Neighbors[Right] = w.GetCellAt(cell.X+1, cell.Y)
		}
	}
}

func (w *Warehouse) Move(direction Direction) {
	recursiveMove(nil, w.GetRobotCell(), direction)
}

func recursiveMove(previous *Cell, current *Cell, direction Direction) (moved bool) {
	if current.Contents == "@" || current.Contents == "O" {
		moved := recursiveMove(current, current.Neighbors[direction], direction)
		if moved && previous != nil {
			current.Contents = previous.Contents
			previous.Contents = "."
			return true
		}
		return moved
	}

	if current.Contents == "#" {
		return false
	}

	if current.Contents == "." {
		current.Contents = previous.Contents
		previous.Contents = "."
		return true
	}

	return false
}

func (w *Warehouse) Print() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			fmt.Print(w.GetCellAt(x, y).Contents)
		}
		fmt.Println()
	}
}

func (w *Warehouse) GetCellAt(x, y int) *Cell {
	for _, cell := range w.Cells {
		if cell.X == x && cell.Y == y {
			return cell
		}
	}
	return nil
}

func (w *Warehouse) GetRobotCell() *Cell {
	for _, cell := range w.Cells {
		if cell.Contents == "@" {
			return cell
		}
	}
	return nil
}

func (w *Warehouse) TotalCoordinates() int {
	total := 0
	for _, cell := range w.Cells {
		if cell.Contents == "O" {
			total += cell.X + 100*cell.Y
		}
	}
	return total
}
