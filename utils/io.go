package utils

import (
	"fmt"
)

// PrintSquare get algebraic notation of square i.e. b2, a6 from array index
func PrintSquare(sq int) string {
	file := FilesBoard[sq]
	rank := RanksBoard[sq]

	// "a"[0] -> returns the byte value of the char 'a' -> convert to int to get ascii value
	// then add the file/rank value to it and convert back to string
	// therefore this automatically translates the files from 0-7 to a-h
	fileStr := string(int("a"[0]) + file)
	rankStr := string(int("1"[0]) + rank)

	squareStr := fileStr + rankStr
	return squareStr
}

// PrintMove prints move in algebraic notation
func PrintMove(move int) string {
	fileFrom := FilesBoard[FromSq(move)]
	rankFrom := RanksBoard[FromSq(move)]
	fileTo := FilesBoard[ToSq(move)]
	rankTo := RanksBoard[ToSq(move)]

	promoted := Promoted(move)

	moveStr := string(int("a"[0])+fileFrom) + string(int("1"[0])+rankFrom) +
		string(int("a"[0])+fileTo) + string(int("1"[0])+rankTo)

	// if this move is a promotion, add char of the piece we promote to at the end of the move string
	// i.e. if a7a8q -> we promote to Queen
	if promoted != 0 {
		pieceChar := "q"
		if IsPieceKnight[promoted] {
			pieceChar = "n"
		} else if IsPieceRookQueen[promoted] && !IsPieceBishopQueen[promoted] {
			pieceChar = "r"
		} else if !IsPieceRookQueen[promoted] && IsPieceBishopQueen[promoted] {
			pieceChar = "b"
		}
		moveStr += pieceChar
	}

	return moveStr
}

// PrintMoveList prints move list
func PrintMoveList(moveList *MoveList) {
	fmt.Println("MoveList:\n", moveList.Count)

	for index := 0; index < moveList.Count; index++ {

		move := moveList.Moves[index].Move
		score := moveList.Moves[index].score

		fmt.Printf("Move:%d > %s (score:%d)\n", index+1, PrintMove(move), score)
	}
	fmt.Printf("MoveList Total %d Moves:\n\n", moveList.Count)
}

// ParseMove parses user move and returns the MOVE int value from the GeneratedMoves for the
// position, that matches the moveStr input. For example if moveStr = 'a2a3'
// loops over all possible moves for the position, finds that move int i.e. 1451231 and returns it
func ParseMove(moveStr string, pos *Board) (move int) {
	// THIS COULD BE DOING BYTE COMPARISON INSTEAD OF INT COMPARISON !!!!!
	// check if files for 'from' and 'to' squares are valid i.e. between 1-8
	if moveStr[1] > "8"[0] || moveStr[1] < "1"[0] {
		return NoMove
	}

	if moveStr[3] > "8"[0] || moveStr[3] < "1"[0] {
		return NoMove
	}

	// check if ranks for 'from' and 'to' squares are valid i.e. between a-h
	if moveStr[0] > "h"[0] || moveStr[0] < "a"[0] {
		return NoMove
	}

	if moveStr[2] > "h"[0] || moveStr[2] < "a"[0] {
		return NoMove
	}

	from := FileRankToSquare(int(moveStr[0]-"a"[0]), int(moveStr[1]-"1"[0]))
	to := FileRankToSquare(int(moveStr[2]-"a"[0]), int(moveStr[3]-"1"[0]))

	// fmt.Printf("Move string: %s, from: %d to: %d\n", moveStr, from, to)

	// // AssertTrue(SquareOnBoard(from) && SquareOnBoard(to))

	var moveList MoveList
	GenerateAllMoves(pos, &moveList)

	for moveNum := 0; moveNum < moveList.Count; moveNum++ {
		move := moveList.Moves[moveNum].Move
		if FromSq(move) == from && ToSq(move) == to {
			promPiece := Promoted(move)
			if promPiece != Empty {
				if IsPieceRookQueen[promPiece] && !IsPieceBishopQueen[promPiece] && moveStr[4] == "r"[0] {
					return move
				} else if !IsPieceRookQueen[promPiece] && IsPieceBishopQueen[promPiece] && moveStr[4] == "b"[0] {
					return move
				} else if IsPieceRookQueen[promPiece] && IsPieceBishopQueen[promPiece] && moveStr[4] == "q"[0] {
					return move
				} else if IsPieceKnight[promPiece] && moveStr[4] == "n"[0] {
					return move
				}
				continue
			}
			// must not be a promotion -> return move
			return move
		}
	}

	return NoMove
}
