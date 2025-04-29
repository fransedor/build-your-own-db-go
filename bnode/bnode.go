package bnode

import (
	"encoding/binary"
	"log"
)

type BNode []byte

// Getters
func (node BNode) Btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) Nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

// Setter
func (node BNode) SetHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// Get pointer based on its index
func (node BNode) GetPointer(idx uint16) uint64 {
	if idx >= node.Nkeys() {
		log.Fatal("Accessing out of range pointer")
	}
	pos := 4 + 8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) SetPointer(idx uint16, val uint64) {
	if idx >= node.Nkeys() {
		log.Fatal("Accessing out of range pointer")
	}
	pos := 4 + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// Get offset based on idx
func (node BNode) GetOffset(idx uint16) uint16 {
	if idx >= node.Nkeys() {
		log.Fatal("Accessing out of range offset")
	}
	if idx == 0 {
		return 0
	}
	pos := 4 + 8*node.Nkeys() + 2*(idx-1)
	return binary.LittleEndian.Uint16(node[pos:])
}

func (node BNode) SetOffset(idx uint16, val uint16) {
	if idx == 0 {
		return
	}
	pos := 4 + 8*node.Nkeys() + 2*(idx-1)
	binary.LittleEndian.PutUint16(node[pos:], val)
}

func (node BNode) KvPos(idx uint16) uint16 {
	if idx >= node.Nkeys() {
		log.Fatal("Accessing out of range position")
	}
	return 4 + 8*node.Nkeys() + 2*node.Nkeys() + node.GetOffset(idx)
}

func (node BNode) GetKey(idx uint16) []byte {
	if idx >= node.Nkeys() {
		log.Fatal("getting out of range key")
	}
	pos := node.KvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:klen]
}

func (node BNode) GetVal(idx uint16) []byte {
	if idx >= node.Nkeys() {
		log.Fatal("getting val from out of range")
	}
	kvPos := node.KvPos(idx)
	klen := binary.LittleEndian.Uint16(node[kvPos:])
	vlen := binary.LittleEndian.Uint16(node[kvPos+2:])
	return node[kvPos+4+klen:][:vlen]
}

func (node BNode) AppendKV(idx uint16, ptr uint64, key []byte, val []byte) {
	node.SetPointer(idx, ptr)
	pos := node.KvPos(idx)

	// Setup KV sizes
	binary.LittleEndian.PutUint16(node[pos:], uint16(len(key)))
	binary.LittleEndian.PutUint16(node[pos+2:], uint16(len(val)))

	// Setup KV data
	copy(node[pos+4:], key)
	copy(node[pos+4+uint16(len(key)):], val)

	// update the offset value
	node.SetOffset(idx+1, node.GetOffset(idx)+4+uint16(len(key)+len(val)))
}

func (node BNode) Nbytes() uint16 {
	return node.KvPos(node.Nkeys())
}
