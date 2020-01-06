package utils

/*
MoveGen(board, moveList)
	loops all pieces
		-> slider loop each dir, add move
			-> AddMove moveList->moves[moveList->count]=move, moveList->count++
*/

/*
The order of the moves searched is important for alpha beta search because
the sooner we find the strongest move for a position the sooner we can skip
checking all the other moves. Because of this we order the moves for a given position
as follows:
1. search the PV move
2. search captures and start with - MvvLVA - most valuable victim least valuable attacker
	this means that if we have a possibility of a pawn capturing a queen -> we search this
	before we search for queen captures pawn
	These captures are still in order, for example:
	PxQ is searched first and then PxN etc.
3. search for killer moves (beta cuttoffs)
4. search history for score
*/

// VictimScore MvvLVA victim scores
var VictimScore = map[int]int{
	Empty:       0,
	WhitePawn:   100,
	WhiteKnight: 200,
	WhiteBishop: 300,
	WhiteRook:   400,
	WhiteQueen:  500,
	WhiteKing:   600,
	BlackPawn:   100,
	BlackKnight: 200,
	BlackBishop: 300,
	BlackRook:   400,
	BlackQueen:  500,
	BlackKing:   600,
}

// MvvLvaScores MvvLVA scores slice
var MvvLvaScores [13][13]int

// InitMvvLva initialize MvvLva slice
func InitMvvLva() {
	for attacker := WhitePawn; attacker <= BlackKing; attacker++ {
		for victim := WhitePawn; victim <= BlackKing; victim++ {
			// if attacker is pawn and victim is queen - score is 505
			// if attacker is queen and victim is pawn - score is 101
			MvvLvaScores[victim][attacker] = VictimScore[victim] + 6 - (VictimScore[attacker] / 100)
		}
	}
}

// GetMoveInt creates and returns a move int from given move information
func GetMoveInt(fromSq, toSq, capturePiece, promotionPiece, flag int) int {
	return fromSq | (toSq << 7) | (capturePiece << 14) | (promotionPiece << 20) | flag
}

// This already exists from validate.go
// // IsSquareOffBoard returns true if square is off board
// func IsSquareOffBoard(sq int) bool {
// 	return (FilesBoard[sq] == OffBoard)
// }

// TODO: COMBINE AND PARAMETERIZE THE CODE BELOW for all moves GEN !!!!!!!!!!

// addQuietMove adds quiet move
func addQuietMove(pos *Board, move int, moveList *MoveList) {
	// // AssertTrue(SquareOnBoard(FromSq(move)))
	// // AssertTrue(SquareOnBoard(ToSq(move)))
	moveList.Moves[moveList.Count].Move = move

	// Set killer moves scores, the best one hasa a higher score
	// i.e. the one at index 0
	if pos.searchKillers[0][pos.ply] == move {
		moveList.Moves[moveList.Count].score = 900000
	} else if pos.searchKillers[1][pos.ply] == move {
		moveList.Moves[moveList.Count].score = 800000
	} else {
		moveList.Moves[moveList.Count].score = pos.searchHistory[pos.Pieces[FromSq(move)]][ToSq(move)]
	}
	moveList.Count++
}

// addCaptureMove adds capture move
func addCaptureMove(pos *Board, move int, moveList *MoveList) {
	// // AssertTrue(SquareOnBoard(FromSq(move)))
	// // AssertTrue(SquareOnBoard(ToSq(move)))
	// // AssertTrue(PieceValid(Captured(move)))
	moveList.Moves[moveList.Count].Move = move
	// the +1000000 makes sure that these moves are searched before killer moves which will have
	// +800000-900000 which will be searche before history moves which can
	// go up to 800000
	moveList.Moves[moveList.Count].score = MvvLvaScores[Captured(move)][pos.Pieces[FromSq(move)]] + 1000000
	moveList.Count++
}

// addEnPassantMove adds quiet move
func addEnPassantMove(pos *Board, move int, moveList *MoveList) {
	moveList.Moves[moveList.Count].Move = move
	moveList.Moves[moveList.Count].score = 105 + 1000000 // based on MvvLVA -> pawn take pawn results in 105
	moveList.Count++
}

