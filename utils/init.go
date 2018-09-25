package utils

import "math/rand"

// AllInit initialize everything
func AllInit() {
	InitSq120To64()
	InitBitMasks()
	InitHashKeys()
	InitFilesRanksBoard()
	InitEvaluationMasks()
	InitMvvLva()
}

// InitEvaluationMasks initialize evaluation masks
func InitEvaluationMasks() {
	// set everything to 0
	for sq := 0; sq < 8; sq++ {
		FileBBMask[sq] = 0
		RankBBMask[sq] = 0
	}

	// For every square, set the corresponding rank/file mask to 1
	for rank := Rank8; rank >= Rank1; rank-- {
		for file := FileA; file <= FileH; file++ {
			sq := uint64(rank*8 + file)
			FileBBMask[file] |= 1 << sq
			RankBBMask[rank] |= 1 << sq
		}
	}

	// set everything to 0
	for sq := 0; sq < 64; sq++ {
		IsolatedMask[sq] = 0
		WhitePassedMask[sq] = 0
		BlackPassedMask[sq] = 0
		WhiteDoubledMask[sq] = 0
	}

	for sq := 0; sq < 64; sq++ {
		targetSq := sq + 8

		for targetSq < 64 {
			WhitePassedMask[sq] |= 1 << uint64(targetSq)
			BlackDoubledMask[sq] |= 1 << uint64(targetSq)
			targetSq += 8
		}

		targetSq = sq - 8
		for targetSq >= 0 {
			BlackPassedMask[sq] |= 1 << uint64(targetSq)
			WhiteDoubledMask[sq] |= 1 << uint64(targetSq)
			targetSq -= 8
		}

		if FilesBoard[Sq120(sq)] > FileA {
			IsolatedMask[sq] |= FileBBMask[FilesBoard[Sq120(sq)]-1]

			targetSq = sq + 7
			for targetSq < 64 {
				WhitePassedMask[sq] |= 1 << uint64(targetSq)
				targetSq += 8
			}

			targetSq = sq - 9
			for targetSq >= 0 {
				BlackPassedMask[sq] |= 1 << uint64(targetSq)
				targetSq -= 8
			}
		}

		if FilesBoard[Sq120(sq)] < FileH {
			IsolatedMask[sq] |= FileBBMask[FilesBoard[Sq120(sq)]+1]

			targetSq = sq + 9
			for targetSq < 64 {
				WhitePassedMask[sq] |= 1 << uint64(targetSq)
				targetSq += 8
			}

			targetSq = sq - 7
			for targetSq >= 0 {
				BlackPassedMask[sq] |= 1 << uint64(targetSq)
				targetSq -= 8
			}
		}
	}

	// for sq := 0; sq < 64; sq++ {
	// 	PrintBitBoard(BlackDoubledMask[sq])
	// }
}

// InitSq120To64 Initialize board covertion arrays
func InitSq120To64() {
	// Set invalid values for all squares in 120Sq array
	for index := 0; index < BoardSquareNum; index++ {
		Sq120ToSq64[index] = 65
	}

	// Set invalid values for all squares in 64Sq array
	for index := 0; index < 64; index++ {
		Sq64ToSq120[index] = 120
	}
	// The above setup is later used for fail safe check that everything is set correctly

	sq64 := 0
	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			sq := FileRankToSquare(file, rank)
			Sq64ToSq120[sq64] = sq
			Sq120ToSq64[sq] = sq64
			sq64++
		}
	}
}

// InitBitMasks initializes bit masks values
func InitBitMasks() {
	// !!!!!
	// consider removing this, the slice will already be initialized to 0
	for index := 0; index < 64; index++ {
		SetMask[index] = 0
		ClearMask[index] = 0
	}

	for index := 0; index < 64; index++ {
		SetMask[index] |= (1 << uint64(index))
		ClearMask[index] = ^SetMask[index] // bitwise complement to SetMask
	}
}

// InitHashKeys initializes hashkeys for all pieces and possible positions, for castling rights, for side to move
func InitHashKeys() {
	for i := 0; i < 13; i++ {
		for j := 0; j < 120; j++ {
			PieceKeys[i][j] = rand.Uint64() // returns a random 64 bit number
		}
	}

	SideKey = rand.Uint64()
	for i := 0; i < 16; i++ {
		CastleKeys[i] = rand.Uint64()
	}
}
