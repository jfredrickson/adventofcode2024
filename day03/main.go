package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	code, err := os.ReadFile("day03/input.txt")
	if err != nil {
		panic(err)
	}

	program := getProgram(string(code))

	fmt.Println("Sum:", program.Run(false))
	fmt.Println("Sum with conditionals:", program.Run(true))
}

func getProgram(s string) *Program {
	program := &Program{Instructions: []Instruction{}}

	instructionRegexp := regexp.MustCompile(`(mul|do|don't)\((\d{1,3},\d{1,3})?\)`)
	matches := instructionRegexp.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		instruction := Instruction{Name: match[1]}
		if match[2] != "" {
			argsRegexp := regexp.MustCompile(`(\d{1,3}),(\d{1,3})`)
			argMatch := argsRegexp.FindStringSubmatch(match[2])
			instruction.Args = []string{argMatch[1], argMatch[2]}
		}
		program.Instructions = append(program.Instructions, instruction)
	}

	return program
}

func (p *Program) Run(useConditionals bool) int {
	sum := 0
	mulEnabled := true

	for _, instruction := range p.Instructions {
		switch instruction.Name {
		case "mul":
			if mulEnabled {
				x, err := strconv.Atoi(instruction.Args[0])
				if err != nil {
					panic(err)
				}

				y, err := strconv.Atoi(instruction.Args[1])
				if err != nil {
					panic(err)
				}

				sum += x * y
			}
		case "do":
			if useConditionals {
				mulEnabled = true
			}
		case "don't":
			if useConditionals {
				mulEnabled = false
			}
		default:
			panic("Shouldn't get here")
		}
	}

	return sum
}

type Program struct {
	Instructions []Instruction
}

type Instruction struct {
	Name string
	Args []string
}
