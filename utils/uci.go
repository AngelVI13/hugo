package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseGo parse UCI go command
// sample go commange is below
//            -white time ms -black time  -b/w increment ms -movetime ms
// go depth 6 wtime 180000 btime 100000 binc 1000 winc 1000 movetime 1000 movestogo 40
func ParseGo(line string, info *SearchInfo, pos *Board) {
	depth := -1
	movesToGo := 30
	moveTime := -1
	timeInt := -1
	inc := 0
	info.TimeSet = false

	if strings.Contains(line, "infinite") {

	}

	if strings.Contains(line, "binc") && pos.side == Black {
		incStr1 := RemoveStringToTheLeftOfMarker(line, "binc ")
		incStr2 := RemoveStringToTheRightOfMarker(incStr1, " ")
		inc, _ = strconv.Atoi(incStr2)
	}

	if strings.Contains(line, "winc") && pos.side == White {
		incStr1 := RemoveStringToTheLeftOfMarker(line, "winc ")
		incStr2 := RemoveStringToTheRightOfMarker(incStr1, " ")
		inc, _ = strconv.Atoi(incStr2)
	}

	if strings.Contains(line, "wtime") && pos.side == White {
		timeStr1 := RemoveStringToTheLeftOfMarker(line, "wtime ")
		timeStr2 := RemoveStringToTheRightOfMarker(timeStr1, " ")
		fmt.Println(timeStr1, "|", timeStr2)
		timeInt, _ = strconv.Atoi(timeStr2)
	}

	if strings.Contains(line, "btime") && pos.side == Black {
		timeStr1 := RemoveStringToTheLeftOfMarker(line, "btime ")
		timeStr2 := RemoveStringToTheRightOfMarker(timeStr1, " ")
		timeInt, _ = strconv.Atoi(timeStr2)
	}

	if strings.Contains(line, "movestogo") {
		movesToGoStr1 := RemoveStringToTheLeftOfMarker(line, "movestogo ")
		movesToGoStr2 := RemoveStringToTheRightOfMarker(movesToGoStr1, " ")
		fmt.Println(movesToGoStr1, "|", movesToGoStr2)
		movesToGo, _ = strconv.Atoi(movesToGoStr2)
	}

	if strings.Contains(line, "movetime") {
		moveTimeStr1 := RemoveStringToTheLeftOfMarker(line, "movetime ")
		moveTimeStr2 := RemoveStringToTheRightOfMarker(moveTimeStr1, " ")
		moveTime, _ = strconv.Atoi(moveTimeStr2)
	}

	if strings.Contains(line, "depth") {
		depthStr1 := RemoveStringToTheLeftOfMarker(line, "depth ")
		depthStr2 := RemoveStringToTheRightOfMarker(depthStr1, " ")
		depth, _ = strconv.Atoi(depthStr2)
	}

	if moveTime != -1 {
		timeInt = moveTime
		movesToGo = 1
	}

	info.StartTime = time.Now()
	info.Depth = depth

	if timeInt != -1 {
		info.TimeSet = true
		timeInt /= movesToGo
		// to be on the safe side we remove 50ms from this value
		timeInt -= 50
		stopTimeInSeconds := (timeInt + inc) // find stop time in miliseconds
		info.StopTime = stopTimeInSeconds
	}

	if depth == -1 {
		info.Depth = MaxDepth
	}

	fmt.Printf("time:%d start:%s stop:%d depth:%d timeset:%t\n", timeInt, info.StartTime, info.StopTime,
		info.Depth, info.TimeSet)

	SearchPosition(pos, info)
}

// ParsePosition parse UCI position
// the expected formats are 'position fen **' or 'position startpos'
func ParsePosition(lineIn string, pos *Board) {
	if strings.Contains(lineIn, "startpos") {
		ParseFen(StartFen, pos)
	} else {
		if strings.Contains(lineIn, "fen") {
			startStr := "fen "
			fen := RemoveStringToTheLeftOfMarker(lineIn, startStr)
			ParseFen(fen, pos)
		} else {
			ParseFen(StartFen, pos)
		}
	}

	movesStr := "moves "
	movesIdx := strings.Index(lineIn, movesStr)
	if movesIdx != -1 {
		fullMovesStr := RemoveStringToTheLeftOfMarker(lineIn, movesStr)
		moveSlice := strings.Split(fullMovesStr, " ")
		for i := range moveSlice {
			move := ParseMove(moveSlice[i], pos)
			if move == NoMove {
				break
			}
			MakeMove(pos, move)
			pos.ply = 0
		}
	}
	// PrintBoard(pos)
}

const (
	// InputBuffer max characters received
	InputBuffer int = 400 * 6
)

// UciLoop main UCI loop
func UciLoop(pos *Board, info *SearchInfo) {
	fmt.Printf("id name %s\n", Name)
	fmt.Printf("id author AngelVI\n")
	fmt.Println("uciok")

	line := ""

	// InitHashTable(&pos.HashTable)

	for {
		line, _ = GetInput("")
		if len(line) < 2 {
			continue
		}
		// Remove leading and trailinig whitespaces
		line = strings.Trim(line, " ")

		if strings.Contains(line, "isready") {
			fmt.Println("readyok")
			continue
		} else if strings.Contains(line, "position") {
			ParsePosition(line, pos)
		} else if strings.Contains(line, "ucinewgame") {
			ParsePosition("position startpos\n", pos)
		} else if strings.Contains(line, "go") {
			ParseGo(line, info, pos)
		} else if strings.Contains(line, "quit") {
			info.Quit = true
			break
		} else if strings.Contains(line, "uci") {
			fmt.Printf("id name %s\n", Name)
			fmt.Printf("id author AngelVI\n")
			fmt.Println("uciok")
		}

		if info.Quit {
			break
		}
	}
}
