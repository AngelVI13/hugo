package utils

import (
	"log"
	"fmt"
	"time"
	"sync"
	// "github.com/jinzhu/copier"
)

const (
	//Infinite infinite value
	Infinite int = 30000
	// IsMate mate value
	IsMate int = Infinite - MaxDepth
)

// CheckUp check if time up or interrupt from GUI
func CheckUp(info *SearchInfo) {
	// check if time up or interrupt from GUI
	// fmt.Println(elapsed, info.StopTime, elapsed.After(info.StopTime))
	elapsedTime := time.Since(info.StartTime).Seconds() * 1000 // get elapsed time in ms
	if info.TimeSet == true && elapsedTime > float64(info.StopTime) {
		info.stopped = true
	}
	// if we received something from the gui -> set stopped/quit to true
	// ReadInput(info)
}

// PickNextMove reorders the moves so that the highest scoring move is picked next
func PickNextMove(moveNum int, moveList *MoveList) {
	var tempMove Move
	bestScore := 0
	bestNum := moveNum

	for i := moveNum; i < moveList.Count; i++ {
		if moveList.Moves[i].score > bestScore {
			bestScore = moveList.Moves[i].score
			bestNum = i
		}
	}
	// swap move with better move
	tempMove = moveList.Moves[moveNum]
	moveList.Moves[moveNum] = moveList.Moves[bestNum]
	moveList.Moves[bestNum] = tempMove
}

// SearchPosition searches a given position
func SearchPosition(pos *Board, info *SearchInfo) int {
	// ... iterative deepening, search init
	// for depth = 1 to maxDepth
	// 		search with alphaBeta if you have enough time left
	// you do not search to maxDepth from the start but instead search
	// with depth 1, then 2, then 3 etc, because you first identify
	// the principle variation or the potentially good moves and in this
	// way when you search again with more depth you can easily eliminate
	// a lot of bad nodes automatically
	// todo good test position for deep mate: 8/8/2K5/4b3/3k4/5Q2/6p1/8 w - - 0 6

	// if we can perform a book move, do that first, otherwise perform search
	bestMove := GetBookMove(pos)
	if bestMove != 0 {
		PerformMove(pos, info, bestMove)
		return 0
	}
	ClearForSearch(pos, info)

	// Create waitgroup to manage goroutines
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)

		// copy pos & info
		pos_copy := *pos
		info_copy := *info

        go func() { 
        	IterativeSearch(&pos_copy, &info_copy)
        	wg.Done()
        }()
	}

	info.IsMainThread = true
	bestMove = IterativeSearch(pos, info)

	// Wait for all goroutines to finish. We are not interested in their output
	// We only run them in order to faster populate the HashTable
	wg.Wait()

	// todo this is only here for debugging
	if bestMove == NoMove {
		panic("No bestMove after IterativeSearch!")
	}

	PerformMove(pos, info, bestMove)

	return 0
}

func IterativeSearch(pos *Board, info *SearchInfo) int {
	// do normal move search
	bestMove := NoMove
	bestScore := -Infinite

	for currentDepth := 1; currentDepth <= info.Depth; currentDepth++ {
		//                    *alpha     *beta
		bestScore = AlphaBeta(-Infinite, Infinite, currentDepth, pos, info, true)

		// make sure to update PV line info into the current board object
		pvMoves := GetPvLine(pos, currentDepth)
		
		if info.stopped == true {
			break
		}

		// we only print info on our main thread
		if !info.IsMainThread {
			continue
		}

		moveTime := int64(time.Since(info.StartTime).Seconds() * 1000) // the UCI protocol expects milliseconds
		if info.GameMode == UciMode {
			fmt.Printf("info score cp %d depth %d nodes %d time %d ", bestScore, currentDepth, info.nodes, moveTime)
			log.Printf("info score cp %d depth %d nodes %d time %d ", bestScore, currentDepth, info.nodes, moveTime)
		} else if info.GameMode == XBoardMode && info.PostThinking == true {
			moveTime *= 10
			fmt.Printf("%d %d %d %d", currentDepth, bestScore, moveTime, info.nodes)
			log.Printf("%d %d %d %d", currentDepth, bestScore, moveTime, info.nodes)
		} else if info.PostThinking == true {
			fmt.Printf("score:%d depth:%d nodes:%d time:%d(ms)", bestScore, currentDepth, info.nodes, moveTime)
			log.Printf("score:%d depth:%d nodes:%d time:%d(ms)", bestScore, currentDepth, info.nodes, moveTime)
		}
		if info.GameMode == UciMode || info.PostThinking == true {
			// Print the principle variation
			pvMoves = GetPvLine(pos, currentDepth)
			fmt.Printf("pv")
			log.Printf("pv")
			for pvNum := 0; pvNum < pvMoves; pvNum++ {
				fmt.Printf(" %s", PrintMove(pos.PvArray[pvNum]))
				log.Printf(" %s", PrintMove(pos.PvArray[pvNum]))
			}
			fmt.Println()
			log.Println()
			// fmt.Printf("Ordering: %.2f\n", info.failHighFirst/info.failHigh)
		}
	}

	bestMove = pos.PvArray[0]
	return bestMove

}

