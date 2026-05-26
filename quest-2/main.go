package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type complexNum struct {
	x int
	y int
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func complexNumFromString(s string) complexNum {
	re := regexp.MustCompile(`\[(-?\d+),(-?\d+)\]`)
	matches := re.FindStringSubmatch(s)
	return complexNum{
		x: aToIIgnoreError(matches[1]),
		y: aToIIgnoreError(matches[2]),
	}
}

func processInput(fileName string) []complexNum {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var complexNums []complexNum

	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		complexNums = append(complexNums, complexNumFromString(line))
	}
	return complexNums
}

func (c complexNum) add(other complexNum) complexNum {
	return complexNum{
		x: c.x + other.x,
		y: c.y + other.y,
	}
}

func (c complexNum) multi(other complexNum) complexNum {
	return complexNum{
		x: (c.x * other.x) - (c.y * other.y),
		y: (c.x * other.y) + (c.y * other.x),
	}
}

func (c complexNum) div(other complexNum) complexNum {
	return complexNum{
		x: c.x / other.x,
		y: c.y / other.y,
	}
}

func (c complexNum) processPartOne(iterations int, A complexNum) complexNum {
	result := c
	divisor := complexNum{x: 10, y: 10}
	for range iterations {
		result = result.multi(result)
		result = result.div(divisor)
		result = result.add(A)
	}
	return result
}

func (c complexNum) isValidEngraving(iterations, boundaryMin, boundaryMax int) bool {
	result := complexNum{x: 0, y: 0}
	divisor := complexNum{x: 100000, y: 100000}
	for range iterations {
		result = result.multi(result)
		result = result.div(divisor)
		result = result.add(c)
		if !result.isWithinBoundary(boundaryMin, boundaryMax) {
			return false
		}
	}
	return true
}

func (c complexNum) isWithinBoundary(boundaryMin, boundaryMax int) bool {
	return c.x >= boundaryMin && c.x <= boundaryMax && c.y >= boundaryMin && c.y <= boundaryMax
}

func (c complexNum) String() string {
	return fmt.Sprintf("[%d,%d]", c.x, c.y)
}

func processGrid(A complexNum, gridSize, cycles, boundaryMin, boundaryMax int) int {
	step := 1000 / gridSize
	totalWithinBoundary := 0
	for row := 0; row <= gridSize; row++ {
		for col := 0; col <= gridSize; col++ {
			current := complexNum{x: A.x + col*step, y: A.y + row*step}
			if current.isValidEngraving(cycles, boundaryMin, boundaryMax) {
				totalWithinBoundary++
			}
		}
	}
	return totalWithinBoundary
}

func main() {
	A := processInput("input-part-1.txt")[0]
	fmt.Println("Part One:", complexNum{x: 0, y: 0}.processPartOne(3, A))

	A2 := processInput("input-part-2.txt")[0]
	fmt.Println("Part Two:", processGrid(A2, 100, 100, -1000000, 1000000))

	A3 := processInput("input-part-3.txt")[0]
	fmt.Println("Part Three:", processGrid(A3, 1000, 100, -1000000, 1000000))

}
