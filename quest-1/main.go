package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type instruction struct {
	raw    string
	offset int
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func generateInstructions(raw string) []instruction {
	var instructions []instruction

	for _, rawInstruction := range strings.Split(raw, ",") {
		direction := rawInstruction[0]
		rawOffset := aToIIgnoreError(string(rawInstruction[1:]))
		if direction == 'L' {
			rawOffset = rawOffset * -1
		}
		instructions = append(instructions, instruction{
			raw:    rawInstruction,
			offset: rawOffset,
		})
	}
	return instructions
}

func processInput(fileName string) ([]string, []instruction) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	return strings.Split(lines[0], ","), generateInstructions(lines[2])
}

func findNameSimple(names []string, instructions []instruction) string {
	currentIndex := 0
	namesLength := len(names)
	for _, instruction := range instructions {
		if (currentIndex + instruction.offset) >= namesLength {
			currentIndex = namesLength - 1
		} else if (currentIndex + instruction.offset) < 0 {
			currentIndex = 0
		} else {
			currentIndex = currentIndex + instruction.offset

		}
	}
	return names[currentIndex]
}

func findNameWithWrapping(names []string, instructions []instruction) string {
	currentIndex := 0
	namesLength := len(names)
	for _, instruction := range instructions {
		// add namesLength to account for negatives
		currentIndex = ((currentIndex+instruction.offset)%namesLength + namesLength) % namesLength
	}
	return names[currentIndex]
}

func findNamesWithSwapping(names []string, instructions []instruction) string {
	namesLength := len(names)
	for _, instruction := range instructions {
		// get the next swap candidate from top (0) then simply swap them
		nextSwapIndex := ((0+instruction.offset)%namesLength + namesLength) % namesLength
		names[0], names[nextSwapIndex] = names[nextSwapIndex], names[0]
	}
	return names[0]
}

func main() {
	names, instructions := processInput("input-part-1.txt")
	names2, instructions2 := processInput("input-part-2.txt")
	names3, instructions3 := processInput("input-part-3.txt")
	fmt.Println("Your part one name is:", findNameSimple(names, instructions))
	fmt.Println("Your part two name is:", findNameWithWrapping(names2, instructions2))
	fmt.Println("Your part two name is:", findNamesWithSwapping(names3, instructions3))
}
