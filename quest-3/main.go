package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type crates []int

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func processInput(fileName string) crates {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var crates crates

	for _, crate := range strings.Split(string(data), ",") {
		if crate == "" {
			continue
		}
		crates = append(crates, aToIIgnoreError(crate))
	}
	return crates
}

func (c crates) uniqueAndSortedDesc() crates {
	seen := make(map[int]struct{}, len(c))
	result := make(crates, 0, len(c))

	for _, n := range c {
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		result = append(result, n)
	}

	slices.SortFunc(result, func(a, b int) int {
		return cmp.Compare(b, a)
	})
	return result
}

func (c crates) sum() int {
	result := 0
	seen := make(map[int]struct{}, len(c))
	for _, n := range c {
		// on smaller crates can go in each, so therefore only track uniques
		if _, ok := seen[n]; !ok {
			result += n
			seen[n] = struct{}{}
		}
	}
	return result
}

func largeSetOfCrates(c crates) int {
	largest := 0
	current := c
	for len(current) > 1 {
		sum := current.sum()
		if sum > largest {
			largest = sum
		}
		current = current[1:]
	}

	return largest
}

// mushroom can only be packed in sets of 20 crates
func smallestMushroomPack(c crates) int {
	smallest := -1
	for i := 0; i <= len(c)-20; i++ {
		window := c[i : i+20]
		sum := window.sum()
		if sum < smallest || smallest == -1 {
			smallest = sum
		}
	}
	return smallest

}

// For part three, we need to find the minimum number of sets of crates needed to pack everything.
// This is equivalent to finding the maximum frequency of any crate size,
// since that determines how many sets we need to accommodate all crates of that size.
func calcSetsToPackEverything(c crates) int {
	freq := make(map[int]int)
	for _, n := range c {
		freq[n]++
	}
	maxFreq := 0
	for _, count := range freq {
		if count > maxFreq {
			maxFreq = count
		}
	}
	return maxFreq
}

func main() {
	crates := processInput("input-part-1.txt").uniqueAndSortedDesc()
	fmt.Println("Part One: The largest set of crates is:", largeSetOfCrates(crates))

	crates2 := processInput("input-part-2.txt").uniqueAndSortedDesc()
	fmt.Println("Part Two: smallest set to pack the mushroom is:", smallestMushroomPack(crates2))

	crates3 := processInput("input-part-3.txt")
	fmt.Println("Part Three: minimum number of sets to pack everything is:", calcSetsToPackEverything(crates3))

}
