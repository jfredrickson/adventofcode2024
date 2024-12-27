package main

import (
	"common"
	"fmt"
	"maps"
	"regexp"
	"strings"
)

func main() {
	circuit := Circuit{
		Wires: make(map[string]*Wire),
		Gates: make([]*Gate, 0),
	}

	processingWires := true
	common.ProcessFile("day24/input.txt", func(line string) {
		if line == "" {
			processingWires = false
			return
		}

		if processingWires {
			wireData := strings.Split(line, ": ")
			wireName := wireData[0]
			wireValue := common.Atoi(wireData[1])
			circuit.Wires[wireName] = &Wire{Name: wireName, Value: wireValue}
		} else {
			connectionData := strings.Split(line, " -> ")
			gateData := strings.Split(connectionData[0], " ")
			outputWire := circuit.GetOrCreateWire(connectionData[1])
			inputWireA := circuit.GetOrCreateWire(gateData[0])
			inputWireB := circuit.GetOrCreateWire(gateData[2])
			gateType := GateType(gateData[1])
			gate := Gate{
				Inputs:     []*Wire{inputWireA, inputWireB},
				OutputWire: outputWire,
				Type:       gateType,
			}
			circuit.Gates = append(circuit.Gates, &gate)
		}
	})

	circuit.Process()

	fmt.Println("Output:", circuit.WiresToDecimal("z"))
}

type GateType string

const (
	AND GateType = "AND"
	OR  GateType = "OR"
	XOR GateType = "XOR"
)

type Gate struct {
	Inputs     []*Wire
	OutputWire *Wire
	Type       GateType
}

type Wire struct {
	Name  string
	Value int
}

type Circuit struct {
	Wires map[string]*Wire
	Gates []*Gate
}

func (c *Circuit) GetOrCreateWire(wireName string) *Wire {
	if wire, exists := (*c).Wires[wireName]; exists {
		return wire
	} else {
		newWire := &Wire{Name: wireName, Value: -1}
		(*c).Wires[wireName] = newWire
		return newWire
	}
}

func (c *Circuit) Process() {
	processing := true
	for processing {
		for _, g := range c.Gates {
			ready := true
			for _, i := range g.Inputs {
				if i.Value == -1 {
					ready = false
					break
				}
			}
			if ready {
				c.Wires[g.OutputWire.Name].Value = g.GetOutput()
			}
		}

		processing = false
		for wire := range maps.Values(c.Wires) {
			if wire.Value == -1 {
				processing = true
			}
		}
	}
}

func (c *Circuit) WiresToDecimal(prefix string) int {
	digits := map[int]int{}
	r := regexp.MustCompile(fmt.Sprintf(`^%s(\d\d)$`, prefix))
	for wireName := range c.Wires {
		if r.MatchString(wireName) {
			place := common.Atoi(strings.TrimLeft(wireName, prefix))
			digits[place] = c.Wires[wireName].Value
		}
	}

	number := 0
	for place, digit := range digits {
		number += digit << place
	}

	return number
}

func (w Wire) String() string {
	return fmt.Sprintf("{%s %d}", w.Name, w.Value)
}

func (g *Gate) GetOutput() int {
	inputA := g.Inputs[0].Value
	inputB := g.Inputs[1].Value

	switch g.Type {
	case AND:
		return inputA & inputB
	case OR:
		return inputA | inputB
	case XOR:
		return inputA ^ inputB
	default:
		return 0
	}
}

func (g Gate) String() string {
	return fmt.Sprintf("%s [%v, %v] -> %v", g.Type, *g.Inputs[0], *g.Inputs[1], *g.OutputWire)
}
