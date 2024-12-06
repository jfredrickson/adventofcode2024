package main

import (
	"common"
	"fmt"
	"slices"
	"strings"
)

func main() {
	processingRules := true
	rules := make([]Rule, 0)
	pageLists := make([][]int, 0)

	common.ProcessFile("day05/input.txt", func(line string) {
		if line == "" {
			processingRules = false
			return
		}

		if processingRules {
			// line is a rule
			pages := common.ToInts(strings.Split(line, "|"))
			rules = append(rules, Rule{Before: pages[0], After: pages[1]})
		} else {
			// line is a list of pages
			pageLists = append(pageLists, common.ToInts(strings.Split(line, ",")))
		}
	})

	correctlyOrderedSum := 0
	incorrectlyOrderedSum := 0

	for _, pageList := range pageLists {
		update := Update{Rules: rules, Pages: pageList}
		result := update.ProcessRules()
		if result > 0 {
			// the page order is correct; add its result to the correct sum
			correctlyOrderedSum += result
		} else {
			// the page order is incorrect; fix it and then add its result to the incorrect sum
			update.Fix()
			incorrectlyOrderedSum += update.ProcessRules()
		}
	}

	fmt.Println("Sum of correctly ordered pages:", correctlyOrderedSum)
	fmt.Println("Sum of incorrectly ordered pages:", incorrectlyOrderedSum)
}

type Rule struct {
	Before int
	After  int
}

type Update struct {
	Rules []Rule
	Pages []int
}

func (u *Update) GetMiddlePage() int {
	return u.Pages[len(u.Pages)/2]
}

// Check if the page order follows the rules, returning the middle page number if so
func (u *Update) ProcessRules() int {
	result := u.GetMiddlePage()

	for _, rule := range u.Rules {
		indexBefore := slices.Index(u.Pages, rule.Before)
		indexAfter := slices.Index(u.Pages, rule.After)
		if indexBefore == -1 || indexAfter == -1 {
			continue
		}
		if indexAfter < indexBefore {
			return 0
		}
	}

	return result
}

// Fix the page order so that it follows the rules
func (u *Update) Fix() {
	// Create a graph of pages
	pageGraph := make(map[int][]int, 0)
	for _, page := range u.Pages {
		pageGraph[page] = make([]int, 0)
	}

	// Build edges in the graph based on rules
	for _, rule := range u.Rules {
		// Process this rule only if the before and after values are both in the set of updated pages
		if !slices.Contains(u.Pages, rule.Before) || !slices.Contains(u.Pages, rule.After) {
			continue
		}
		pageGraph[rule.After] = append(pageGraph[rule.After], rule.Before)
	}

	orderedPages := make([]int, 0)

	// Topological sort
	for len(pageGraph) > 0 {
		// Find a node with no incoming edges
		for pageNumber, node := range pageGraph {
			if len(node) == 0 {
				// If node has no incoming edges, it's next in the ordered list
				orderedPages = append(orderedPages, pageNumber)
				// Remove that next node from the graph
				for otherPageNumber, otherNode := range pageGraph {
					if slices.Contains(otherNode, pageNumber) {
						index := slices.Index(otherNode, pageNumber)
						pageGraph[otherPageNumber] = slices.Delete(pageGraph[otherPageNumber], index, index+1)
					}
				}
				delete(pageGraph, pageNumber)
			}
		}
	}

	// Replace the incorrectly ordered pages with the correctly ordered pages
	u.Pages = orderedPages
}
