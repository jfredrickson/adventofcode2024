package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	r := NewRacetrack("day20/example.txt")

	baselineTime, _ := r.FastestTime()

	cheatTimes := make(map[int]int, 0)
	for _, wp := range r.GetInnerWallPairs() {
		t1 := r.GetTileAt(wp[0].X, wp[0].Y)
		t2 := r.GetTileAt(wp[1].X, wp[1].Y)
		t1OriginalType := t1.Type
		t2OriginalType := t2.Type
		t1.Type = Track
		t2.Type = Track
		t1.CheatSequence = 1
		t2.CheatSequence = 2

		time, _ := r.FastestTime()
		eliminate := false

		if time == -1 {
			eliminate = true
		}
		if time >= baselineTime {
			eliminate = true
		}

		if !eliminate {
			cheatTimes[baselineTime-time] += 1

		}

		t1.Type = t1OriginalType
		t2.Type = t2OriginalType
		t1.CheatSequence = 0
		t2.CheatSequence = 0
	}

	fmt.Println("Baseline path time:", baselineTime)

	for time, count := range cheatTimes {
		fmt.Println("Cheat times saving", time, "picoseconds:", count)
	}
}

type TileType string

const (
	Track TileType = "."
	Wall  TileType = "#"
)

type Tile struct {
	X, Y          int
	Neighbors     []*Tile
	Type          TileType
	CheatSequence int
}

type Racetrack struct {
	Tiles         []*Tile
	Start, End    *Tile
	Height, Width int
}

func NewRacetrack(filename string) *Racetrack {
	tileData := make([][]string, 0)
	var width int
	common.ProcessFile(filename, func(line string) {
		width = len(line)
		tileData = append(tileData, strings.Split(line, ""))
	})

	r := Racetrack{
		Width:  width,
		Height: len(tileData),
	}

	// Create the individual tiles
	for y, line := range tileData {
		for x, char := range line {
			var tileType TileType
			if char == "#" {
				tileType = Wall
			} else {
				tileType = Track
			}

			t := Tile{
				X:         x,
				Y:         y,
				Neighbors: make([]*Tile, 0),
				Type:      tileType,
			}

			r.Tiles = append(r.Tiles, &t)

			if char == "S" {
				r.Start = &t
			}
			if char == "E" {
				r.End = &t
			}
		}
	}

	// Connect the graph
	for _, t := range r.Tiles {
		if north := r.GetTileAt(t.X, t.Y-1); north != nil {
			t.Neighbors = append(t.Neighbors, north)
		}
		if south := r.GetTileAt(t.X, t.Y+1); south != nil {
			t.Neighbors = append(t.Neighbors, south)
		}
		if west := r.GetTileAt(t.X-1, t.Y); west != nil {
			t.Neighbors = append(t.Neighbors, west)
		}
		if east := r.GetTileAt(t.X+1, t.Y); east != nil {
			t.Neighbors = append(t.Neighbors, east)
		}
	}

	return &r
}

func (r *Racetrack) FastestTime() (int, []*Tile) {
	// Use BFS to find the shortest path
	queue := []*Tile{r.Start}
	distance := make(map[*Tile]int)
	previous := make(map[*Tile]*Tile)
	distance[r.Start] = 0

	for len(queue) > 0 {
		currentTile := queue[0]
		queue = queue[1:]

		for _, neighbor := range currentTile.Neighbors {
			// Don't bother walking into walls
			if neighbor.Type == Wall {
				continue
			}

			// Don't go the wrong way through cheat tiles
			if currentTile.CheatSequence == 2 && neighbor.CheatSequence == 1 {
				continue
			}
			if currentTile.CheatSequence == 0 && neighbor.CheatSequence == 2 {
				continue
			}
			if currentTile.CheatSequence == 1 && neighbor.CheatSequence == 0 {
				continue
			}

			// Ensure we haven't already visited this tile
			if _, exists := distance[neighbor]; !exists {
				distance[neighbor] = distance[currentTile] + 1
				previous[neighbor] = currentTile

				// Check if we've reached the end
				if neighbor == r.End {
					// Recreate the path
					path := make([]*Tile, 0)
					for tile := r.End; tile != nil; tile = previous[tile] {
						path = append([]*Tile{tile}, path...)
					}
					return distance[neighbor], path
				}

				// Add the neighbor to the queue
				queue = append(queue, neighbor)
			}
		}
	}

	// No path found
	return -1, nil
}

func (r *Racetrack) GetInnerWallPairs() [][]*Tile {
	pairs := make([][]*Tile, 0)

	for _, tile := range r.Tiles {
		if tile.Type == Wall {
			// Skip outer wall tiles
			if tile.X == 0 || tile.X == r.Width-1 || tile.Y == 0 || tile.Y == r.Height-1 {
				continue
			}

			for _, neighbor := range tile.Neighbors {
				// Skip outer wall tiles
				if neighbor.X == 0 || neighbor.X == r.Width-1 || neighbor.Y == 0 || neighbor.Y == r.Height-1 {
					continue
				}

				pairs = append(pairs, []*Tile{tile, neighbor})
			}
		}
	}

	return pairs
}

func (r *Racetrack) GetTileAt(x, y int) *Tile {
	for _, tile := range r.Tiles {
		if tile.X == x && tile.Y == y {
			return tile
		}
	}
	return nil
}

func (r *Racetrack) String() string {
	var s string
	for y := 0; y < r.Height; y++ {
		for x := 0; x < r.Width; x++ {
			tile := r.GetTileAt(x, y)
			icon := tile.Type
			if tile == r.Start {
				icon = "S"
			}
			if tile == r.End {
				icon = "E"
			}
			if tile.CheatSequence == 1 {
				icon = "1"
			}
			if tile.CheatSequence == 2 {
				icon = "2"
			}
			s += fmt.Sprintf("%v", icon)
		}
		s += "\n"
	}
	return s
}

func (t *Tile) String() string {
	var s string
	s += fmt.Sprintf("(%d,%d)", t.X, t.Y)
	// for _, neighbor := range t.Neighbors {
	// 	s += fmt.Sprintf("  %v(%d,%d)", neighbor.Type, neighbor.X, neighbor.Y)
	// }
	return s
}
