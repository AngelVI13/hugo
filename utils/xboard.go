package utils

import (
	"fmt"
	inpututils "local/input-utils"
	stringutils "local/string-utils"
	"strconv"
	"strings"
	"time"
)

// ThreeFoldRepetition Detects how many repetitions for a given position
func ThreeFoldRepetition(pos *Board) int {
	r := 0

	for i := 0; i < pos.histPly; i++ {
		if pos.history[i].posKey == pos.posKey {
			r++
		}
	}
	return r
}

// DrawMaterial determine if position is a draw
func DrawMaterial(pos *Board) bool {
	// if there are pawns on the board the one of the sides can get mated
	if pos.pieceNum[WhitePawn] != 0 || pos.pieceNum[BlackPawn] != 0 {
		return false
	}
	// if there are major pieces on the board the one of the sides can get mated
	if pos.pieceNum[WhiteQueen] != 0 || pos.pieceNum[BlackQueen] != 0 || pos.pieceNum[WhiteRook] != 0 || pos.pieceNum[BlackRook] != 0 {
		return false
	}
	if pos.pieceNum[WhiteBishop] > 1 || pos.pieceNum[BlackBishop] > 1 {
		return false
	}
	if pos.pieceNum[WhiteKnight] > 1 || pos.pieceNum[BlackKnight] > 1 {
		return false
	}
	if pos.pieceNum[WhiteKnight] != 0 && pos.pieceNum[WhiteBishop] != 0 {
		return false
	}
	if pos.pieceNum[BlackKnight] != 0 && pos.pieceNum[BlackBishop] != 0 {
		return false
	}

	return true
}

// CheckResult is called everytime a move is made this functio is called to check if the game is a draw
func CheckResult(pos *Board) bool {

	if pos.fiftyMove > 100 {
		fmt.Printf("1/2-1/2 {fifty move rule (claimed by Hugo)}\n")
		return true
	}

	if ThreeFoldRepetition(pos) >= 2 {
		fmt.Printf("1/2-1/2 {3-fold repetition (claimed by Hugo)}\n")
		return true
	}

	if DrawMaterial(pos) == true {
		fmt.Printf("1/2-1/2 {insufficient material (claimed by Hugo)}\n")
		return true
	}

	var moveList MoveList
	GenerateAllMoves(pos, &moveList)

	found := 0
	for moveNum := 0; moveNum < moveList.Count; moveNum++ {

		if !MakeMove(pos, moveList.Moves[moveNum].Move) {
			continue
		}
		found++
		TakeMove(pos)
		break
	}

	// we have legal moves -> game is not over
	if found != 0 {
		return false
	}

	InCheck := IsSquareAttacked(pos.kingSquare[pos.side], pos.side^1, pos)

	if InCheck == true {
		if pos.side == White {
			fmt.Printf("0-1 {black mates (claimed by Hugo)}\n")
			return true
		}
		fmt.Printf("0-1 {white mates (claimed by Hugo)}\n")
		return true
	}
	// not in check but no legal moves left -> stalemate
	fmt.Printf("\n1/2-1/2 {stalemate (claimed by Hugo)}\n")
	return true

}

