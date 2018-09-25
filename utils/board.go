package utils

import (
	"fmt"
	"strconv"
)

// ResetBoard resets a given board
func ResetBoard(pos *Board) {
	// Set all board positions to OffBoard
	for i := 0; i < BoardSquareNum; i++ {
		pos.Pieces[i] = OffBoard
	}

	// Set all real board positions to Empty
	for i := 0; i < 64; i++ {
		pos.Pieces[Sq120(i)] = Empty
	}

	// Reset number of big pieces, major pieces, minor pieces and pawns
	for i := 0; i < 2; i++ {
		pos.bigPieceNum[i] = 0
		pos.majorPieceNum[i] = 0
		pos.minorPieceNum[i] = 0
		pos.material[i] = 0
	}

	// The pawns slice contains information for White,Black & Both
	for i := 0; i < 3; i++ {
		pos.Pawns[i] = 0
	}

	// Reset piece number
	for i := 0; i < 13; i++ {
		pos.pieceNum[i] = 0
	}

	pos.kingSquare[White] = NoSquare
	pos.kingSquare[Black] = NoSquare

	pos.side = Both
	pos.enPas = NoSquare
	pos.fiftyMove = 0

	pos.ply = 0
	pos.histPly = 0

	pos.castlePerm = 0

	pos.posKey = 0
}

// ParseFen parse fen position string and setup a position accordingly
// !!! FIX ERROR HANDLING, now it simply returns a non-zero int whenever there is an error
func ParseFen(fen string, pos *Board) int {
	// // AssertTrue(fen != "")
	// // AssertTrue(pos != nil)

	rank := Rank8 // we start from rank 8 since the notation starts from rank 8
	file := FileA
	piece := 0
	count := 0 // number of empty squares declared inside fen string
	sq64 := 0
	sq120 := 0

	ResetBoard(pos)
	char := 0

	for (rank >= Rank1) && char < len(fen) {
		count = 1
		switch t := string(fen[char]); t {
		case "p", "r", "n", "b", "k", "q", "P", "R", "N", "B", "K", "Q":
			// If we have a piece related char -> set the piece to corresponding value, i.e p -> BlackPawn
			piece = PieceNotationMap[t]
		case "1", "2", "3", "4", "5", "6", "7", "8":
			// otherwise it must be a count of a number of empty squares
			piece = Empty
			count, _ = strconv.Atoi(t) // get number of empty squares and store in count
		case "/", " ":
			// if we have / or space then we are either at the end of the rank or at the end of the piece list
			// -> reset variables and continue the while loop
			rank--
			file = FileA
			char++
			continue
		default:
			fmt.Println("FEN error")
			return -1
		}

		// This loop, skips over all empty positions in a rank
		// When it comes to a piece that is different that "1"-"8" it places it on the corresponding square
		for i := 0; i < count; i++ {
			sq64 = rank*8 + file
			sq120 = Sq120(sq64)
			if piece != Empty {
				pos.Pieces[sq120] = piece
			}
			file++
		}
		char++
	}

	newChar := ""
	// newChar should be set to the side to move part of the FEN string here
	newChar = string(fen[char])
	// // AssertTrue(newChar == "w" || newChar == "b")

	if newChar == "w" {
		pos.side = White
	} else {
		pos.side = Black
	}

	// move character pointer 2 characters further and it should now point to the start of the castling permissions part of FEN
	char += 2

	// Iterate over the next 4 chars - they show if white is allowed to castle king or quenside and the same for black
	for i := 0; i < 4; i++ {
		newChar = string(fen[char])
		if newChar == " " {
			// when we hit a space, it means there are no more castling permissions => break
			break
		}
		switch newChar { // Depending on the char, enable the corresponding castling permission related bit
		case "K":
			pos.castlePerm |= WhiteKingCastling
		case "Q":
			pos.castlePerm |= WhiteQueenCastling
		case "k":
			pos.castlePerm |= BlackKingCastling
		case "q":
			pos.castlePerm |= BlackQueenCastling
		default:
			break
		}
		char++
	}

	// // AssertTrue(pos.castlePerm >= 0 && pos.castlePerm <= 15)
	// move to the en passant square related part of FEN
	char++
	newChar = string(fen[char])

	if newChar != "-" {
		file := FileNotationMap[newChar]
		char++
		rank, _ := strconv.Atoi(string(fen[char])) // get rank number
		rank--                                     // decrement rank to match our indexes, i.e. Rank1 == 0

		// // AssertTrue(file >= FileA && file <= FileH)
		// // AssertTrue(rank >= Rank1 && rank <= Rank8)

		pos.enPas = FileRankToSquare(file, rank)
	}

	pos.posKey = GeneratePosKey(pos) // generate pos key for new position

	UpdateListsMaterial(pos)

	return 0
}

