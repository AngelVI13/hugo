package utils

import (
	"fmt"
)

// BitTable BitTable with default values
var BitTable = [64]int{63, 30, 3, 32, 25, 41, 22, 33, 15, 50, 42, 13, 11, 53, 19, 34, 61, 29, 2, 51, 21, 43, 45, 10, 18, 47, 1, 54, 9, 57, 0, 35, 62, 31, 40, 4, 49, 5, 52, 26, 60, 6, 23, 44, 46, 27, 56, 16, 7, 39, 48, 24, 59, 14, 12, 55, 38, 28, 58, 20, 37, 17, 36, 8}

// CountBits Counts number of 1's in b
func CountBits(b uint64) int {
	r := 0
	for b != 0 {
		b &= b - 1
		r++
	}
	return r
}

// PopBit Pops bit from bitboard
func PopBit(bb *uint64) int {
	b := *bb ^ (*bb - 1)
	fold := ((b & 0xffffffff) ^ (b >> 32))
	*bb &= (*bb - 1)

	// for some reason in the c-version the result of the multiplication is truncated to 32bits -> have to cast here as well
	return BitTable[uint32(fold*0x783a9b23)>>26]
}

// PrintBitBoard prints a give bitboard to screen
func PrintBitBoard(bb uint64) {
	var shiftMe uint64 = 1

	fmt.Println()
	for rank := Rank8; rank >= Rank1; rank-- {
		for file := FileA; file <= FileH; file++ {
			sq := FileRankToSquare(file, rank) // 120 based index
			sq64 := Sq64(sq)                   // get 64 based index

			// check if position has a value of 1 and print X
			if (shiftMe<<uint64(sq64))&bb != 0 {
				fmt.Print("X")
			} else {
				fmt.Print("-")
			}
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println()
}
