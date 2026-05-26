package main

import (
	"bufio"
	"fmt"
	"os"
)

type direction struct {
	rowOffset int
	colOffset int
	label     string
}

type cell struct {
	row int
	col int
}

const (
	soundSource = '@'
	vocalBone   = '#'
	empty       = '.'
)

var (
	directions = []direction{
		// order matters here for the problem at hand, we want to check up first, then right, then down, then left
		{rowOffset: -1, colOffset: 0, label: "up"},
		{rowOffset: 0, colOffset: 1, label: "right"},
		{rowOffset: 1, colOffset: 0, label: "down"},
		{rowOffset: 0, colOffset: -1, label: "left"},
	}
)

func processInput(fileName string) (cell, cell, [][]rune) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]rune
	var soundSourceCell cell
	var vocalBoneCell cell

	scanner := bufio.NewScanner(file)

	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		gridRow := []rune{}
		for col, char := range line {
			switch char {
			case soundSource:
				soundSourceCell = cell{row: row, col: col}
			case vocalBone:
				vocalBoneCell = cell{row: row, col: col}
			}
			gridRow = append(gridRow, char)
		}
		grid = append(grid, gridRow)
		row++
	}
	return soundSourceCell, vocalBoneCell, grid
}

func (c cell) isWithinGrid(grid [][]rune) bool {
	return c.row >= 0 && c.row < len(grid) && c.col >= 0 && c.col < len(grid[0])
}

func partOne(start, end cell, grid [][]rune) int {
	totalSteps := 0
	current := start
	visited := map[cell]bool{current: true}
	dirIndex := 0
	for current.row != end.row || current.col != end.col {
		dir := directions[dirIndex%4]
		dirIndex++
		next := cell{row: current.row + dir.rowOffset, col: current.col + dir.colOffset}
		if next.isWithinGrid(grid) && !visited[next] {
			current = next
			visited[current] = true
			totalSteps++
		}
	}
	return totalSteps
}

func main() {
	soundSourceCell, vocalBoneCell, grid := processInput("input-part-1.txt")
	fmt.Println("Part One steps:", partOne(soundSourceCell, vocalBoneCell, grid))
}
