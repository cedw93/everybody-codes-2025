package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type gear struct {
	teeth int
	turns int
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func processInput(fileName string) []*gear {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	gears := []*gear{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		gears = append(gears, &gear{
			teeth: aToIIgnoreError(line),
			turns: 0,
		})
	}

	return gears
}

func calcTurnsOfLastGearAfterTurnsOfFirstGear(gears []*gear, turns int) int {
	return (turns * gears[0].teeth) / gears[len(gears)-1].teeth
}

func calcMinimumTurnsOfFirstFromLastNturns(gears []*gear, turns int) int {
	// Force integer division to round up by adding the divisor - 1 to the dividend before dividing
	return (turns*gears[len(gears)-1].teeth + gears[0].teeth - 1) / gears[0].teeth
}

func main() {
	gears := processInput("input-part-1.txt")
	fmt.Println("After 2025 turns of gear 1 the last gear will have turned: ", calcTurnsOfLastGearAfterTurnsOfFirstGear(gears, 2025), "times")
	gears2 := processInput("input-part-2.txt")
	fmt.Println("After 2025 turns of gear 1 the last gear will have turned: ", calcMinimumTurnsOfFirstFromLastNturns(gears2, 10_000_000_000_000), "times")
}
