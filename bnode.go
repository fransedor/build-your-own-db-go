package main

import (
	"encoding/binary"
	"log"
)

// Node Format: (B == bytes)
// | type | nkeys |  pointers  |  offsets   | key-values | unused |
// |  2B  |   2B  | nkeys × 8B | nkeys × 2B |     ...    |        |

//	key-values
//	| key_size | val_size | key | val |
//	|    2B    |    2B    | ... | ... |

//For example, a leaf node {"k1":"hi", "k3":"hello"} is encoded as:

// | type | nkeys | pointers | offsets |            key-values           | unused |
// |   2  |   2   | nil nil  |  8 19   | 2 2 "k1" "hi"  2 5 "k3" "hello" |        |
// |  2B  |  2B   |   2×8B   |  2×2B   | 4B + 2B + 2B + 4B + 2B + 5B     |        |
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

func (node BNode) getKey(idx uint16) []byte {
	if idx >= node.nkeys() {
		log.Fatal("getting out of range key")
	}

	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])

	return node[pos+4:][:klen]
}

func (node BNode) getVal(idx uint16) []byte {
	if idx >= node.nkeys() {
		log.Fatal("getting val from out of range")
	}

	kvPos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[kvPos:])
	vlen := binary.LittleEndian.Uint16(node[kvPos+2:])
	return node[kvPos+4+klen:][:vlen]
}
