package main

import (
	"common"
	"fmt"
	"math"
	"strings"
)

func main() {
	m := NewMaze("day16/input.txt")

	fmt.Println("Best path score:", m.FindBestPath())
	fmt.Println(len(m.FindBestPaths()), "best paths")
}

type Direction string

const (
	North Direction = "N"
	East  Direction = "E"
	South Direction = "S"
	West  Direction = "W"
)

type Tile struct {
	X, Y      int
	Neighbors map[Direction]*Tile
}

type Maze struct {
	Tiles      []*Tile
	Start, End *Tile
}

func NewMaze(filename string) *Maze {
	tileData := make([][]string, 0)
	common.ProcessFile(filename, func(line string) {
		tileData = append(tileData, strings.Split(line, ""))
	})

	m := Maze{}

	// Create the individual tiles
	for y, line := range tileData {
		for x, char := range line {
			if char == "#" {
				continue
			}
			t := Tile{
				X:         x,
				Y:         y,
				Neighbors: make(map[Direction]*Tile, 0),
			}
			m.Tiles = append(m.Tiles, &t)
			if char == "S" {
				m.Start = &t
			}
			if char == "E" {
				m.End = &t
			}
		}
	}

	// Connect the graph
	for _, t := range m.Tiles {
		if north := m.GetTileAt(t.X, t.Y-1); north != nil {
			t.Neighbors[North] = north
		}
		if south := m.GetTileAt(t.X, t.Y+1); south != nil {
			t.Neighbors[South] = south
		}
		if west := m.GetTileAt(t.X-1, t.Y); west != nil {
			t.Neighbors[West] = west
		}
		if east := m.GetTileAt(t.X+1, t.Y); east != nil {
			t.Neighbors[East] = east
		}
	}

	return &m
}

func (m *Maze) GetTileAt(x, y int) *Tile {
	for _, tile := range m.Tiles {
		if tile.X == x && tile.Y == y {
			return tile
		}
	}
	return nil
}

func (m *Maze) FindBestPath() int {
	type State struct {
		Tile        *Tile
		Orientation Direction
		Score       int
	}

	// Priority queue for Dijkstra's algorithm
	pq := make([]State, 0)
	pq = append(pq, State{Tile: m.Start, Orientation: East, Score: 0})

	// Visited states: map tile -> orientation -> lowest score
	visited := make(map[*Tile]map[Direction]int)

	// Initialize visited map, setting all costs to infinity (or at least MaxInt)
	for _, tile := range m.Tiles {
		visited[tile] = make(map[Direction]int)
		for _, direction := range []Direction{North, East, South, West} {
			visited[tile][direction] = math.MaxInt
		}
	}

	// Use Dijkstra's algorithm
	lowestScore := math.MaxInt
	for len(pq) > 0 {
		// Get the lowest score state out of the priority queue
		current := pq[0]
		pq = pq[1:]

		// If we reached the end, update the lowest score
		if current.Tile == m.End {
			if current.Score < lowestScore {
				lowestScore = current.Score
			}
			continue
		}

		// If we've already found a path to this state with a lower score, skip this
		if current.Score >= visited[current.Tile][current.Orientation] {
			continue
		}

		// Mark this state as visited
		visited[current.Tile][current.Orientation] = current.Score

		// Check all neighbors
		for direction, neighbor := range current.Tile.Neighbors {
			newScore := current.Score + getCost(current.Tile, neighbor, current.Orientation)
			newState := State{
				Tile:        neighbor,
				Orientation: direction,
				Score:       newScore,
			}
			pq = append(pq, newState)
		}
	}

	return lowestScore
}

func (m *Maze) FindBestPaths() [][]*Tile {
	type State struct {
		Tile        *Tile
		Orientation Direction
		Score       int
		Path        []*Tile
	}

	// Priority queue for Dijkstra's algorithm
	pq := make([]State, 0)
	pq = append(pq, State{Tile: m.Start, Orientation: East, Score: 0, Path: make([]*Tile, 0)})

	// Visited states: map tile -> orientation -> lowest score
	visited := make(map[*Tile]map[Direction]int)

	// Initialize visited map, setting all costs to infinity (or at least to MaxInt)
	for _, tile := range m.Tiles {
		visited[tile] = make(map[Direction]int)
		for _, direction := range []Direction{North, East, South, West} {
			visited[tile][direction] = math.MaxInt
		}
	}

	// Track best paths and lowest score
	bestPaths := [][]*Tile{}
	lowestScore := math.MaxInt

	// Use Dijkstra's algorithm
	for len(pq) > 0 {
		// Get the lowest score state out of the priority queue
		current := pq[0]
		pq = pq[1:]

		// If we reached the end, update the lowest score
		if current.Tile == m.End {
			if current.Score < lowestScore {
				// If we found a better path score, reset the list of best paths and add this one
				lowestScore = current.Score
				bestPaths = [][]*Tile{current.Path}
			} else if current.Score == lowestScore {
				// This path is equal to the current best path, add it
				bestPaths = append(bestPaths, current.Path)
			}
			continue
		}

		// If we've already found a path to this state with a lower score, skip this
		if current.Score >= visited[current.Tile][current.Orientation] {
			continue
		}

		// Mark this state as visited
		visited[current.Tile][current.Orientation] = current.Score

		// Check all neighbors
		for direction, neighbor := range current.Tile.Neighbors {
			newScore := current.Score + getCost(current.Tile, neighbor, current.Orientation)
			newPath := append([]*Tile{}, current.Path...)
			newPath = append(newPath, neighbor)
			newState := State{
				Tile:        neighbor,
				Orientation: direction,
				Score:       newScore,
				Path:        newPath,
			}
			pq = append(pq, newState)
		}
	}

	return bestPaths
}

func (m *Maze) String() string {
	var s string
	for _, t := range m.Tiles {
		s += fmt.Sprintf("%d,%d ->", t.X, t.Y)
		for direction, neighbor := range t.Neighbors {
			s += fmt.Sprintf(" %s:%d,%d", direction, neighbor.X, neighbor.Y)
		}
		s += "\n"
	}
	return s
}

// Get the cost of moving from one tile to another, given an orientation
func getCost(from, to *Tile, orientation Direction) int {
	cost := 1 // Cost is at least 1 for all moves
	for direction, neighbor := range from.Neighbors {
		if neighbor == to && direction != orientation {
			cost += 1000 // We have to turn to go to this neighbor, which costs 1000
		}
	}
	return cost
}
