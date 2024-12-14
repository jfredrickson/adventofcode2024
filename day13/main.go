package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	machines := loadMachineDefinitions("day13/input.txt")

	cost := 0
	for _, machine := range machines {
		a, b := findButtonPresses(machine)
		cost += (a * 3) + (b * 1)
	}
	fmt.Println("Cost of tokens:", cost)
}

type Point struct {
	X, Y int
}

type Machine struct {
	A     Point
	B     Point
	Prize Point
}

func findButtonPresses(machine Machine) (a, b int) {
	// Figure out the maximum number of times we need to iterate to find the best combination
	iterations := max(machine.Prize.X/machine.A.X, machine.Prize.X/machine.B.X)
	// Start with the maximum number of B presses
	for i := iterations; i >= 0; i-- {
		remaining := machine.Prize.X - (i * machine.B.X)
		// Check if the remainder of necessary presses is attainable by presses of A
		if remaining%machine.A.X == 0 {
			// If so, we found a potential multiplier based on the X axis
			a := remaining / machine.A.X
			b := i
			// verify the multiplier is valid for the Y axis as well
			if a*machine.A.Y+b*machine.B.Y == machine.Prize.Y {
				return a, b
			}
		}
	}
	return 0, 0
}

func loadMachineDefinitions(filename string) []Machine {
	data, _ := os.ReadFile(filename)
	machines := make([]Machine, 0)
	for _, machineData := range strings.Split(string(data), "\n\n") {
		machines = append(machines, parseMachineData(machineData))
	}
	return machines
}

func parseMachineData(machineData string) Machine {
	machine := Machine{}
	propertyData := strings.Split(machineData, "\n")
	machine.A = parseProperty(propertyData[0])
	machine.B = parseProperty(propertyData[1])
	machine.Prize = parseProperty(propertyData[2])
	return machine
}

func parseProperty(propertyData string) Point {
	r := regexp.MustCompile(`.+: X[+=](\d+), Y[+=](\d+)`)
	p := r.FindStringSubmatch(propertyData)
	x, _ := strconv.Atoi(p[1])
	y, _ := strconv.Atoi(p[2])
	return Point{X: x, Y: y}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
