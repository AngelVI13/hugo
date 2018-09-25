package utils

// SquareOnBoard returns true if square is on board
func SquareOnBoard(sq int) bool {
	if FilesBoard[sq] == OffBoard {
		return false
	}
	return true
}

// SideValid returns true if side is valid
func SideValid(side int) bool {
	if side == White || side == Black {
		return true
	}
	return false
}

// FileRankValid returns true if a file/rank is valid
func FileRankValid(fr int) bool {
	if fr >= 0 && fr <= 7 {
		return true
	}
	return false
}

// PieceValidEmpty returns true if a piece is valid (empty is also included as valid)
func PieceValidEmpty(piece int) bool {
	if piece >= Empty && piece <= BlackKing {
		return true
	}
	return false
}

// PieceValid returns true if piece is valid (empty state not included)
func PieceValid(piece int) bool {
	if piece >= WhitePawn && piece <= BlackKing {
		return true
	}
	return false
}
