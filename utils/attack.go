package utils

// KnightDir Squares increment to find places where the knight will be attacking the current piece
// For example if we want to check if square 55 (e4) is attacked. We need to check if there is a
// opposite coloured knight on square 55-8 = 47, 55-19=36 etc.
var knightDir = [8]int{-8, -19, -21, -12, 8, 19, 21, 12}
var rookDir = [4]int{-1, -10, 1, 10} // horizontal and vertical direction from a given pos
var bishopDir = [4]int{-9, -11, 11, 9}
var kingDir = [8]int{-1, -10, 1, 10, -9, -11, 11, 9}

// IsSquareAttacked determines if a given square is attacked from the opponent
func IsSquareAttacked(sq, side int, pos *Board) bool {
	// side here is the attacking side !!!!!!!!!!!!!!!!!

	// // AssertTrue(SquareOnBoard(sq))
	// // AssertTrue(SideValid(side))
	// // AssertTrue(// CheckBoard(pos))

	// pawns
	// if attacking side is white and there are pawns infornt to the left and right of us, then we are attacked
	if side == White {
		if pos.Pieces[sq-11] == WhitePawn || pos.Pieces[sq-9] == WhitePawn {
			return true
		}
	} else {
		if pos.Pieces[sq+11] == BlackPawn || pos.Pieces[sq+9] == BlackPawn {
			return true
		}
	}

	// knights
	// Loop through 8 directions
	for index := 0; index < 8; index++ {
		// find what piece is in that direction
		pce := pos.Pieces[sq+knightDir[index]]
		// if there is a knight of the opposite side at that piece -> return true
		if pce != OffBoard && IsPieceKnight[pce] && PieceColour[pce] == side {
			return true
		}
	}

	// rooks, queens
	for index := 0; index < 4; index++ {
		dir := rookDir[index]  // get current direction
		tSq := sq + dir        // take the first square
		pce := pos.Pieces[tSq] // see what piece is there
		for pce != OffBoard {  // while the piece is not OffBoard
			if pce != Empty { // if we hit a piece
				if IsPieceRookQueen[pce] && PieceColour[pce] == side { // if that piece is a rook or queen from the opposite side
					return true // our square is under attack -> return true
				}
				break // otherwise we hit a piece that is not an attacker -> try another direction
			}
			tSq += dir            // increment new piece square and perform check again
			pce = pos.Pieces[tSq] // get new piece
		}
	}

	// bishops, queens
	for index := 0; index < 4; index++ { // !!!!! could be rewriten as for _, dir := range bishopDir !!!!!!
		dir := bishopDir[index]
		tSq := sq + dir
		pce := pos.Pieces[tSq]
		for pce != OffBoard {
			if pce != Empty {
				if IsPieceBishopQueen[pce] && PieceColour[pce] == side {
					return true
				}
				break
			}
			tSq += dir
			pce = pos.Pieces[tSq]
		}
	}

	// kings
	for index := 0; index < 8; index++ {
		pce := pos.Pieces[sq+kingDir[index]]
		if pce != OffBoard && IsPieceKing[pce] && PieceColour[pce] == side {
			return true
		}
	}

	return false
}
