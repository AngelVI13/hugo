package main


import (
	utils "hugo/utils"
	"testing"
)

func BenchmarkGetAllMoves(b *testing.B) {
	utils.AllInit()
	var board utils.Board
	utils.ParseFen(utils.StartFen, &board)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var moveList utils.MoveList
		utils.GenerateAllMoves(&board, &moveList)
	}
	b.StopTimer()
}