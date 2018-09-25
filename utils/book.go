package utils

import (
	"fmt"
	ioutils "local/io-utils"
	"local/string-utils"
	"math/rand"
	"strings"
	"time"
)

// GetBookMove returns a move from the opening book
func GetBookMove(pos *Board) int {
	if pos.histPly > 25 {
		return 0
	}

	book, err := ioutils.ScanFile(BookFile)
	if err != nil {
		fmt.Println("Book error")
		return 0
	}

	currentLine := ""
	for i := 0; i < pos.histPly; i++ {
		currentLine += PrintMove(pos.history[i].move) + " "
	}

	fmt.Println(currentLine)
	bookMoves := make([]int, 0)

	for _, bookLine := range book {
		if strings.Contains(bookLine, currentLine) {
			nextMovesStr := stringutils.RemoveStringToTheLeftOfMarker(bookLine, currentLine)
			nextMoveStr := stringutils.RemoveStringToTheRightOfMarker(nextMovesStr, " ")

			if len(nextMoveStr) > 5 {
				fmt.Println("Book move parsing error")
				continue // parsing eror
			} else {
				bookMoves = append(bookMoves, ParseMove(nextMoveStr, pos))
			}
		}
	}

	numberOfBookMoves := len(bookMoves)
	if len(bookMoves) > 0 {
		rand.Seed(time.Now().Unix())
		return bookMoves[rand.Intn(numberOfBookMoves)]
	}
	return 0
}
