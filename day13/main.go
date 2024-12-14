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

	correctedCost := 0
	for _, machine := range machines {
		machine.Prize.X += 10000000000000
		machine.Prize.Y += 10000000000000
		a, b := findButtonPresses(machine)
		correctedCost += (a * 3) + (b * 1)
	}
	fmt.Println("Corrected cost of tokens:", correctedCost)
}

func findButtonPresses(m Machine) (a, b int) {
	// So this seems to be a system of equations
	// aPresses * mAX + bPresses * mBX = mPX
	// aPresses * mAY + bPresses * mBY = mPY

	// Need to solve for either aPresses or bPresses
	// aPresses = (mPX - bPresses * mBX) / mAX
	// aPresses = (mPY - bPresses * mBY) / mAY

	// Now can eliminate aPresses because the two right hand sides are equal
	// (mPX - bPresses * mBX) / mAX = (mPY - bPresses * mBY) / mAY
	// (mPX - bPresses * mBX) * mAY = (mPY - bPresses * mBY) * mAX

	// Now isolate bPresses on one side to get the potential number of B button presses
	// mPX * mAY - bPresses * mBX * mAY = mPY * mAX - bPresses * mBY * mAX
	// mPX * mAY - mPY * mAX = bPresses * mBX * MAY - bPresses * mBY * mAX
	// mPX * mAY - mPY * mAX = bPresses (mBX * mAY - mBY * mAX)
	// bPresses = (mPX * mAY - mPY * mAX) / (mBX * mAY - mBY * mAX)

	// Use that equation to solve for B presses, then solve for A presses in turn
	b = (m.Prize.X*m.A.Y - m.Prize.Y*m.A.X) / (m.B.X*m.A.Y - m.B.Y*m.A.X)
	a = (m.Prize.X - b*m.B.X) / m.A.X

	// Verify that the combination of A and B presses matches the prize coordinates
	if a*m.A.X+b*m.B.X == m.Prize.X && a*m.A.Y+b*m.B.Y == m.Prize.Y {
		return a, b
	}

	return 0, 0
}

type Point struct {
	X, Y int
}

type Machine struct {
	A     Point
	B     Point
	Prize Point
}

// Approach via iteration, which worked fine for part 1 but not part 2
// func findButtonPresses(machine Machine) (a, b int) {
// 	// Figure out the maximum number of times we need to iterate to find the best combination
// 	iterations := max(machine.Prize.X/machine.A.X, machine.Prize.X/machine.B.X)
// 	// Start with the maximum number of B presses
// 	for i := iterations; i >= 0; i-- {
// 		remaining := machine.Prize.X - (i * machine.B.X)
// 		// Check if the remainder of necessary presses is attainable by presses of A
// 		if remaining%machine.A.X == 0 {
// 			// If so, we found a potential multiplier based on the X axis
// 			a := remaining / machine.A.X
// 			b := i
// 			// verify the multiplier is valid for the Y axis as well
// 			if a*machine.A.Y+b*machine.B.Y == machine.Prize.Y {
// 				return a, b
// 			}
// 		}
// 	}
// 	return 0, 0
// }

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	} else {
// 		return b
// 	}
// }

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
