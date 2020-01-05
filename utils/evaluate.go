package utils

// PawnTable pawn table
var PawnTable = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 0, -10, -10, 0, 10, 10,
	5, 0, 0, 5, 5, 0, 0, 5,
	0, 0, 10, 20, 20, 10, 0, 0,
	5, 5, 5, 10, 10, 5, 5, 5,
	10, 10, 10, 20, 20, 10, 10, 10,
	20, 20, 20, 30, 30, 20, 20, 20,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// KnightTable knight table
var KnightTable = [64]int{
	0, -10, 0, 0, 0, 0, -10, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 10, 10, 10, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	5, 10, 15, 20, 20, 15, 10, 5,
	5, 10, 10, 20, 20, 10, 10, 5,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// BishopTable bishop table
var BishopTable = [64]int{
	0, 0, -10, 0, 0, -10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// RookTable rook table
var RookTable = [64]int{
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	25, 25, 25, 25, 25, 25, 25, 25,
	0, 0, 5, 10, 10, 5, 0, 0,
}

// KingE king endgame table
var KingE = [64]int{
	-50, -10, 0, 0, 0, 0, -10, -50,
	-10, 0, 10, 10, 10, 10, 0, -10,
	0, 10, 15, 15, 15, 15, 10, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 10, 15, 15, 15, 15, 10, 0,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-50, -10, 0, 0, 0, 0, -10, -50,
}

// KingO king opening/middle game table
var KingO = [64]int{
	0, 10, 10, -10, -10, 0, 20, 10,
	-30, -30, -30, -30, -30, -30, -30, -30,
	-50, -50, -50, -50, -50, -50, -50, -50,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
}

// Mirror64 slice that is used to get a mirror version of the tables for black's evaluation
var Mirror64 = [64]int{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

// PawnPassed passed pawn bonuses depending on how far down the board it is
var PawnPassed = [8]int{0, 5, 10, 20, 35, 60, 100, 200}

const (
	// PawnIsolated isolated pawn bonus
	PawnIsolated = -10
	// PawnDoubled doubled pawn bonus
	PawnDoubled = -10
	// RookOpenFile rook on open file bonus
	RookOpenFile = 10
	// RookSemiOpenfile rook on semi-open file bonus
	RookSemiOpenfile = 5
	// QueenOpenFile queen on open file bonus
	QueenOpenFile = 5
	// QueenSemiOpenFile queen on semi-open file bonus
	QueenSemiOpenFile = 3
	// BishopPair bonus
	BishopPair = 30
	// KingNearOpenFile king on or near open file bonus
	KingNearOpenFile = -10
)

// EndGameMaterial defines the boundary limit for the endgame
var EndGameMaterial = 1*PieceVal[WhiteRook] + 2*PieceVal[WhiteKnight] + 2*PieceVal[WhitePawn] + PieceVal[WhiteKing]

// EvalPosition evaluate position and return value
func EvalPosition(pos *Board) int {
	score := pos.material[White] - pos.material[Black] // get current score

	// !!! FIX THIS !!!!
	// if MaterialDraw(pos) == true && pos.pieceNum[WhitePawn] == 0 && pos.pieceNum[BlackPawn] == 0 {
	// 	return 0
	// }

	for sq := 0; sq < BoardSquareNum; sq++ {
		switch pos.Pieces[sq] {
		case OffBoard, Empty:
			continue
		case WhitePawn:
			evalWhitePawn(pos, sq, &score)
		case BlackPawn:
			evalBlackPawn(pos, sq, &score)
		case WhiteKnight:
			score += KnightTable[Sq64(sq)]
		case BlackKnight:
			score -= KnightTable[Mirror64[Sq64(sq)]]
		case WhiteBishop:
			score += BishopTable[Sq64(sq)]
		case BlackBishop:
			score -= BishopTable[Mirror64[Sq64(sq)]]
		case WhiteRook:
			evalWhiteRook(pos, sq, &score)
		case BlackRook:
			evalBlackRook(pos, sq, &score)
		case WhiteQueen:
			evalWhiteQueen(pos, sq, &score)
		case BlackQueen:
			evalBlackQueen(pos, sq, &score)
		case WhiteKing:
			evalWhiteKing(pos, sq, &score)
		case BlackKing:
			evalBlackKing(pos, sq, &score)
		default:
			panic("Unknown square type")
		}
	}

	if pos.pieceNum[WhiteBishop] >= 2 {
		score += BishopPair
	}
	if pos.pieceNum[BlackBishop] >= 2 {
		score -= BishopPair
	}

	if pos.side == White {
		return score
	}
	return -score
}

// evalWhitePawn evaluate white pawn at square (120 based)
func evalWhitePawn(pos *Board, sq int, score *int) {
	*score += PawnTable[Sq64(sq)]

	// if there is no same colour pawn on a neighbouring square -> isolated pawn
	if IsolatedMask[Sq64(sq)]&pos.Pawns[White] == 0 {
		// fmt.Printf("WhitePawn Iso:%s\n", PrintSquare(sq))
		*score += PawnIsolated
	}

	if WhitePassedMask[Sq64(sq)]&pos.Pawns[Black] == 0 {
		// fmt.Printf("WhitePawn passed:%s\n", PrintSquare(sq))
		*score += PawnPassed[RanksBoard[sq]]
	}

	if WhiteDoubledMask[Sq64(sq)]&pos.Pawns[White] != 0 {
		// fmt.Printf("WhitePawn doubled: %s\n", PrintSquare(sq))
		*score += PawnDoubled
	}
}

// evalBlackPawn evaluate black pawn at square (120 based)
func evalBlackPawn(pos *Board, sq int, score *int) {
	*score -= PawnTable[Mirror64[Sq64(sq)]]

	// if there is no same colour pawn on a neighbouring square -> isolated pawn
	if IsolatedMask[Sq64(sq)]&pos.Pawns[Black] == 0 {
		// fmt.Printf("BlackPawn Iso:%s\n", PrintSquare(sq))
		*score -= PawnIsolated
	}

	// if there are no opposite color pawns in front or to the side of the black pawn -> passed pawn
	if BlackPassedMask[Sq64(sq)]&pos.Pawns[White] == 0 {
		// fmt.Printf("BlackPawn passed:%s\n", PrintSquare(sq))
		*score -= PawnPassed[7-RanksBoard[sq]] // need to subtract from 7 since black pawns go down the rank from 7->0
	}

	if BlackDoubledMask[Sq64(sq)]&pos.Pawns[Black] != 0 {
		// fmt.Printf("BlackPawn doubled: %s\n", PrintSquare(sq))
		*score -= PawnDoubled
	}
}

// evalWhiteRook evaluate white rook at square (120 based)
func evalWhiteRook(pos *Board, sq int, score *int) {
	*score += RookTable[Sq64(sq)]

	if (pos.Pawns[Both] & FileBBMask[FilesBoard[sq]]) == 0 {
		*score += RookOpenFile
	} else if (pos.Pawns[White] & FileBBMask[FilesBoard[sq]]) == 0 {
		// none of our pawns on the file -> semiopen file bonus
		*score += RookSemiOpenfile
	}
}

// evalBlackRook evaluate black rook at square (120 based)
func evalBlackRook(pos *Board, sq int, score *int) {
	*score -= RookTable[Mirror64[Sq64(sq)]]

	if (pos.Pawns[Both] & FileBBMask[FilesBoard[sq]]) == 0 {
		*score -= RookOpenFile
	} else if (pos.Pawns[Black] & FileBBMask[FilesBoard[sq]]) == 0 {
		// none of our pawns on the file -> semiopen file bonus
		*score -= RookSemiOpenfile
	}
}

// evalWhiteQueen evaluate white queen at square (120 based)
func evalWhiteQueen(pos *Board, sq int, score *int) {
	if (pos.Pawns[Both] & FileBBMask[FilesBoard[sq]]) == 0 {
		*score += QueenOpenFile
	} else if (pos.Pawns[White] & FileBBMask[FilesBoard[sq]]) == 0 {
		// none of our pawns on the file -> semiopen file bonus
		*score += QueenSemiOpenFile
	}
}

// evalBlackQueen evaluate black queen at square (120 based)
func evalBlackQueen(pos *Board, sq int, score *int) {
	if (pos.Pawns[Both] & FileBBMask[FilesBoard[sq]]) == 0 {
		*score -= QueenOpenFile
	} else if (pos.Pawns[Black] & FileBBMask[FilesBoard[sq]]) == 0 {
		// none of our pawns on the file -> semiopen file bonus
		*score -= QueenSemiOpenFile
	}
}

// evalWhiteKing evaluate white king at square (120 based)
func evalWhiteKing(pos *Board, sq int, score *int) {
	// if we are in the endgame add endgame king square score
	if pos.material[Black] <= EndGameMaterial {
		*score += KingE[Sq64(sq)]
	} else {
		*score += KingO[Sq64(sq)]
		*score += evalWhiteKingSafety(pos)
	}
}

// evalBlackKing evaluate white king at square (120 based)
func evalBlackKing(pos *Board, sq int, score *int) {
	if pos.material[White] <= EndGameMaterial {
		*score -= KingE[Mirror64[Sq64(sq)]]
	} else {
		*score -= KingO[Mirror64[Sq64(sq)]]
		*score -= evalBlackKingSafety(pos)
	}
}

func evalWhiteKingSafety(pos *Board) int {
	safetyScore := 0
	kingSq := pos.kingSquare[White]

	// if king is castled short
	if pos.kingSquare[White] > E1 && pos.kingSquare[White] <= H1 {
		safetyScore += evalWhiteKingPawn(pos, H2)
		safetyScore += evalWhiteKingPawn(pos, G2)
		safetyScore += evalWhiteKingPawn(pos, F2) / 2 // F and C pawns are not that severe
	} else if pos.kingSquare[White] < D1 && pos.kingSquare[White] >= A1 {
		safetyScore += evalWhiteKingPawn(pos, A2)
		safetyScore += evalWhiteKingPawn(pos, B2)
		safetyScore += evalWhiteKingPawn(pos, C2) / 2 // F and C pawns are not that severe
	} else {
		// if the king is on an open file
		if (pos.Pawns[Both] & FileBBMask[FilesBoard[kingSq]]) == 0 {
			// fmt.Printf("White King file is open (%s)\n", PrintSquare(kingSq))
			safetyScore += KingNearOpenFile
		}
		// if file to the left of the king is open
		if SquareOnBoard(kingSq - 1) {
			if (pos.Pawns[Both] & FileBBMask[FilesBoard[kingSq-1]]) == 0 {
				// fmt.Printf("White King file-1 is open (%s)\n", PrintSquare(kingSq-1))
				safetyScore += KingNearOpenFile
			}
		}
		// if file to the right of the king is open
		if SquareOnBoard(kingSq + 1) {
			if (pos.Pawns[Both] & FileBBMask[FilesBoard[kingSq+1]]) == 0 {
				// fmt.Printf("White King file+1 is open (%s)\n", PrintSquare(kingSq+1))
				safetyScore += KingNearOpenFile
			}
		}
	}

	// scale the king safety value according to the opponent's material;
	// the premise is that your king safety can only be bad if the
	// opponent has enough pieces to attack you
	// safetyScore *= pos.material[Black]
	// safetyScore /= 3100

	return safetyScore
}

func evalWhiteKingPawn(pos *Board, pawn int) int {
	safetyScore := 0

	// if pawn hasnt moved
	if (pos.Pawns[White] & (1 << uint64(Sq64(pawn)))) != 0 {
		// everything is okay
		// fmt.Printf("White king pawn hasnt moved (%s)\n", PrintSquare(pawn))
	} else if (pos.Pawns[White] & (1 << uint64(Sq64(pawn+10)))) != 0 {
		// if pawn moved one square forward add -10
		// fmt.Printf("White king pawn moved one square forward (%s)\n", PrintSquare(pawn+10))
		safetyScore -= 10
	} else if (pos.Pawns[White] & BlackDoubledMask[Sq64(pawn)]) != 0 {
		// pawn moved more than one square
		// fmt.Printf("White king pawn moved more than one square forward (%s)\n", PrintSquare(pawn))
		safetyScore -= 20
	} else {
		// fmt.Printf("No white king pawn (%s).\n", PrintSquare(pawn))
		// no pawn on this file
		safetyScore -= 25
	}

	//check for opposing pawns close to the white king pawns
	if (pos.Pawns[Black] & (1 << uint64(Sq64(pawn+10)))) != 0 {
		// enemy pawn on the 3rd rank
		// fmt.Printf("Enemy pawn on 3rd rank (%s)\n", PrintSquare(pawn+10))
		safetyScore -= 10
	} else if (pos.Pawns[Black] & (1 << uint64(Sq64(pawn+20)))) != 0 {
		// enemy pawn on the 4th rank
		// fmt.Printf("Enemy pawn on 4th rank (%s)\n", PrintSquare(pawn+20))
		safetyScore -= 5
	}

	return safetyScore
}

func evalBlackKingSafety(pos *Board) int {
	safetyScore := 0
	kingSq := pos.kingSquare[Black]

	// if king is castled short
	if kingSq > E8 && kingSq <= H8 {
		safetyScore += evalBlackKingPawn(pos, H7)
		safetyScore += evalBlackKingPawn(pos, G7)
		safetyScore += evalBlackKingPawn(pos, F7) / 2
	} else if kingSq < D8 && kingSq >= A8 {
		// if king is castled
		safetyScore += evalBlackKingPawn(pos, A7)
		safetyScore += evalBlackKingPawn(pos, B7)
		safetyScore += evalBlackKingPawn(pos, C7) / 2 // F and C pawns are not that severe
	} else {
		// if the king is on an open file
		if (pos.Pawns[Both] & FileBBMask[FilesBoard[kingSq]]) == 0 {
			// fmt.Printf("Black King file is open (%s)\n", PrintSquare(kingSq))
			safetyScore += KingNearOpenFile
		}
		// if file to the left of the king is open
		if SquareOnBoard(kingSq - 1) {
			if (pos.Pawns[Both] & FileBBMask[FilesBoard[kingSq-1]]) == 0 {
				// fmt.Printf("Black King file-1 is open (%s)\n", PrintSquare(kingSq-1))
				safetyScore += KingNearOpenFile
			}
		}
		// if file to the right of the king is open
		if SquareOnBoard(kingSq + 1) {
			if (pos.Pawns[Both] & FileBBMask[FilesBoard[kingSq+1]]) == 0 {
				// fmt.Printf("Black King file+1 is open (%s)\n", PrintSquare(kingSq+1))
				safetyScore += KingNearOpenFile
			}
		}
	}

	// scale the king safety value according to the opponent's material;
	// the premise is that your king safety can only be bad if the
	// opponent has enough pieces to attack you
	// safetyScore *= pos.material[White]
	// safetyScore /= 3100

	return safetyScore
}

func evalBlackKingPawn(pos *Board, pawn int) int {
	safetyScore := 0

	// if pawn hasnt moved
	if (pos.Pawns[Black] & (1 << uint64(Sq64(pawn)))) != 0 {
		// everything is okay
		// fmt.Printf("Black king pawn hasnt moved (%s)\n", PrintSquare(pawn))
	} else if (pos.Pawns[Black] & (1 << uint64(Sq64(pawn-10)))) != 0 {
		// if pawn moved one square forward add -10
		// fmt.Printf("Black king pawn moved one square forward (%s)\n", PrintSquare(pawn-10))
		safetyScore -= 10
	} else if (pos.Pawns[Black] & WhiteDoubledMask[Sq64(pawn)]) != 0 {
		// pawn moved more than one square
		// fmt.Printf("Black king pawn moved more than one square forward (%s)\n", PrintSquare(pawn))
		safetyScore -= 20
	} else {
		// fmt.Printf("No Black king pawn (%s).\n", PrintSquare(pawn))
		// no pawn on this file
		safetyScore -= 25
	}

	//check for opposing pawns close to the white king pawns
	if (pos.Pawns[White] & (1 << uint64(Sq64(pawn-10)))) != 0 {
		// enemy pawn on the 3rd rank
		// fmt.Printf("Enemy pawn on 3rd rank (%s)\n", PrintSquare(pawn-10))
		safetyScore -= 10
	} else if (pos.Pawns[White] & (1 << uint64(Sq64(pawn-20)))) != 0 {
		// enemy pawn on the 4th rank
		// fmt.Printf("Enemy pawn on 4th rank (%s)\n", PrintSquare(pawn-20))
		safetyScore -= 5
	}

	return safetyScore
}

// MirrorBoard takes in a position and modifies it to be the mirrored version of it
func MirrorBoard(pos *Board) {
	var swapPiece = [13]int{Empty, BlackPawn, BlackKnight, BlackBishop, BlackRook, BlackQueen, BlackKing, WhitePawn, WhiteKnight, WhiteBishop, WhiteRook, WhiteQueen, WhiteKing}
	var tempPiecesSlice [64]int
	var tempSide = pos.side ^ 1
	var tempCastlePerm = 0
	var tempEnPassant = NoSquare

	if pos.castlePerm&WhiteKingCastling != 0 {
		tempCastlePerm |= BlackKingCastling
	}

	if pos.castlePerm&WhiteQueenCastling != 0 {
		tempCastlePerm |= BlackQueenCastling
	}

	if pos.castlePerm&BlackKingCastling != 0 {
		tempCastlePerm |= WhiteKingCastling
	}

	if pos.castlePerm&BlackQueenCastling != 0 {
		tempCastlePerm |= WhiteQueenCastling
	}

	if pos.enPas != NoSquare {
		tempEnPassant = Sq120(Mirror64[Sq64(pos.enPas)])
	}

	for sq := 0; sq < 64; sq++ {
		tempPiecesSlice[sq] = pos.Pieces[Sq120(Mirror64[sq])]
	}

	// clear board
	ResetBoard(pos)

	// write mirrored information to all relevant arrays
	for sq := 0; sq < 64; sq++ {
		tempPiece := swapPiece[tempPiecesSlice[sq]]
		pos.Pieces[Sq120(sq)] = tempPiece
	}

	pos.side = tempSide
	pos.castlePerm = tempCastlePerm
	pos.enPas = tempEnPassant

	pos.posKey = GeneratePosKey(pos)

	UpdateListsMaterial(pos)

	// // AssertTrue(CheckBoard(pos))
}

// abs local method to compute absolute value of int without needing to convert to float
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// MaterialDraw Determines if given the available pieces the position is a material draw, based on sjeng
func MaterialDraw(pos *Board) bool {
	if pos.pieceNum[WhiteRook] == 0 && pos.pieceNum[BlackRook] == 0 && pos.pieceNum[WhiteQueen] == 0 && pos.pieceNum[BlackQueen] == 0 {
		if pos.pieceNum[BlackBishop] == 0 && pos.pieceNum[WhiteBishop] == 0 {
			if pos.pieceNum[WhiteKnight] < 3 && pos.pieceNum[BlackKnight] < 3 {
				return true
			}
		} else if pos.pieceNum[WhiteKnight] == 0 && pos.pieceNum[BlackKnight] == 0 {
			if abs(pos.pieceNum[WhiteBishop]-pos.pieceNum[BlackBishop]) < 2 {
				return true
			}
		} else if (pos.pieceNum[WhiteKnight] < 3 && pos.pieceNum[WhiteBishop] == 0) || (pos.pieceNum[BlackBishop] == 1 && pos.pieceNum[WhiteKnight] == 0) {
			if (pos.pieceNum[BlackKnight] < 3 && pos.pieceNum[BlackBishop] == 0) || (pos.pieceNum[BlackBishop] == 1 && pos.pieceNum[BlackKnight] == 0) {
				return true
			}
		}
	} else if pos.pieceNum[WhiteQueen] == 0 && pos.pieceNum[BlackQueen] == 0 {
		if pos.pieceNum[WhiteRook] == 1 && pos.pieceNum[BlackRook] == 1 {
			if (pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) < 2 && (pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) < 2 {
				return true
			}
		} else if pos.pieceNum[WhiteRook] == 1 && pos.pieceNum[BlackRook] == 0 {
			if (pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) == 0 && ((pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) == 1 || (pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) == 2) {
				return true
			}
		} else if pos.pieceNum[WhiteRook] == 0 && pos.pieceNum[BlackRook] == 1 {
			if (pos.pieceNum[BlackKnight]+pos.pieceNum[BlackBishop]) == 0 && ((pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) == 1 || (pos.pieceNum[WhiteKnight]+pos.pieceNum[WhiteBishop]) == 2) {
				return true
			}
		}
	}
	return false
}
