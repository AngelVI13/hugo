package utils

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ScanFile reads file and returns []slice with all lines
func ScanFile(filename string) ([]string, error) {
	lines := make([]string, 0)

	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Printf("open file error: %v", err)
		return lines, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic("Couldn't close file")
		}
	}()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return nil, err
	}

	return lines, nil
}

// GetBookMove returns a move from the opening book
func GetBookMove(pos *Board) int {
	if pos.histPly > 25 {
		return 0
	}

	book, err := ScanFile(BookFile)
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
			nextMovesStr := RemoveStringToTheLeftOfMarker(bookLine, currentLine)
			nextMoveStr := RemoveStringToTheRightOfMarker(nextMovesStr, " ")

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

// RemoveStringToTheLeftOfMarker Removes substring to the left of a marker. If markers not in string -> return unchanged string.
func RemoveStringToTheLeftOfMarker(s, marker string) (result string) {
	markerIdx := strings.Index(s, marker)
	if markerIdx == -1 {
		return s
	}

	markerIdx = markerIdx + len(marker)

	textForRemoval := s[0:markerIdx]
	resultStr := strings.Replace(s, textForRemoval, "", -1)
	return resultStr
}

// RemoveStringToTheRightOfMarker Removes substring to the right of a marker. If markers not in string -> return unchanged string.
func RemoveStringToTheRightOfMarker(s, marker string) (result string) {
	markerIdx := strings.Index(s, marker)
	if markerIdx == -1 {
		return s
	}

	// markerIdx = markerIdx + len(marker)

	textForRemoval := s[markerIdx:]
	resultStr := strings.Replace(s, textForRemoval, "", -1)
	return resultStr
}