// XBoardLoop main loop
func XBoardLoop(pos *Board, info *SearchInfo) {
	// before we do anything else, we first send back the features that this engine supports
	fmt.Printf("feature ping=1 setboard=1 colors=0 usermove=1\n")
	fmt.Printf("feature done=1\n")

	// initialize variables
	info.GameMode = XBoardMode
	info.PostThinking = true

	depth := -1
	movesToGo := [2]int{30, 30}
	moveTime := -1
	timeInt := -1
	inc := 0
	engineSide := Both
	timeLeft := 0
	seconds := 0
	movesPerSession := 0
	move := NoMove
	// score := 0

	command := ""

	// Setting default values
	engineSide = Black
	ParseFen(StartFen, pos)
	depth = -1
	timeInt = -1

	for {

		if pos.side == engineSide && CheckResult(pos) == false {
			info.StartTime = time.Now()
			info.Depth = depth

			if timeInt != -1 {
				info.TimeSet = true
				timeInt /= movesToGo[pos.side]
				timeInt -= 50
				stopTimeInSeconds := (timeInt + inc) // find stop time in miliseconds
				info.StopTime = stopTimeInSeconds
			}

			if depth == -1 || depth > MaxDepth {
				info.Depth = MaxDepth
			}

			fmt.Printf("time:%d start:%s stop:%d depth:%d timeset:%t movestogo:%d mps:%d\n", timeInt, info.StartTime,
				info.StopTime, info.Depth, info.TimeSet, movesToGo[pos.side], movesPerSession)
			// think
			SearchPosition(pos, info)

			if movesPerSession != 0 {
				movesToGo[pos.side^1]--
				if movesToGo[pos.side^1] < 1 {
					movesToGo[pos.side^1] = movesPerSession
				}
			}
		}

		command, _ = inpututils.GetInput("")
		if len(command) < 2 {
			continue
		}

		// if (!fgets(inBuf, 80, stdin))
		// continue;

		if strings.Contains(command, "quit") {
			break
		}

		if strings.Contains(command, "force") {
			engineSide = Both
			continue
		}

		if strings.Contains(command, "protover") {
			fmt.Printf("feature ping=1 setboard=1 colors=0 usermove=1\n")
			fmt.Printf("feature done=1\n")
			continue
		}

		if strings.Contains(command, "sd") {
			depthStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "sd ")
			depthStr2 := stringutils.RemoveStringToTheRightOfMarker(depthStr1, " ")
			depth, _ = strconv.Atoi(depthStr2)
			fmt.Printf("DEBUG depth: %d\n", depth)
			continue
		}

		if strings.Contains(command, "st") {
			moveTimeStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "st ")
			moveTimeStr2 := stringutils.RemoveStringToTheRightOfMarker(moveTimeStr1, " ")
			moveTime, _ = strconv.Atoi(moveTimeStr2)
			fmt.Printf("DEBUG movetime: %d\n", moveTime)
			continue
		}

		if strings.Contains(command, "time") {
			timeIntStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "time ")
			timeIntStr2 := stringutils.RemoveStringToTheRightOfMarker(timeIntStr1, " ")
			timeInt, _ = strconv.Atoi(timeIntStr2)
			timeInt *= 10
			fmt.Printf("DEBUG time:%d\n", timeInt)
		}

		if strings.Contains(command, "level") {
			seconds = 0
			moveTime = -1
			inputStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "level ")
			inputStrSlice := strings.Split(inputStr1, " ")
			if len(inputStrSlice) != 3 {
				fmt.Println(command)
				panic("Incorrect level command!")
			}

			if strings.Contains(command, ":") {
				movesPerSession, _ = strconv.Atoi(inputStrSlice[0])
				timeInfoSlice := strings.Split(inputStrSlice[1], ":")
				if len(timeInfoSlice) != 2 {
					fmt.Println(command)
					panic("Incorrect level command!")
				}
				timeLeft, _ = strconv.Atoi(timeInfoSlice[0])
				seconds, _ = strconv.Atoi(timeInfoSlice[1])
				inc, _ = strconv.Atoi(inputStrSlice[2])
				fmt.Printf("DEBUG level with :\n")
			} else {
				movesPerSession, _ = strconv.Atoi(inputStrSlice[0])
				timeLeft, _ = strconv.Atoi(inputStrSlice[1])
				inc, _ = strconv.Atoi(inputStrSlice[2])
				fmt.Printf("DEBUG level without :\n")
			}

			timeLeft *= 60000
			timeLeft += seconds * 1000
			movesToGo[0] = 30
			movesToGo[1] = 30
			if movesPerSession != 0 {
				movesToGo[0] = movesPerSession
				movesToGo[1] = movesPerSession
			}
			timeInt = -1
			fmt.Printf("DEBUG level timeLeft:%d movesToGo:%d inc:%d mps:%d\n", timeLeft, movesToGo[0], inc, movesPerSession)
			continue
		}

		// the protocol might send ping 3 and we need to reply with pong 3
		if strings.Contains(command, "ping") {
			pingStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "ping ")
			pingStr2 := stringutils.RemoveStringToTheRightOfMarker(pingStr1, " ")
			pingNum, _ := strconv.Atoi(pingStr2)
			fmt.Printf("pong %d\n", pingNum)
			continue
		}

		if strings.Contains(command, "new") {
			engineSide = Black
			ParseFen(StartFen, pos)
			depth = -1
			continue
		}

		if strings.Contains(command, "setboard") {
			engineSide = Both
			startStr := "setboard "
			fen := stringutils.RemoveStringToTheLeftOfMarker(command, startStr)
			ParseFen(fen, pos)
			continue
		}

		if strings.Contains(command, "go") {
			engineSide = pos.side
			continue
		}

		if strings.Contains(command, "usermove") {
			movesToGo[pos.side]--
			moveStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "usermove ")
			moveStr2 := stringutils.RemoveStringToTheRightOfMarker(moveStr1, " ")
			move = ParseMove(moveStr2, pos)
			if move == NoMove {
				continue
			}
			MakeMove(pos, move)
			pos.ply = 0
		}
	}
}

