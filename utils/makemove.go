package utils

// --- Hashing 'macros' ---
func hashPiece(piece, sq int, pos *Board) {
	pos.posKey ^= PieceKeys[piece][sq]
}

func hashCastlePerm(pos *Board) {
	pos.posKey ^= CastleKeys[pos.castlePerm]
}

func hashSide(pos *Board) {
	pos.posKey ^= SideKey
}

func hashEnPass(pos *Board) {
	pos.posKey ^= PieceKeys[Empty][pos.enPas]
}

// ------------------------

// CastlePerm used to simplify hashing castle permissions
// Everytime we make a move we will take pos.castlePerm &= CastlePerm[sq]
// in this way if any of the rooks or the king moves, the castle permission will be
// disabled for that side. In any other move, the castle permissions will remain the
// same, since 15 is the max number associated with all possible castling permissions
// for both sides
var CastlePerm = [BoardSquareNum]int{
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 13, 15, 15, 15, 12, 15, 15, 14, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15,  7, 15, 15, 15,  3, 15, 15, 11, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
}

func clearPiece(sq int, pos *Board) {

	// // AssertTrue(SquareOnBoard(sq))

	pce := pos.Pieces[sq]

	// // AssertTrue(PieceValid(pce))

	colour := PieceColour[pce]
	tempPieceNum := -1

	hashPiece(pce, sq, pos)

	pos.Pieces[sq] = Empty
	pos.material[colour] -= PieceVal[pce]

	// Decrement number of pieces from respective arrays
	if PieceBig[pce] {
		pos.bigPieceNum[colour]--
		if PieceMaj[pce] {
			pos.majorPieceNum[colour]--
		} else {
			pos.minorPieceNum[colour]--
		}
	} else {
		// if piece is a pawn, remove it from the board of same coloured pawns and also from board of both coloured pawns
		ClearBit(&pos.Pawns[colour], Sq120ToSq64[sq])
		ClearBit(&pos.Pawns[Both], Sq120ToSq64[sq])
	}

	/*
		pos.pceNum[wP] == 5 Looping from 0 to 4
		pos.pList[pce][0] == sq0
		pos.pList[pce][1] == sq1
		pos.pList[pce][2] == sq2
		pos.pList[pce][3] == sq3
		pos.pList[pce][4] == sq4

		sq==sq3 so t_pceNum = 3;
	*/

	// Loop over all available WhitePawns for example, and on of the white pawns
	// sould have the same square number as the one passed to the function to be cleared
	// find that one
	for index := 0; index < pos.pieceNum[pce]; index++ {
		if pos.pieceList[pce][index] == sq {
			tempPieceNum = index
			break
		}
	}

	// // AssertTrue(tempPieceNum != -1)

	pos.pieceNum[pce]--
	// pos.pceNum[wP] == 4

	pos.pieceList[pce][tempPieceNum] = pos.pieceList[pce][pos.pieceNum[pce]]
	//pos.pList[wP][3]	= pos.pList[wP][4] = sq4
	/*
		pos.pceNum[wP] == 4 Looping from 0 to 3
		pos.pList[pce][0] == sq0
		pos.pList[pce][1] == sq1
		pos.pList[pce][2] == sq2
		pos.pList[pce][3] == sq4
	*/

}

func addPiece(sq, pce int, pos *Board) {

	// // AssertTrue(PieceValid(pce))
	// // AssertTrue(SquareOnBoard(sq))

	colour := PieceColour[pce]

	hashPiece(pce, sq, pos)

	pos.Pieces[sq] = pce

	if PieceBig[pce] {
		pos.bigPieceNum[colour]++
		if PieceMaj[pce] {
			pos.majorPieceNum[colour]++
		} else {
			pos.minorPieceNum[colour]++
		}
	} else {
		SetBit(&pos.Pawns[colour], Sq120ToSq64[sq])
		SetBit(&pos.Pawns[Both], Sq120ToSq64[sq])
	}

	pos.material[colour] += PieceVal[pce]
	pos.pieceList[pce][pos.pieceNum[pce]] = sq
	pos.pieceNum[pce]++
}

