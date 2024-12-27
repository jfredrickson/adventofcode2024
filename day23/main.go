package main

import (
	"common"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

func main() {
	network := Network{}

	common.ProcessFile("day23/example.txt", func(line string) {
		computerData := strings.Split(line, "-")
		c1, c2 := Computer(computerData[0]), Computer(computerData[1])
		network.Add(c1, c2)
	})

	triples := network.FindTriples()

	matches := 0
	for _, t := range triples {
		if t.Match(`^t.$`) {
			matches++
		}
	}

	fmt.Println("Interconnected sets containing a computer name starting with 't':", matches)
}

type Computer string
type Network map[Computer]Interconnection
type Interconnection map[Computer]bool

func (n Network) Add(c1, c2 Computer) {
	if _, exists := n[c1]; !exists {
		n[c1] = Interconnection{}
	}
	if _, exists := n[c2]; !exists {
		n[c2] = Interconnection{}
	}
	n[c1][c2] = true
	n[c2][c1] = true
}

func (n Network) String() string {
	var s string
	for c, interconnections := range n {
		otherComputers := make([]string, 0, len(interconnections))
		for otherComputer := range interconnections {
			otherComputers = append(otherComputers, string(otherComputer))
		}
		s += fmt.Sprintf("{%v: %v}", c, strings.Join(otherComputers, ","))
	}
	return s
}

func (n Network) FindTriples() []Interconnection {
	triples := []Interconnection{}
	for computerA, otherComputers := range n {
		for computerB := range otherComputers {
			for computerC := range otherComputers {
				if n[computerB][computerC] {
					tripleExists := false
					for _, triple := range triples {
						if triple[computerA] && triple[computerB] && triple[computerC] {
							tripleExists = true
							break
						}
					}
					if !tripleExists {
						triple := Interconnection{}
						triple.Add(computerA)
						triple.Add(computerB)
						triple.Add(computerC)
						triples = append(triples, triple)
					}
				}
			}
		}
	}
	return triples
}

func (i Interconnection) Add(c Computer) {
	i[c] = true
}

func (i Interconnection) String() string {
	return fmt.Sprintf("%v", strings.Join(i.SortedNames(), ","))
}

func (i Interconnection) Match(regex string) bool {
	r := regexp.MustCompile(regex)
	for computer := range i {
		if r.MatchString(string(computer)) {
			return true
		}
	}
	return false
}

func (i Interconnection) SortedNames() []string {
	names := make([]string, 0, len(i))
	for name := range i {
		names = append(names, string(name))
	}
	slices.Sort(names)
	return names
}