// addWhitePawnCaptureMove add white pawn capture move
func addWhitePawnCaptureMove(pos *Board, from, to, cap int, moveList *MoveList) {
	// // AssertTrue(PieceValidEmpty(cap))
	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	if RanksBoard[from] == Rank7 {
		// add all promotion with capture related moves
		addCaptureMove(pos, GetMoveInt(from, to, cap, WhiteQueen, 0), moveList)
		addCaptureMove(pos, GetMoveInt(from, to, cap, WhiteRook, 0), moveList)
		addCaptureMove(pos, GetMoveInt(from, to, cap, WhiteBishop, 0), moveList)
		addCaptureMove(pos, GetMoveInt(from, to, cap, WhiteKnight, 0), moveList)
	} else {
		// add normal capture moves without promotion
		addCaptureMove(pos, GetMoveInt(from, to, cap, Empty, 0), moveList)
	}
}

// addWhitePawnMove add white pawn normal moves
func addWhitePawnMove(pos *Board, from, to int, moveList *MoveList) {
	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	if RanksBoard[from] == Rank7 {
		// add normal promotion without capture
		addQuietMove(pos, GetMoveInt(from, to, Empty, WhiteQueen, 0), moveList)
		addQuietMove(pos, GetMoveInt(from, to, Empty, WhiteRook, 0), moveList)
		addQuietMove(pos, GetMoveInt(from, to, Empty, WhiteBishop, 0), moveList)
		addQuietMove(pos, GetMoveInt(from, to, Empty, WhiteKnight, 0), moveList)
	} else {
		addQuietMove(pos, GetMoveInt(from, to, Empty, Empty, 0), moveList)
	}
}

// addBlackPawnCaptureMove add black pawn capture move
func addBlackPawnCaptureMove(pos *Board, from, to, cap int, moveList *MoveList) {
	// // AssertTrue(PieceValidEmpty(cap))
	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	if RanksBoard[from] == Rank2 {
		// add all promotion with capture related moves
		addCaptureMove(pos, GetMoveInt(from, to, cap, BlackQueen, 0), moveList)
		addCaptureMove(pos, GetMoveInt(from, to, cap, BlackRook, 0), moveList)
		addCaptureMove(pos, GetMoveInt(from, to, cap, BlackBishop, 0), moveList)
		addCaptureMove(pos, GetMoveInt(from, to, cap, BlackKnight, 0), moveList)
	} else {
		// add normal capture moves without promotion
		addCaptureMove(pos, GetMoveInt(from, to, cap, Empty, 0), moveList)
	}
}

// addBlackPawnMove add black pawn normal moves
func addBlackPawnMove(pos *Board, from, to int, moveList *MoveList) {
	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	if RanksBoard[from] == Rank2 {
		// add normal promotion without capture
		addQuietMove(pos, GetMoveInt(from, to, Empty, BlackQueen, 0), moveList)
		addQuietMove(pos, GetMoveInt(from, to, Empty, BlackRook, 0), moveList)
		addQuietMove(pos, GetMoveInt(from, to, Empty, BlackBishop, 0), moveList)
		addQuietMove(pos, GetMoveInt(from, to, Empty, BlackKnight, 0), moveList)
	} else {
		addQuietMove(pos, GetMoveInt(from, to, Empty, Empty, 0), moveList)
	}
}