// PerformMove performs the best found move from search or book
func PerformMove(pos *Board, info *SearchInfo, bestMove int) {
	if info.GameMode == UciMode {
		fmt.Printf("bestmove %s\n", PrintMove(bestMove))
		log.Printf("bestmove %s\n", PrintMove(bestMove))
	} else if info.GameMode == XBoardMode {
		fmt.Printf("move %s\n", PrintMove(bestMove))
		log.Printf("move %s\n", PrintMove(bestMove))
		MakeMove(pos, bestMove)
	} else {
		fmt.Printf("\n\n***!! Hugo makes move %s !!***\n\n", PrintMove(bestMove))
		log.Printf("\n\n***!! Hugo makes move %s !!***\n\n", PrintMove(bestMove))
		MakeMove(pos, bestMove)
		PrintBoard(pos)
	}
}

// ClearForSearch clear all info for search
func ClearForSearch(pos *Board, info *SearchInfo) {
	for i := 0; i < 13; i++ {
		for j := 0; j < BoardSquareNum; j++ {
			pos.searchHistory[i][j] = 0
		}
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < MaxDepth; j++ {
			pos.searchKillers[i][j] = 0
		}
	}

	

	pos.ply = 0
	info.stopped = false
	info.nodes = 0
	info.failHigh = 0
	info.failHighFirst = 0
	info.IsMainThread = false
}

// Quiescence searches untill a quiet position in order to eliminate the horizon effect
func Quiescence(alpha, beta int, pos *Board, info *SearchInfo) int {
	// // AssertTrue(// CheckBoard(pos))

	// perform check on every 2047 nodes
	if (info.nodes & 2047) == 0 {
		CheckUp(info)
	}

	info.nodes++

	if IsRepetition(pos) || pos.fiftyMove >= 100 {
		return 0
	}

	if pos.ply > MaxDepth-1 {
		return EvalPosition(pos)
	}

	score := EvalPosition(pos)

	// if score is already bigger than beta then we can return
	// since we wont make any score bigger than beta
	// -- because the player has not yet moved and the score is already
	// more than beta, then very likely the score will be much more after
	// the move is made so its much more greater than beta then....
	// so thats why we can assume that beta cutoff here is a good estimate
	if score >= beta {
		return beta
	}

	// if score is greater than alpha then update alpha
	if score > alpha {
		alpha = score
	}

	var moveList MoveList
	GenerateAllCaptures(pos, &moveList)

	legalMoves := 0
	score = -Infinite
	// PvMove := ProbeHashEntry(pos)

	for moveNum := 0; moveNum < moveList.Count; moveNum++ {

		PickNextMove(moveNum, &moveList)

		if !MakeMove(pos, moveList.Moves[moveNum].Move) {
			continue
		}

		legalMoves++
		score = -Quiescence(-beta, -alpha, pos, info)
		TakeMove(pos)

		if info.stopped == true {
			return 0
		}

		if score > alpha {
			if score >= beta {
				if legalMoves == 1 {
					info.failHighFirst++
				}
				info.failHigh++
				return beta
			}
			alpha = score
		}
	}
	return alpha
}

