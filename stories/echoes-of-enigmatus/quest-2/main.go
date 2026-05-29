package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type instruction struct {
	op    string
	id    int
	left  node
	right node
}

type node struct {
	rank  int
	value rune
	left  *node
	right *node
}

func (i instruction) String() string {
	return fmt.Sprintf("%s id=%d left=[%d,%c] right=[%d,%c]", i.op, i.id, i.left.rank, i.left.value, i.right.rank, i.right.value)
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func (n *node) swap(other *node) {
	n.rank, other.rank = other.rank, n.rank
	n.value, other.value = other.value, n.value
}

func (n *node) String() string {
	return fmt.Sprintf("rank=%d value=%c", n.rank, n.value)
}

func lineToInstruction(line string) instruction {
	re := regexp.MustCompile(`^(\w+)\s+id=(\d+)\s+left=\[(\d+),(.)\]\s+right=\[(\d+),(.)\]$`)
	swapRe := regexp.MustCompile(`^SWAP\s+(\d+)$`)
	matches := re.FindStringSubmatch(line)
	swapMatched := swapRe.FindStringSubmatch(line)
	if matches == nil && swapMatched == nil {
		fmt.Println("no match:", line)
		panic("invalid instruction format")
	}

	if swapMatched != nil {
		return instruction{
			op: "SWAP",
			id: aToIIgnoreError(swapMatched[1]),
		}
	}

	return instruction{
		op: matches[1],
		id: aToIIgnoreError(matches[2]),
		left: node{
			rank:  aToIIgnoreError(matches[3]),
			value: rune(matches[4][0]),
		},
		right: node{
			rank:  aToIIgnoreError(matches[5]),
			value: rune(matches[6][0]),
		},
	}
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
		instruction := lineToInstruction(line)
		instructions = append(instructions, instruction)
	}
	return instructions
}

func (n *node) insert(candidate *node) {
	if candidate.rank < n.rank {
		if n.left == nil {
			n.left = candidate
		} else {
			n.left.insert(candidate)
		}
	} else {
		if n.right == nil {
			n.right = candidate
		} else {
			n.right.insert(candidate)
		}
	}
}

func (n *node) walk(levels map[int][]rune, depth int) {
	if n == nil {
		return
	}
	levels[depth] = append(levels[depth], n.value)
	n.left.walk(levels, depth+1)
	n.right.walk(levels, depth+1)
}

func wordAtMostFrequentDepth(root *node) string {
	levels := map[int][]rune{}
	root.walk(levels, 0)

	maxCount := 0
	maxDepth := 0
	for depth, nodes := range levels {
		if len(nodes) > maxCount {
			maxCount = len(nodes)
			maxDepth = depth
		}
	}
	return string(levels[maxDepth])
}

func partOne(instructions []instruction) string {
	startPoint := instructions[0]
	instructions = instructions[1:]

	leftTree := &node{
		rank:  startPoint.left.rank,
		value: startPoint.left.value,
		left:  nil,
		right: nil,
	}
	rightTree := &node{
		rank:  startPoint.right.rank,
		value: startPoint.right.value,
		left:  nil,
		right: nil,
	}

	for _, instr := range instructions {
		leftTree.insert(&node{
			rank:  instr.left.rank,
			value: instr.left.value,
		})
		rightTree.insert(&node{
			rank:  instr.right.rank,
			value: instr.right.value,
		})
	}

	return wordAtMostFrequentDepth(leftTree) + wordAtMostFrequentDepth(rightTree)
}

func partTwo(instructions []instruction) string {
	startPoint := instructions[0]
	instructions = instructions[1:]

	leftTree := &node{
		rank:  startPoint.left.rank,
		value: startPoint.left.value,
		left:  nil,
		right: nil,
	}
	rightTree := &node{
		rank:  startPoint.right.rank,
		value: startPoint.right.value,
		left:  nil,
		right: nil,
	}

	nodesByID := map[int][]*node{
		startPoint.id: {leftTree, rightTree},
	}

	for _, instr := range instructions {
		if instr.op == "SWAP" {
			if swapNodes, ok := nodesByID[instr.id]; ok {
				// fmt.Printf("Performing swap for id=%d: %v\n", instr.id, swapNodes)
				swapNodes[0].swap(swapNodes[1])
			}
			continue
		}

		leftCanidate := &node{
			rank:  instr.left.rank,
			value: instr.left.value,
		}
		rightCanidate := &node{
			rank:  instr.right.rank,
			value: instr.right.value,
		}
		nodesByID[instr.id] = append(nodesByID[instr.id], leftCanidate)
		nodesByID[instr.id] = append(nodesByID[instr.id], rightCanidate)
		leftTree.insert(leftCanidate)
		rightTree.insert(rightCanidate)
	}

	return wordAtMostFrequentDepth(leftTree) + wordAtMostFrequentDepth(rightTree)
}

func main() {
	instructions := processInput("input-part-1.txt")
	fmt.Println("Part One largest result:", partOne(instructions))
	instructions2 := processInput("input-part-2.txt")
	fmt.Println("Part One largest result:", partTwo(instructions2))
}
