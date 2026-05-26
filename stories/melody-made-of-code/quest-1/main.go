package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type dragonScale struct {
	raw         string
	red         string
	blue        string
	green       string
	shine       string
	redBinary   int
	greenBinary int
	blueBinary  int
	shineBinary int
	binary      int
	value       int
	group       string
}

const (
	redMatte          = "red-matte"
	redGlossy         = "red-shiny"
	blueMatte         = "blue-matte"
	blueGlossy        = "blue-shiny"
	greenMatte        = "green-matte"
	greenGlossy       = "green-shiny"
	noGroup           = "---"
	matteThresholdMax = 0b011110
	shinyThresholdMin = 0b100001
)

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func stringToBinary(s string) int {
	result := 0b0000000
	for idx, char := range s {
		if char >= 'A' && char <= 'Z' {
			result |= 1 << (len(s) - 1 - idx)
		}
	}
	return result
}

func (d *dragonScale) determineGroup() {
	if d.shineBinary > matteThresholdMax && d.shineBinary < shinyThresholdMin {
		d.group = noGroup
	} else if d.redBinary > d.greenBinary && d.redBinary > d.blueBinary {
		if d.shineBinary > matteThresholdMax {
			d.group = redGlossy
		} else {
			d.group = redMatte
		}
	} else if d.greenBinary > d.redBinary && d.greenBinary > d.blueBinary {
		if d.shineBinary > matteThresholdMax {
			d.group = greenGlossy
		} else {
			d.group = greenMatte
		}
	} else if d.blueBinary > d.redBinary && d.blueBinary > d.greenBinary {
		if d.shineBinary > matteThresholdMax {
			d.group = blueGlossy
		} else {
			d.group = blueMatte
		}
	}
}

func noteToDragonScale(note string) dragonScale {
	parts := strings.Split(note, ":")
	value := parts[0]
	colours := strings.Split(parts[1], " ")
	redBinary := stringToBinary(colours[0])
	greenBinary := stringToBinary(colours[1])
	blueBinary := stringToBinary(colours[2])
	shine := ""
	shineBinary := 0b0
	if len(colours) > 3 {
		shine = colours[3]
		shineBinary = stringToBinary(colours[3])
	}
	scale := dragonScale{
		raw:         note,
		red:         colours[0],
		green:       colours[1],
		blue:        colours[2],
		shine:       shine,
		redBinary:   redBinary,
		greenBinary: greenBinary,
		blueBinary:  blueBinary,
		shineBinary: shineBinary,
		value:       aToIIgnoreError(value),
	}
	(&scale).determineGroup()
	return scale
}

func (d dragonScale) String() string {
	return fmt.Sprintf("Value: %d, Red: %s (%06b), Green: %s (%06b), Blue: %s (%06b), Binary: %06b, Group: %s",
		d.value, d.red, d.redBinary, d.green, d.greenBinary, d.blue, d.blueBinary, d.binary, d.group)
}

func processInput(fileName string) []dragonScale {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scales := []dragonScale{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		scales = append(scales, noteToDragonScale(line))
	}
	return scales
}

func countDominantGreens(scales []dragonScale) int {
	result := 0
	for _, scale := range scales {
		if scale.greenBinary > scale.redBinary && scale.greenBinary > scale.blueBinary {
			result += scale.value
		}
	}
	return result
}

func (d dragonScale) sumColours() int {
	return d.redBinary + d.greenBinary + d.blueBinary
}

func findDarkestShine(scales []dragonScale) dragonScale {
	HighestShineValue := 0b0
	var result dragonScale
	for _, scale := range scales {
		if result == (dragonScale{}) || scale.shineBinary > HighestShineValue {
			HighestShineValue = scale.shineBinary
			result = scale
		} else if scale.shineBinary == HighestShineValue {
			if scale.sumColours() < result.sumColours() {
				result = scale
			}
		}
	}
	return result
}

func largestGrouping(scales []dragonScale) (string, int) {
	groups := make(map[string][]dragonScale)
	largestGroupSize := 0
	largestGroupKey := ""
	result := 0
	for _, scale := range scales {
		groups[scale.group] = append(groups[scale.group], scale)
	}

	for groupName, groupScales := range groups {
		if len(groupScales) > largestGroupSize {
			largestGroupSize = len(groupScales)
			largestGroupKey = groupName
		}
	}

	for _, scale := range groups[largestGroupKey] {
		result += scale.value
	}

	return largestGroupKey, result
}

func main() {
	scales := processInput("input-part-1.txt")
	fmt.Println("Dominant greens count:", countDominantGreens(scales))
	scalesWithShine := processInput("input-part-2.txt")
	fmt.Println("Darkest scale amongst those with the brightest shine:", findDarkestShine(scalesWithShine).value)
	scalesWithGroups := processInput("input-part-3.txt")
	groupKey, groupValue := largestGrouping(scalesWithGroups)
	fmt.Println("Value of the largest similar scales:", groupValue, "from group", groupKey)
}
