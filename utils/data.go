package utils

// PieceChar string with piece characters
var PieceChar = ".PNBRQKpnbrqk"

// SideChar string with side characters
var SideChar = "wb-"

// RankChar string with rank characters
var RankChar = "12345678"

// FileChar string with file characters
var FileChar = "abcdefgh"

// PieceBig A map used to identify if a piece is considered "Big"
var PieceBig = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: true,
	WhiteBishop: true,
	WhiteRook:   true,
	WhiteQueen:  true,
	WhiteKing:   true,
	BlackPawn:   false,
	BlackKnight: true,
	BlackBishop: true,
	BlackRook:   true,
	BlackQueen:  true,
	BlackKing:   true,
}

// PieceMaj A map used to identify if a piece is considered "Major"
var PieceMaj = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: false,
	WhiteRook:   true,
	WhiteQueen:  true,
	WhiteKing:   true,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: false,
	BlackRook:   true,
	BlackQueen:  true,
	BlackKing:   true,
}

// PieceMin A map used to identify if a piece is considered "Minor"
var PieceMin = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: true,
	WhiteBishop: true,
	WhiteRook:   false,
	WhiteQueen:  false,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: true,
	BlackBishop: true,
	BlackRook:   false,
	BlackQueen:  false,
	BlackKing:   false,
}

// PieceVal A map used to identify a piece's value
var PieceVal = map[int]int{
	Empty:       0,
	WhitePawn:   100,
	WhiteKnight: 325,
	WhiteBishop: 325,
	WhiteRook:   550,
	WhiteQueen:  1000,
	WhiteKing:   50000,
	BlackPawn:   100,
	BlackKnight: 325,
	BlackBishop: 325,
	BlackRook:   550,
	BlackQueen:  1000,
	BlackKing:   50000,
}

// PieceColour A map used to identify a piece's colour
var PieceColour = map[int]int{
	Empty:       Both,
	WhitePawn:   White,
	WhiteKnight: White,
	WhiteBishop: White,
	WhiteRook:   White,
	WhiteQueen:  White,
	WhiteKing:   White,
	BlackPawn:   Black,
	BlackKnight: Black,
	BlackBishop: Black,
	BlackRook:   Black,
	BlackQueen:  Black,
	BlackKing:   Black,
}

// IsPieceKnight holds information if a given piece is a knight
var IsPieceKnight = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: true,
	WhiteBishop: false,
	WhiteRook:   false,
	WhiteQueen:  false,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: true,
	BlackBishop: false,
	BlackRook:   false,
	BlackQueen:  false,
	BlackKing:   false,
}

// IsPieceKing holds information if a given piece is a king
var IsPieceKing = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: false,
	WhiteRook:   false,
	WhiteQueen:  false,
	WhiteKing:   true,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: false,
	BlackRook:   false,
	BlackQueen:  false,
	BlackKing:   true,
}

// IsPieceRookQueen holds information if a given piece is a rook or queen
var IsPieceRookQueen = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: false,
	WhiteRook:   true,
	WhiteQueen:  true,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: false,
	BlackRook:   true,
	BlackQueen:  true,
	BlackKing:   false,
}

// IsPieceBishopQueen holds information if a given piece is a bishop or queen
var IsPieceBishopQueen = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: true,
	WhiteRook:   false,
	WhiteQueen:  true,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: true,
	BlackRook:   false,
	BlackQueen:  true,
	BlackKing:   false,
}

// IsPiecePawn holds information if a given piece is a pawn
var IsPiecePawn = map[int]bool{
	Empty:       false,
	WhitePawn:   true,
	WhiteKnight: false,
	WhiteBishop: false,
	WhiteRook:   false,
	WhiteQueen:  false,
	WhiteKing:   false,
	BlackPawn:   true,
	BlackKnight: false,
	BlackBishop: false,
	BlackRook:   false,
	BlackQueen:  false,
	BlackKing:   false,
}

// PieceSlides holds information if a given piece slides
var PieceSlides = map[int]bool{
	Empty:       false,
	WhitePawn:   false,
	WhiteKnight: false,
	WhiteBishop: true,
	WhiteRook:   true,
	WhiteQueen:  true,
	WhiteKing:   false,
	BlackPawn:   false,
	BlackKnight: false,
	BlackBishop: true,
	BlackRook:   true,
	BlackQueen:  true,
	BlackKing:   false,
}

// LoopSlidePiece sliding pieces slice used for looping
var LoopSlidePiece = [...]int{WhiteBishop, WhiteRook, WhiteQueen, 0, BlackBishop, BlackRook, BlackQueen, 0}

// LoopSlideIndex sliding pieces index slice to index where
// the white pieces start in the above LoopSlidePiece, and where black
var LoopSlideIndex = [...]int{0, 4}

// LoopNonSlidePiece non-sliding pieces slice used for looping
var LoopNonSlidePiece = [...]int{WhiteKnight, WhiteKing, 0, BlackKnight, BlackKing, 0}

// LoopNonSlideIndex non-sliding pieces index slice to index where
// the white pieces start in the above LoopSlidePiece, and where black
var LoopNonSlideIndex = [...]int{0, 3}

// PiececeDir squares increment for each direction
var PiececeDir = map[int][]int{
	Empty:       {0, 0, 0, 0, 0, 0, 0},
	WhitePawn:   {0, 0, 0, 0, 0, 0, 0},
	WhiteKnight: {-8, -19, -21, -12, 8, 19, 21, 12},
	WhiteBishop: {-9, -11, 11, 9, 0, 0, 0, 0},
	WhiteRook:   {-1, -10, 1, 10, 0, 0, 0, 0},
	WhiteQueen:  {-1, -10, 1, 10, -9, -11, 11, 9},
	WhiteKing:   {-1, -10, 1, 10, -9, -11, 11, 9},
	BlackPawn:   {0, 0, 0, 0, 0, 0, 0},
	BlackKnight: {-8, -19, -21, -12, 8, 19, 21, 12},
	BlackBishop: {-9, -11, 11, 9, 0, 0, 0, 0},
	BlackRook:   {-1, -10, 1, 10, 0, 0, 0, 0},
	BlackQueen:  {-1, -10, 1, 10, -9, -11, 11, 9},
	BlackKing:   {-1, -10, 1, 10, -9, -11, 11, 9},
}

// NumberOfDir number of directions in which each piece can move
var NumberOfDir = map[int]int{
	Empty:       0,
	WhitePawn:   0,
	WhiteKnight: 8,
	WhiteBishop: 4,
	WhiteRook:   4,
	WhiteQueen:  8,
	WhiteKing:   8,
	BlackPawn:   0,
	BlackKnight: 8,
	BlackBishop: 4,
	BlackRook:   4,
	BlackQueen:  8,
	BlackKing:   8,
}
