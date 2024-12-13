package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	data, _ := os.ReadFile("day12/input.txt")

	input := bytes.Split(data, []byte("\n"))
	plots := make([][]*Plot, 0)
	boundary := Plot{Plant: '#'}

	// Create the garden
	for y, row := range input {
		plotRow := make([]*Plot, 0, len(row))
		for x, cell := range row {
			plot := Plot{Plant: cell, X: x, Y: y}
			plotRow = append(plotRow, &plot)
		}
		plots = append(plots, plotRow)
	}

	// Populate neighbors
	for y, row := range plots {
		for x, plot := range row {
			if y > 0 {
				plot.Neighbors = append(plot.Neighbors, plots[y-1][x])
			} else {
				plot.Neighbors = append(plot.Neighbors, &boundary)
			}
			if y < len(plots)-1 {
				plot.Neighbors = append(plot.Neighbors, plots[y+1][x])
			} else {
				plot.Neighbors = append(plot.Neighbors, &boundary)
			}
			if x > 0 {
				plot.Neighbors = append(plot.Neighbors, plots[y][x-1])
			} else {
				plot.Neighbors = append(plot.Neighbors, &boundary)
			}
			if x < len(row)-1 {
				plot.Neighbors = append(plot.Neighbors, plots[y][x+1])
			} else {
				plot.Neighbors = append(plot.Neighbors, &boundary)
			}
		}
	}

	// the region key is the first plot visited in that region
	regions := make(map[*Plot]*Region)

	// Populate regions and their areas and perimeters
	visitedPlots := make(map[*Plot]bool)
	for _, row := range plots {
		for _, plot := range row {
			search(plot, nil, visitedPlots, regions)
		}
	}

	// Count sides by actually counting corners, since it's easier to identify corners
	// In a polygon, corners equal the number of sides
	neighborOffsets := [][]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} // North, east, south, west
	for _, region := range regions {
		for _, plot := range region.Plots {
			for i := 0; i < len(neighborOffsets); i++ {
				var n1, n2, nd *Plot

				// Get offsets of a pair of neighbors plus their corresponding diagonal
				o1 := neighborOffsets[i]
				o2 := neighborOffsets[(i+1)%len(neighborOffsets)]
				od := []int{o1[0] + o2[0], o1[1] + o2[1]}

				// Calculate resulting location of the pair of neighbors and the diagonal plot
				n1x := plot.X + o1[0]
				n1y := plot.Y + o1[1]
				n2x := plot.X + o2[0]
				n2y := plot.Y + o2[1]
				ndx := plot.X + od[0]
				ndy := plot.Y + od[1]

				// Check if any of those locations would be out of bounds, and assign a dummy boundary struct to it
				if n1x < 0 || n1x >= len(plots[0]) || n1y < 0 || n1y >= len(plots) {
					n1 = &boundary
				} else {
					n1 = plots[n1y][n1x]
				}
				if n2x < 0 || n2x >= len(plots[0]) || n2y < 0 || n2y >= len(plots) {
					n2 = &boundary
				} else {
					n2 = plots[n2y][n2x]
				}
				if ndx < 0 || ndx >= len(plots[0]) || ndy < 0 || ndy >= len(plots) {
					nd = &boundary
				} else {
					nd = plots[ndy][ndx]
				}

				// If both neighbors are in the same region and the diagonal is not, this is a corner
				if n1.Plant == plot.Plant && n2.Plant == plot.Plant && nd.Plant != plot.Plant {
					region.Sides++
				}

				// If both neighbors are not part of this region, this is a corner
				if n1.Plant != plot.Plant && n2.Plant != plot.Plant {
					region.Sides++
				}
			}
		}
	}

	// Calculate total fencing price
	price := 0
	for _, region := range regions {
		price += region.Area * region.Perimeter
	}
	fmt.Println("Total fencing price:", price)

	// Calculate total fencing price with discount
	discountedPrice := 0
	for _, region := range regions {
		discountedPrice += region.Area * region.Sides
	}
	fmt.Println("Discounted fencing price:", discountedPrice)
}

func search(plot *Plot, region *Region, visitedPlots map[*Plot]bool, regions map[*Plot]*Region) {
	if visitedPlots[plot] {
		return
	}

	if region == nil {
		region = &Region{Area: 0, Perimeter: 0}
		regions[plot] = region
	}

	visitedPlots[plot] = true
	region.Plots = append(region.Plots, plot)
	region.Area++

	for _, neighbor := range plot.Neighbors {
		if neighbor.Plant == plot.Plant {
			search(neighbor, region, visitedPlots, regions)
		} else {
			region.Perimeter++
		}
	}
}

type Plot struct {
	X, Y      int
	Plant     byte
	Neighbors []*Plot
}

type Region struct {
	Plots     []*Plot
	Area      int
	Perimeter int
	Sides     int
}