// PrintBoard prints board for a given position
func PrintBoard(pos *Board) {
	fmt.Printf("\nGame Board:\n\n")

	for rank := Rank8; rank >= Rank1; rank-- {
		fmt.Printf("%d  ", rank+1)
		for file := FileA; file <= FileH; file++ {
			sq := FileRankToSquare(file, rank)
			piece := pos.Pieces[sq]
			fmt.Printf("%3c", PieceChar[piece])
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n   ")
	for file := FileA; file <= FileH; file++ {
		fmt.Printf("%3c", 'a'+file)
	}
	fmt.Printf("\n")
	fmt.Printf("side:%c\n", SideChar[pos.side])
	fmt.Printf("enPas:%d\n", pos.enPas)

	// Compute castling permissions
	wKCA := "-"
	if pos.castlePerm&WhiteKingCastling != 0 {
		wKCA = "K"
	}

	wQCA := "-"
	if pos.castlePerm&WhiteQueenCastling != 0 {
		wQCA = "Q"
	}

	bKCA := "-"
	if pos.castlePerm&BlackKingCastling != 0 {
		bKCA = "k"
	}

	bQCA := "-"
	if pos.castlePerm&BlackQueenCastling != 0 {
		bQCA = "q"
	}

	fmt.Printf("castle:%s%s%s%s\n", wKCA, wQCA, bKCA, bQCA)
	fmt.Printf("PosKey:%X\n", pos.posKey)
}

// UpdateListsMaterial updates all material related piece lists
func UpdateListsMaterial(pos *Board) {
	for index := 0; index < BoardSquareNum; index++ {
		piece := pos.Pieces[index]
		if piece != OffBoard && piece != Empty {
			colour := PieceColour[piece]

			if PieceBig[piece] == true {
				pos.bigPieceNum[colour]++
			}
			if PieceMin[piece] == true {
				pos.minorPieceNum[colour]++
			}
			if PieceMaj[piece] == true {
				pos.majorPieceNum[colour]++
			}

			pos.material[colour] += PieceVal[piece]

			// piece list [13][10] i.e.
			// pieceList[WhitePawn][0] = a1
			// pieceList[WhitePawn][0] = a2 etc
			pos.pieceList[piece][pos.pieceNum[piece]] = index
			pos.pieceNum[piece]++ // increment piece number

			if piece == WhiteKing || piece == BlackKing {
				pos.kingSquare[colour] = index
			}

			// If we have a pawn, set the bit corresponding to the board position of the pawn in the related pawn bitboard
			if piece == WhitePawn {
				SetBit(&pos.Pawns[White], Sq64(index))
				SetBit(&pos.Pawns[Both], Sq64(index))
			} else if piece == BlackPawn {
				SetBit(&pos.Pawns[Black], Sq64(index))
				SetBit(&pos.Pawns[Both], Sq64(index))
			}
		}
	}
}

// InitFilesRanksBoard initialize arrays that hold information about which rank & file a square is on the board
func InitFilesRanksBoard() {
	// Set all square to OffBoard
	for index := 0; index < BoardSquareNum; index++ {
		FilesBoard[index] = OffBoard
		RanksBoard[index] = OffBoard
	}

	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			sq := FileRankToSquare(file, rank)
			FilesBoard[sq] = file
			RanksBoard[sq] = rank
		}
	}
}

// // CheckBoard makes a checkboard
// func CheckBoard(pos *Board) bool {
// 	// setup temporary variables
// 	tPceNum := [13]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
// 	tBigPce := [2]int{0, 0}
// 	tMajPce := [2]int{0, 0}
// 	tMinPce := [2]int{0, 0}
// 	tMaterial := [2]int{0, 0}
// 	tPawns := [3]uint64{0, 0, 0} // !!! Here initializing to 0 is not needed, you can straight away initialize to pos.Pawns[...]

// 	tPawns[White] = pos.Pawns[White]
// 	tPawns[Black] = pos.Pawns[Black]
// 	tPawns[Both] = pos.Pawns[Both]