// GenerateAllMoves generates all moves for a position
func GenerateAllMoves(pos *Board, moveList *MoveList) {

	// // AssertTrue(// CheckBoard(pos))

	moveList.Count = 0

	side := pos.side

	// TODO: COMBINE AND PARAMETERIZE THE CODE BELOW for PAWN GEN !!!!!!!!!!
	if side == White {
		// loop through all of the pawns on the board
		for pceNum := 0; pceNum < pos.pieceNum[WhitePawn]; pceNum++ {
			sq := pos.pieceList[WhitePawn][pceNum] // find the square of the pawn
			// // AssertTrue(SquareOnBoard(sq))          // check if pawn is on the board

			// add simple pawn move forward if next sq is empty
			if pos.Pieces[sq+10] == Empty {
				addWhitePawnMove(pos, sq, sq+10, moveList)
				// if we are on the second rank, generate a double pawn move if 4th rank sq is empty
				if RanksBoard[sq] == Rank2 && pos.Pieces[sq+20] == Empty {
					// don't forget to set the flag for PAWN START
					addQuietMove(pos, GetMoveInt(sq, (sq+20), Empty, Empty, MoveFlagPawnStart), moveList)
				}
			}

			// Capture to the left and right
			// check if the square that we are capturing on is on the board and that it has a black piece on it
			if SquareOnBoard(sq+9) && PieceColour[pos.Pieces[sq+9]] == Black {
				addWhitePawnCaptureMove(pos, sq, sq+9, pos.Pieces[sq+9], moveList)
			}
			// check if the square that we are capturing on is on the board and that it has a black piece on it
			if SquareOnBoard(sq+11) && PieceColour[pos.Pieces[sq+11]] == Black {
				addWhitePawnCaptureMove(pos, sq, sq+11, pos.Pieces[sq+11], moveList)
			}

			if pos.enPas != NoSquare {
				// check if the sq+9 square is equal to the enpassant square that we have stored in our pos
				if sq+9 == pos.enPas {
					addEnPassantMove(pos, GetMoveInt(sq, sq+9, Empty, Empty, MoveFlagEnPass), moveList)
				}
				if sq+11 == pos.enPas {
					addEnPassantMove(pos, GetMoveInt(sq, sq+11, Empty, Empty, MoveFlagEnPass), moveList)
				}
			}
		}
		// if the position allows white king castling
		// here we do not check if square G1 (final square after castling) is attacked
		// this will be handled at the end of the function where we will verify that all generated
		// moves are legal
		if (pos.castlePerm & WhiteKingCastling) != 0 {
			if pos.Pieces[F1] == Empty && pos.Pieces[G1] == Empty {
				if !IsSquareAttacked(E1, Black, pos) && !IsSquareAttacked(F1, Black, pos) {
					addQuietMove(pos, GetMoveInt(E1, G1, Empty, Empty, MoveFlagCastle), moveList)
				}
			}
		}

		if (pos.castlePerm & WhiteQueenCastling) != 0 {
			if pos.Pieces[D1] == Empty && pos.Pieces[C1] == Empty && pos.Pieces[B1] == Empty {
				if !IsSquareAttacked(E1, Black, pos) && !IsSquareAttacked(D1, Black, pos) {
					addQuietMove(pos, GetMoveInt(E1, C1, Empty, Empty, MoveFlagCastle), moveList)
				}
			}
		}
	} else {
		for pceNum := 0; pceNum < pos.pieceNum[BlackPawn]; pceNum++ {
			sq := pos.pieceList[BlackPawn][pceNum]
			// // AssertTrue(SquareOnBoard(sq))

			// add simple pawn move forward if next sq is empty
			if pos.Pieces[sq-10] == Empty {
				addBlackPawnMove(pos, sq, sq-10, moveList)
				// if we are on the second rank, generate a double pawn move if 4th rank sq is empty
				if RanksBoard[sq] == Rank7 && pos.Pieces[sq-20] == Empty {
					// don't forget to set the flag for PAWN START
					addQuietMove(pos, GetMoveInt(sq, (sq-20), Empty, Empty, MoveFlagPawnStart), moveList)
				}
			}

			// Capture to the left and right
			// check if the square that we are capturing on is on the board and that it has a black piece on it
			if SquareOnBoard(sq-9) && PieceColour[pos.Pieces[sq-9]] == White {
				addBlackPawnCaptureMove(pos, sq, sq-9, pos.Pieces[sq-9], moveList)
			}
			// check if the square that we are capturing on is on the board and that it has a black piece on it
			if SquareOnBoard(sq-11) && PieceColour[pos.Pieces[sq-11]] == White {
				addBlackPawnCaptureMove(pos, sq, sq-11, pos.Pieces[sq-11], moveList)
			}

			if pos.enPas != NoSquare {
				// check if the sq-9 square is equal to the enpassant square that we have stored in our pos
				if sq-9 == pos.enPas {
					addEnPassantMove(pos, GetMoveInt(sq, sq-9, Empty, Empty, MoveFlagEnPass), moveList)
				}
				if sq-11 == pos.enPas {
					addEnPassantMove(pos, GetMoveInt(sq, sq-11, Empty, Empty, MoveFlagEnPass), moveList)
				}
			}
		}
		// castling
		if (pos.castlePerm & BlackKingCastling) != 0 {
			if pos.Pieces[F8] == Empty && pos.Pieces[G8] == Empty {
				if !IsSquareAttacked(E8, White, pos) && !IsSquareAttacked(F8, White, pos) {
					addQuietMove(pos, GetMoveInt(E8, G8, Empty, Empty, MoveFlagCastle), moveList)
				}
			}
		}

		if (pos.castlePerm & BlackQueenCastling) != 0 {
			if pos.Pieces[D8] == Empty && pos.Pieces[C8] == Empty && pos.Pieces[B8] == Empty {
				if !IsSquareAttacked(E8, White, pos) && !IsSquareAttacked(D8, White, pos) {
					addQuietMove(pos, GetMoveInt(E8, C8, Empty, Empty, MoveFlagCastle), moveList)
				}
			}
		}
	}

	// Loop for slide pieces
	pieceIndex := LoopSlideIndex[side]
	piece := LoopSlidePiece[pieceIndex]
	pieceIndex++

	for piece != 0 {
		// // AssertTrue(PieceValid(piece))

		for pceNum := 0; pceNum < pos.pieceNum[piece]; pceNum++ {
			sq := pos.pieceList[piece][pceNum]
			// // AssertTrue(SquareOnBoard(sq))

			for index := 0; index < NumberOfDir[piece]; index++ {
				dir := PiececeDir[piece][index]
				targetSq := sq + dir

				// while we are still on the board, take a sliding piece and add a possible move
				// untill we see another piece or we hit the edge of the board
				for SquareOnBoard(targetSq) {
					// BLACK ^ 1 == WHITE       WHITE ^ 1 == BLACK
					if pos.Pieces[targetSq] != Empty {
						if PieceColour[pos.Pieces[targetSq]] == side^1 {
							addCaptureMove(pos, GetMoveInt(sq, targetSq, pos.Pieces[targetSq], Empty, 0), moveList)
						}
						break // if we hit a non-empty square, we break from this direction
					}
					addQuietMove(pos, GetMoveInt(sq, targetSq, Empty, Empty, 0), moveList)
					targetSq += dir
				}
			}
		}

		piece = LoopSlidePiece[pieceIndex]
		pieceIndex++
	}

	/* Loop for non slide */
	pieceIndex = LoopNonSlideIndex[side]
	piece = LoopNonSlidePiece[pieceIndex]
	pieceIndex++

	for piece != 0 {
		// // AssertTrue(PieceValid(piece))

		for pceNum := 0; pceNum < pos.pieceNum[piece]; pceNum++ {
			sq := pos.pieceList[piece][pceNum]
			// // AssertTrue(SquareOnBoard(sq))

			for index := 0; index < NumberOfDir[piece]; index++ {
				dir := PiececeDir[piece][index]
				targetSq := sq + dir

				if !SquareOnBoard(targetSq) {
					continue
				}

				// BLACK ^ 1 == WHITE       WHITE ^ 1 == BLACK
				if pos.Pieces[targetSq] != Empty {
					if PieceColour[pos.Pieces[targetSq]] == side^1 {
						addCaptureMove(pos, GetMoveInt(sq, targetSq, pos.Pieces[targetSq], Empty, 0), moveList)
					}
					continue
				}
				addQuietMove(pos, GetMoveInt(sq, targetSq, Empty, Empty, 0), moveList)
			}
		}

		piece = LoopNonSlidePiece[pieceIndex]
		pieceIndex++
	}
}

