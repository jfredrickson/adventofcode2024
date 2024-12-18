package main

import (
	"common"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func main() {
	cpu := Reset()
	program := make([]int, 0)
	var programData string

	processingRegisters := true
	common.ProcessFile("day17/input.txt", func(line string) {
		if line == "" {
			processingRegisters = false
			return
		}

		if processingRegisters {
			r := regexp.MustCompile(`^Register (.): (\d+)$`)
			match := r.FindStringSubmatch(line)
			if match[1] == "A" {
				cpu.A = common.Atoi(match[2])
			}
			if match[1] == "B" {
				cpu.B = common.Atoi(match[2])
			}
			if match[1] == "C" {
				cpu.C = common.Atoi(match[2])
			}
		} else {
			programData = strings.Split(line, ": ")[1]
			program = common.ToInts(strings.Split(programData, ","))
		}
	})

	// Part 1

	cpu.Run(program)
	fmt.Println("Output:", cpu.Sprint())

	// Part 2

	// This is a brute force method that just won't work in the scale of the real input.
	// Correct approach might be to reverse the program, running it backwards?
	// Example of reversed operations:
	// - XOR is basically the inverse of itself, so the XOR operations can just be performed as is
	// - Division becomes multiplication
	// - 'out' instruction basically pops the last element off the output string
	// - Modulo - not sure yet what to do about that, but keep in mind this is a 3-bit system

	// Set up a function to run as a goroutine
	run := func(a int) int {
		c := Reset()
		c.A = a

		// A debug version of Run() that compares the output with the programData
		// This saves about 25% of runtime, not enough to matter for the real input
		for c.Pointer < len(program)-1 {
			// If the current instruction is 'out', see if the output matches so far
			if program[c.Pointer] == 5 && strings.Index(programData, c.Sprint()) != 0 {
				// If it doesn't match, no need to continue with this program
				return -1
			}
			c.Opcodes[program[c.Pointer]](program[c.Pointer+1])
			c.Pointer += 2
		}

		if c.Sprint() == programData {
			return a
		}
		return -1
	}

	// Prepare channels and wait group
	maxA := math.MaxInt
	maxConcurrency := 16

	// Run the program with batches of A values
	found := false
	lowestA := math.MaxInt
	for i := 0; i <= maxA && !found; i += maxConcurrency {
		results := make(chan int, maxConcurrency)
		var wg sync.WaitGroup
		for j := 0; j < maxConcurrency; j++ {
			a := i + j
			if a > maxA {
				break
			}
			wg.Add(1)
			go func(a int) {
				defer wg.Done()
				results <- run(a)
			}(a)
		}
		wg.Wait()
		close(results)

		// Find the lowest value of A among this batch
		for result := range results {
			if result > -1 && result < lowestA {
				lowestA = result
			}
		}

		if lowestA < maxA {
			found = true
		}
	}

	fmt.Println("Quine register A:", lowestA)
}

type Opcode func(int)

type CPU struct {
	A, B, C int
	Opcodes map[int]Opcode
	Pointer int
	Output  []int
}

func Reset() *CPU {
	c := CPU{
		A:       0,
		B:       0,
		C:       0,
		Opcodes: make(map[int]Opcode),
		Pointer: 0,
		Output:  make([]int, 0),
	}
	c.Opcodes[0] = c.adv
	c.Opcodes[1] = c.bxl
	c.Opcodes[2] = c.bst
	c.Opcodes[3] = c.jnz
	c.Opcodes[4] = c.bxc
	c.Opcodes[5] = c.out
	c.Opcodes[6] = c.bdv
	c.Opcodes[7] = c.cdv
	return &c
}

func (c *CPU) Run(program []int) {
	for c.Pointer < len(program)-1 {
		// fn := strings.Split(runtime.FuncForPC(reflect.ValueOf(c.Opcodes[program[c.Pointer]]).Pointer()).Name(), ".")[2][:3]
		// fmt.Println("Instruction pointer", c.Pointer)
		// fmt.Println("  Opcode", program[c.Pointer], fn)
		// fmt.Println("  Operand", program[c.Pointer+1])
		// fmt.Println("  A:", c.A, "  B:", c.B, "  C:", c.C)
		// fmt.Println("  Output:", c.Output)
		c.Opcodes[program[c.Pointer]](program[c.Pointer+1])
		c.Pointer += 2
	}
}

func (c *CPU) Sprint() string {
	out := make([]string, 0)
	for _, o := range c.Output {
		out = append(out, strconv.Itoa(o))
	}
	return strings.Join(out, ",")
}

func (c *CPU) adv(combo int) {
	c.A = c.A / common.Pow(2, c.parseCombo(combo))
}

func (c *CPU) bxl(literal int) {
	c.B = c.B ^ literal
}

func (c *CPU) bst(combo int) {
	c.B = c.parseCombo(combo) % 8
}

func (c *CPU) jnz(literal int) {
	if c.A != 0 {
		c.Pointer = literal - 2
	}
}

func (c *CPU) bxc(_ int) {
	c.B = c.B ^ c.C
}

func (c *CPU) out(combo int) {
	c.Output = append(c.Output, c.parseCombo(combo)%8)
}

func (c *CPU) bdv(combo int) {
	c.B = c.A / common.Pow(2, c.parseCombo(combo))
}

func (c *CPU) cdv(combo int) {
	c.C = c.A / common.Pow(2, c.parseCombo(combo))
}

func (c *CPU) parseCombo(operand int) int {
	if operand < 4 {
		return operand
	}
	if operand == 4 {
		return c.A
	}
	if operand == 5 {
		return c.B
	}
	if operand == 6 {
		return c.C
	}
	return -1
}