// ConsoleLoop console loop for playing through the console
func ConsoleLoop(pos *Board, info *SearchInfo) {
	fmt.Printf("Welcome to Hugo In Console Mode!\n")
	fmt.Printf("Type help for commands\n\n")

	info.GameMode = ConsoleMode
	info.PostThinking = true

	depth := MaxDepth
	moveTime := 3000 // 3 seconds move time
	engineSide := Both
	move := NoMove

	engineSide = Black
	ParseFen(StartFen, pos)

	command := ""

	for {
		if pos.side == engineSide && CheckResult(pos) == false {
			info.StartTime = time.Now()
			info.Depth = depth

			if moveTime != 0 {
				info.TimeSet = true
				info.StopTime = moveTime
			}

			SearchPosition(pos, info)
		}

		command, _ = inpututils.GetInput("\nHugo > ")
		if len(command) < 2 {
			continue
		}

		if strings.Contains(command, "help") {
			fmt.Printf("Commands:\n")
			fmt.Printf("quit - quit game\n")
			fmt.Printf("force - computer will not think\n")
			fmt.Printf("print - show board\n")
			fmt.Printf("post - show thinking\n")
			fmt.Printf("nopost - do not show thinking\n")
			fmt.Printf("new - start new game\n")
			fmt.Printf("mirror - prints the current position and then it's mirrored image.\n")
			fmt.Printf("setboard x - set position to fen x\n")
			fmt.Printf("go - set computer thinking\n")
			fmt.Printf("depth x - set depth to x\n")
			fmt.Printf("time x - set thinking time to x seconds (depth still applies if set)\n")
			fmt.Printf("view - show current depth and moveTime settings\n")
			fmt.Printf("showline - show current move line so far\n")
			fmt.Printf("** note ** - to reset time and depth, set to 0\n")
			fmt.Printf("enter moves using b7b8q notation\n\n\n")
			continue
		}

		if strings.Contains(command, "mirror") {
			PrintBoard(pos)
			fmt.Printf("Eval:%d\n", EvalPosition(pos))
			MirrorBoard(pos)
			PrintBoard(pos)
			fmt.Printf("Eval:%d\n", EvalPosition(pos))
			MirrorBoard(pos)
			continue
		}

		if strings.Contains(command, "setboard") {
			engineSide = Both
			startStr := "setboard "
			fen := stringutils.RemoveStringToTheLeftOfMarker(command, startStr)
			ParseFen(fen, pos)
			continue
		}

		if strings.Contains(command, "quit") {
			info.Quit = true
			break
		}

		if strings.Compare(command, "post") == 0 {
			info.PostThinking = true
			continue
		}

		if strings.Contains(command, "print") {
			PrintBoard(pos)
			continue
		}

		if strings.Compare(command, "nopost") == 0 {
			info.PostThinking = false
			continue
		}

		if strings.Contains(command, "force") {
			engineSide = Both
			continue
		}

		if strings.Contains(command, "view") {
			if depth == MaxDepth {
				fmt.Printf("depth not set ")
			} else {
				fmt.Printf("depth %d", depth)
			}

			if moveTime != 0 {
				fmt.Printf(" moveTime %ds\n", moveTime/1000)
			} else {
				fmt.Printf(" moveTime not set\n")
			}
			continue
		}

		if strings.Contains(command, "showline") {
			fmt.Println(GetBookMove(pos))
			continue
		}

		if strings.Contains(command, "depth") {
			depthStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "depth ")
			depthStr2 := stringutils.RemoveStringToTheRightOfMarker(depthStr1, " ")
			depth, _ = strconv.Atoi(depthStr2)
			if depth == 0 {
				depth = MaxDepth
			}
			continue
		}

		if strings.Contains(command, "time") {
			moveTimeStr1 := stringutils.RemoveStringToTheLeftOfMarker(command, "time ")
			moveTimeStr2 := stringutils.RemoveStringToTheRightOfMarker(moveTimeStr1, " ")
			moveTime, _ = strconv.Atoi(moveTimeStr2)
			moveTime *= 1000
			continue
		}

		if strings.Contains(command, "new") {
			engineSide = Black
			ParseFen(StartFen, pos)
			continue
		}

		if strings.Contains(command, "go") {
			engineSide = pos.side
			continue
		}

		move = ParseMove(command, pos)
		if move == NoMove {
			fmt.Printf("Command unknown:%s\n", command)
			continue
		}
		MakeMove(pos, move)
		pos.ply = 0
	}
}
