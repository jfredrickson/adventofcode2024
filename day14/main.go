package main

import (
	"common"
	"fmt"
	"regexp"
)

func main() {
	bathroom := Bathroom{}
	bathroom.Load("day14/input.txt", 101, 103)

	safetyFactorSeconds := 100

	safetyFactor := -1
	easterEgg := -1
	tick := 0
	for safetyFactor < 0 || easterEgg < 0 {
		// Look for safety factor at the desired number of seconds
		if tick == safetyFactorSeconds {
			safetyFactor = bathroom.CalculateSafetyFactor()
		}

		// Look for easter egg by looking for a state in which all robots are in unique locations
		uniqueLocations := true
		for y := 0; y < bathroom.Height; y++ {
			for x := 0; x < bathroom.Width; x++ {
				if len(bathroom.RobotsAt(x, y)) > 1 {
					uniqueLocations = false
				}
			}
		}

		if uniqueLocations {
			easterEgg = tick
			bathroom.Print()
		}

		bathroom.Tick(1)
		tick++
	}

	fmt.Println("Safety factor:", safetyFactor)
	fmt.Println("Easter egg found at:", easterEgg)
}

type Point struct {
	X, Y int
}

type Robot struct {
	Position Point
	Velocity Point
}

type Bathroom struct {
	Robots        []*Robot
	Width, Height int
}

func (b *Bathroom) Load(filename string, width, height int) {
	b.Width = width
	b.Height = height
	common.ProcessFile(filename, func(line string) {
		r := regexp.MustCompile(`p=(-*\d+),(-*\d+) v=(-*\d+),(-*\d+)`)
		match := r.FindStringSubmatch(line)
		p := Point{X: common.Atoi(match[1]), Y: common.Atoi(match[2])}
		v := Point{X: common.Atoi(match[3]), Y: common.Atoi(match[4])}
		b.Robots = append(b.Robots, &Robot{Position: p, Velocity: v})
	})
}

func (b *Bathroom) RobotsAt(x, y int) []*Robot {
	robots := make([]*Robot, 0)
	for _, robot := range b.Robots {
		if robot.Position.X == x && robot.Position.Y == y {
			robots = append(robots, robot)
		}
	}
	return robots
}

func (b *Bathroom) Tick(seconds int) {
	for _, robot := range b.Robots {
		robot.Position.X = (robot.Position.X + (robot.Velocity.X * seconds)) % b.Width
		robot.Position.Y = (robot.Position.Y + (robot.Velocity.Y * seconds)) % b.Height
		if robot.Position.X < 0 {
			robot.Position.X += b.Width
		}
		if robot.Position.Y < 0 {
			robot.Position.Y += b.Height
		}
	}
}

func (b *Bathroom) CalculateSafetyFactor() int {
	wd := b.Width / 2
	hd := b.Height / 2

	var nw, ne, sw, se int
	for y := 0; y < hd; y++ {
		for x := 0; x < wd; x++ {
			nw += len(b.RobotsAt(x, y))
		}
		for x := wd + 1; x < b.Width; x++ {
			ne += len(b.RobotsAt(x, y))
		}
	}
	for y := hd + 1; y < b.Height; y++ {
		for x := 0; x < wd; x++ {
			sw += len(b.RobotsAt(x, y))
		}
		for x := wd + 1; x < b.Width; x++ {
			se += len(b.RobotsAt(x, y))
		}
	}

	return nw * ne * sw * se
}

func (b *Bathroom) Print() {
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			robotCount := len(b.RobotsAt(x, y))
			if robotCount > 0 {
				fmt.Print(robotCount)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
