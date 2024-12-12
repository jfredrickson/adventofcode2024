package main

import (
	"common"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, _ := os.ReadFile("day11/input.txt")
	input := common.ToInts(strings.Split(string(data), " "))

	fmt.Println("Number of stones after 25 blinks:", count(input, 25))
	fmt.Println("Number of stones after 75 blinks:", count(input, 75))
}

func count(input []int, numBlinks int) int {
	// Create a map of the quantity of each stone number
	stones := make(map[int]int, 0)
	for _, num := range input {
		stones[num]++
	}

	// Apply the rules to get an updated map for each blink
	for range numBlinks {
		stones = applyRules(stones)
	}

	// Sum up the resulting unique stone numbers
	count := 0
	for _, quantity := range stones {
		count += quantity
	}
	return count
}

func applyRules(stones map[int]int) map[int]int {
	changedStones := make(map[int]int, 0)

	for number, quantity := range stones {
		if number == 0 {
			changedStones[1] += quantity
		} else if countDigits(number)%2 == 0 {
			left := atoi(itoa(number)[:len(itoa(number))/2])
			right := atoi(itoa(number)[len(itoa(number))/2:])
			changedStones[left] += quantity
			changedStones[right] += quantity
		} else {
			changedStones[number*2024] += quantity
		}
	}

	return changedStones
}

/* A few utility functions just to make things more readable in the main code */

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

func countDigits(n int) int {
	return len(itoa(n))
}