func movePiece(from, to int, pos *Board) {

	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	pce := pos.Pieces[from]
	colour := PieceColour[pce]

	// hash the piece out of the from square and then later hash it back in to the new square
	hashPiece(pce, from, pos)
	pos.Pieces[from] = Empty

	hashPiece(pce, to, pos)
	pos.Pieces[to] = pce

	if !PieceBig[pce] {
		ClearBit(&pos.Pawns[colour], Sq120ToSq64[from])
		ClearBit(&pos.Pawns[Both], Sq120ToSq64[from])
		SetBit(&pos.Pawns[colour], Sq120ToSq64[to])
		SetBit(&pos.Pawns[Both], Sq120ToSq64[to])
	}

	// Update the square value for the given piece
	// i.e. change its destination from 'from sq' -> 'to sq'
	for index := 0; index < pos.pieceNum[pce]; index++ {
		if pos.pieceList[pce][index] == from {
			pos.pieceList[pce][index] = to

			break
		}
	}
}

// MakeMove perform a move
// return false if the side to move has left themselves in check after the move i.e. illegal move
func MakeMove(pos *Board, move int) bool {

	// // AssertTrue(// CheckBoard(pos))

	from := FromSq(move)
	to := ToSq(move)
	side := pos.side

	// Make sure all input info is valid
	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))
	// // AssertTrue(SideValid(side))
	// // AssertTrue(PieceValid(pos.Pieces[from]))

	// Store has value before we do any hashing in/out of pieces etc
	pos.history[pos.histPly].posKey = pos.posKey

	// if this is an en passant move
	if move&MoveFlagEnPass != 0 {
		// if the side thats making the capture is white
		// then we need to remove the black pawn right behind the new position of the white piece
		// i.e. new_pos - 10 -> translated to array index
		if side == White {
			clearPiece(to-10, pos)
		} else {
			clearPiece(to+10, pos)
		}
	} else if move&MoveFlagCastle != 0 {
		// if its a castlign move, based on the TO square, make the appopriate move, otherwise assert false
		switch to {
		case C1:
			movePiece(A1, D1, pos)
		case C8:
			movePiece(A8, D8, pos)
		case G1:
			movePiece(H1, F1, pos)
		case G8:
			movePiece(H8, F8, pos)
		default:
			// // AssertTrue(false)
		}
	}

	// If the current enpassant square is SET, then we hash in the poskey
	if pos.enPas != NoSquare {
		hashEnPass(pos)
	}
	hashCastlePerm(pos) // hash out the castling permissions

	// store information to the history array about this move
	pos.history[pos.histPly].move = move
	pos.history[pos.histPly].fiftyMove = pos.fiftyMove
	pos.history[pos.histPly].enPas = pos.enPas
	pos.history[pos.histPly].castlePerm = pos.castlePerm

	// if a rook or king has moved the remove the respective castling permission from castlePerm
	pos.castlePerm &= CastlePerm[from]
	pos.castlePerm &= CastlePerm[to]
	pos.enPas = NoSquare // set enpassant square to no square

	hashCastlePerm(pos) // hash back in the castling perm

	pos.fiftyMove++ // increment firfty move rule

	// get what piece, if any, was captured in the move and if somethig was actually captured
	// i.e. captured piece is not empty remove captured piece and reset fifty move rule
	if captured := Captured(move); captured != Empty {
		// // AssertTrue(PieceValid(captured))
		clearPiece(to, pos)
		pos.fiftyMove = 0
	}

	// increase halfmove counter and ply counter values
	pos.histPly++
	pos.ply++

	// check if we need to set a new en passant square i.e. if this is a pawn start
	// then depending on the side find the piece just behind the new pawn destination
	// i.e. A4 -> compute A3 and set that as a possible enpassant capture square
	if IsPiecePawn[pos.Pieces[from]] {
		pos.fiftyMove = 0
		if move&MoveFlagPawnStart != 0 {
			if side == White {
				pos.enPas = from + 10
				// // AssertTrue(RanksBoard[pos.enPas] == Rank3)
			} else {
				pos.enPas = from - 10
				// // AssertTrue(RanksBoard[pos.enPas] == Rank6)
			}
			hashEnPass(pos) // hash in the enpass
		}
	}

	movePiece(from, to, pos)

	// get promoted piece and if its not empty, clear old piece (pawn)
	// and add new piece (whatever was the selected promotion piece)
	if promotedPiece := Promoted(move); promotedPiece != Empty {
		// // AssertTrue(PieceValid(promotedPiece) && !IsPiecePawn[promotedPiece])
		clearPiece(to, pos)
		addPiece(to, promotedPiece, pos)
	}

	// if we move the king -> update king square
	if IsPieceKing[pos.Pieces[to]] {
		pos.kingSquare[pos.side] = to
	}

	pos.side ^= 1 // change side to move
	hashSide(pos) // hash in the new side

	// // AssertTrue(// CheckBoard(pos))

	// check if after this move, our king is in check -> if yes -> illegal move
	if IsSquareAttacked(pos.kingSquare[side], pos.side, pos) {
		TakeMove(pos)
		return false
	}

	return true
}