// AlphaBeta performs alpha beta search
func AlphaBeta(alpha, beta, depth int, pos *Board, info *SearchInfo, doNull bool) int {
	// // AssertTrue(// CheckBoard(pos))

	if depth <= 0 {
		return Quiescence(alpha, beta, pos, info)
	}

	// perform check on every 2047 nodes
	if (info.nodes & 2047) == 0 {
		CheckUp(info)
	}

	info.nodes++

	// Detect draw cases
	if (IsRepetition(pos) || pos.fiftyMove >= 100) && (pos.ply != 0) {
		return 0
	}

	// If we are at max depth, return eval
	if pos.ply > MaxDepth-1 {
		return EvalPosition(pos)
	}

	// If we are in check, increase the search depth in order to make sure we dont
	// get mated in the next couple of moves
	inCheck := IsSquareAttacked(pos.kingSquare[pos.side], pos.side^1, pos)
	if inCheck == true {
		depth++
	}

	score := -Infinite
	pvMove := NoMove

	if ProbeHashEntry(pos, &pvMove, &score, alpha, beta, depth) == true {
		return score
	}

	// if normal alphabeta was called rather than the null move call
	// if we are not in check (otherwise the king can get captured)
	// if also we have made at least 1 move into the search
	// bigPieceNum protects against a zugzwang situation in which this might evaluate to a draw whereas we are actually loosing
	// only perform null moves at depth >= 4
	if doNull == true && inCheck == false && pos.ply != 0 && pos.bigPieceNum[pos.side] > 0 && depth >= 4 {
		MakeNullMove(pos)
		score = -AlphaBeta(-beta, -beta+1, depth-4, pos, info, false)
		TakeNullMove(pos)
		if info.stopped == true {
			return 0
		}
		if score >= beta && abs(score) < IsMate {
			info.nullCut++
			return beta
		}
	}

	var moveList MoveList
	GenerateAllMoves(pos, &moveList)

	legalMoves := 0
	oldAlpha := alpha
	bestMove := NoMove
	bestScore := -Infinite
	score = -Infinite // reset score here after null move pruning

	// if we have a PvMove then score it to 2mil so that it is searched first
	if pvMove != NoMove {
		for moveNum := 0; moveNum < moveList.Count; moveNum++ {
			if moveList.Moves[moveNum].Move == pvMove {
				moveList.Moves[moveNum].score = 2000000
				break
			}
		}
	}

	for moveNum := 0; moveNum < moveList.Count; moveNum++ {
		// todo investigate sorting the whole movelist instead of finding the best move and returning it
		PickNextMove(moveNum, &moveList)

		if !MakeMove(pos, moveList.Moves[moveNum].Move) {
			continue
		}

		legalMoves++
		// search next depth with flipped alpha and beta since the sides
		// will change
		score = -AlphaBeta(-beta, -alpha, depth-1, pos, info, true)
		TakeMove(pos)

		if info.stopped == true {
			return 0
		}

		if score > bestScore {
			bestScore = score
			bestMove = moveList.Moves[moveNum].Move
			if score > alpha {
				if score >= beta {
					// beta cutoff found -> return beta
					if legalMoves == 1 {
						info.failHighFirst++
					}
					info.failHigh++

					// finds non captures that are causing beta cuttoffs
					if (moveList.Moves[moveNum].Move & MoveFlagCapture) == 0 {
						// move previous killer move to 1 index and set
						// new move to be the first element of the killers slice
						pos.searchKillers[1][pos.ply] = pos.searchKillers[0][pos.ply]
						pos.searchKillers[0][pos.ply] = moveList.Moves[moveNum].Move
					}

					StoreHashEntry(pos, bestMove, beta, HFBeta, depth)

					return beta
				}
				alpha = score

				if (moveList.Moves[moveNum].Move & MoveFlagCapture) == 0 {
					// move previous killer move to 1 index and set
					// new move to be the first element of the killers slice
					pos.searchHistory[pos.Pieces[FromSq(bestMove)]][ToSq(bestMove)] += depth // prioritizes depth
				}
			}
		}
	}

	if legalMoves == 0 {
		// if the enemy side is attacking our king and we don't have any legal
		// moves left -> mate found
		if inCheck == true {
			// this returns the distance to mate from root
			return -Infinite + pos.ply
		}
		// otherwise its a stalemate -> draw
		return 0
	}

	if alpha != oldAlpha {
		StoreHashEntry(pos, bestMove, bestScore, HFExact, depth)
	} else {
		StoreHashEntry(pos, bestMove, alpha, HFAlpha, depth)
	}

	return alpha
}

// IsRepetition Determine if position is a repetition
func IsRepetition(pos *Board) bool {
	// Loop through moves and check if the current position is equal to the posKey of any previous positions
	// since when we have a capture or a pawn move i.e. when we reset the 50 move counter, we can not go back
	// to the same position - we cannot unmove a pawn or uncapture a piece -> limit the search to the last time
	// the fifty move was reset
	for i := pos.histPly - pos.fiftyMove; i < pos.histPly-1; i++ {

		// // AssertTrue(i >= 0 && i < MaxGameMoves)
		if pos.posKey == pos.history[i].posKey {
			return true
		}
	}

	return false
}
