package main

import (
	"common"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	topoMap := NewTopoMap("day10/input.txt")
	scores, ratings := topoMap.Totals()

	fmt.Println("Trailhead score total:", scores)
	fmt.Println("Trailhead rating total:", ratings)
}

type Step struct {
	X, Y   int
	Height int
}

type TopoMap map[*Step][]*Step

func NewTopoMap(filename string) TopoMap {
	t := TopoMap{}

	var input [][]string
	common.ProcessFile(filename, func(line string) {
		input = append(input, strings.Split(line, ""))
	})

	for y := 0; y < len(input); y++ {
		for x := 0; x < len(input[y]); x++ {
			height, _ := strconv.Atoi(input[y][x])
			step := &Step{
				X:      x,
				Y:      y,
				Height: height,
			}
			t[step] = []*Step{}
		}
	}

	for step := range t {
		// Calculate all possible neighbors relative to this step
		neighbors := [][]int{
			{step.X, step.Y - 1}, {step.X - 1, step.Y}, {step.X + 1, step.Y}, {step.X, step.Y + 1},
		}

		for _, neighbor := range neighbors {
			nx, ny := neighbor[0], neighbor[1]
			// Check if in bounds
			if ny < 0 || nx < 0 || ny >= len(input) || nx >= len(input[ny]) {
				continue
			}

			// Get the neighbor from the map based on coordinates
			var neighbor *Step
			for s := range t {
				if s.X == nx && s.Y == ny {
					neighbor = s
					break
				}
			}

			// If the neighbor is a valid next step, add it to the neighbors list
			if neighbor.Height == step.Height+1 {
				t[step] = append(t[step], neighbor)
			}
		}
	}

	return t
}

func (t *TopoMap) AllTrailheads() []*Step {
	trailheads := make([]*Step, 0)
	for step := range *t {
		if step.Height == 0 {
			trailheads = append(trailheads, step)
		}
	}
	return trailheads
}

func (t *TopoMap) AllEnds() []*Step {
	ends := make([]*Step, 0)
	for step := range *t {
		if step.Height == 9 {
			ends = append(ends, step)
		}
	}
	return ends
}

func (t *TopoMap) Totals() (int, int) {
	score := 0
	rating := 0

	for _, trailhead := range t.AllTrailheads() {
		for _, end := range t.AllEnds() {
			allPaths := make([][]*Step, 0)
			t.findAllPaths(trailhead, end, map[*Step]bool{}, []*Step{}, &allPaths)
			if len(allPaths) > 0 {
				score++
			}
			rating += len(allPaths)
		}
	}

	return score, rating
}

// Return all paths from start to end
func (t *TopoMap) findAllPaths(start, end *Step, visited map[*Step]bool, path []*Step, allPaths *[][]*Step) {
	if visited[start] {
		return
	}

	visited[start] = true
	path = append(path, start)

	if start == end {
		*allPaths = append(*allPaths, append([]*Step{}, path...))
	} else {
		for _, neighbor := range (*t)[start] {
			t.findAllPaths(neighbor, end, visited, path, allPaths)
		}
	}

	visited[start] = false
}
