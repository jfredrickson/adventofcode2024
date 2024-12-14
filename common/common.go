package common

import (
	"bufio"
	"os"
	"strconv"
)

type Number interface {
	int
}

// Read a file line by line, calling the process function on each line
func ProcessFile(filename string, process func(string)) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		process(line)
	}
	if err = scanner.Err(); err != nil {
		panic(err)
	}
}

// Convert a slice of strings to a slice of ints
func ToInts(s []string) []int {
	ints := make([]int, len(s))
	for i, val := range s {
		intval, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		ints[i] = intval
	}
	return ints
}

func Abs[N Number](n N) N {
	if n < 0 {
		return -n
	}
	return n
}

func Atoi[N Number](s string) N {
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return N(val)
}
