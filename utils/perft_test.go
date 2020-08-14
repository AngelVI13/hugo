package utils

import (
	"testing"
)


func TestPerftStartingPosition(t *testing.T) {
	AllInit()
	var board Board
	ParseFen(StartFen, &board)
	LeafNodes = 0

	perft(5, &board)
	if LeafNodes != 4865609 {
		t.Errorf("Expected 4865609 moves at depth 5 from starting position, got %d\n", LeafNodes)
	}
}


func BenchmarkPerftStartingPositionDepth3(b *testing.B) {
	AllInit()
	var board Board
	ParseFen(StartFen, &board)
	LeafNodes = 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		perft(3, &board)
	}
	b.StopTimer()

	// fmt.Println(LeafNodes)
}
