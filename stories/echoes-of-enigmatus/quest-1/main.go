package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type instruction struct {
	A      int
	B      int
	C      int
	X      int
	Y      int
	Z      int
	M      int
	result int
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func (i instruction) calcEni() int {
	return eni(i.A, i.X, i.M) + eni(i.B, i.Y, i.M) + eni(i.C, i.Z, i.M)
}

func (i instruction) calcEniWithRemainders(maxRemainders int) int {
	return eniWithNRemainders(i.A, i.X, i.M, maxRemainders) + eniWithNRemainders(i.B, i.Y, i.M, maxRemainders) + eniWithNRemainders(i.C, i.Z, i.M, maxRemainders)
}

func lineToInstruction(line string) instruction {
	re := regexp.MustCompile(`([A-Z])=(\d+)`)
	matches := re.FindAllStringSubmatch(line, -1)
	var instr instruction
	for _, match := range matches {
		switch match[1] {
		case "A":
			instr.A = aToIIgnoreError(match[2])
		case "B":
			instr.B = aToIIgnoreError(match[2])
		case "C":
			instr.C = aToIIgnoreError(match[2])
		case "X":
			instr.X = aToIIgnoreError(match[2])
		case "Y":
			instr.Y = aToIIgnoreError(match[2])
		case "Z":
			instr.Z = aToIIgnoreError(match[2])
		case "M":
			instr.M = aToIIgnoreError(match[2])
		}
	}
	return instr
}

func (i instruction) String() string {
	return fmt.Sprintf("A=%d, B=%d, C=%d, X=%d, Y=%d, M=%d", i.A, i.B, i.C, i.X, i.Y, i.M)
}

func processInput(fileName string) []instruction {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	instructions := []instruction{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		instructions = append(instructions, lineToInstruction(line))
	}
	return instructions
}

func eni(n, exp, mod int) int {
	remainders := []string{}
	score := 1
	for range exp {
		remainder := (score * n) % mod
		score = remainder
		remainders = append([]string{strconv.Itoa(remainder)}, remainders...)
	}

	return aToIIgnoreError(strings.Join(remainders, ""))
}

func eniWithNRemainders(n, exp, mod, maxRemainders int) int {
	// Find the cycle in the score sequence
	score := 1
	seen := map[int]int{} // score -> iteration index
	scores := []int{}     // all remainders in order

	cycleStart := -1
	cycleLen := 0
	for i := 0; i < exp; i++ {
		remainder := (score * n) % mod
		score = remainder
		if prev, ok := seen[score]; ok {
			cycleStart = prev
			cycleLen = i - prev
			break
		}
		seen[score] = i
		scores = append(scores, remainder)
	}

	if cycleStart == -1 {
		// No cycle found within exp iterations, use scores directly
		result := []string{}
		start := len(scores) - maxRemainders
		if start < 0 {
			start = 0
		}
		for _, s := range scores[start:] {
			result = append([]string{strconv.Itoa(s)}, result...)
		}
		return aToIIgnoreError(strings.Join(result, ""))
	}

	// Use cycle to determine the last maxRemainders remainders
	// We need the remainders at positions exp-maxRemainders .. exp-1
	if maxRemainders > exp {
		maxRemainders = exp
	}
	lastN := []string{}
	for i := 0; i < maxRemainders; i++ {
		pos := exp - maxRemainders + i
		var val int
		if pos < cycleStart {
			val = scores[pos]
		} else {
			val = scores[cycleStart+(pos-cycleStart)%cycleLen]
		}
		lastN = append([]string{strconv.Itoa(val)}, lastN...)
	}

	return aToIIgnoreError(strings.Join(lastN, ""))
}

func partOne(instructions []instruction) int {
	currentMax := 0
	for _, instr := range instructions {
		result := instr.calcEni()
		if result > currentMax {
			currentMax = result
		}
	}

	return currentMax
}

func partTwo(instructions []instruction, maxRemainders int) int {
	currentMax := 0
	for _, instr := range instructions {
		result := instr.calcEniWithRemainders(maxRemainders)
		if result > currentMax {
			currentMax = result
		}
	}

	return currentMax
}

func eniSumRemainders(n, exp, mod int) int {
	score := 1
	seen := map[int]int{}
	scores := []int{}

	cycleStart := -1
	cycleLen := 0
	for i := 0; i < exp; i++ {
		remainder := (score * n) % mod
		score = remainder
		if prev, ok := seen[score]; ok {
			cycleStart = prev
			cycleLen = i - prev
			break
		}
		seen[score] = i
		scores = append(scores, remainder)
	}

	if cycleStart == -1 {
		// No cycle, just sum all scores
		total := 0
		for _, s := range scores {
			total += s
		}
		return total
	}

	// Example: eni(2,7,5)
	// scores = [2, 4, 3, 1] then cycle detected (score=2 seen at index 0)
	// cycleStart=0, cycleLen=4
	//
	// Sequence:  2, 4, 3, 1, | 2, 4, 3, 1, | 2, 4, 3, ...
	//            ^-----------^  ^-----------^
	//            pre-cycle=[]   cycle repeats
	//
	// preCycleSum = 0 (nothing before index 0)
	// cycleSum = 2+4+3+1 = 10
	// remaining = 7-0 = 7, fullCycles = 7/4 = 1, leftover = 7%4 = 3
	// leftoverSum = 2+4+3 = 9
	// total = 0 + 1*10 + 9 = 19 ✓

	// Sum the pre-cycle part
	preCycleSum := 0
	for i := 0; i < cycleStart; i++ {
		preCycleSum += scores[i]
	}

	// Sum of one full cycle
	cycleSum := 0
	for i := cycleStart; i < cycleStart+cycleLen; i++ {
		cycleSum += scores[i]
	}

	// How many iterations remain after the pre-cycle
	// e.g. eni(2,7,5): remaining = 7-0 = 7 iterations to account for
	// fullCycles = 7/4 = 1 complete repetition of [2,4,3,1]
	// leftover = 7%4 = 3 extra values from the start of the next cycle [2,4,3]
	remaining := exp - cycleStart
	fullCycles := remaining / cycleLen
	leftover := remaining % cycleLen

	// Sum of the leftover partial cycle
	// These are the first few values of an incomplete final cycle
	leftoverSum := 0
	for i := 0; i < leftover; i++ {
		leftoverSum += scores[cycleStart+i]
	}

	return preCycleSum + fullCycles*cycleSum + leftoverSum
}

func (i instruction) calcEniSumRemainders() int {
	return eniSumRemainders(i.A, i.X, i.M) + eniSumRemainders(i.B, i.Y, i.M) + eniSumRemainders(i.C, i.Z, i.M)
}

func partThree(instructions []instruction, maxRemainders int) int {
	currentMax := 0
	for _, instr := range instructions {
		result := instr.calcEniSumRemainders()
		if result > currentMax {
			currentMax = result
		}
	}
	return currentMax
}

func main() {
	instructions := processInput("input-part-1.txt")
	fmt.Println("Part One largest result:", partOne(instructions))
	instructionsTwo := processInput("input-part-2.txt")
	fmt.Println("Part Two largest result with 5 remainders:", partTwo(instructionsTwo, 5))
	instructionsThree := processInput("input-part-3.txt")
	fmt.Println("Part Three largest sum of all remainders:", partThree(instructionsThree, 5))
}