// 	// check piece lists
// 	// Iterate over all pieces i.e. WhitePawn, and then iterate on every one of the 8 white pawns and check if their positions
// 	// stored in the pieceList array, matches what is stored in the Pieces array (where the whole board should be represented)
// 	for tPiece := WhitePawn; tPiece <= BlackKing; tPiece++ {
// 		for tPceNum := 0; tPceNum < pos.pieceNum[tPiece]; tPceNum++ {
// 			sq120 := pos.pieceList[tPiece][tPceNum]
// 			// AssertTrue(pos.Pieces[sq120] == tPiece)
// 		}
// 	}

// 	// check piece count and other counters
// 	// similar to update material func
// 	for sq64 := 0; sq64 < 64; sq64++ {
// 		sq120 := Sq120(sq64)
// 		tPiece := pos.Pieces[sq120]
// 		tPceNum[tPiece]++

// 		// No need to update piece numbers if piece is Empty
// 		if tPiece == Empty {
// 			continue
// 		}

// 		colour := PieceColour[tPiece]

// 		if PieceBig[tPiece] == true {
// 			tBigPce[colour]++
// 		}
// 		if PieceMin[tPiece] == true {
// 			tMinPce[colour]++
// 		}
// 		if PieceMaj[tPiece] == true {
// 			tMajPce[colour]++
// 		}
// 		tMaterial[colour] += PieceVal[tPiece]
// 	}

// 	// Check that the number of pieces we found on the board and stored in tPceNum,
// 	// equals the number of pieces that our position says that we have
// 	for tPiece := WhitePawn; tPiece <= BlackKing; tPiece++ {
// 		// AssertTrue(tPceNum[tPiece] == pos.pieceNum[tPiece])
// 	}

// 	// check bitboards count
// 	pcount := CountBits(tPawns[White])
// 	// AssertTrue(pcount == pos.pieceNum[WhitePawn])
// 	pcount = CountBits(tPawns[Black])
// 	// AssertTrue(pcount == pos.pieceNum[BlackPawn])
// 	pcount = CountBits(tPawns[Both])
// 	// AssertTrue(pcount == (pos.pieceNum[BlackPawn] + pos.pieceNum[WhitePawn]))

// 	// check bitboards squares
// 	for tPawns[White] != 0 {
// 		// pop removes a bit from a bitboard and returns a 64bit square index value where that pawn is
// 		sq64 := PopBit(&tPawns[White])
// 		// AssertTrue(pos.Pieces[Sq120(sq64)] == WhitePawn)
// 	}

// 	for tPawns[Black] != 0 {
// 		sq64 := PopBit(&tPawns[Black])
// 		// AssertTrue(pos.Pieces[Sq120(sq64)] == BlackPawn)
// 	}

// 	for tPawns[Both] != 0 {
// 		sq64 := PopBit(&tPawns[Both])
// 		// AssertTrue((pos.Pieces[Sq120(sq64)] == BlackPawn) || (pos.Pieces[Sq120(sq64)] == WhitePawn))
// 	}

// 	// Check material and num of major, minor and big pieces match
// 	// AssertTrue(tMaterial[White] == pos.material[White] && tMaterial[Black] == pos.material[Black])
// 	// AssertTrue(tMinPce[White] == pos.minorPieceNum[White] && tMinPce[Black] == pos.minorPieceNum[Black])
// 	// AssertTrue(tMajPce[White] == pos.majorPieceNum[White] && tMajPce[Black] == pos.majorPieceNum[Black])
// 	// AssertTrue(tBigPce[White] == pos.bigPieceNum[White] && tBigPce[Black] == pos.bigPieceNum[Black])

// 	// AssertTrue(pos.side == White || pos.side == Black)
// 	// AssertTrue(GeneratePosKey(pos) == pos.posKey)

// 	// Check if en passant square is correct, it can only happen when white on 3rd rank and when black on 6th rank
// 	// AssertTrue(pos.enPas == NoSquare || (RanksBoard[pos.enPas] == Rank6 && pos.side == White) || (RanksBoard[pos.enPas] == Rank3 && pos.side == Black))

// 	// Check if at king positions in our position array, there is actualy a King piece
// 	// AssertTrue(pos.Pieces[pos.kingSquare[White]] == WhiteKing)
// 	// AssertTrue(pos.Pieces[pos.kingSquare[Black]] == BlackKing)

// 	return true
// }
