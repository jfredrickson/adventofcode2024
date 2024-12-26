package main

import (
	"common"
	"fmt"
)

const priceChanges = 2000

func main() {
	buyers := []Buyer{}

	common.ProcessFile("day22/input.txt", func(line string) {
		secret := common.Atoi(line)
		price := secret % 10
		buyer := Buyer{{Secret: secret, Price: price, Change: 0}}
		for range priceChanges {
			lastPrice := price
			secret = nextSecret(secret)
			price = secret % 10
			buyer = append(buyer, State{Secret: secret, Price: price, Change: price - lastPrice})
		}
		buyers = append(buyers, buyer)
	})

	sum := 0
	for _, b := range buyers {
		sum += b[len(b)-1].Secret
	}

	fmt.Printf("Sum of %dth secret numbers: %d\n", priceChanges, sum)
	fmt.Println("Most bananas possible:", findMostBananas(buyers))
}

type Sequence [4]int

type Buyer []State

type State struct {
	Secret, Price, Change int
}

func nextSecret(s int) int {
	s = ((s * 64) ^ s) % 16777216
	s = ((s / 32) ^ s) % 16777216
	s = ((s * 2048) ^ s) % 16777216
	return s
}

func findMostBananas(buyers []Buyer) int {
	bananas := make(map[Sequence]int, 0)

	for _, buyer := range buyers {
		// Track the bananas for the current buyer
		currentBananas := make(map[Sequence]int, 0)

		// Iterate through 4-number sequences in the changes
		for i := 0; i < len(buyer)-4; i++ {
			sequence := Sequence{}
			for j := 0; j < 4; j++ {
				sequence[j] = buyer[i+j+1].Change
			}
			// Keep only the first occurrence of this sequence in the current buyer
			if _, exists := currentBananas[sequence]; !exists {
				currentBananas[sequence] += buyer[i+4].Price
			}
		}

		for seq, amount := range currentBananas {
			bananas[seq] += amount
		}
	}

	max := 0
	for _, amount := range bananas {
		if amount > max {
			max = amount
		}
	}

	return max
}
