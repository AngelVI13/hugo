package utils

import (
	"fmt"
)

// pvtable.go - contains methods used to extract the principle variation i.e. the best line found for the position
const (
	HashSize int = 0x100000 * 10 // magic number for number of entries in the pv table (appox 16MB in the C-version)
)

// InitHashTable initializes pv table
func InitHashTable(table *TranspositionTable) {
	// given a const size, and the size of the pEntry struct, compute how many entries can fit into the given PvSize
	// -> create a slice of PVEntries with that size
	table.numEntries = HashSize
	table.pTable = make([]HashEntry, table.numEntries)
	ClearHashTable(table)
	fmt.Printf("HashTable init complete with %d entries\n", table.numEntries)
}

// ClearHashTable clears a given pvtable
func ClearHashTable(table *TranspositionTable) {
	for hashEntry := range table.pTable {
		table.pTable[hashEntry].posKey = 0
		table.pTable[hashEntry].move = NoMove
		table.pTable[hashEntry].depth = 0
		table.pTable[hashEntry].score = 0
		table.pTable[hashEntry].flags = 0
	}
}

// StoreHashEntry store a principle variation move
func StoreHashEntry(pos *Board, move, score, flags, depth int) {
	// this returns a number between 0 and numentries - 1
	index := int(pos.posKey % uint64(HashTable.numEntries))

	// AssertTrue(index >= 0 && index <= HashTable.numEntries-1)
	// AssertTrue(depth >= 1 && depth < MaxDepth)
	// AssertTrue(flags >= HFAlpha && flags <= HFExact)
	// AssertTrue(score >= -Infinite && score <= Infinite)
	// AssertTrue(pos.ply >= 0 && pos.ply < MaxDepth)

	if score > IsMate {
		score += pos.ply
	} else if score < -IsMate {
		score -= pos.ply
	}

	HashTable.pTable[index].move = move
	HashTable.pTable[index].posKey = pos.posKey
	HashTable.pTable[index].flags = flags
	HashTable.pTable[index].score = score
	HashTable.pTable[index].depth = depth
}

// ProbeHashEntry probe pv table. Return the principle variation move for a given position from the PV table
func ProbeHashEntry(pos *Board, move, score *int, alpha, beta, depth int) bool {
	// this returns a number between 0 and numentries - 1
	index := int(pos.posKey % uint64(HashTable.numEntries))

	// AssertTrue(index >= 0 && index <= HashTable.numEntries-1)
	// AssertTrue(depth >= 1 && depth < MaxDepth)
	// AssertTrue(alpha < beta)
	// AssertTrue(alpha >= -Infinite && alpha <= Infinite)
	// AssertTrue(beta >= -Infinite && beta <= Infinite)
	// AssertTrue(pos.ply >= 0 && pos.ply < MaxDepth)

	if HashTable.pTable[index].posKey == pos.posKey {
		*move = HashTable.pTable[index].move
		if HashTable.pTable[index].depth >= depth {

			// AssertTrue(HashTable.pTable[index].depth >= 1 && HashTable.pTable[index].depth < MaxDepth)
			// AssertTrue(HashTable.pTable[index].flags >= HFAlpha && HashTable.pTable[index].flags <= HFExact)

			*score = HashTable.pTable[index].score
			if *score > IsMate {
				*score -= pos.ply
			} else if *score < -IsMate {
				*score += pos.ply
			}

			switch HashTable.pTable[index].flags {
			case HFAlpha:
				if *score <= alpha {
					*score = alpha
					return true
				}
			case HFBeta:
				if *score >= beta {
					*score = beta
					return true
				}
			case HFExact:
				return true
			}
		}
	}

	return false
}

// GetPvLine get the principle variation line i.e. the best line found during the search algorithm
func GetPvLine(pos *Board, depth int) int {
	// // AssertTrue(depth < MaxDepth)

	move := ProbePvMove(pos)
	count := 0

	for move != NoMove && count < depth {
		// // AssertTrue(count < MaxDepth)

		if MoveExists(pos, move) {
			MakeMove(pos, move)
			pos.PvArray[count] = move
			count++
		} else {
			break
		}
		move = ProbePvMove(pos)
	}

	// take back all the moves we have made while probing the pv table so that we dont interfere with the position
	for pos.ply > 0 {
		TakeMove(pos)
	}

	return count
}

// ProbePvMove probe for a pv move
func ProbePvMove(pos *Board) int {
	index := int(pos.posKey % uint64(HashTable.numEntries))

	// AssertTrue(index >= 0 && index <= uint64(pos.HashTable.numEntries-1))

	if HashTable.pTable[index].posKey == pos.posKey {
		return HashTable.pTable[index].move
	}

	return NoMove
}
