package main

import (
	"common"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	c := NewCalibration("day07/input.txt")
	fmt.Println("Total calibration sum:", c.Total([]Operator{add, multiply}))
	fmt.Println("New total calibration sum:", c.Total([]Operator{add, multiply, concat}))

}

type Equation struct {
	TestValue int
	Numbers   []int
	Solvable  bool
}

type Operator func(int, int) int

type Calibration struct {
	Equations []Equation
	Operators []Operator
}

func NewCalibration(filename string) Calibration {
	c := Calibration{
		Equations: []Equation{},
	}

	common.ProcessFile(filename, func(line string) {
		components := strings.Split(line, ": ")

		testValue, err := strconv.Atoi(components[0])
		if err != nil {
			panic(err)
		}

		numbers := common.ToInts(strings.Split(components[1], " "))

		c.Equations = append(c.Equations, Equation{
			TestValue: testValue,
			Numbers:   numbers,
			Solvable:  false,
		})
	})

	return c
}

func (c *Calibration) Total(operators []Operator) int {
	total := 0

	for _, e := range c.Equations {
		positions := len(e.Numbers) - 1
		total += trySequences(positions, e, operators, []Operator{})
	}

	return total
}

func trySequences(positions int, equation Equation, availableOps []Operator, sequence []Operator) int {
	if len(sequence) == positions {
		accumulator := equation.Numbers[0]
		for i, op := range sequence {
			accumulator = op(accumulator, equation.Numbers[i+1])
		}
		if accumulator == equation.TestValue {
			return accumulator
		} else {
			return 0
		}
	}

	for _, op := range availableOps {
		result := trySequences(positions, equation, availableOps, append(sequence, op))
		if result != 0 {
			return result
		}
	}

	return 0
}

func add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

func concat(a, b int) int {
	concatenated := strconv.Itoa(a) + strconv.Itoa(b)
	result, err := strconv.Atoi(concatenated)
	if err != nil {
		panic(err)
	}
	return result
}
