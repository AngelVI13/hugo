package main

import (
	"fmt"
	utils "local/hugo/utils"
	inpututils "local/input-utils"
	"strings"
)

const (
	fen1     string = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	fen2     string = "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2"
	fen3     string = "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"
	fen4     string = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	fen5     string = "8/3q1p2/8/5P2/4Q3/8/8/8 w - - 0 2"
	fen6     string = "rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1"
	fen7     string = "rnbqkbnr/p1p1p3/3p3p/1p1p4/2P1Pp2/8/PP1P1PpP/RNBQKB1R b KQkq e3 0 1"
	fen8     string = "5k2/1n6/4n3/6N1/8/3N4/8/5K2 b - - 0 1"
	fen9     string = "6k1/8/5r2/8/1nR5/5N2/8/6K1 b - - 0 1"
	fen10    string = "6k1/8/4nq2/8/1nQ5/5N2/1N6/6K1 b - - 0 1 "
	fen11    string = "6k1/1b6/4n3/8/1n4B1/1B3N2/1N6/2b3K1 b - - 0 1 "
	fen12    string = "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1"
	fen13    string = "3rk2r/8/8/8/8/8/6p1/R3K2R b KQk - 0 1"
	fen14    string = "r3k2r/pPppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	perftFen string = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	fen15    string = "2rr3k/pp3pp1/1nnqbN1p/3pN3/2pP4/2P3Q1/PPB4P/R4RK1 w - -"
	fen16    string = "r1b1k2r/ppppnppp/2n2q2/2b5/3NP3/2P1B3/PP3PPP/RN1QKB1R w KQkq - 0 1"
)

func showSqAtBySide(side int, pos *utils.Board) {

	fmt.Printf("\n\nSquares attacked by:%c\n", utils.SideChar[side])
	for rank := utils.Rank8; rank >= utils.Rank1; rank-- {
		for file := utils.FileA; file <= utils.FileH; file++ {
			sq := utils.FileRankToSquare(file, rank)
			if utils.IsSquareAttacked(sq, side, pos) {
				fmt.Printf("X")
			} else {
				fmt.Printf("-")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n\n")

}

func main() {
	utils.AllInit()

	var board utils.Board
	var info utils.SearchInfo
	utils.InitHashTable(&board.HashTable)

	line := ""
	fmt.Printf("Welcome to Hugo! Type 'hugo' for console mode...\n")

	for {
		line, _ = inpututils.GetInput("")
		if len(line) < 2 {
			continue
		}

		if strings.Contains(line, "uci") {
			utils.UciLoop(&board, &info)
			if info.Quit == true {
				break
			}
			continue
		} else if strings.Contains(line, "xboard") {
			utils.XBoardLoop(&board, &info)
			if info.Quit == true {
				break
			}
			continue
		} else if strings.Contains(line, "hugo") {
			utils.ConsoleLoop(&board, &info)
			if info.Quit == true {
				break
			}
			continue
		} else if strings.Contains(line, "quit") {
			break
		}
	}
}