// GenerateAllCaptures generates all moves for a position
func GenerateAllCaptures(pos *Board, moveList *MoveList) {
	// // AssertTrue(// CheckBoard(pos))

	moveList.Count = 0

	side := pos.side

	// TODO: COMBINE AND PARAMETERIZE THE CODE BELOW for PAWN GEN !!!!!!!!!!
	if side == White {
		// loop through all of the pawns on the board
		for pceNum := 0; pceNum < pos.pieceNum[WhitePawn]; pceNum++ {
			sq := pos.pieceList[WhitePawn][pceNum] // find the square of the pawn
			// // AssertTrue(SquareOnBoard(sq))          // check if pawn is on the board

			// Capture to the left and right
			// check if the square that we are capturing on is on the board and that it has a black piece on it
			if SquareOnBoard(sq+9) && PieceColour[pos.Pieces[sq+9]] == Black {
				addWhitePawnCaptureMove(pos, sq, sq+9, pos.Pieces[sq+9], moveList)
			}
			// check if the square that we are capturing on is on the board and that it has a black piece on it
			if SquareOnBoard(sq+11) && PieceColour[pos.Pieces[sq+11]] == Black {
				addWhitePawnCaptureMove(pos, sq, sq+11, pos.Pieces[sq+11], moveList)
			}

			if pos.enPas != NoSquare {
				// check if the sq+9 square is equal to the enpassant square that we have stored in our pos
				if sq+9 == pos.enPas {
					addEnPassantMove(pos, GetMoveInt(sq, sq+9, Empty, Empty, MoveFlagEnPass), moveList)
				}
				if sq+11 == pos.enPas {
					addEnPassantMove(pos, GetMoveInt(sq, sq+11, Empty, Empty, MoveFlagEnPass), moveList)
				}
			}
		}
	} else {
		for pceNum := 0; pceNum < pos.pieceNum[BlackPawn]; pceNum++ {
			sq := pos.pieceList[BlackPawn][pceNum]
			// // AssertTrue(SquareOnBoard(sq))

			// Capture to the left and right
			// check if the square that we are capturing on is on the board and that it has a black piece on it
			if SquareOnBoard(sq-9) && PieceColour[pos.Pieces[sq-9]] == White {
				addBlackPawnCaptureMove(pos, sq, sq-9, pos.Pieces[sq-9], moveList)
			}
			// check if the square that we are capturing on is on the board and that it has a black piece on it
			if SquareOnBoard(sq-11) && PieceColour[pos.Pieces[sq-11]] == White {
				addBlackPawnCaptureMove(pos, sq, sq-11, pos.Pieces[sq-11], moveList)
			}

			if pos.enPas != NoSquare {
				// check if the sq-9 square is equal to the enpassant square that we have stored in our pos
				if sq-9 == pos.enPas {
					addEnPassantMove(pos, GetMoveInt(sq, sq-9, Empty, Empty, MoveFlagEnPass), moveList)
				}
				if sq-11 == pos.enPas {
					addEnPassantMove(pos, GetMoveInt(sq, sq-11, Empty, Empty, MoveFlagEnPass), moveList)
				}
			}
		}
	}

	// Loop for slide pieces
	pieceIndex := LoopSlideIndex[side]
	piece := LoopSlidePiece[pieceIndex]
	pieceIndex++

	for piece != 0 {
		// // AssertTrue(PieceValid(piece))

		for pceNum := 0; pceNum < pos.pieceNum[piece]; pceNum++ {
			sq := pos.pieceList[piece][pceNum]
			// // AssertTrue(SquareOnBoard(sq))

			for index := 0; index < NumberOfDir[piece]; index++ {
				dir := PiececeDir[piece][index]
				targetSq := sq + dir

				// while we are still on the board, take a sliding piece and add a possible move
				// untill we see another piece or we hit the edge of the board
				for SquareOnBoard(targetSq) {
					// BLACK ^ 1 == WHITE       WHITE ^ 1 == BLACK
					if pos.Pieces[targetSq] != Empty {
						if PieceColour[pos.Pieces[targetSq]] == side^1 {
							addCaptureMove(pos, GetMoveInt(sq, targetSq, pos.Pieces[targetSq], Empty, 0), moveList)
						}
						break // if we hit a non-empty square, we break from this direction
					}
					targetSq += dir
				}
			}
		}

		piece = LoopSlidePiece[pieceIndex]
		pieceIndex++
	}

	/* Loop for non slide */
	pieceIndex = LoopNonSlideIndex[side]
	piece = LoopNonSlidePiece[pieceIndex]
	pieceIndex++

	for piece != 0 {
		// // AssertTrue(PieceValid(piece))

		for pceNum := 0; pceNum < pos.pieceNum[piece]; pceNum++ {
			sq := pos.pieceList[piece][pceNum]
			// // AssertTrue(SquareOnBoard(sq))

			for index := 0; index < NumberOfDir[piece]; index++ {
				dir := PiececeDir[piece][index]
				targetSq := sq + dir

				if !SquareOnBoard(targetSq) {
					continue
				}

				// BLACK ^ 1 == WHITE       WHITE ^ 1 == BLACK
				if pos.Pieces[targetSq] != Empty {
					if PieceColour[pos.Pieces[targetSq]] == side^1 {
						addCaptureMove(pos, GetMoveInt(sq, targetSq, pos.Pieces[targetSq], Empty, 0), moveList)
					}
					continue
				}
			}
		}

		piece = LoopNonSlidePiece[pieceIndex]
		pieceIndex++
	}
}

// MoveExists returns true if a move is legal i.e. exists
func MoveExists(pos *Board, move int) bool {
	var moveList MoveList

	GenerateAllMoves(pos, &moveList) // generate all moves for the position

	// loop through all moves
	for moveNum := 0; moveNum < moveList.Count; moveNum++ {
		if !MakeMove(pos, moveList.Moves[moveNum].Move) {
			continue
		}

		// take back the made move from the 'if' above and if that move is legal and is equal to the move
		// we passed in -> move exists -> return true
		TakeMove(pos)

		if moveList.Moves[moveNum].Move == move {
			return true
		}
	}

	return false
}
