package main

import (
	"common"
	"fmt"
)

func main() {
	sum := 0
	common.ProcessFile("day22/input.txt", func(line string) {
		s := common.Atoi(line)
		for range 2000 {
			s = next(s)
		}
		sum += s
	})

	fmt.Println("Sum of 2000th secret numbers:", sum)
}

func next(s int) int {
	s = ((s * 64) ^ s) % 16777216
	s = ((s / 32) ^ s) % 16777216
	s = ((s * 2048) ^ s) % 16777216
	return s
}
