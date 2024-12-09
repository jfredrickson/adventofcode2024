package main

import (
	"common"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	filename := "day09/input.txt"

	fs1 := FileSystem{}
	fs2 := FileSystem{}
	fs1.Populate(filename)
	fs2.Populate(filename)

	fs1.Compact()
	fs2.Defragment()

	fmt.Println("Compacted checksum:", fs1.Checksum())
	fmt.Println("Defragmented checksum:", fs2.Checksum())
}

type FileSystem []int

// Populate the filesystem with data from the input file
func (fs *FileSystem) Populate(filename string) {
	input, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	diskMap := common.ToInts(strings.Split(string(input), ""))

	fileId := 0
	isFile := true
	for _, data := range diskMap {
		out := -1
		if isFile {
			out = fileId
			fileId++
		}
		for i := 0; i < data; i++ {
			*fs = append(*fs, out)
		}
		isFile = !isFile
	}
}

// Find the index of the next empty section of the given size
func (fs FileSystem) NextEmptySection(size int) int {
	for i := 0; i < len(fs); i++ {
		// Check if we'd go past the end of the filesystem
		if i+size > len(fs) {
			return -1
		}
		// If all blocks are empty (max value is -1), this section is available
		if slices.Max(fs[i:i+size]) == -1 {
			return i
		}
	}

	return -1
}

// Move blocks of data from one index to another
func (fs *FileSystem) Move(fromIndex, toIndex, size int) {
	for i := 0; i < size; i++ {
		(*fs)[toIndex+i] = (*fs)[fromIndex+i]
		(*fs)[fromIndex+i] = -1
	}
}

// Move all data blocks to the beginning of the filesystem
func (fs *FileSystem) Compact() {
	// Work backwards from the end of the filesystem
	var ei int
	for di := len(*fs) - 1; di >= 0; di-- {
		// If block is empty, skip
		if (*fs)[di] == -1 {
			continue
		}

		// Find next available empty block
		ei = fs.NextEmptySection(1)
		if ei == -1 {
			panic("No empty block found")
		}

		// Check if the filesystem is done
		if slices.Max((*fs)[ei:]) == -1 {
			break
		}

		// Move the data block to the empty block
		fs.Move(di, ei, 1)
	}
}

// Move contiguous blocks of data to the beginning of the filesystem
func (fs *FileSystem) Defragment() {
	// Work backwards from the end of the filesystem
	for di := len(*fs) - 1; di >= 0; di-- {
		// If block is empty, skip
		if (*fs)[di] == -1 {
			continue
		}

		// Find the size of this data block
		size := 1
		fileId := (*fs)[di]
		for i := di - 1; i >= 0; i-- {
			if (*fs)[i] == fileId {
				di--
				size++
			} else {
				break
			}
		}

		// Find next available empty section
		ei := fs.NextEmptySection(size)
		if ei == -1 || ei >= di {
			// There are no empty sections or the only empty sections are closer to the end of the filesystem
			continue
		}

		fs.Move(di, ei, size)
	}
}

func (fs FileSystem) Checksum() int {
	sum := 0
	for i, v := range fs {
		if v == -1 {
			continue
		}
		sum += v * i
	}
	return sum
}
