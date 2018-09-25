package utils

import (
	"fmt"
	"time"
)

// LeafNodes global leaf nodes counter
var LeafNodes int64

func perft(depth int, pos *Board) {

	// // AssertTrue(// CheckBoard(pos))

	if depth == 0 {
		LeafNodes++
		return
	}

	var moveList MoveList
	GenerateAllMoves(pos, &moveList)

	for moveNum := 0; moveNum < moveList.Count; moveNum++ {
		if !MakeMove(pos, moveList.Moves[moveNum].Move) {
			continue
		}
		perft(depth-1, pos)
		TakeMove(pos)
	}

	return
}

// PerftTest run perft test
func PerftTest(depth int, pos *Board) {

	// // AssertTrue(// CheckBoard(pos))

	PrintBoard(pos)
	fmt.Printf("\nStarting Test To Depth:%d\n", depth)
	LeafNodes = 0

	start := time.Now()

	var moveList MoveList
	GenerateAllMoves(pos, &moveList)

	for moveNum := 0; moveNum < moveList.Count; moveNum++ {
		move := moveList.Moves[moveNum].Move
		if !MakeMove(pos, move) {
			continue
		}
		cumulativeNodes := LeafNodes
		perft(depth-1, pos)
		TakeMove(pos)
		oldNodes := LeafNodes - cumulativeNodes
		fmt.Printf("move %d : %s : %d\n", moveNum+1, PrintMove(move), oldNodes)
	}

	elapsed := time.Since(start)

	fmt.Printf("\nTest Complete : %d nodes visited in: %s \n", LeafNodes, elapsed)

	return
}