// TakeMove revert move, opposite to MakeMove()
func TakeMove(pos *Board) {

	// // AssertTrue(// CheckBoard(pos))

	pos.histPly--
	pos.ply--

	move := pos.history[pos.histPly].move
	from := FromSq(move)
	to := ToSq(move)

	// // AssertTrue(SquareOnBoard(from))
	// // AssertTrue(SquareOnBoard(to))

	if pos.enPas != NoSquare {
		hashEnPass(pos)
	}
	hashCastlePerm(pos)

	pos.castlePerm = pos.history[pos.histPly].castlePerm
	pos.fiftyMove = pos.history[pos.histPly].fiftyMove
	pos.enPas = pos.history[pos.histPly].enPas

	if pos.enPas != NoSquare {
		hashEnPass(pos)
	}
	hashCastlePerm(pos)

	pos.side ^= 1
	hashSide(pos)

	if MoveFlagEnPass&move != 0 {
		if pos.side == White {
			addPiece(to-10, BlackPawn, pos)
		} else {
			addPiece(to+10, WhitePawn, pos)
		}
	} else if MoveFlagCastle&move != 0 {
		switch to {
		case C1:
			movePiece(D1, A1, pos)
		case C8:
			movePiece(D8, A8, pos)
		case G1:
			movePiece(F1, H1, pos)
		case G8:
			movePiece(F8, H8, pos)
		default:
			// // AssertTrue(false)
		}
	}

	movePiece(to, from, pos)

	if IsPieceKing[pos.Pieces[from]] {
		pos.kingSquare[pos.side] = from
	}

	if captured := Captured(move); captured != Empty {
		// // AssertTrue(PieceValid(captured))
		addPiece(to, captured, pos)
	}

	if promoted := Promoted(move); promoted != Empty {
		// // AssertTrue(PieceValid(Promoted(move)) && !IsPiecePawn[Promoted(move)])
		clearPiece(from, pos)
		if PieceColour[Promoted(move)] == White {
			addPiece(from, WhitePawn, pos)
		} else {
			addPiece(from, BlackPawn, pos)
		}
	}

	// // AssertTrue(// CheckBoard(pos))
}

// MakeNullMove make null move, i.e. take the opponnent a free move and examine
// the position after he makes it
func MakeNullMove(pos *Board) {
	// AssertTrue(CheckBoard(pos))
	// make sure we are not in check
	// AssertTrue(!IsSquareAttacked(pos.kingSquare[pos.side], pos.side^1, pos))

	pos.ply++
	pos.history[pos.histPly].posKey = pos.posKey

	if pos.enPas != NoSquare {
		hashEnPass(pos)
	}

	pos.history[pos.histPly].move = NoMove
	pos.history[pos.histPly].fiftyMove = pos.fiftyMove
	pos.history[pos.histPly].enPas = pos.enPas
	pos.history[pos.histPly].castlePerm = pos.castlePerm
	pos.enPas = NoSquare

	pos.side ^= 1
	pos.histPly++
	hashSide(pos)

	// AssertTrue(CheckBoard(pos))
}

// TakeNullMove take back null move
func TakeNullMove(pos *Board) {
	// AssertTrue(CheckBoard(pos))

	pos.histPly--
	pos.ply--

	if pos.enPas != NoSquare {
		hashEnPass(pos)
	}

	pos.castlePerm = pos.history[pos.histPly].castlePerm
	pos.fiftyMove = pos.history[pos.histPly].fiftyMove
	pos.enPas = pos.history[pos.histPly].enPas

	if pos.enPas != NoSquare {
		hashEnPass(pos)
	}

	pos.side ^= 1
	hashSide(pos)

	// AssertTrue(CheckBoard(pos))
}
