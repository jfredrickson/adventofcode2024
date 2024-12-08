package main

import (
	"common"
	"fmt"
)

func main() {
	m, maxX, maxY := NewAntennaMap("day08/input.txt")

	fmt.Println("Antinodes:", m.GetAntinodeCount(maxX, maxY))
	fmt.Println("Resonant antinodes:", m.GetResonantAntinodeCount(maxX, maxY))
}

type Vector struct {
	X, Y int
}

type AntennaPair struct {
	A, B Vector
}

type AntennaMap map[string][]Vector

// Add another vector to this one, returning a new Vector
func (v Vector) Add(other *Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y}
}

// Make an antenna map from a file, also returning the dimensions of the map (maxX and maxY)
func NewAntennaMap(filename string) (AntennaMap, int, int) {
	m := AntennaMap{}
	row := 0
	var lineLength int

	common.ProcessFile(filename, func(line string) {
		lineLength = len(line)
		for col, char := range line {
			if char != '.' {
				m[string(char)] = append(m[string(char)], Vector{col, row})
			}
		}
		row++
	})

	return m, lineLength, row
}

// Calculate the location of all antinodes given map dimensions (maxX and maxY)
func (m AntennaMap) GetAntinodeCount(maxX, maxY int) int {
	antinodeSet := make(map[Vector]bool)

	for antennaType := range m {
		pairs := m.getDirectionalPairs(antennaType)
		for _, pair := range pairs {
			potentialAntinode := pair.B.Add(pair.Distance())
			// Ensure potential antinode is within bounds
			if potentialAntinode.X >= 0 && potentialAntinode.X < maxX && potentialAntinode.Y >= 0 && potentialAntinode.Y < maxY {
				antinodeSet[potentialAntinode] = true
			}
		}
	}

	return len(antinodeSet)
}

// Calculate the location of all resonant antinodes given map dimensions (maxX and maxY)
func (m AntennaMap) GetResonantAntinodeCount(maxX, maxY int) int {
	antinodeSet := make(map[Vector]bool)

	for antennaType := range m {
		pairs := m.getDirectionalPairs(antennaType)
		for _, pair := range pairs {
			potentialAntinode := pair.B.Add(pair.Distance())
			// Ensure potential antinode is within bounds
			for potentialAntinode.X >= 0 && potentialAntinode.X < maxX && potentialAntinode.Y >= 0 && potentialAntinode.Y < maxY {
				antinodeSet[potentialAntinode] = true
				potentialAntinode = potentialAntinode.Add(pair.Distance())
			}
			// Antenna pair points themselves are antinodes
			antinodeSet[pair.A] = true
			antinodeSet[pair.B] = true
		}
	}

	return len(antinodeSet)
}

// Get all directional pairs for a given antenna type
func (m AntennaMap) getDirectionalPairs(antennaType string) []AntennaPair {
	pairs := []AntennaPair{}

	// First pass, get simple unordered pairs
	for i := 0; i < len(m[antennaType]); i++ {
		for j := i + 1; j < len(m[antennaType]); j++ {
			pairs = append(pairs, AntennaPair{m[antennaType][i], m[antennaType][j]})
		}
	}

	// Second pass, get reversed pairs
	for _, pair := range pairs {
		reversedPair := AntennaPair{pair.B, pair.A}
		pairs = append(pairs, reversedPair)
	}

	return pairs
}

func (p *AntennaPair) Distance() *Vector {
	return &Vector{
		X: p.B.X - p.A.X,
		Y: p.B.Y - p.A.Y,
	}
}
