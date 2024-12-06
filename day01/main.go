package main

import (
	"common"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func main() {
	search := Search{
		Locations1: []int{},
		Locations2: []int{},
	}

	search.LoadFile("day01/input.txt")

	fmt.Println("Total distance:", search.TotalDistance())
	fmt.Println("Similarity score:", search.SimilarityScore())
}

type Search struct {
	Locations1 []int
	Locations2 []int
}

func (s *Search) LoadFile(filename string) {
	common.ProcessFile(filename, func(line string) {
		first, second, _ := strings.Cut(line, " ")
		first = strings.TrimSpace(first)
		second = strings.TrimSpace(second)

		loc1, err := strconv.Atoi(first)
		if err != nil {
			panic(err)
		}

		loc2, err := strconv.Atoi(second)
		if err != nil {
			panic(err)
		}

		s.Locations1 = append(s.Locations1, loc1)
		s.Locations2 = append(s.Locations2, loc2)
	})
}

func (s *Search) TotalDistance() int {
	slices.Sort(s.Locations1)
	slices.Sort(s.Locations2)

	sum := 0

	for i := 0; i < len(s.Locations1); i++ {
		sum += common.Abs(s.Locations1[i] - s.Locations2[i])
	}

	return sum
}

func (s *Search) SimilarityScore() int {
	sum := 0
	counts := make(map[int]int)

	for i := 0; i < len(s.Locations1); i++ {
		for j := 0; j < len(s.Locations2); j++ {
			if s.Locations1[i] == s.Locations2[j] {
				counts[s.Locations1[i]]++
			}
		}
	}

	for num, count := range counts {
		sum += num * count
	}

	return sum
}
