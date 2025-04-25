package main

import (
	"encoding/binary"
	"log"
)

// Node Format: (B == bytes)
// | type | nkeys |  pointers  |  offsets   | key-values | unused |
// |  2B  |   2B  | nkeys × 8B | nkeys × 2B |     ...    |        |

type BNode []byte

// getters
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

// setter
func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// Get pointer based on its index
// Pointers are stored as array of 8bytes
func (node BNode) getPointer(idx uint16) uint64 {
	if idx >= node.nkeys() {
		log.Fatal("Accessing out of range pointer")
	}
	// 4 represents header bytes, 8 represents each pointer byte size
	pos := 4 + 8*idx

	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPointer(idx uint16, val uint64) {
	if idx >= node.nkeys() {
		log.Fatal("Accessing out of range pointer")
	}

	pos := 4 + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// Get offset based on idx
func (node BNode) getOffset(idx uint16) uint16 {
	if idx >= node.nkeys() {
		log.Fatal("Accessing out of range offset")
	}

	if idx == 0 {
		return 0
	}

	// We use idx - 1 because we do not store the first offset because it's always 0
	pos := 4 + 8*node.nkeys() + 2*(idx-1)
	return binary.LittleEndian.Uint16(node[pos:])
}

func (node BNode) kvPos(idx uint16) uint16 {
	if idx >= node.nkeys() {
		log.Fatal("Accessing out of range position")
	}

	return 4 + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)

}
