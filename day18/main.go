package main

import (
	"common"
	"fmt"
)

func main() {
	ms := NewMemorySpace(71) // example=7 input=71
	corruption := loadCorruption("day18/input.txt")

	// Apply first 1024 corrupted memory positions
	for i := range 1024 {
		fmt.Println("applying corruption", i, corruption[i])
		ms.Positions[corruption[i]].Safe = false
	}

	fmt.Println(ms)

	startKey, endKey := keyFor(0, 0), keyFor(ms.Size-1, ms.Size-1)
	start := ms.Positions[startKey]
	end := ms.Positions[endKey]
	path := ms.FindPath(start, end)
	steps := len(path) - 1 // -1 because we don't count the start

	fmt.Println("Steps to end:", steps)

	for i := 1024; i < len(corruption); i++ {
		ms.Positions[corruption[i]].Safe = false
		path = ms.FindPath(start, end)
		if len(path) == 0 {
			fmt.Println("Cut off at:", corruption[i])
			break
		}
	}
}

type Position struct {
	X, Y      int
	Safe      bool
	Neighbors []*Position
}

func (p *Position) String() string {
	var safety, neighbors string

	if p.Safe {
		safety = "safe"
	} else {
		safety = "unsafe"
	}

	for _, n := range p.Neighbors {
		neighbors += fmt.Sprintf("(%d,%d) ", n.X, n.Y)
	}
	return fmt.Sprintf("(%d,%d) %s -> %s", p.X, p.Y, safety, neighbors)
}

type MemorySpace struct {
	Positions map[string]*Position
	Size      int
}

func NewMemorySpace(size int) *MemorySpace {
	ms := MemorySpace{
		Positions: map[string]*Position{},
		Size:      size,
	}

	// Populate positions
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			ms.Positions[keyFor(x, y)] = &Position{X: x, Y: y, Safe: true}
		}
	}

	// Build the graph
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kCurrent := keyFor(x, y)
			// Up
			if y > 0 {
				kUp := keyFor(x, y-1)
				ms.Positions[kCurrent].Neighbors = append(ms.Positions[kCurrent].Neighbors, ms.Positions[kUp])
			}
			// Down
			if y < size-1 {
				kDown := keyFor(x, y+1)
				ms.Positions[kCurrent].Neighbors = append(ms.Positions[kCurrent].Neighbors, ms.Positions[kDown])
			}
			// Left
			if x > 0 {
				kLeft := keyFor(x-1, y)
				ms.Positions[kCurrent].Neighbors = append(ms.Positions[kCurrent].Neighbors, ms.Positions[kLeft])
			}
			// Right
			if x < size-1 {
				kRight := keyFor(x+1, y)
				ms.Positions[kCurrent].Neighbors = append(ms.Positions[kCurrent].Neighbors, ms.Positions[kRight])
			}
		}
	}

	return &ms
}

func (ms *MemorySpace) FindPath(start, end *Position) []*Position {
	// BFS
	queue := [][]*Position{}
	visited := map[*Position]bool{}

	// Initialize the queue and visited positions with the starting position
	queue = append(queue, []*Position{start})
	visited[start] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		currentPosition := path[len(path)-1]

		// Check if we're already at the end
		if currentPosition == end {
			return path
		}

		// Check all neighbors
		for _, neighbor := range currentPosition.Neighbors {
			// If we haven't visited this neighbor and it's a valid one (aka safe), add it to the current path
			if neighbor.Safe && !visited[neighbor] {
				visited[neighbor] = true
				newPath := append([]*Position{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	return nil
}

func (ms *MemorySpace) String() string {
	s := ""
	for y := 0; y < ms.Size; y++ {
		for x := 0; x < ms.Size; x++ {
			if ms.Positions[keyFor(x, y)].Safe {
				s += "."
			} else {
				s += "#"
			}
		}
		s += "\n"
	}
	return s
}

func keyFor(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func loadCorruption(filename string) (positions []string) {
	common.ProcessFile(filename, func(line string) {
		positions = append(positions, line)
	})
	return positions
}
